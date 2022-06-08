package internal

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
)

var (
	ErrImgDoesntExist = errors.New("no img record found")
	ErrImgRecExists   = errors.New("img record already exists")
	ErrUpdatingDB     = errors.New("error updating db")
	ErrDeletingDB     = errors.New("error deleting from")
	ErrNotMatching    = errors.New("volumes have differed")
	ErrVolMissingImg  = errors.New("volume is missing an img")
)

func getId(vol, path string) string {
	return fmt.Sprintf("%s%s", vol, path)
}

func (kv *KVImg) GetFile(key []byte) (*ImgRecord, int, error) {
	rec := kv.GetImgRecord(key)
	if rec.Deleted == DELETED {
		return rec, -1, ErrImgDoesntExist
	}

	// Check the vols
	if len(rec.Locations) > 0 {
		if !compareLocations(rec.Locations, keyToVol(key, kv.servers, kv.replicas)) {
			return rec, -1, ErrNotMatching
		}
	}

	// Check to see if the image exists on a vol
	// TODO: Check all of them ? Or keep a cache of last time checked
	path := keyToPath(key)
	found := -1

	for srv := range rand.Perm(len(rec.Locations)) {
		if _, err := send_get(getId(rec.Locations[srv], path)); err == nil {
			found = srv
			break
		}
	}

	if found == -1 {
		return rec, -1, ErrVolMissingImg
	}

	return rec, found, nil
}

func (kv *KVImg) UploadFile(key []byte, data io.Reader, length int64) error {
	rec := kv.GetImgRecord(key)
	if rec.Deleted == EXISTS {
		return ErrImgRecExists
	}
	path := keyToPath(key)
	vols := keyToVol(key, kv.GetServers(), kv.replicas)
	buf := &bytes.Buffer{}
	body := io.TeeReader(data, buf)
	for i, vol := range vols {
		if i != 0 {
			body = bytes.NewReader(buf.Bytes())
		}
		if _, err := send_put(getId(vol, path), body); err != nil {
			log.Printf("error while trying to upload to %s", vol)
		}
	}
	rec.Locations = vols
	rec.Hash = fmt.Sprintf("%x", md5.Sum(buf.Bytes()))

	if !kv.SaveImgRecord(key, rec) {
		return ErrUpdatingDB
	}

	log.Printf("uploaded to %s %v", path, vols)

	return nil
}

func (kv *KVImg) DeleteFile(key []byte) error {
	rec := kv.GetImgRecord(key)
	if rec.Deleted == DELETED {
		return ErrImgDoesntExist
	}

	path := keyToPath(key)

	for _, srv := range rec.Locations {
		if _, err := send_delete(getId(srv, path)); err != nil {
			fmt.Printf("error trying to delete from %s", srv)
		}

	}

	if !kv.DeleteImgRecord(key) {
		return ErrDeletingDB
	}

	return nil
}
