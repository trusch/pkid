package storage

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
	"github.com/trusch/pkid/types"
	"github.com/trusch/storage"
	"github.com/trusch/storage/engines/meta"
)

//StorageImpl is an implementation of the Storage interface
type StorageImpl struct {
	store storage.Storage
}

const (
	clientBucket string = "pkid-clients"
	caBucket            = "pkid-cas"
	serverBucket        = "pkid-servers"
)

// New returnes a new pki storage using github.com/trusch/storage
func New(uri string, token ...string) (*StorageImpl, error) {
	t := ""
	if len(token) > 0 {
		t = token[0]
	}
	store, err := meta.NewStorage(uri, t)
	if err != nil {
		return nil, err
	}
	if err = store.CreateBucket(clientBucket); err != nil {
		return nil, err
	}
	if err = store.CreateBucket(serverBucket); err != nil {
		return nil, err
	}
	if err = store.CreateBucket(caBucket); err != nil {
		return nil, err
	}
	return &StorageImpl{store}, nil
}

// GetID returns a new uuid
func (s *StorageImpl) GetID() string {
	u1 := uuid.NewV4()
	return u1.String()
}

// SaveCA saves a CA to backend
func (s *StorageImpl) SaveCA(ca *types.CAEntity) error {
	bs, err := json.Marshal(ca)
	if err != nil {
		return err
	}
	return s.store.Put(caBucket, ca.ID, bs)
}

// SaveClient saves a client to backend
func (s *StorageImpl) SaveClient(client *types.Entity) error {
	bs, err := json.Marshal(client)
	if err != nil {
		return err
	}
	return s.store.Put(clientBucket, client.ID, bs)
}

// SaveServer saves a Server to backend
func (s *StorageImpl) SaveServer(server *types.Entity) error {
	bs, err := json.Marshal(server)
	if err != nil {
		return err
	}
	return s.store.Put(serverBucket, server.ID, bs)
}

// LoadCA loads a CA from backend
func (s *StorageImpl) LoadCA(id string) (*types.CAEntity, error) {
	bs, err := s.store.Get(caBucket, id)
	if err != nil {
		return nil, err
	}
	entity := &types.CAEntity{}
	err = json.Unmarshal(bs, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// LoadClient loads a Client from backend
func (s *StorageImpl) LoadClient(clientID string) (*types.Entity, error) {
	bs, err := s.store.Get(clientBucket, clientID)
	if err != nil {
		return nil, err
	}
	entity := &types.Entity{}
	err = json.Unmarshal(bs, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// LoadServer loads a Server from backend
func (s *StorageImpl) LoadServer(serverID string) (*types.Entity, error) {
	bs, err := s.store.Get(serverBucket, serverID)
	if err != nil {
		return nil, err
	}
	entity := &types.Entity{}
	err = json.Unmarshal(bs, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
