package generator

import (
	"crypto/x509"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trusch/pkid/types"
)

func TestGenerateSelfSignedRSA(t *testing.T) {
	options := &Options{
		Name:    "my-cert",
		RsaBits: 2048,
	}
	entity, err := Generate(nil, options)
	assert.NotEmpty(t, entity)
	assert.NoError(t, err)
}

func TestGenerateSelfSignedECDSA(t *testing.T) {
	options := &Options{
		Name:  "my-cert",
		Curve: "P521",
	}
	entity, err := Generate(nil, options)
	assert.NotEmpty(t, entity)
	assert.NoError(t, err)
}

func TestGenerateSubFoo(t *testing.T) {
	options := &Options{
		Name: "my-ca",
		IsCA: true,
	}
	entity, err := Generate(nil, options)
	assert.NotEmpty(t, entity)
	assert.NoError(t, err)
	caEntity := &types.CAEntity{
		Entity: entity,
		Serial: big.NewInt(1),
	}
	options = &Options{
		Name:  "my-server",
		Usage: x509.ExtKeyUsageServerAuth,
	}
	serverEntity, err := Generate(caEntity, options)
	assert.NotEmpty(t, serverEntity)
	assert.NoError(t, err)
}
