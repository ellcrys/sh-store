package db

import "github.com/ellcrys/util"

// Account represents a entity
type Account struct {
	SN               int64  `json:"-" structs:"-" mapstructure:"-" gorm:"primary_key"`
	ID               string `json:"id,omitempty" structs:"id,omitempty" mapstructure:"id,omitempty" gorm:"type:varchar(36)"`
	FirstName        string `json:"first_name,omitempty" structs:"first_name,omitempty" mapstructure:"first_name,omitempty" gorm:""`
	LastName         string `json:"last_name,omitempty" structs:"last_name,omitempty" mapstructure:"last_name,omitempty" gorm:""`
	Email            string `json:"email,omitempty" structs:"email,omitempty" mapstructure:"email,omitempty" gorm:"index:idx_acct_email"`
	Password         string `json:"password,omitempty" structs:"password,omitempty" mapstructure:"password,omitempty" gorm:"index:idx_acct_pass"`
	Confirmed        bool   `json:"confirmed" structs:"confirmed,omitempty" mapstructure:"confirmed,omitempty" gorm:""`
	ConfirmationCode string `json:"confirmation_code,omitempty" structs:"confirmation_code,omitempty" mapstructure:"confirmation_code,omitempty" gorm:""`
	CreatedAt        string `json:"created_at,omitempty" structs:"created_at,omitempty" mapstructure:"created_at,omitempty" gorm:"type:timestamp"`
	UpdatedAt        string `json:"updated_at,omitempty" structs:"updated_at,omitempty" mapstructure:"updated_at,omitempty" gorm:"type:timestamp"`
}

// Contract defines the structure of a contract app
type Contract struct {
	SN           int64  `json:"-" structs:"-" mapstructure:"-" gorm:"primary_key"`
	ID           string `json:"id,omitempty" structs:"id,omitempty" mapstructure:"id,omitempty" gorm:"type:varchar(36);index:idx_cntr_id"`
	Creator      string `json:"creator,omitempty" structs:"creator,omitempty" mapstructure:"creator,omitempty" gorm:"type:varchar(36);index:idx_cntr_creator"`
	Name         string `json:"name,omitempty" structs:"name,omitempty" mapstructure:"name,omitempty" gorm:"index:idx_cntr_name;unique_index"`
	ClientID     string `json:"client_id,omitempty" structs:"client_id,omitempty" mapstructure:"client_id,omitempty" gorm:"unique_index"`
	ClientSecret string `json:"client_secret,omitempty" structs:"client_secret,omitempty" mapstructure:"client_secret,omitempty" gorm:"index:idx_client_secret"`
	CreatedAt    string `json:"created_at,omitempty" structs:"created_at,omitempty" mapstructure:"created_at,omitempty" gorm:"type:timestamp"`
	UpdatedAt    string `json:"updated_at,omitempty" structs:"updated_at,omitempty" mapstructure:"updated_at,omitempty" gorm:"type:timestamp"`
}

// NewContract creates a new contract. ID and CreatedAt are assigned values.
func NewContract() *Contract {
	return &Contract{ID: util.UUID4()}
}

// NewAccount creates a new account. An ID is assigned.
func NewAccount() *Account {
	return &Account{ID: util.UUID4()}
}

// Bucket represents a logical group for objects
type Bucket struct {
	SN        int64  `json:"-" structs:"-" mapstructure:"-" gorm:"primary_key"`
	ID        string `json:"id,omitempty" structs:"id,omitempty" mapstructure:"id,omitempty" gorm:"type:varchar(36)"`
	Name      string `json:"name,omitempty" structs:"name,omitempty" mapstructure:"name,omitempty" gorm:"type:varchar(36);unique_index"`
	Creator   string `json:"creator,omitempty" structs:"creator,omitempty" mapstructure:"creator,omitempty"`
	Immutable bool   `json:"immutable,omitempty" structs:"immutable,omitempty" mapstructure:"immutable,omitempty"`
	CreatedAt string `json:"created_at,omitempty" structs:"created_at,omitempty" mapstructure:"created_at,omitempty" gorm:"type:timestamp"`
	UpdatedAt string `json:"updated_at,omitempty" structs:"updated_at,omitempty" mapstructure:"updated_at,omitempty" gorm:"type:timestamp"`
}

// NewBucket creates a new bucket. An ID is assigned.
func NewBucket() *Bucket {
	return &Bucket{ID: util.UUID4()}
}

// Mapping represents
type Mapping struct {
	SN        int64  `json:"-" structs:"-" mapstructure:"-" gorm:"primary_key"`
	ID        string `json:"id,omitempty" structs:"id,omitempty" mapstructure:"id,omitempty" gorm:"type:varchar(36)"`
	Bucket    string `json:"bucket,omitempty" structs:"bucket,omitempty" mapstructure:"bucket,omitempty"`
	Creator   string `json:"creator,omitempty" structs:"creator,omitempty" mapstructure:"creator,omitempty"`
	Name      string `json:"name,omitempty" structs:"name,omitempty" mapstructure:"name,omitempty" gorm:"type:varchar(36);unique_index"`
	Mapping   string `json:"mapping,omitempty" structs:"mapping,omitempty" mapstructure:"mapping,omitempty" gorm:"type:varchar(256)"`
	CreatedAt string `json:"created_at,omitempty" structs:"created_at,omitempty" mapstructure:"created_at,omitempty" gorm:"type:timestamp"`
	UpdatedAt string `json:"updated_at,omitempty" structs:"updated_at,omitempty" mapstructure:"updated_at,omitempty" gorm:"type:timestamp"`
}

// NewMapping creates a new mapping. An ID is assigned.
func NewMapping() *Mapping {
	return &Mapping{ID: util.UUID4()}
}

// Object represents a bucket object
type Object struct {
	SN          int64       `json:"-" structs:"-" mapstructure:"-" gorm:"primary_key"`
	ID          string      `json:"id,omitempty" structs:"id,omitempty" mapstructure:"id,omitempty" gorm:"type:varchar(36);unique_index"`
	OwnerID     string      `json:"owner_id,omitempty" structs:"owner_id,omitempty" mapstructure:"owner_id,omitempty" gorm:"index:idx_owner_id"`
	CreatorID   string      `json:"creator_id,omitempty" structs:"creator_id,omitempty" mapstructure:"creator_id,omitempty" gorm:"index:idx_creator_id"`
	Bucket      string      `json:"bucket,omitempty" structs:"bucket,omitempty" mapstructure:"bucket,omitempty" gorm:"index:idx_created_at"`
	Key         string      `json:"key,omitempty" structs:"key,omitempty" mapstructure:"key,omitempty" gorm:"type:varchar(64);index:idx_key"`
	Value       string      `json:"value,omitempty" structs:"value,omitempty" mapstructure:"value,omitempty" gorm:"type:varchar(64000);index:idx_value"`
	CreatedAt   string      `json:"created_at,omitempty" structs:"created_at,omitempty" mapstructure:"created_at,omitempty" gorm:"type:timestamp"`
	UpdatedAt   string      `json:"updated_at,omitempty" structs:"updated_at,omitempty" mapstructure:"updated_at,omitempty" gorm:"type:timestamp"`
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
	return &Object{ID: util.UUID4()}
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
