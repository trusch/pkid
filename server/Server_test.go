package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trusch/pkid/manager"
	"github.com/trusch/pkid/storage"
	"github.com/trusch/pkid/types"
)

type ServerSuite struct {
	suite.Suite
	srv *Server
}

func (suite *ServerSuite) SetupTest() {
	store, err := storage.NewMetaStorage("leveldb://test-store")
	suite.NoError(err)
	suite.NotEmpty(store)
	mgr := manager.NewThreadSafeManager(store)
	suite.NotEmpty(mgr)
	suite.srv = New(":8080", mgr)
	go suite.srv.ListenAndServe()
}

func (suite *ServerSuite) TearDownTest() {
	suite.srv.Stop()
	os.RemoveAll("test-store")
}

func (suite *ServerSuite) request(method, path string) (string, error) {
	var (
		resp *http.Response
		err  error
	)
	switch method {
	case "GET":
		resp, err = http.Get(fmt.Sprintf("http://localhost:8080%v", path))
	case "POST":
		resp, err = http.Post(fmt.Sprintf("http://localhost:8080%v", path), "", nil)

	}
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return string(body), fmt.Errorf("%v", resp.StatusCode)
	}
	return string(body), nil
}

func (suite *ServerSuite) TestCreateSelfSignedCA() {
	rootID, err := suite.request("POST", "/ca?name=root")
	suite.NoError(err)
	suite.NotEmpty(rootID)
	key, err := suite.request("GET", fmt.Sprintf("/ca/%v/key", rootID))
	suite.NoError(err)
	suite.NotEmpty(key)
	cert, err := suite.request("GET", fmt.Sprintf("/ca/%v/cert", rootID))
	suite.NoError(err)
	suite.NotEmpty(cert)
}

func (suite *ServerSuite) TestList() {
	rootID, err := suite.request("POST", "/ca?name=root")
	suite.NoError(err)
	suite.NotEmpty(rootID)
	resp, err := suite.request("GET", fmt.Sprintf("/ca/%v/ca", rootID))
	suite.NoError(err)
	data := map[string]string{}
	err = json.Unmarshal([]byte(resp), &data)
	suite.NoError(err)
	suite.Equal(0, len(data))
	resp, err = suite.request("GET", fmt.Sprintf("/ca/%v/client", rootID))
	suite.NoError(err)
	data = map[string]string{}
	err = json.Unmarshal([]byte(resp), &data)
	suite.NoError(err)
	data = map[string]string{}
	suite.Equal(0, len(data))
	resp, err = suite.request("GET", fmt.Sprintf("/ca/%v/server", rootID))
	suite.NoError(err)
	data = map[string]string{}
	err = json.Unmarshal([]byte(resp), &data)
	suite.NoError(err)
	suite.Equal(0, len(data))
	subID, err := suite.request("POST", fmt.Sprintf("/ca/%v/ca?name=subca", rootID))
	suite.NoError(err)
	suite.NotEmpty(subID)
	subID, err = suite.request("POST", fmt.Sprintf("/ca/%v/client?name=subclient", rootID))
	suite.NoError(err)
	suite.NotEmpty(subID)
	subID, err = suite.request("POST", fmt.Sprintf("/ca/%v/server?name=subserver", rootID))
	suite.NoError(err)
	suite.NotEmpty(subID)
	resp, err = suite.request("GET", fmt.Sprintf("/ca/%v/ca", rootID))
	suite.NoError(err)
	data = map[string]string{}
	err = json.Unmarshal([]byte(resp), &data)
	suite.NoError(err)
	suite.Equal(1, len(data))
	log.Print(data)
	resp, err = suite.request("GET", fmt.Sprintf("/ca/%v/client", rootID))
	suite.NoError(err)
	data = map[string]string{}
	err = json.Unmarshal([]byte(resp), &data)
	suite.NoError(err)
	suite.Equal(1, len(data))
	log.Print(data)
	resp, err = suite.request("GET", fmt.Sprintf("/ca/%v/server", rootID))
	suite.NoError(err)
	data = map[string]string{}
	err = json.Unmarshal([]byte(resp), &data)
	suite.NoError(err)
	suite.Equal(1, len(data))
	log.Print(data)
}

