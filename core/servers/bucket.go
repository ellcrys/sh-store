package servers

import (
	"github.com/ellcrys/elldb/core/servers/common"
	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/ellcrys/elldb/core/servers/proto_rpc"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

// CreateBucket creates a bucket
func (s *RPC) CreateBucket(ctx context.Context, req *proto_rpc.CreateBucketMsg) (*proto_rpc.Bucket, error) {

	if errs := validateBucket(req); len(errs) > 0 {
		return nil, common.Errors(400, errs)
	}

	contract := ctx.Value(CtxContract).(*db.Contract)

	// check if bucket with matching name exists
	err := s.db.Where("name = ?", req.Name).First(&db.Bucket{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, common.ServerError
	} else if err == nil {
		return nil, common.Error(400, "", "/name", "name has been used")
	}

	var bucket = db.NewBucket()
	bucket.Creator = contract.ID
	bucket.Name = req.Name
	bucket.Immutable = req.Immutable
	if err = db.CreateBucket(s.db, bucket); err != nil {
		return nil, common.ServerError
	}

	var resp proto_rpc.Bucket
	copier.Copy(&resp, bucket)

	return &resp, nil
}

// getBucket fetches a bucket
func (s *RPC) getBucket(name string) (*db.Bucket, error) {
	var b db.Bucket
	err := s.db.Where("name = ?", name).First(&b).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.Error(404, "", "/bucket", "bucket not found")
		}
		return nil, common.ServerError
	}
	return &b, err
}
