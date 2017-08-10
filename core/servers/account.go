package servers

import (
	"github.com/asaskevich/govalidator"
	"github.com/ellcrys/elldb/core/servers/common"
	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/ellcrys/elldb/core/servers/proto_rpc"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

// getAccountByEmail fetches an account by email
func (s *RPC) getAccountByEmail(email string) (*db.Account, error) {
	var account db.Account
	err := s.db.Where("email = ?", email).First(&account).Error
	return &account, err
}

// GetAccount fetches an account object
func (s *RPC) GetAccount(ctx context.Context, req *proto_rpc.GetAccountMsg) (*proto_rpc.Account, error) {

	var account db.Account
	var q = s.db.Where("id = ?", req.ID)

	if govalidator.IsEmail(req.ID) {
		q = s.db.Where("email = ?", req.ID)
	}

	if err := q.First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.Error(404, "", "id", "account not found")
		}
		return nil, common.ServerError
	}

	var resp proto_rpc.Account
	copier.Copy(&resp, account)

	return &resp, nil
}

// getAccount returns an account or common.APIError if not found
func (s *RPC) getAccount(id string) (*db.Account, error) {
	var m db.Account
	if err := s.db.Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.Error(404, common.CodeInvalidParam, "", "account not found")
		}
		return nil, common.ServerError
	}
	return &m, nil
}