func (suite *ServerSuite) TestCreateSignedCA() {
	rootID, err := suite.request("POST", "/ca?name=root")
	suite.NoError(err)
	suite.NotEmpty(rootID)
	subID, err := suite.request("POST", fmt.Sprintf("/ca/%v/ca?name=sub", rootID))
	suite.NoError(err)
	suite.NotEmpty(subID)
	key, err := suite.request("GET", fmt.Sprintf("/ca/%v/ca/%v/key", rootID, subID))
	suite.NoError(err)
	suite.NotEmpty(key)
	cert, err := suite.request("GET", fmt.Sprintf("/ca/%v/ca/%v/cert", rootID, subID))
	suite.NoError(err)
	suite.NotEmpty(cert)
}

func (suite *ServerSuite) TestCreateSignedClient() {
	rootID, err := suite.request("POST", "/ca?name=root")
	suite.NoError(err)
	suite.NotEmpty(rootID)
	subID, err := suite.request("POST", fmt.Sprintf("/ca/%v/client?name=sub", rootID))
	suite.NoError(err)
	suite.NotEmpty(subID)
	key, err := suite.request("GET", fmt.Sprintf("/ca/%v/client/%v/key", rootID, subID))
	suite.NoError(err)
	suite.NotEmpty(key)
	cert, err := suite.request("GET", fmt.Sprintf("/ca/%v/client/%v/cert", rootID, subID))
	suite.NoError(err)
	suite.NotEmpty(cert)
}

func (suite *ServerSuite) TestCreateSignedServer() {
	rootID, err := suite.request("POST", "/ca?name=root")
	suite.NoError(err)
	suite.NotEmpty(rootID)
	subID, err := suite.request("POST", fmt.Sprintf("/ca/%v/server?name=sub", rootID))
	suite.NoError(err)
	suite.NotEmpty(subID)
	key, err := suite.request("GET", fmt.Sprintf("/ca/%v/server/%v/key", rootID, subID))
	suite.NoError(err)
	suite.NotEmpty(key)
	cert, err := suite.request("GET", fmt.Sprintf("/ca/%v/server/%v/cert", rootID, subID))
	suite.NoError(err)
	suite.NotEmpty(cert)
}

func (suite *ServerSuite) TestGetCA() {
	rootID, err := suite.request("POST", "/ca?name=root")
	suite.NoError(err)
	suite.NotEmpty(rootID)
	caID, err := suite.request("POST", fmt.Sprintf("/ca/%v/ca?name=subca", rootID))
	suite.NoError(err)
	suite.NotEmpty(caID)
	clientID, err := suite.request("POST", fmt.Sprintf("/ca/%v/client?name=subclient", rootID))
	suite.NoError(err)
	suite.NotEmpty(clientID)
	serverID, err := suite.request("POST", fmt.Sprintf("/ca/%v/server?name=subserver", rootID))
	suite.NoError(err)
	suite.NotEmpty(serverID)
	resp, err := suite.request("GET", fmt.Sprintf("/ca/%v", rootID))
	suite.NoError(err)
	suite.NotEmpty(resp)
	ca := &types.CAEntity{}
	err = json.Unmarshal([]byte(resp), ca)
	suite.NoError(err)
	suite.Empty(ca.Key)
	suite.Empty(ca.Cert)
	suite.Equal(1, len(ca.CAs))
	suite.Equal(1, len(ca.Clients))
	suite.Equal(1, len(ca.Servers))
	suite.Equal(0, len(ca.Revoked))
}

func (suite *ServerSuite) TestRevoke() {
	rootID, err := suite.request("POST", "/ca?name=root")
	suite.NoError(err)
	suite.NotEmpty(rootID)
	caID, err := suite.request("POST", fmt.Sprintf("/ca/%v/ca?name=subca", rootID))
	suite.NoError(err)
	suite.NotEmpty(caID)
	_, err = suite.request("POST", fmt.Sprintf("/ca/%v/ca/%v/revoke", rootID, caID))
	suite.NoError(err)
	resp, err := suite.request("GET", fmt.Sprintf("/ca/%v", rootID))
	suite.NoError(err)
	suite.NotEmpty(resp)
	ca := &types.CAEntity{}
	err = json.Unmarshal([]byte(resp), ca)
	suite.NoError(err)
	suite.Equal(1, len(ca.Revoked))
}

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
