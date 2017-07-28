package db

import (
	"time"

	"github.com/ellcrys/util"
)

// Identity represents a entity
type Identity struct {
	SN           int64  `json:"-" structs:"-" mapstructure:"-" gorm:"primary_key"`
	ID           string `json:"id,omitempty" structs:"id,omitempty" mapstructure:"id,omitempty" gorm:"type:varchar(36)"`
	FirstName    string `json:"first_name,omitempty" structs:"first_name,omitempty" mapstructure:"first_name,omitempty" gorm:""`
	LastName     string `json:"last_name,omitempty" structs:"last_name,omitempty" mapstructure:"last_name,omitempty" gorm:""`
	Email        string `json:"email,omitempty" structs:"email,omitempty" mapstructure:"email,omitempty" gorm:""`
	Password     string `json:"password,omitempty" structs:"password,omitempty" mapstructure:"password,omitempty" gorm:""`
	Developer    bool   `json:"developer,omitempty" structs:"developer,omitempty" mapstructure:"developer,omitempty" gorm:""`
	Confirmed    bool   `json:"confirmed" structs:"confirmed,omitempty" mapstructure:"confirmed,omitempty" gorm:""`
	ClientID     string `json:"client_id,omitempty" structs:"client_id,omitempty" mapstructure:"client_id,omitempty" gorm:""`
	ClientSecret string `json:"client_secret,omitempty" structs:"client_secret,omitempty" mapstructure:"client_secret,omitempty" gorm:""`
	CreatedAt    int64  `json:"created_at,omitempty" structs:"created_at,omitempty" mapstructure:"created_at,omitempty"`
}

// NewIdentity creates a new identity. An ID is assigned.
func NewIdentity() *Identity {
	return &Identity{ID: util.UUID4(), CreatedAt: time.Now().UnixNano()}
}

// Bucket represents a logical group for objects
type Bucket struct {
	SN        int64  `json:"-" structs:"-" mapstructure:"-" gorm:"primary_key"`
	ID        string `json:"id,omitempty" structs:"id,omitempty" mapstructure:"id,omitempty" gorm:"type:varchar(36)"`
	Name      string `json:"name,omitempty" structs:"name,omitempty" mapstructure:"name,omitempty" gorm:"type:varchar(36)"`
	Identity  string `json:"identity,omitempty" structs:"identity,omitempty" mapstructure:"identity,omitempty"`
	Immutable bool   `json:"immutable,omitempty" structs:"immutable,omitempty" mapstructure:"immutable,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty" structs:"created_at,omitempty" mapstructure:"created_at,omitempty"`
}

// NewBucket creates a new bucket. An ID is assigned.
func NewBucket() *Bucket {
	return &Bucket{ID: util.UUID4(), CreatedAt: time.Now().UnixNano()}
}

// Mapping represents
type Mapping struct {
	SN        int64  `json:"-" structs:"-" mapstructure:"-" gorm:"primary_key"`
	ID        string `json:"id,omitempty" structs:"id,omitempty" mapstructure:"id,omitempty" gorm:"type:varchar(36)"`
	Bucket    string `json:"bucket,omitempty" structs:"bucket,omitempty" mapstructure:"bucket,omitempty"`
	Identity  string `json:"identity,omitempty" structs:"identity,omitempty" mapstructure:"identity,omitempty"`
	Name      string `json:"name,omitempty" structs:"name,omitempty" mapstructure:"name,omitempty" gorm:"type:varchar(36)"`
	Mapping   string `json:"mapping,omitempty" structs:"mapping,omitempty" mapstructure:"mapping,omitempty" gorm:"type:varchar(256)"`
	CreatedAt int64  `json:"created_at,omitempty" structs:"created_at,omitempty" mapstructure:"created_at,omitempty"`
}

// NewMapping creates a new mapping. An ID is assigned.
func NewMapping() *Mapping {
	return &Mapping{ID: util.UUID4(), CreatedAt: time.Now().UnixNano()}
}

