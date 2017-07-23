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

// getIdentityByEmail fetches an identity by email
func (s *RPC) getIdentityByEmail(email string) (*db.Identity, error) {
	var identity db.Identity
	err := s.db.Where("email = ?", email).First(&identity).Error
	return &identity, err
}

// CreateIdentity creates a new identity for object
func (s *RPC) CreateIdentity(ctx context.Context, req *proto_rpc.CreateIdentityMsg) (*proto_rpc.GetObjectResponse, error) {

	if errs := validateIdentity(req); len(errs) > 0 {
		return nil, common.NewMultiAPIErr(400, "validation errors", errs)
	}

	// check if email has been used
	_, err := s.getIdentityByEmail(req.Email)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, common.ServerError
		}
	} else {
		return nil, common.NewSingleAPIErr(400, common.CodeInvalidParam, "email", "Email is not available", nil)
	}

	var newIdentity = db.NewIdentity()
	copier.Copy(&newIdentity, &req)

	// generate client id and secret for developer identity
	if newIdentity.Developer {
		newIdentity.ClientID = util.RandString(32)
		newIdentity.ClientSecret = util.RandString(util.RandNum(28, 32))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, common.ServerError
	}

	newIdentity.Password = string(hashedPassword)
	if err = s.db.Create(&newIdentity).Error; err != nil {
		return nil, common.ServerError
	}

	return &proto_rpc.GetObjectResponse{
		Object: util.MustStringify(newIdentity),
	}, nil
}

// GetIdentity fetches an identity object
func (s *RPC) GetIdentity(ctx context.Context, req *proto_rpc.GetIdentityMsg) (*proto_rpc.GetObjectResponse, error) {

	var identity db.Identity
	var q = s.db.Where("id = ?", req.ID)

	if govalidator.IsEmail(req.ID) {
		q = s.db.Where("email = ?", req.ID)
	}

	if err := q.First(&identity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "Identity not found", nil)
		}
		return nil, common.ServerError
	}

	return &proto_rpc.GetObjectResponse{
		Object: util.MustStringify(identity),
	}, nil
}

// getIdentity returns an identity or common.APIError if not found
func (s *RPC) getIdentity(id string) (*db.Identity, error) {
	var m db.Identity
	if err := s.db.Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.NewSingleAPIErr(404, common.CodeInvalidParam, "identity", "identity not found", nil)
		}
		return nil, common.ServerError
	}
	return &m, nil
}
