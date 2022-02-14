package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountType int

const (
	Personal AccountType = iota
	Business
)

type Customer struct {
	ID           *primitive.ObjectID `bson:"_id"`
	Username     *string             `bson:"userName" json:"userName" validate:"userName,required"`
	FullName     *string             `bson:"fullName" json:"fullName"  validate:"required, min=3, max=150"`
	Email        *string             `bson:"email" json:"email" validate:"email,required"`
	Password     *string             `bson:"password,omitempty" json:"password,omitempty" validate:"min=6"`
	AccountInfo  *AccountInfo        `bson:"accountInfo,inline" json:"accountInfo"  validate:"required"`
	Created_at   *time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"date,required"`
	Updated_at   *time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty" validate:"date,required"`
	Deleted_at   *time.Time          `bson:"deletedAt,omitempty" json:"deletedAt,omitempty" validate:"date"`
	Token        *string             `bson:"token" json:"token"`
	SECRET_KEY   *string             `bson:"secretKey" json:"secretKey"  validate:"required"`
	Transactions *[]ThetaTransaction `bson:"transactions,omitempty" json:"transactions,omitempty"`
}

type AccountInfo struct {
	AccountID   *string      `bson:"accountId" json:"accountId"  validate:"required"`
	PinCode     *int         `bson:"pincode" json:"pincode"  validate:"required, min=4,max=6"`
	Balance     *float64     `bson:"balance" json:"balance"  validate:"required"`
	AccountType *AccountType `bson:"accountType" json:"accountType" validate:"required, max=1"`
}
