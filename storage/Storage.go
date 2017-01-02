package storage

import "github.com/trusch/pkid/types"

// Storage Interface
type Storage interface {
	GetID() string
	SaveCA(ca *types.CAEntity) error
	SaveClient(client *types.Entity) error
	SaveServer(server *types.Entity) error
	LoadCA(id string) (*types.CAEntity, error)
	LoadClient(clientID string) (*types.Entity, error)
	LoadServer(serverID string) (*types.Entity, error)
}