// Object represents a bucket object
type Object struct {
	SN          int64       `json:"-" structs:"-" mapstructure:"-" gorm:"primary_key"`
	ID          string      `json:"id,omitempty" structs:"id,omitempty" mapstructure:"id,omitempty" gorm:"type:varchar(36)"`
	OwnerID     string      `json:"owner_id,omitempty" structs:"owner_id,omitempty" mapstructure:"owner_id,omitempty" gorm:"index:idx_owner_id"`
	CreatorID   string      `json:"creator_id,omitempty" structs:"creator_id,omitempty" mapstructure:"creator_id,omitempty" gorm:"index:idx_creator_id"`
	Bucket      string      `json:"bucket,omitempty" structs:"bucket,omitempty" mapstructure:"bucket,omitempty" gorm:"index:idx_created_at"`
	Key         string      `json:"key,omitempty" structs:"key,omitempty" mapstructure:"key,omitempty" gorm:"type:varchar(64);index:idx_key"`
	Value       string      `json:"value,omitempty" structs:"value,omitempty" mapstructure:"value,omitempty" gorm:"type:varchar(64000);index:idx_value"`
	CreatedAt   int64       `json:"created_at,omitempty" structs:"created_at,omitempty" mapstructure:"created_at,omitempty" gorm:"index:idx_created_at"`
	UpdatedAt   int64       `json:"updated_at,omitempty" structs:"updated_at,omitempty" mapstructure:"updated_at,omitempty" gorm:"index:idx_created_at"`
	Ref1        string      `json:"ref1,omitempty" structs:"ref1,omitempty" mapstructure:"ref1,omitempty" gorm:"type:varchar(64);index:idx_ref1"`
	Ref2        string      `json:"ref2,omitempty" structs:"ref2,omitempty" mapstructure:"ref2,omitempty" gorm:"type:varchar(64);index:idx_ref2"`
	Ref3        string      `json:"ref3,omitempty" structs:"ref3,omitempty" mapstructure:"ref3,omitempty" gorm:"type:varchar(64);index:idx_ref3"`
	Ref4        string      `json:"ref4,omitempty" structs:"ref4,omitempty" mapstructure:"ref4,omitempty" gorm:"type:varchar(64);index:idx_ref4"`
	Ref5        string      `json:"ref5,omitempty" structs:"ref5,omitempty" mapstructure:"ref5,omitempty" gorm:"type:varchar(64);index:idx_ref5"`
	Ref6        string      `json:"ref6,omitempty" structs:"ref6,omitempty" mapstructure:"ref6,omitempty" gorm:"type:varchar(64);index:idx_ref6"`
	Ref7        string      `json:"ref7,omitempty" structs:"ref7,omitempty" mapstructure:"ref7,omitempty" gorm:"type:varchar(64);index:idx_ref7"`
	Ref8        string      `json:"ref8,omitempty" structs:"ref8,omitempty" mapstructure:"ref8,omitempty" gorm:"type:varchar(64);index:idx_ref8"`
	Ref9        string      `json:"ref9,omitempty" structs:"ref9,omitempty" mapstructure:"ref9,omitempty" gorm:"type:varchar(64);index:idx_ref9"`
	Ref10       string      `json:"ref10,omitempty" structs:"ref10,omitempty" mapstructure:"ref10,omitempty" gorm:"type:varchar(64);index:idx_ref10"`
	QueryParams QueryParams `json:"-" structs:"-" mapstructure:"-" gorm:"-"`
}

// NewObject creates a new object. An ID is assigned.
func NewObject() *Object {
	return &Object{ID: util.UUID4(), CreatedAt: time.Now().UnixNano(), UpdatedAt: time.Now().UnixNano()}
}

// Query represents an object for querying objects
type Query interface {
	// GetQueryParams returns the special query parameters
	GetQueryParams() *QueryParams
}

// Expr describes a query expression
type Expr struct {
	Expr string
	Args []interface{}
}

// QueryParams represents object query options
type QueryParams struct {
	Expr    Expr   `json:"-" structs:"-" mapstructure:"-" gorm:"-"`
	OrderBy string `json:"-" structs:"-" mapstructure:"-" gorm:"-"`
	Limit   int    `json:"-" structs:"-" mapstructure:"-" gorm:"-"`
}

// GetQueryParams returns the query parameters attached to the object
func (o *Object) GetQueryParams() *QueryParams {
	return &o.QueryParams
}
