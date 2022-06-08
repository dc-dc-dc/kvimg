package internal

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type KVImg struct {
	servers  map[string]interface{}
	db       *leveldb.DB
	replicas int
}

func smallerInt(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func NewKVImg(servers []string, db *leveldb.DB) *KVImg {
	_servers := make(map[string]interface{})
	for _, s := range servers {
		_servers[s] = true
	}

	return &KVImg{
		servers:  _servers,
		db:       db,
		replicas: smallerInt(3, len(_servers)),
	}
}

func (kv *KVImg) GetServers() map[string]interface{} {
	return kv.servers
}

func (kv *KVImg) RemoveServer(server string) {
	delete(kv.servers, server)
}

func (kv *KVImg) AddServer(server string) {
	kv.servers[server] = true
}
