package servers

import (
	"fmt"

	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/ellcrys/util"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

// MaxGetAllMappingLimit is the maximum number of mapping to return for GetAllMapping()
var MaxGetAllMappingLimit = int32(50)

// CreateMapping creates a mapping for a bucket
func (s *RPC) CreateMapping(ctx context.Context, req *proto_rpc.CreateMappingMsg) (*proto_rpc.Mapping, error) {

	if len(req.Name) == 0 {
		req.Name = util.UUID4()
	}

	if len(req.Bucket) == 0 {
		return nil, common.Error(400, "", "", "bucket is required")
	}

	contract := ctx.Value(CtxContract).(*db.Contract)

	bucket, err := s.getBucket(req.Bucket)
	if err != nil {
		return nil, err
	} else if bucket.Creator != contract.ID {
		return nil, common.Error(401, "", "", "You do not have permission to perform this action")
	}

	var mapping map[string]string
	if err := util.FromJSON([]byte(req.Mapping), &mapping); err != nil {
		return nil, common.Error(400, "", "", "malformed mapping")
	}

	errs := validateMapping(mapping)
	if len(errs) > 0 {
		return nil, common.Errors(400, errs)
	}

	// check if an existing mapping with same name exists
	err = s.db.Where("name = ? AND creator = ?", req.Name, contract.ID).First(&db.Mapping{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, common.ServerError
	} else if err == nil {
		return nil, common.Error(400, "", "", fmt.Sprintf("mapping with name '%s' already exist", req.Name))
	}

	var newMapping = db.NewMapping()
	newMapping.Creator = contract.ID
	newMapping.Name = req.Name
	newMapping.Bucket = req.Bucket
	newMapping.Mapping = string(req.Mapping)

	if err := s.db.Create(&newMapping).Error; err != nil {
		return nil, common.ServerError
	}

	return &proto_rpc.Mapping{
		Bucket:  req.Bucket,
		Name:    req.Name,
		ID:      newMapping.ID,
		Creator: contract.ID,
		Mapping: req.Mapping,
	}, nil
}

// GetMapping fetches a mapping belonging to the logged in developer
func (s *RPC) GetMapping(ctx context.Context, req *proto_rpc.GetMappingMsg) (*proto_rpc.Mapping, error) {

	if len(req.Name) == 0 {
		return nil, common.Error(400, "", "/id", "name is required")
	}

	contract := ctx.Value(CtxContract).(*db.Contract)

	var mapping db.Mapping
	if err := s.db.Where("name = ? AND creator = ?", req.Name, contract.ID).Last(&mapping).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.Error(404, "", "", "mapping not found")
		}
		return nil, common.ServerError
	}

	return &proto_rpc.Mapping{
		Bucket:  mapping.Bucket,
		Name:    mapping.Name,
		Creator: mapping.Creator,
		Mapping: string(util.MustStringify(mapping.Mapping)),
	}, nil
}

// GetAllMapping fetches the most recent mappings belonging to the logged in developer
func (s *RPC) GetAllMapping(ctx context.Context, req *proto_rpc.GetAllMappingMsg) (*proto_rpc.Mappings, error) {

	contract := ctx.Value(CtxContract).(*db.Contract)

	if req.Limit == 0 || req.Limit > MaxGetAllMappingLimit {
		req.Limit = MaxGetAllMappingLimit
	}

	var mappings []*proto_rpc.Mapping
	if err := s.db.Where("creator = ?", contract.ID).Order("created_at DESC").Limit(req.Limit).Find(&mappings).Error; err != nil {
		return nil, common.ServerError
	}

	return &proto_rpc.Mappings{
		Mappings: mappings,
	}, nil
}

// getMapping returns a mapping or common.APIError
func (s *RPC) getMapping(name, creator string) (*db.Mapping, error) {
	var m db.Mapping
	if err := s.db.Where("name = ? AND creator = ?", name, creator).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.Error(404, common.CodeInvalidParam, "", "mapping not found")
		}
		return nil, common.ServerError
	}
	return &m, nil
}
