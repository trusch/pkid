package entity

import (
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityRSA(t *testing.T) {
	entity, err := NewEntityFromFile("test-rsa.crt", "test-rsa.key")
	assert.Nil(t, err)
	certPem, err := entity.GetCertAsPEM()
	assert.Nil(t, err)
	keyPem, err := entity.GetKeyAsPEM()
	assert.Nil(t, err)
	ioutil.WriteFile("test-rsa-copy.crt", certPem, 0755)
	ioutil.WriteFile("test-rsa-copy.key", keyPem, 0755)
	test1 := exec.Command("diff", "test-rsa.crt", "test-rsa-copy.crt")
	err = test1.Run()
	assert.Nil(t, err)
	test2 := exec.Command("diff", "test-rsa.crt", "test-rsa-copy.crt")
	err = test2.Run()
	assert.Nil(t, err)
}

func TestEntityEC(t *testing.T) {
	entity, err := NewEntityFromFile("test-ec.crt", "test-ec.key")
	assert.Nil(t, err)
	certPem, err := entity.GetCertAsPEM()
	assert.Nil(t, err)
	keyPem, err := entity.GetKeyAsPEM()
	assert.Nil(t, err)
	ioutil.WriteFile("test-ec-copy.crt", certPem, 0755)
	ioutil.WriteFile("test-ec-copy.key", keyPem, 0755)
	test1 := exec.Command("diff", "test-ec.crt", "test-ec-copy.crt")
	err = test1.Run()
	assert.Nil(t, err)
	test2 := exec.Command("diff", "test-ec.crt", "test-ec-copy.crt")
	err = test2.Run()
	assert.Nil(t, err)
}
