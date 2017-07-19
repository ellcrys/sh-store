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

// CreateMapping creates a mapping for an identity
func (s *RPC) CreateMapping(ctx context.Context, req *proto_rpc.CreateMappingMsg) (*proto_rpc.CreateMappingResponse, error) {

	if len(req.Name) == 0 {
		req.Name = util.UUID4()
	}

	var mapping map[string]string
	if err := util.FromJSON(req.Mapping, &mapping); err != nil {
		return nil, common.NewSingleAPIErr(400, "", "", "malformed mapping", nil)
	}

	errs := validateMapping(mapping)
	if len(errs) > 0 {
		return nil, common.NewMultiAPIErr(400, "mapping error", errs)
	}

	developerID := ctx.Value(CtxIdentity).(string)

	// check if an existing mapping with same name exists
	err := s.db.Where("name = ? AND identity = ?", req.Name, developerID).First(&db.Mapping{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, common.ServerError
	} else if err == nil {
		return nil, common.NewSingleAPIErr(400, "", "", fmt.Sprintf("mapping with name '%s' already exist", req.Name), nil)
	}

	var newMapping = db.NewMapping()
	newMapping.Identity = developerID
	newMapping.Name = req.Name
	newMapping.Mapping = string(req.Mapping)

	if err := s.db.Create(&newMapping).Error; err != nil {
		return nil, common.ServerError
	}

	return &proto_rpc.CreateMappingResponse{
		Name: req.Name,
		ID:   newMapping.ID,
	}, nil
}

// GetMapping fetches a mapping belonging to the logged in developer
func (s *RPC) GetMapping(ctx context.Context, req *proto_rpc.GetMappingMsg) (*proto_rpc.GetMappingResponse, error) {

	if len(req.Name) == 0 {
		return nil, common.NewSingleAPIErr(400, "", "id", "name is required", nil)
	}

	developerID := ctx.Value(CtxIdentity).(string)

	var mapping db.Mapping
	if err := s.db.Where("name = ? AND identity = ?", req.Name, developerID).Last(&mapping).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "mapping not found", nil)
		}
		return nil, common.ServerError
	}

	return &proto_rpc.GetMappingResponse{
		Mapping: util.MustStringify(mapping),
	}, nil
}

// GetAllMapping fetches the most recent mappings belonging to the logged in developer
func (s *RPC) GetAllMapping(ctx context.Context, req *proto_rpc.GetAllMappingMsg) (*proto_rpc.GetAllMappingResponse, error) {

	developerID := ctx.Value(CtxIdentity).(string)

	if req.Limit == 0 || req.Limit > MaxGetAllMappingLimit {
		req.Limit = MaxGetAllMappingLimit
	}

	var mappings []db.Mapping
	if err := s.db.Where("identity = ?", developerID).Order("created_at DESC").Limit(req.Limit).Find(&mappings).Error; err != nil {
		return nil, common.ServerError
	}

	return &proto_rpc.GetAllMappingResponse{
		Mappings: util.MustStringify(mappings),
	}, nil
}
