package internal

import "github.com/syndtr/goleveldb/leveldb"

type KVImg struct {
	servers map[string]interface{}
	db      *leveldb.DB
}

func NewKVImg(servers []string, db *leveldb.DB) *KVImg {
	_servers := make(map[string]interface{})
	for _, s := range servers {
		_servers[s] = true
	}

	return &KVImg{
		servers: _servers,
		db:      db,
	}
}

func (kv *KVImg) RemoveServer(server string) {
	delete(kv.servers, server)
}

func (kv *KVImg) AddServer(server string) {
	kv.servers[server] = true
}
