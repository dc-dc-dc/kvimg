package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dc-dc-dc/KVImg/internal"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	action  = flag.String("action", "server", "<server|rebuild> what action to preform")
	dbFile  = flag.String("db", "./db", "location of db file")
	port    = flag.Int("port", 3000, "the port to bind too")
	servers = flag.String("servers", "", "comma seperated webdav servers")
)

func init() {
	log.SetFlags(log.LstdFlags)
	flag.Parse()
}

func main() {
	if *action != "server" && *action != "rebuild" {
		flag.PrintDefaults()
		os.Exit(2)
	}
	if *servers == "" {
		log.Printf("servers param is required, got: %s", *servers)
		flag.PrintDefaults()
		os.Exit(2)
	}

	db, err := leveldb.OpenFile(*dbFile, nil)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	kvImg := internal.NewKVImg(strings.Split(*servers, ","), db)
	if err := server(kvImg, *port); err != nil {
		panic(err)
	}
}

func server(kv *internal.KVImg, port int) error {
	log.Printf("starting server on port %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), kv)
}

// func testLocalFile(KVImg *internal.KVImg) {
// 	key := []byte("gitignore")
// 	log.Print("finding")
// 	rec, _, err := KVImg.GetFile(key)
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Printf("%d, %s, %v", rec.Deleted, rec.Hash, rec.Locations)
// 	log.Print("deleting")
// 	if err = KVImg.DeleteFile(key); err != nil {
// 		log.Print(err.Error())
// 	}
// 	file, err := os.Open("./.gitignore")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fileinfo, err := file.Stat()
// 	if err != nil {
// 		panic(err)
// 	}
// 	data := make([]byte, fileinfo.Size())
// 	_, err = file.Read(data)
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Print("uploading")

// 	if err = KVImg.UploadFile(key, bytes.NewReader(data), fileinfo.Size()); err != nil {
// 		panic(err)
// 	}
// }
