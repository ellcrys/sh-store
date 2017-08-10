package servers

import (
	"fmt"

	"github.com/ellcrys/jsonapi"

	"github.com/asaskevich/govalidator"
	"github.com/ellcrys/elldb/core/servers/common"
	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/ellcrys/elldb/core/servers/proto_rpc"
	"github.com/ellcrys/util"
	"github.com/jinzhu/gorm"
	mv "github.com/ncodes/mapvalidator"
)

// validateContract checks whether a contract message is valid
func validateContract(req *proto_rpc.CreateContractMsg) []*jsonapi.ErrorObject {
	var errs []*jsonapi.ErrorObject
	if govalidator.IsNull(req.Name) {
		errs = append(errs, &jsonapi.ErrorObject{Code: common.CodeInvalidParam, Source: "/name", Detail: "name is required"})
	} else if !govalidator.Matches(req.Name, "^[a-z0-9_]+$") {
		errs = append(errs, &jsonapi.ErrorObject{Code: common.CodeInvalidParam, Source: "/name", Detail: "name is not valid. Only alphanumeric and underscore characters are allowed"})
	}
	return errs
}

// validateBucket validates a bucket to be created
func validateBucket(req *proto_rpc.CreateBucketMsg) []*jsonapi.ErrorObject {
	var errs []*jsonapi.ErrorObject
	if govalidator.IsNull(req.Name) {
		errs = append(errs, &jsonapi.ErrorObject{Code: common.CodeInvalidParam, Source: "/name", Detail: "name is required"})
	}
	return errs
}

// validateObjects validates a slice of objects returning all errors found.
// - empty object not allowed
// - max number of object must not be exceeded
// - owner id is required
// - owner id of all objects must be the same
// - owner id must exist
// - key is required
// Invalid field will be replaced with their mappings if mapping is provided
func (s *RPC) validateObjects(objs []map[string]interface{}, mapping map[string]string) ([]*jsonapi.ErrorObject, error) {

	var errs []*jsonapi.ErrorObject

	// Returns the custom field of a standard db.Object field
	// using the mapping provided. If a standard object field does not
	// have a custom field, the standard field is returned.
	var getCustomFieldFromMapping = func(field string) string {
		for custom, _field := range mapping {
			if field == _field {
				return custom
			}
		}
		return field
	}

	if len(objs) == 0 {
		errs = append(errs, &jsonapi.ErrorObject{Detail: "no object provided. At least one object is required"})
		return errs, nil
	}

	if len(objs) > MaxObjectPerPut {
		errs = append(errs, &jsonapi.ErrorObject{Detail: fmt.Sprintf("too many objects. Only a maximum of %d can be created at once", MaxObjectPerPut)})
		return errs, nil
	}

	prevOwnerID := ""
	for i, o := range objs {
		rules := []mv.Rule{
			mv.RequiredWithMsg("owner_id", fmt.Sprintf("object %d: %s is required", i, getCustomFieldFromMapping("owner_id"))),
			mv.TypeWithMsg("owner_id", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getCustomFieldFromMapping("owner_id"))),
			mv.RequiredWithMsg("creator_id", fmt.Sprintf("object %d: %s is required", i, getCustomFieldFromMapping("creator_id"))),
			mv.TypeWithMsg("creator_id", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, "creator_id")),
			mv.RequiredWithMsg("key", fmt.Sprintf("object %d: %s is required", i, getCustomFieldFromMapping("key"))),
			mv.TypeWithMsg("key", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getCustomFieldFromMapping("key"))),
			mv.TypeWithMsg("value", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getCustomFieldFromMapping("value"))),
			mv.TypeWithMsg("ref1", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getCustomFieldFromMapping("ref1"))),
			mv.TypeWithMsg("ref2", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getCustomFieldFromMapping("ref2"))),
			mv.TypeWithMsg("ref3", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getCustomFieldFromMapping("ref3"))),
			mv.TypeWithMsg("ref4", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getCustomFieldFromMapping("ref4"))),
			mv.TypeWithMsg("ref5", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getCustomFieldFromMapping("ref5"))),
			mv.TypeWithMsg("ref6", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getCustomFieldFromMapping("ref6"))),
			mv.TypeWithMsg("ref7", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getCustomFieldFromMapping("ref7"))),
			mv.TypeWithMsg("ref8", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getCustomFieldFromMapping("ref8"))),
			mv.TypeWithMsg("ref9", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getCustomFieldFromMapping("ref9"))),
			mv.TypeWithMsg("ref10", mv.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getCustomFieldFromMapping("ref10"))),

			// same 'owner_id' must be shared on all objects
			mv.Custom("owner_id", fmt.Sprintf("object %d: all objects must share the same '%s' (owner_id) value", i, getCustomFieldFromMapping("owner_id")), func(fieldValue interface{}, fullMap map[string]interface{}) bool {
				if fv, ok := fieldValue.(string); ok {
					if prevOwnerID == "" {
						prevOwnerID = fv
					}
					return prevOwnerID == fv
				}
				return true
			}),

			// 'key' cannot start with '$' character
			mv.Custom("key", fmt.Sprintf("object %d: %s cannot start with a '$' character", i, getCustomFieldFromMapping("key")), func(fieldValue interface{}, fullMap map[string]interface{}) bool {
				if fv, ok := fieldValue.(string); ok {
					return fv[0] != '$'
				}
				return true
			}),
		}

		vErrs := mv.Validate(o, rules...)
		for _, err := range vErrs {
			errs = append(errs, &jsonapi.ErrorObject{Code: common.CodeInvalidParam, Detail: err.Error()})
		}
	}

	if len(errs) > 0 {
		return errs, nil
	}

	// assuming the owner is an account, check if the account it exist
	if err := s.objectOwnerExists(objs[0]["owner_id"].(string)); err != nil {
		if err == gorm.ErrRecordNotFound {
			errs = append(errs, &jsonapi.ErrorObject{Code: common.CodeInvalidParam, Detail: "owner of object(s) does not exist"})
		} else {
			return nil, err
		}
	}

	return errs, nil
}

// validateMapping checks whether a map is valid
func validateMapping(mapping map[string]string) []*jsonapi.ErrorObject {
	var errs []*jsonapi.ErrorObject
	validObjectFields := db.GetValidObjectFields()

	if len(mapping) == 0 {
		errs = append(errs, &jsonapi.ErrorObject{
			Source: "",
			Detail: "requires at least one custom field",
		})
	}

	for customField, field := range mapping {
		if len(field) == 0 {
			errs = append(errs, &jsonapi.ErrorObject{
				Source: "/" + customField,
				Detail: "field name is required",
			})
			continue
		}
		if !util.InStringSlice(validObjectFields, field) {
			errs = append(errs, &jsonapi.ErrorObject{
				Source: "/" + customField,
				Detail: fmt.Sprintf("field value '%s' is not a valid object field", field),
			})
		}
	}
	return errs
}
