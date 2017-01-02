package manager

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/trusch/pkid/entity"
	"github.com/trusch/pkid/generator"
	"github.com/trusch/pkid/storage"
	"github.com/trusch/pkid/types"
)

type BasicManager struct {
	store storage.Storage
}

func NewBasicManager(store storage.Storage) Manager {
	return &BasicManager{store}
}

func (mgr *BasicManager) GetCA(id string) (*types.CAEntity, error) {
	return mgr.store.LoadCA(id)
}

func (mgr *BasicManager) GetClient(id string) (*types.Entity, error) {
	return mgr.store.LoadClient(id)
}

func (mgr *BasicManager) GetServer(id string) (*types.Entity, error) {
	return mgr.store.LoadServer(id)
}

func (mgr *BasicManager) CreateCA(caID string, options *generator.Options) (string, error) {
	ca, _ := mgr.store.LoadCA(caID)
	options.IsCA = true
	entity, err := generator.Generate(ca, options)
	if err != nil {
		return "", err
	}
	entity.ID = mgr.store.GetID()

	newCaEntity := &types.CAEntity{Entity: entity, Serial: big.NewInt(1)}
	err = mgr.store.SaveCA(newCaEntity)
	if err != nil {
		return "", err
	}
	if ca != nil {
		ca.Serial.Add(ca.Serial, big.NewInt(1))
		if ca.CAs == nil {
			ca.CAs = make(map[string]string)
		}
		ca.CAs[newCaEntity.ID] = newCaEntity.Name
		err = mgr.store.SaveCA(ca)
		if err != nil {
			return "", err
		}
	}
	return newCaEntity.ID, nil
}

func (mgr *BasicManager) CreateClient(caID string, options *generator.Options) (string, error) {
	ca, _ := mgr.store.LoadCA(caID)
	options.Usage = x509.ExtKeyUsageClientAuth
	entity, err := generator.Generate(ca, options)
	if err != nil {
		return "", err
	}
	entity.ID = mgr.store.GetID()
	err = mgr.store.SaveClient(entity)
	if err != nil {
		return "", err
	}
	if ca != nil {
		ca.Serial.Add(ca.Serial, big.NewInt(1))
		if ca.Clients == nil {
			ca.Clients = make(map[string]string)
		}
		ca.Clients[entity.ID] = entity.Name
		err = mgr.store.SaveCA(ca)
		if err != nil {
			return "", err
		}
	}
	return entity.ID, nil
}

func (mgr *BasicManager) CreateServer(caID string, options *generator.Options) (string, error) {
	ca, _ := mgr.store.LoadCA(caID)
	options.Usage = x509.ExtKeyUsageServerAuth
	entity, err := generator.Generate(ca, options)
	if err != nil {
		return "", err
	}
	entity.ID = mgr.store.GetID()
	err = mgr.store.SaveServer(entity)
	if err != nil {
		return "", err
	}
	if ca != nil {
		ca.Serial.Add(ca.Serial, big.NewInt(1))
		if ca.Servers == nil {
			ca.Servers = make(map[string]string)
		}
		ca.Servers[entity.ID] = entity.Name
		err = mgr.store.SaveCA(ca)
		if err != nil {
			return "", err
		}
	}
	return entity.ID, nil
}

func (mgr *BasicManager) RevokeCA(caID, id string) error {
	ca, err := mgr.GetCA(caID)
	if err != nil {
		return err
	}
	subCa, err := mgr.GetCA(id)
	if err != nil {
		return err
	}
	serial, err := mgr.getSerialFromEntity(subCa.Entity)
	if err != nil {
		return err
	}
	subCa.IsRevoked = true
	err = mgr.store.SaveCA(subCa)
	if err != nil {
		return err
	}
	ca.Revoked = append(ca.Revoked, serial)
	return mgr.store.SaveCA(ca)
}

func (mgr *BasicManager) RevokeClient(caID, id string) error {
	ca, err := mgr.GetCA(caID)
	if err != nil {
		return err
	}
	client, err := mgr.GetClient(id)
	if err != nil {
		return err
	}
	serial, err := mgr.getSerialFromEntity(client)
	if err != nil {
		return err
	}
	client.IsRevoked = true
	err = mgr.store.SaveClient(client)
	if err != nil {
		return err
	}
	ca.Revoked = append(ca.Revoked, serial)
	return mgr.store.SaveCA(ca)
}

func (mgr *BasicManager) RevokeServer(caID, id string) error {
	ca, err := mgr.GetCA(caID)
	if err != nil {
		return err
	}
	server, err := mgr.GetServer(id)
	if err != nil {
		return err
	}
	serial, err := mgr.getSerialFromEntity(server)
	if err != nil {
		return err
	}
	server.IsRevoked = true
	err = mgr.store.SaveServer(server)
	if err != nil {
		return err
	}
	ca.Revoked = append(ca.Revoked, serial)
	return mgr.store.SaveCA(ca)
}

func (mgr *BasicManager) GetCRL(caID string) (string, error) {
	caEntity, err := mgr.GetCA(caID)
	if err != nil {
		return "", err
	}
	revokedCerts := make([]pkix.RevokedCertificate, len(caEntity.Revoked))
	now := time.Now()
	for idx, serial := range caEntity.Revoked {
		revokedCerts[idx] = pkix.RevokedCertificate{
			SerialNumber:   serial,
			RevocationTime: now,
		}
	}
	entity, err := entity.NewEntityFromPEM([]byte(caEntity.Cert), []byte(caEntity.Key))
	if err != nil {
		return "", err
	}
	derCRL, err := entity.Cert.CreateCRL(rand.Reader, entity.Key, revokedCerts, now, now.Add(365*24*time.Hour))
	if err != nil {
		return "", err
	}
	out := &bytes.Buffer{}
	err = pem.Encode(out, &pem.Block{Type: "X509 CRL", Bytes: derCRL})
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func (mgr *BasicManager) getSerialFromEntity(e *types.Entity) (*big.Int, error) {
	parsed, err := entity.NewEntityFromPEM([]byte(e.Cert), []byte(e.Key))
	if err != nil {
		return nil, err
	}
	return parsed.Cert.SerialNumber, nil
}
