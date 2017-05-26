package servers

import (
	"github.com/ellcrys/util"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"github.com/ncodes/patchain/cockroach/tables"
)

// NewObjectResponse creates a standard object response
func NewObjectResponse(oType string, o *tables.Object, links map[string]string) (*proto_rpc.ObjectResponse, error) {
	var obj proto_rpc.Object
	util.CopyToStruct(&obj, o)
	return &proto_rpc.ObjectResponse{
		Data: &proto_rpc.ObjectResponseData{
			Type:       oType,
			ID:         o.ID,
			Attributes: &obj,
		},
		Links: links,
	}, nil
}

// NewMultiObjectResponse creates a multi object response from a slice of objects
func NewMultiObjectResponse(oType string, objects []*tables.Object) (*proto_rpc.MultiObjectResponse, error) {
	var objResponses []*proto_rpc.ObjectResponseData

	for _, obj := range objects {
		var _obj proto_rpc.Object
		util.CopyToStruct(&_obj, obj)
		objResponses = append(objResponses, &proto_rpc.ObjectResponseData{
			Type:       oType,
			ID:         _obj.ID,
			Attributes: &_obj,
		})
	}

	return &proto_rpc.MultiObjectResponse{
		Data: objResponses,
	}, nil
}
