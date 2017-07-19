package servers

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/ellcrys/util"
	"github.com/jinzhu/gorm"
	mv "github.com/ncodes/mapvalidator"
)

// validateIdentity validates an identity to be created
func validateIdentity(req *proto_rpc.CreateIdentityMsg) []common.Error {
	var errs []common.Error
	if govalidator.IsNull(req.FirstName) {
		errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "First Name is required", Field: "first_name"})
	}
	if govalidator.IsNull(req.LastName) {
		errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "Last Name is required", Field: "last_name"})
	}
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

// validateBucket validates a bucket to be created
func validateBucket(req *proto_rpc.CreateBucketMsg) []common.Error {
	var errs []common.Error
	if govalidator.IsNull(req.Name) {
		errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "Name is required", Field: "name"})
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
func (s *RPC) validateObjects(objs []map[string]interface{}, mapping map[string]string) ([]common.Error, error) {

	var errs []common.Error

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
		errs = append(errs, common.Error{Message: "no object provided. At least one object is required"})
		return errs, nil
	}

	if len(objs) > MaxObjectPerPut {
		errs = append(errs, common.Error{Message: fmt.Sprintf("too many objects. Only a maximum of %d can be created at once", MaxObjectPerPut)})
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
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: err.Error()})
		}
	}

	if len(errs) > 0 {
		return errs, nil
	}

	// check whether the owner of the objects exist
	err := s.db.Where("id = ?", objs[0]["owner_id"]).Last(&db.Identity{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "owner of object(s) does not exist"})
		} else {
			return nil, err
		}
	}

	return errs, nil
}

// validateMapping checks whether a map is valid
func validateMapping(mapping map[string]string) []common.Error {
	var errs []common.Error
	validObjectFields := db.GetValidObjectFields()

	if len(mapping) == 0 {
		errs = append(errs, common.Error{
			Field:   "mapping",
			Message: "requires at least one custom field",
		})
	}

	for customField, field := range mapping {
		if len(field) == 0 {
			errs = append(errs, common.Error{
				Field:   "mapping." + customField,
				Message: fmt.Sprintf("field name is required"),
			})
			continue
		}
		if !util.InStringSlice(validObjectFields, field) {
			errs = append(errs, common.Error{
				Field:   "mapping." + customField,
				Message: fmt.Sprintf("field value '%s' is not a valid object field", field),
			})
		}
	}
	return errs
}
