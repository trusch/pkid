package manager

import (
	"github.com/trusch/pkid/generator"
	"github.com/trusch/pkid/types"
)

type Manager interface {
	GetCA(id string) (*types.CAEntity, error)
	GetClient(id string) (*types.Entity, error)
	GetServer(id string) (*types.Entity, error)
	CreateCA(caID string, options *generator.Options) (string, error)
	CreateClient(caID string, options *generator.Options) (string, error)
	CreateServer(caID string, options *generator.Options) (string, error)
	RevokeCA(caID, id string) error
	RevokeClient(caID, id string) error
	RevokeServer(caID, id string) error
	GetCRL(caID string) (string, error)
}
