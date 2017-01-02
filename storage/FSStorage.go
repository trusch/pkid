package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/satori/go.uuid"
	"github.com/trusch/pkid/types"
)

type FSStorage struct {
	baseDir string
}

func NewFSStorage(baseDir string) (*FSStorage, error) {
	fs := &FSStorage{baseDir}
	return fs, fs.assertDirectory()
}

func (fs *FSStorage) GetID() string {
	u1 := uuid.NewV4()
	return u1.String()
}

func (fs *FSStorage) SaveCA(ca *types.CAEntity) error {
	if err := fs.assertDirectory("cas"); err != nil {
		return err
	}
	bs, err := yaml.Marshal(ca)
	if err != nil {
		return err
	}
	if err := fs.writeFile(bs, "cas", ca.ID+".yaml"); err != nil {
		return err
	}
	return nil
}

func (fs *FSStorage) SaveClient(client *types.Entity) error {
	if err := fs.assertDirectory("clients"); err != nil {
		return err
	}
	bs, err := yaml.Marshal(client)
	if err != nil {
		return err
	}
	if err := fs.writeFile(bs, "clients", client.ID+".yaml"); err != nil {
		return err
	}
	return nil
}

func (fs *FSStorage) SaveServer(server *types.Entity) error {
	if err := fs.assertDirectory("servers"); err != nil {
		return err
	}
	bs, err := yaml.Marshal(server)
	if err != nil {
		return err
	}
	if err := fs.writeFile(bs, "servers", server.ID+".yaml"); err != nil {
		return err
	}
	return nil
}

func (fs *FSStorage) LoadCA(id string) (*types.CAEntity, error) {
	ca := &types.CAEntity{}
	caBs, err := fs.loadFile("cas", id+".yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(caBs, ca)
	return ca, err
}

func (fs *FSStorage) LoadClient(clientID string) (*types.Entity, error) {
	client := &types.Entity{}
	bs, err := fs.loadFile("clients", clientID+".yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bs, client)
	return client, err
}

func (fs *FSStorage) LoadServer(serverID string) (*types.Entity, error) {
	server := &types.Entity{}
	bs, err := fs.loadFile("servers", serverID+".yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bs, server)
	return server, err
}

func (fs *FSStorage) assertDirectory(relDir ...string) error {
	args := append([]string{fs.baseDir}, relDir...)
	return os.MkdirAll(filepath.Join(args...), 0700)
}

func (fs *FSStorage) writeFile(content []byte, file ...string) error {
	args := append([]string{fs.baseDir}, file...)
	return ioutil.WriteFile(filepath.Join(args...), content, 0700)
}

func (fs *FSStorage) loadFile(file ...string) ([]byte, error) {
	args := append([]string{fs.baseDir}, file...)
	return ioutil.ReadFile(filepath.Join(args...))
}
