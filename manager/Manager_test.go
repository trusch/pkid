package manager

import (
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trusch/pkid/generator"
	"github.com/trusch/pkid/storage"
)

type ManagerSuite struct {
	suite.Suite
	manager Manager
}

func NewManagerSuite(manager Manager) *ManagerSuite {
	res := &ManagerSuite{}
	res.manager = manager
	return res
}

func (suite *ManagerSuite) TearDownTest() {
	os.RemoveAll("./test-store")
}

func (suite *ManagerSuite) TestCreateSelfSigned() {
	caID, err := suite.manager.CreateCA("", &generator.Options{Name: "my-ca"})
	suite.NoError(err)
	suite.NotEmpty(caID)
	serverID, err := suite.manager.CreateServer("", &generator.Options{Name: "my-server"})
	suite.NoError(err)
	suite.NotEmpty(serverID)
	clientID, err := suite.manager.CreateClient("", &generator.Options{Name: "my-client"})
	suite.NoError(err)
	suite.NotEmpty(clientID)
}

func (suite *ManagerSuite) TestCreateSigned() {
	rootCaID, err := suite.manager.CreateCA("", &generator.Options{Name: "root-ca"})
	suite.NoError(err)
	suite.NotEmpty(rootCaID)
	caID, err := suite.manager.CreateCA(rootCaID, &generator.Options{Name: "my-ca"})
	suite.NoError(err)
	suite.NotEmpty(caID)
	serverID, err := suite.manager.CreateServer(rootCaID, &generator.Options{Name: "my-server"})
	suite.NoError(err)
	suite.NotEmpty(serverID)
	clientID, err := suite.manager.CreateClient(rootCaID, &generator.Options{Name: "my-client"})
	suite.NoError(err)
	suite.NotEmpty(clientID)
}

func (suite *ManagerSuite) TestCreateAndGet() {
	rootCaID, err := suite.manager.CreateCA("", &generator.Options{Name: "root-ca"})
	suite.NoError(err)
	suite.NotEmpty(rootCaID)
	caID, err := suite.manager.CreateCA(rootCaID, &generator.Options{Name: "my-ca"})
	suite.NoError(err)
	suite.NotEmpty(caID)
	serverID, err := suite.manager.CreateServer(rootCaID, &generator.Options{Name: "my-server"})
	suite.NoError(err)
	suite.NotEmpty(serverID)
	clientID, err := suite.manager.CreateClient(rootCaID, &generator.Options{Name: "my-client"})
	suite.NoError(err)
	suite.NotEmpty(clientID)
	ca, err := suite.manager.GetCA(caID)
	suite.NoError(err)
	suite.NotNil(ca)
	client, err := suite.manager.GetClient(clientID)
	suite.NoError(err)
	suite.NotNil(client)
	server, err := suite.manager.GetServer(serverID)
	suite.NoError(err)
	suite.NotNil(server)
}

func (suite *ManagerSuite) TestRevoke() {
	rootCaID, err := suite.manager.CreateCA("", &generator.Options{Name: "root-ca"})
	suite.NoError(err)
	suite.NotEmpty(rootCaID)
	caID, err := suite.manager.CreateCA(rootCaID, &generator.Options{Name: "my-ca"})
	suite.NoError(err)
	suite.NotEmpty(caID)
	serverID, err := suite.manager.CreateServer(rootCaID, &generator.Options{Name: "my-server"})
	suite.NoError(err)
	suite.NotEmpty(serverID)
	clientID, err := suite.manager.CreateClient(rootCaID, &generator.Options{Name: "my-client"})
	suite.NoError(err)
	suite.NotEmpty(clientID)

	err = suite.manager.RevokeCA(rootCaID, caID)
	suite.NoError(err)
	err = suite.manager.RevokeServer(rootCaID, serverID)
	suite.NoError(err)
	err = suite.manager.RevokeClient(rootCaID, clientID)
	suite.NoError(err)

	ca, err := suite.manager.GetCA(rootCaID)
	suite.NoError(err)
	suite.Equal([]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, ca.Revoked)

	crl, err := suite.manager.GetCRL(rootCaID)
	suite.NoError(err)
	suite.NotEmpty(crl)
}

func TestBasicManager(t *testing.T) {
	store, _ := storage.NewFSStorage("./test-store")
	mgr := NewBasicManager(store)
	suite.Run(t, NewManagerSuite(mgr))
}

func TestThreadSafeManager(t *testing.T) {
	store, _ := storage.NewFSStorage("./test-store")
	mgr := NewThreadSafeManager(store)
	suite.Run(t, NewManagerSuite(mgr))
}
