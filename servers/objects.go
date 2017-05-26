package servers

import (
	"fmt"

	"strings"

	"github.com/jinzhu/copier"
	"github.com/ncodes/safehold/servers/common"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"github.com/ncodes/patchain"
	"github.com/ncodes/patchain/cockroach/tables"
	"github.com/ncodes/patchain/object"
	"golang.org/x/net/context"
)

var (
	// MaxObjectPerPut is the maximum number of objects allowed in a single PUT operation
	MaxObjectPerPut = 25
)

// validateObjects validates a slice of objects returning all errors found.
// - owner id is required
// - owner id of all objects must be the same
// - owner id must exist
// - key is required
func (s *RPC) validateObjects(obj []*proto_rpc.Object) []common.Error {
	var errs []common.Error
	var ownerID = obj[0].OwnerID

	for i, o := range obj {

		if len(strings.TrimSpace(o.OwnerID)) == 0 {
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: fmt.Sprintf("object %d: owner id is required", i), Field: "objects"})
			continue
		}

		if ownerID != o.OwnerID {
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: fmt.Sprintf("object %d: has a different owner. All objects must be owned by a single identity", i), Field: "objects"})
			continue
		}

		_, err := s.object.GetLast(&tables.Object{ID: o.OwnerID})
		if err != nil {
			if err == patchain.ErrNotFound {
				errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: fmt.Sprintf("object %d: owner id does not exist", i), Field: "objects"})
				continue
			}
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: fmt.Sprintf("object %d: server error", i), Field: "objects"})
			continue
		}

		if len(strings.TrimSpace(o.Key)) == 0 {
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: fmt.Sprintf("object %d: key is required", i), Field: "objects"})
			continue
		}
	}
	return errs
}

// CreateObjects creates one or more objects
func (s *RPC) CreateObjects(ctx context.Context, req *proto_rpc.CreateObjectsMsg) (*proto_rpc.MultiObjectResponse, error) {

	if len(req.Objects) == 0 {
		return nil, common.NewSingleAPIErr(400, "", "", "no object provided. At least one object is required", nil)
	}
	if len(req.Objects) > MaxObjectPerPut {
		return nil, common.NewSingleAPIErr(400, "", "", "too many objects. Only a maximum of 25 can be created at once", nil)
	}

	if errs := s.validateObjects(req.Objects); len(errs) > 0 {
		return nil, common.NewMultiAPIErr(400, "validation errors", errs)
	}

	// ensure caller has permission to PUT objects on behalf of the object owner
	developerID := ctx.Value(CtxIdentity).(string)
	if developerID != req.Objects[0].OwnerID {
		// TODO: check object owner ACL to see if object owner has granted developer permission
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you are not authorized to create objects for the owner", nil)
	}

	objHandler := object.NewObject(s.db)

	var objs []*tables.Object
	copier.Copy(&objs, req.Objects)
	if err := objHandler.Put(objs); err != nil {
		logRPC.Errorf("%+v", err)
		return nil, common.NewSingleAPIErr(500, common.CodePutError, "objects", err.Error(), nil)
	}

	resp, _ := NewMultiObjectResponse("object", objs)
	return resp, nil
}
