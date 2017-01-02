package storage

import (
	"errors"
	"net/url"

	"github.com/trusch/pkid/types"
)

type MetaStorage struct {
	base Storage
}

func NewMetaStorage(uriStr string) (*MetaStorage, error) {
	uri, err := url.Parse(uriStr)
	if err != nil {
		return nil, err
	}
	var base Storage
	switch uri.Scheme {
	case "file":
		base, err = NewFSStorage(uri.Host + uri.Path)
	case "leveldb":
		base, err = NewLevelDBStorage(uri.Host + uri.Path)
	default:
		err = errors.New("unknown uri scheme, try file:// or leveldb://")
	}
	if err != nil {
		return nil, err
	}
	return &MetaStorage{base}, nil
}

func (store *MetaStorage) GetID() string {
	return store.base.GetID()
}
func (store *MetaStorage) SaveCA(ca *types.CAEntity) error {
	return store.base.SaveCA(ca)
}
func (store *MetaStorage) SaveClient(client *types.Entity) error {
	return store.base.SaveClient(client)
}
func (store *MetaStorage) SaveServer(server *types.Entity) error {
	return store.base.SaveServer(server)
}
func (store *MetaStorage) LoadCA(id string) (*types.CAEntity, error) {
	return store.base.LoadCA(id)
}
func (store *MetaStorage) LoadClient(clientID string) (*types.Entity, error) {
	return store.base.LoadClient(clientID)
}
func (store *MetaStorage) LoadServer(serverID string) (*types.Entity, error) {
	return store.base.LoadServer(serverID)
}
