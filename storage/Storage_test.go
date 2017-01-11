package storage

//func (suite *StorageSuite) TestStorage(t *testing.T) {

//}
// Basic imports
import (
	"crypto/x509"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trusch/pkid/generator"
	"github.com/trusch/pkid/types"
)

type StorageSuite struct {
	suite.Suite
	store Storage
}

func (suite *StorageSuite) TearDownSuite() {
	os.RemoveAll("./test-store.db")
}

func (suite *StorageSuite) TestSaveLoadCA() {
	entity, err := generator.Generate(nil, &generator.Options{Name: "test-ca", IsCA: true, RsaBits: 2048})
	suite.NoError(err)
	entity.ID = suite.store.GetID()
	caEntity := &types.CAEntity{
		Entity: entity,
		Serial: big.NewInt(1),
	}
	err = suite.store.SaveCA(caEntity)
	suite.NoError(err)
	restoredEntity, err := suite.store.LoadCA(caEntity.ID)
	suite.NoError(err)
	suite.Equal(caEntity.ID, restoredEntity.ID)
	suite.Equal(caEntity.Name, restoredEntity.Name)
	suite.Equal(caEntity.Cert, restoredEntity.Cert)
	suite.Equal(caEntity.Key, restoredEntity.Key)
	suite.Equal(caEntity.Serial, restoredEntity.Serial)
}

func (suite *StorageSuite) TestSaveLoadClient() {
	entity, err := generator.Generate(nil, &generator.Options{Name: "test-ca", IsCA: true})
	suite.NoError(err)
	entity.ID = suite.store.GetID()
	caEntity := &types.CAEntity{
		Entity: entity,
		Serial: big.NewInt(1),
	}
	entity, err = generator.Generate(caEntity, &generator.Options{Name: "test-client", Usage: x509.ExtKeyUsageClientAuth})
	suite.NoError(err)
	entity.ID = suite.store.GetID()
	err = suite.store.SaveClient(entity)
	suite.NoError(err)
	restoredEntity, err := suite.store.LoadClient(entity.ID)
	suite.NoError(err)
	suite.Equal(entity.ID, restoredEntity.ID)
	suite.Equal(entity.Name, restoredEntity.Name)
	suite.Equal(entity.Cert, restoredEntity.Cert)
	suite.Equal(entity.Key, restoredEntity.Key)
}

func (suite *StorageSuite) TestSaveLoadServer() {
	entity, err := generator.Generate(nil, &generator.Options{Name: "test-ca", IsCA: true})
	suite.NoError(err)
	entity.ID = suite.store.GetID()
	caEntity := &types.CAEntity{
		Entity: entity,
		Serial: big.NewInt(1),
	}
	entity, err = generator.Generate(caEntity, &generator.Options{Name: "test-server", Usage: x509.ExtKeyUsageServerAuth})
	suite.NoError(err)
	entity.ID = suite.store.GetID()
	err = suite.store.SaveServer(entity)
	suite.NoError(err)
	restoredEntity, err := suite.store.LoadServer(entity.ID)
	suite.NoError(err)
	suite.Equal(entity.ID, restoredEntity.ID)
	suite.Equal(entity.Name, restoredEntity.Name)
	suite.Equal(entity.Cert, restoredEntity.Cert)
	suite.Equal(entity.Key, restoredEntity.Key)
}

// func TestStorageImplWithLevelDB(t *testing.T) {
// 	store, err := New("leveldb://test-store.db")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, store)
// 	s := new(StorageSuite)
// 	s.store = store
// 	suite.Run(t, s)
// }

func TestStorageImplWithFile(t *testing.T) {
	store, err := New("file://test-store.db")
	assert.NoError(t, err)
	assert.NotNil(t, store)
	s := new(StorageSuite)
	s.store = store
	suite.Run(t, s)
}
