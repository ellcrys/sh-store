package servers

import (
	val "github.com/asaskevich/govalidator"
	"github.com/ncodes/safehold/servers/common"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"github.com/ncodes/patchain"
	"github.com/ncodes/patchain/cockroach/tables"
	"github.com/ncodes/patchain/object"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

// validateIdentity validates a requested identity
func validateIdentity(req *proto_rpc.CreateIdentityMsg) []common.Error {
	var errs []common.Error
	if val.IsNull(req.Email) {
		errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "Email is required", Field: "email"})
	}
	if !val.IsEmail(req.Email) {
		errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "Email is not valid", Field: "email"})
	}
	if val.IsNull(req.Password) {
		errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "Password is required", Field: "password"})
	}
	if len(req.Password) < 6 {
		errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "Password is too short. Must be at least 6 characters long", Field: "password"})
	}
	return errs
}

// GetSystemIdentity fetches the system identity
func (s *RPC) GetSystemIdentity() (*tables.Object, error) {
	return s.object.GetLast(&tables.Object{
		Key: object.MakeIdentityKey(systemEmail),
	})
}

// CreateIdentity creates a new identity for object
func (s *RPC) CreateIdentity(ctx context.Context, req *proto_rpc.CreateIdentityMsg) (*proto_rpc.ObjectResponse, error) {

	// validate request
	if errs := validateIdentity(req); len(errs) > 0 {
		return nil, common.NewMultiAPIErr(400, "validation errors", errs)
	}

	systemIdentity, err := s.GetSystemIdentity()
	if err != nil {
		logRPC.Errorf("%+v", err)
		if err == patchain.ErrNotFound {
			return nil, common.NewSingleAPIErr(400, "", "", "system identity not found", nil)
		}
		return nil, common.ServerError
	}

	objHandler := object.NewObject(s.db)

	// ensure identity does not already exists
	_, err = objHandler.GetLast(&tables.Object{OwnerID: systemIdentity.ID, QueryParams: patchain.KeyStartsWith(object.MakeIdentityKey(req.Email))})
	if err != nil {
		if err != patchain.ErrNotFound {
			return nil, common.ServerError
		}
	}

	if err == nil {
		return nil, common.NewSingleAPIErr(400, "used_email", "email", "Email is not available", nil)
	}

	var numPartitions = int64(5)
	var newIdentity *tables.Object

	// create the identity and partitions
	var do = func() error {
		return s.db.Transact(func(txDB patchain.DB, commit patchain.CommitFunc, rollback patchain.RollbackFunc) error {

			var err error
			dbOpt := &patchain.UseDBOption{DB: txDB}
			passwordHash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

			if req.Developer {
				newIdentity = object.MakeDeveloperIdentityObject(systemIdentity.ID, systemIdentity.ID, req.Email, string(passwordHash), true)
			} else {
				newIdentity = object.MakeIdentityObject(systemIdentity.ID, systemIdentity.ID, req.Email, string(passwordHash), true)
			}

			err = s.object.Put(newIdentity, dbOpt)
			if err != nil {
				return err
			}

			_, err = s.object.CreatePartitions(numPartitions, newIdentity.ID, systemIdentity.ID, dbOpt)
			if err != nil {
				return err
			}

			return nil
		})
	}

	err = objHandler.Retry(func(stopRetry func()) error { return do() })
	if err != nil {
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	resp, _ := NewObjectResponse("identity", newIdentity, nil)
	return resp, nil
}
