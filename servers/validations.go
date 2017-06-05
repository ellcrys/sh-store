package servers

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/ncodes/mapvalidator"
	"github.com/ncodes/patchain"
	"github.com/ncodes/patchain/cockroach/tables"
	"github.com/ncodes/safehold/servers/common"
	"github.com/ncodes/safehold/servers/proto_rpc"
)

// validateObjects validates a slice of objects returning all errors found.
// - empty object not allowed
// - max number of object must not be exceeded
// - owner id is required
// - owner id of all objects must be the same
// - owner id must exist
// - key is required
// Invalid field will be replaced with their mappings if mapping is provided
func (s *RPC) validateObjects(objs []map[string]interface{}, mapping map[string]string) ([]common.Error, error) {

	var errs []common.Error

	// get the mapped field for an object field or return the object field if it has no mapped field
	var getMappedField = func(objectField string) string {
		for mappedField, _objectField := range mapping {
			if objectField == _objectField {
				return mappedField
			}
		}
		return objectField
	}

	if len(objs) == 0 {
		errs = append(errs, common.Error{Message: "no object provided. At least one object is required"})
		return errs, nil
	}

	if len(objs) > MaxObjectPerPut {
		errs = append(errs, common.Error{Message: fmt.Sprintf("too many objects. Only a maximum of %d can be created at once", MaxObjectPerPut)})
		return errs, nil
	}

	var prevOwnerID string
	for i, o := range objs {
		rules := []mapvalidator.Rule{
			mapvalidator.RequiredWithMsg("owner_id", fmt.Sprintf("object %d: %s is required", i, getMappedField("owner_id"))),
			mapvalidator.TypeWithMsg("owner_id", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("owner_id"))),
			mapvalidator.RequiredWithMsg("key", fmt.Sprintf("object %d: %s is required", i, getMappedField("key"))),
			mapvalidator.TypeWithMsg("key", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("key"))),
			mapvalidator.TypeWithMsg("value", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("value"))),
			mapvalidator.TypeWithMsg("ref1", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref1"))),
			mapvalidator.TypeWithMsg("ref2", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref2"))),
			mapvalidator.TypeWithMsg("ref3", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref3"))),
			mapvalidator.TypeWithMsg("ref4", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref4"))),
			mapvalidator.TypeWithMsg("ref5", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref5"))),
			mapvalidator.TypeWithMsg("ref6", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref6"))),
			mapvalidator.TypeWithMsg("ref7", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref7"))),
			mapvalidator.TypeWithMsg("ref8", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref8"))),
			mapvalidator.TypeWithMsg("ref9", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref9"))),
			mapvalidator.TypeWithMsg("ref10", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref10"))),

			mapvalidator.Custom("owner_id", fmt.Sprintf("object %d: all objects must share the same %s", i, getMappedField("owner_id")), func(fieldValue interface{}, fullMap map[string]interface{}) bool {
				if fv, ok := fieldValue.(string); ok {
					if prevOwnerID == "" {
						prevOwnerID = fv
					}
					return prevOwnerID == fv
				}
				return true
			}),

			mapvalidator.Custom("key", fmt.Sprintf("object %d: %s cannot start with a '$' character", i, getMappedField("key")), func(fieldValue interface{}, fullMap map[string]interface{}) bool {
				if fv, ok := fieldValue.(string); ok {
					return fv[0] != '$'
				}
				return true
			}),
		}

		vErrs := mapvalidator.Validate(o, rules...)
		for _, err := range vErrs {
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: err.Error()})
		}
	}

	if len(errs) > 0 {
		return errs, nil
	}

	_, err := s.object.GetLast(&tables.Object{ID: objs[0]["owner_id"].(string)})
	if err != nil {
		if err == patchain.ErrNotFound {
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "owner of object(s) does not exist"})
		} else {
			return nil, err
		}
	}

	return errs, nil
}

// validateIdentity validates a requested identity
func validateIdentity(req *proto_rpc.CreateIdentityMsg) []common.Error {
	var errs []common.Error
	if govalidator.IsNull(req.Email) {
		errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "Email is required", Field: "email"})
	}
	if !govalidator.IsEmail(req.Email) {
		errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "Email is not valid", Field: "email"})
	}
	if govalidator.IsNull(req.Password) {
		errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "Password is required", Field: "password"})
	}
	if len(req.Password) < 6 {
		errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "Password is too short. Must be at least 6 characters long", Field: "password"})
	}
	return errs
}
