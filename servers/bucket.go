package servers

import (
	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

// CreateBucket creates a bucket
func (s *RPC) CreateBucket(ctx context.Context, req *proto_rpc.CreateBucketMsg) (*proto_rpc.Bucket, error) {

	if errs := validateBucket(req); len(errs) > 0 {
		return nil, common.NewMultiAPIErr(400, "validation errors", errs)
	}

	// check if bucket with matching name exists
	err := s.db.Where("name = ?", req.Name).First(&db.Bucket{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, common.ServerError
	} else if err == nil {
		return nil, common.NewSingleAPIErr(400, "", "name", "Bucket name has been used", nil)
	}

	var bucket = db.NewBucket()
	copier.Copy(&bucket, &req)
	bucket.Identity = ctx.Value(CtxIdentity).(string)
	if err = db.CreateBucket(s.db, bucket); err != nil {
		return nil, common.ServerError
	}

	var resp proto_rpc.Bucket
	copier.Copy(&resp, bucket)
	return &resp, nil
}
