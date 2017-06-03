package servers

import (
	"fmt"

	"github.com/ellcrys/util"
	"github.com/labstack/gommon/log"
	"github.com/ncodes/patchain"
	"github.com/ncodes/patchain/cockroach"
	"github.com/ncodes/patchain/cockroach/tables"
	"github.com/ncodes/patchain/object"
	"github.com/ncodes/safehold/servers/common"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"github.com/ncodes/safehold/session"
	"golang.org/x/net/context"
)

// validateMapping checks whether a map is valid
func validateMapping(mapping map[string]interface{}) []common.Error {
	var errs []common.Error
	validObjectFields := cockroach.NewDB().GetValidObjectFields()

	if len(mapping) == 0 {
		errs = append(errs, common.Error{
			Field:   "mapping",
			Message: "requires at least one mapped field",
		})
	}

	for mapField, fieldData := range mapping {
		switch data := fieldData.(type) {
		case string:
			if len(data) == 0 {
				errs = append(errs, common.Error{
					Field:   mapField,
					Message: fmt.Sprintf("column name is required"),
				})
				continue
			}
			if !util.InStringSlice(validObjectFields, data) {
				errs = append(errs, common.Error{
					Field:   mapField,
					Message: fmt.Sprintf("column name '%s' is unknown", data),
				})
			}
		default:
			errs = append(errs, common.Error{
				Field:   mapField,
				Message: "invalid value type. Expected string type",
			})
		}
	}
	return errs
}

// CreateMapping creates a mapping for an identity
func (s *RPC) CreateMapping(ctx context.Context, req *proto_rpc.CreateMappingMsg) (*proto_rpc.CreateMappingResponse, error) {

	developerID := ctx.Value(CtxIdentity).(string)

	if len(req.Name) == 0 {
		req.Name = util.UUID4()
	}

	var mapping map[string]interface{}
	if err := util.FromJSON(req.Mapping, &mapping); err != nil {
		return nil, common.NewSingleAPIErr(404, "", "", "malformed mapping. Expected json object.", nil)
	}

	// validate mapping
	errs := validateMapping(mapping)
	if len(errs) > 0 {
		return nil, common.NewMultiAPIErr(400, "mapping error", errs)
	}

	fullSessionID := makeDBSessionID(developerID, util.UUID4())
	s.dbSession.CreateUnregisteredSession(fullSessionID, developerID)
	defer s.dbSession.CommitEnd(fullSessionID)

	mappingObj := object.MakeMappingObject(developerID, req.Name, string(req.Mapping))
	err := s.dbSession.SendOp(fullSessionID, &session.Op{
		OpType: session.OpPutObjects,
		Data:   mappingObj,
		Done:   make(chan struct{}),
	})
	if err != nil {
		s.dbSession.RollbackEnd(fullSessionID)
		return nil, common.ServerError
	}

	return &proto_rpc.CreateMappingResponse{
		Name: req.Name,
	}, nil
}

// GetMapping fetches a mapping belonging to the logged in developer
func (s *RPC) GetMapping(ctx context.Context, req *proto_rpc.GetMappingMsg) (*proto_rpc.GetMappingResponse, error) {

	developerID := ctx.Value(CtxIdentity).(string)

	if len(req.Name) == 0 {
		return nil, common.NewSingleAPIErr(400, "", "id", "name is required", nil)
	}

	fullSessionID := makeDBSessionID(developerID, util.UUID4())
	s.dbSession.CreateUnregisteredSession(fullSessionID, developerID)
	defer s.dbSession.CommitEnd(fullSessionID)

	var mapping tables.Object
	err := s.dbSession.SendOp(fullSessionID, &session.Op{
		OpType: session.OpGetObjects,
		Data: `{ 
			"key": "` + object.MakeMappingKey(req.Name) + `", 
			"owner_id": "` + developerID + `" 
		}`,
		OrderBy: "timestamp desc",
		Limit:   1,
		Out:     &mapping,
		Done:    make(chan struct{}),
	})
	if err != nil {
		if err == patchain.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "mapping not found", nil)
		}
		log.Errorf("%+v", err)
		return nil, common.ServerError
	}

	return &proto_rpc.GetMappingResponse{
		Mapping: []byte(mapping.Value),
	}, nil
}

// GetAllMapping fetches the most recent mappings belonging to the logged in developer
func (s *RPC) GetAllMapping(ctx context.Context, req *proto_rpc.GetAllMappingMsg) (*proto_rpc.GetAllMappingResponse, error) {
	developerID := ctx.Value(CtxIdentity).(string)

	if req.Limit == 0 {
		req.Limit = 50
	}

	fullSessionID := makeDBSessionID(developerID, util.UUID4())
	s.dbSession.CreateUnregisteredSession(fullSessionID, developerID)
	defer s.dbSession.CommitEnd(fullSessionID)

	var expr = `{ "$sw": "` + object.MappingPrefix + `" }`
	if len(req.Name) > 0 {
		expr = `{ "$eq": "` + object.MakeMappingKey(req.Name) + `" }`
	}

	var mappings []tables.Object
	err := s.dbSession.SendOp(fullSessionID, &session.Op{
		OpType: session.OpGetObjects,
		Data: `{ 
			"key":` + expr + `, 
			"owner_id": "` + developerID + `" 
		}`,
		OrderBy: "timestamp desc",
		Limit:   int(req.Limit),
		Out:     &mappings,
		Done:    make(chan struct{}),
	})
	if err != nil {
		if err == patchain.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "mapping not found", nil)
		}
		log.Errorf("%+v", err)
		return nil, common.ServerError
	}

	var result []*proto_rpc.Mapping
	for _, m := range mappings {
		result = append(result, &proto_rpc.Mapping{
			Name:    m.Key,
			Mapping: []byte(m.Value),
		})
	}

	return &proto_rpc.GetAllMappingResponse{
		Mappings: result,
	}, nil
}
