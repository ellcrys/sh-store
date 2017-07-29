package servers

import (
	"github.com/asaskevich/govalidator"
	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/ellcrys/util"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

// getAccountByEmail fetches an account by email
func (s *RPC) getAccountByEmail(email string) (*db.Account, error) {
	var account db.Account
	err := s.db.Where("email = ?", email).First(&account).Error
	return &account, err
}

// CreateAccount creates a new account for object
func (s *RPC) CreateAccount(ctx context.Context, req *proto_rpc.CreateAccountMsg) (*proto_rpc.GetObjectResponse, error) {

	if errs := validateAccount(req); len(errs) > 0 {
		return nil, common.NewMultiAPIErr(400, "validation errors", errs)
	}

	// check if email has been used
	_, err := s.getAccountByEmail(req.Email)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, common.ServerError
		}
	} else {
		return nil, common.NewSingleAPIErr(400, common.CodeInvalidParam, "email", "Email is not available", nil)
	}

	var newAccount = db.NewAccount()
	copier.Copy(&newAccount, &req)

	// generate client id and secret for developer account
	if newAccount.Developer {
		newAccount.ClientID = util.RandString(32)
		newAccount.ClientSecret = util.RandString(util.RandNum(28, 32))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, common.ServerError
	}

	newAccount.Password = string(hashedPassword)
	if err = s.db.Create(&newAccount).Error; err != nil {
		return nil, common.ServerError
	}

	return &proto_rpc.GetObjectResponse{
		Object: util.MustStringify(newAccount),
	}, nil
}

// GetAccount fetches an account object
func (s *RPC) GetAccount(ctx context.Context, req *proto_rpc.GetAccountMsg) (*proto_rpc.GetObjectResponse, error) {

	var account db.Account
	var q = s.db.Where("id = ?", req.ID)

	if govalidator.IsEmail(req.ID) {
		q = s.db.Where("email = ?", req.ID)
	}

	if err := q.First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "Account not found", nil)
		}
		return nil, common.ServerError
	}

	return &proto_rpc.GetObjectResponse{
		Object: util.MustStringify(account),
	}, nil
}

// getAccount returns an account or common.APIError if not found
func (s *RPC) getAccount(id string) (*db.Account, error) {
	var m db.Account
	if err := s.db.Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.NewSingleAPIErr(404, common.CodeInvalidParam, "account", "account not found", nil)
		}
		return nil, common.ServerError
	}
	return &m, nil
}
