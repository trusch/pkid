package types

// EntityType is the storage entity type
import "math/big"

type EntityType int

// List of valid EntityType's
const (
	CA EntityType = iota
	Server
	Client
)

// Entity is ID with a Cert and a Key (both pem encoded)
type Entity struct {
	ID        string
	Name      string
	Cert      string
	Key       string
	IsRevoked bool
}

// A CAEntity is a Entity with a serial number (used for next issued cert)
type CAEntity struct {
	*Entity
	Serial  *big.Int
	Revoked []*big.Int
	Clients map[string]string
	Servers map[string]string
	CAs     map[string]string
}
