package servers

import (
	"fmt"

	"github.com/ellcrys/util"
	"github.com/ncodes/patchain/cockroach"
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
		case map[string]interface{}:
			if data["field"] == nil {
				errs = append(errs, common.Error{
					Field:   mapField,
					Message: "'field' property is required",
				})
				continue
			}
			field, ok := data["field"].(string)
			if !ok {
				errs = append(errs, common.Error{
					Field:   fmt.Sprintf("%s.field", mapField),
					Message: "invalid value type. Expected string type",
				})
				continue
			}
			if len(field) == 0 {
				errs = append(errs, common.Error{
					Field:   fmt.Sprintf("%s.field", mapField),
					Message: fmt.Sprintf("column name is required"),
				})
				continue
			}
			if !util.InStringSlice(validObjectFields, field) {
				errs = append(errs, common.Error{
					Field:   fmt.Sprintf("%s.field", mapField),
					Message: fmt.Sprintf("column name '%s' is unknown", field),
				})
				continue
			}
			if data["protected"] != nil {
				if _, ok := data["protected"].(bool); !ok {
					errs = append(errs, common.Error{
						Field:   fmt.Sprintf("%s.protected", mapField),
						Message: "invalid value type. Expected boolean type",
					})
				}
				continue
			}
		default:
			errs = append(errs, common.Error{
				Field:   mapField,
				Message: "invalid value type. Expected string or json type",
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
		return nil, common.ServerError
	}

	return &proto_rpc.CreateMappingResponse{
		Name: req.Name,
	}, nil
}
