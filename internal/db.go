package internal

import "github.com/syndtr/goleveldb/leveldb"

func (kv *KVImg) GetImgRecord(key []byte) *ImgRecord {
	data, err := kv.db.Get(key, nil)
	rec := &ImgRecord{}
	if err != leveldb.ErrNotFound {
		parseImgRecordData(rec, data)
	}
	return rec
}

func (kv *KVImg) SaveImgRecord(key []byte, img *ImgRecord) bool {
	return kv.db.Put(key, img.ToBytes(), nil) == nil
}

func (kv *KVImg) DeleteImgRecord(key []byte) bool {
	return kv.db.Delete(key, nil) == nil
}
