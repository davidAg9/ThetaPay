package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountType int

const (
	Personal = iota
	Business
)

type Customers struct {
	ID          *primitive.ObjectID `bson:"_id"`
	Username    *string
	FullName    *string      `bson:"fullName" json:"fullName"  validate:"required, min=3, max=150"`
	Email       *string      `bson:"email" json:"email" validate:"email,required"`
	Password    *string      `bson:"password,omitempty" json:"password,omitempty" validate:"min=6"`
	AccountInfo *AccountInfo `bson:"accountInfo" json:"accountInfo"  validate:"required"`
	Created_at  *time.Time   `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"date,required"`
	Updated_at  *time.Time   `bson:"updatedAt,omitempty" json:"updatedAt,omitempty" validate:"date,required"`
	Deleted_at  *time.Time   `bson:"deletedAt,omitempty" json:"deletedAt,omitempty" validate:"date"`
	Token       *string      `bson:"token" json:"token"`
	SECRET_KEY  *string      `bson:"secretKey" json:"secretKey"  validate:"required"`
}

type AccountInfo struct {
	AccountID   *string      `bson:"accountId" json:"accountId"  validate:"required"`
	PinCode     *int         `bson:"pincode" json:"pincode"  validate:"required, min=4,max=6"`
	Balance     *float64     `bson:"balance" json:"balance"  validate:"required"`
	AccountType *AccountType `bson:"accountType" json:"accountType" validate:"required, max=1"`
}
