package main

import (
	"flag"
	"log"

	"github.com/trusch/pkid/manager"
	"github.com/trusch/pkid/server"
	"github.com/trusch/pkid/storage"
)

var storagePath = flag.String("storage", "leveldb:///usr/share/pkid/datastore", "storage backend uri")
var listenAddr = flag.String("listen", ":80", "listen address")
var token = flag.String("token", "", "bearer authorization token for secure storaged backend")

func main() {
	flag.Parse()
	store, err := storage.New(*storagePath, *token)
	if err != nil {
		log.Fatal(err)
	}
	mgr := manager.NewThreadSafeManager(store)
	srv := server.New(*listenAddr, mgr)
	log.Fatal(srv.ListenAndServe())
}
