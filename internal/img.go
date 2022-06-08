package internal

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
)

type deleted int

// Structure
// If 	D
// Else H<hash 32><Comma seperated locations>
type ImgRecord struct {
	Deleted   deleted  `json:"deleted"`
	Hash      string   `json:"hash"`
	Locations []string `json:"locations"`
}

type hashVal struct {
	sum []byte
	srv string
}

type hashVals []hashVal

const (
	DELETED deleted = 0
	EXISTS  deleted = 1
)

// TODO: Remove this was only for testing
func (kv *KVImg) KeyToServer(key []byte) (string, []string) {
	return keyToPath(key), keyToVol(key, kv.GetServers(), kv.replicas)
}

func keyToPath(key []byte) string {
	mkey := md5.Sum(key)
	b64key := base64.StdEncoding.EncodeToString(key)
	return fmt.Sprintf("/%02x/%02x/%s", mkey[0], mkey[1], b64key)
}

func compareLocations(saved, generated []string) bool {
	if len(saved) != len(generated) {
		return false
	}

	for i := range saved {
		if saved[i] != generated[i] {
			return false
		}
	}

	return true
}

func keyToVol(key []byte, servers map[string]interface{}, count int) []string {
	vals := make(hashVals, len(servers))
	index := 0
	for server := range servers {
		hash := md5.New()
		hash.Write(key)
		hash.Write([]byte(server))
		vals[index] = hashVal{srv: server, sum: hash.Sum(nil)}
		index++
	}
	sort.Stable(vals)
	vols := make([]string, count)
	for i, val := range vals {
		if i == count {
			break
		}
		vols[i] = val.srv
	}

	return vols
}

func (img *ImgRecord) ToBytes() []byte {
	res := ""
	if len(img.Hash) == 32 {
		res = "H" + img.Hash
	}

	return []byte(res + strings.Join(img.Locations, ","))
}

func parseImgRecordData(dest *ImgRecord, data []byte) {
	s := string(data)
	// Start off with a deleted state
	dest.Deleted = DELETED
	if s[0:1] == "H" {
		dest.Hash = s[1:33]
		dest.Deleted = EXISTS
		s = s[33:]
	}
	dest.Locations = strings.Split(s, ",")
}

func (hv hashVals) Swap(i, j int) {
	hv[i], hv[j] = hv[j], hv[i]
}

func (hv hashVals) Len() int {
	return len(hv)
}

func (hv hashVals) Less(i, j int) bool {
	return bytes.Compare(hv[i].sum, hv[j].sum) == 1
}
