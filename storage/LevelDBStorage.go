package storage

import (
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/trusch/pkid/types"
)

type LevelDBStorage struct {
	db *leveldb.DB
}

func NewLevelDBStorage(baseDir string) (*LevelDBStorage, error) {
	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}
	db, err := leveldb.OpenFile(baseDir, o)
	if err != nil {
		return nil, err
	}
	store := &LevelDBStorage{db}
	return store, nil
}

func (store *LevelDBStorage) GetID() string {
	u1 := uuid.NewV4()
	return u1.String()
}

func (store *LevelDBStorage) SaveCA(ca *types.CAEntity) error {
	bs, err := json.Marshal(ca)
	if err != nil {
		return err
	}
	key := []byte(fmt.Sprintf("ca::%v", ca.ID))
	return store.db.Put(key, bs, nil)
}

func (store *LevelDBStorage) SaveClient(client *types.Entity) error {
	bs, err := json.Marshal(client)
	if err != nil {
		return err
	}
	key := []byte(fmt.Sprintf("client::%v", client.ID))
	return store.db.Put(key, bs, nil)
}

func (store *LevelDBStorage) SaveServer(server *types.Entity) error {
	bs, err := json.Marshal(server)
	if err != nil {
		return err
	}
	key := []byte(fmt.Sprintf("server::%v", server.ID))
	return store.db.Put(key, bs, nil)
}

func (store *LevelDBStorage) LoadCA(id string) (*types.CAEntity, error) {
	key := []byte(fmt.Sprintf("ca::%v", id))
	bs, err := store.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	res := &types.CAEntity{}
	return res, json.Unmarshal(bs, res)
}

func (store *LevelDBStorage) LoadClient(clientID string) (*types.Entity, error) {
	key := []byte(fmt.Sprintf("client::%v", clientID))
	bs, err := store.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	res := &types.Entity{}
	return res, json.Unmarshal(bs, res)
}

func (store *LevelDBStorage) LoadServer(serverID string) (*types.Entity, error) {
	key := []byte(fmt.Sprintf("server::%v", serverID))
	bs, err := store.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	res := &types.Entity{}
	return res, json.Unmarshal(bs, res)
}
