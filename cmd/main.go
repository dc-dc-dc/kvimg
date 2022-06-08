package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dc-dc-dc/KVImg/internal"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	servers := []string{"http://192.168.1.231:3000"}
	db, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	kvImg := internal.NewKVImg(servers, db)
	log.SetFlags(log.Ldate | log.LstdFlags)
	if err := testServer(kvImg); err != nil {
		panic(err)
	}
}

func testServer(kv *internal.KVImg) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", 3000), kv)
}

func testLocalFile(KVImg *internal.KVImg) {
	key := []byte("gitignore")
	log.Print("finding")
	rec, _, err := KVImg.GetFile(key)
	if err != nil {
		panic(err)
	}
	log.Printf("%d, %s, %v", rec.Deleted, rec.Hash, rec.Locations)
	log.Print("deleting")
	if err = KVImg.DeleteFile(key); err != nil {
		log.Print(err.Error())
	}
	file, err := os.Open("./.gitignore")
	if err != nil {
		panic(err)
	}
	fileinfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	data := make([]byte, fileinfo.Size())
	_, err = file.Read(data)
	if err != nil {
		panic(err)
	}
	log.Print("uploading")

	if err = KVImg.UploadFile(key, bytes.NewReader(data), fileinfo.Size()); err != nil {
		panic(err)
	}
}
