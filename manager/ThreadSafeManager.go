package manager

import (
	"github.com/trusch/pkid/generator"
	"github.com/trusch/pkid/storage"
	"github.com/trusch/pkid/types"
	"github.com/trusch/transaction"
)

type ThreadSafeManager struct {
	basic       *BasicManager
	transaction *transaction.Manager
}

func NewThreadSafeManager(store storage.Storage) Manager {
	mgr := &ThreadSafeManager{
		basic:       &BasicManager{store},
		transaction: transaction.NewManager(nil),
	}
	return mgr
}

func (mgr *ThreadSafeManager) GetCA(id string) (*types.CAEntity, error) {
	v, e := mgr.transaction.Transaction(func(context interface{}) (interface{}, error) {
		return mgr.basic.GetCA(id)
	})
	return v.(*types.CAEntity), e
}

func (mgr *ThreadSafeManager) GetClient(id string) (*types.Entity, error) {
	v, e := mgr.transaction.Transaction(func(context interface{}) (interface{}, error) {
		return mgr.basic.GetClient(id)
	})
	return v.(*types.Entity), e
}

func (mgr *ThreadSafeManager) GetServer(id string) (*types.Entity, error) {
	v, e := mgr.transaction.Transaction(func(context interface{}) (interface{}, error) {
		return mgr.basic.GetServer(id)
	})
	return v.(*types.Entity), e
}

func (mgr *ThreadSafeManager) CreateCA(caID string, options *generator.Options) (string, error) {
	v, e := mgr.transaction.Transaction(func(context interface{}) (interface{}, error) {
		return mgr.basic.CreateCA(caID, options)
	})
	return v.(string), e
}

func (mgr *ThreadSafeManager) CreateClient(caID string, options *generator.Options) (string, error) {
	v, e := mgr.transaction.Transaction(func(context interface{}) (interface{}, error) {
		return mgr.basic.CreateClient(caID, options)
	})
	return v.(string), e
}

func (mgr *ThreadSafeManager) CreateServer(caID string, options *generator.Options) (string, error) {
	v, e := mgr.transaction.Transaction(func(context interface{}) (interface{}, error) {
		return mgr.basic.CreateServer(caID, options)
	})
	return v.(string), e
}

func (mgr *ThreadSafeManager) RevokeCA(caID, id string) error {
	_, e := mgr.transaction.Transaction(func(context interface{}) (interface{}, error) {
		return nil, mgr.basic.RevokeCA(caID, id)
	})
	return e
}

func (mgr *ThreadSafeManager) RevokeClient(caID, id string) error {
	_, e := mgr.transaction.Transaction(func(context interface{}) (interface{}, error) {
		return nil, mgr.basic.RevokeClient(caID, id)
	})
	return e
}

func (mgr *ThreadSafeManager) RevokeServer(caID, id string) error {
	_, e := mgr.transaction.Transaction(func(context interface{}) (interface{}, error) {
		return nil, mgr.basic.RevokeServer(caID, id)
	})
	return e
}

func (mgr *ThreadSafeManager) GetCRL(caID string) (string, error) {
	v, e := mgr.transaction.Transaction(func(context interface{}) (interface{}, error) {
		return mgr.basic.GetCRL(caID)
	})
	return v.(string), e
}
