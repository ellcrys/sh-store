package servers

import (
	"github.com/ellcrys/util"
	"github.com/jinzhu/copier"

	"github.com/jinzhu/gorm"

	"github.com/ellcrys/elldb/core/servers/common"
	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/ellcrys/elldb/core/servers/proto_rpc"
	"golang.org/x/net/context"
)

// CreateContract creates a contract
func (s *RPC) CreateContract(ctx context.Context, req *proto_rpc.CreateContractMsg) (*proto_rpc.Contract, error) {

	if errs := validateContract(req); len(errs) > 0 {
		return nil, common.Errors(400, errs)
	}

	account := ctx.Value(CtxAccount).(*db.Account)

	// check if a contract with a matching name exists
	err := s.db.Where("name = ?", req.Name).First(&db.Contract{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, common.ServerError
	} else if err == nil {
		return nil, common.Error(400, "", "/name", "name is not available")
	}

	var contract = db.NewContract()
	contract.Creator = account.ID
	contract.Name = req.Name
	contract.ClientID = util.RandString(32)
	contract.ClientSecret = util.RandString(util.RandNum(27, 32))
	if err := s.db.Create(contract).Error; err != nil {
		return nil, common.ServerError
	}

	var resp proto_rpc.Contract
	copier.Copy(&resp, contract)

	return &resp, nil
}
