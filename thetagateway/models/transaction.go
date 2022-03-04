package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TxnType string

const (
	Deposit = "deposit"
	Tranfer = "tranfer"
)

type MomoTransaction struct {
	TxnID       *primitive.ObjectID `bson:"_id" json:"txnId,omitempty"  `
	Amount      *float64            `bson:"amount" json:"amount" validate:"required"`
	Email       *string             `bson:"email" json:"email" validate:"email,required"`
	Description *string             `bson:"description,omitempty" json:"description,omitempty"`
	MerchantId  *string             `bson:"merchantId" json:"merchantId,omitempty" `
	Number      *string             `bson:"number" json:"number" validate:"required" `
	Created_at  time.Time           `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"date"`
	Updated_at  time.Time           `bson:"updatedAt,omitempty" json:"-" validate:"date"`
	Deleted_at  time.Time           `bson:"deletedAt,omitempty" json:"-" validate:"date"`
	Trans_Type  TxnType             `bson:"txnType" json:"txnType" validate:"required"`
}

type VisaTransaction struct {
	TxnID       *primitive.ObjectID `bson:"_id" json:"txnId,omitempty"  `
	Amount      *float64            `bson:"amount" json:"amount" validate:"required"`
	Email       *string             `bson:"email" json:"email" validate:"email,required"`
	Description *string             `bson:"description,omitempty" json:"description,omitempty"`
	MerchantId  *string             `bson:"merchantId" json:"merchantId,omitempty" `
	CardNo      *string             `bson:"cardNo" json:"cardNo" validate:"required" `
	Cvv         *string             `bson:"-" json:"cvv"  validate:"required"`
	ExpiryDate  time.Time           `bson:"-" json:"expiryDate"  validate:"required"`
	Created_at  time.Time           `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"date"`
	Updated_at  time.Time           `bson:"updatedAt,omitempty" json:"-" validate:"date"`
	Deleted_at  time.Time           `bson:"deletedAt,omitempty" json:"-" validate:"date"`
	Trans_Type  TxnType             `bson:"txnType" json:"txnType" validate:"required"`
}
type ThetaTransaction struct {
	TxnID         *primitive.ObjectID `bson:"_id" json:"txnId,omitempty"  `
	Amount        *float64            `bson:"amount" json:"amount" validate:"required"`
	Email         *string             `bson:"email" json:"email" validate:"email,required"`
	Description   *string             `bson:"description,omitempty" json:"description,omitempty"`
	AcountId      *string             `bson:"accountId" json:"accountId,omitempty" validate:"required"`
	MerchantId    *string             `bson:"merchantId" json:"merchantId,omitempty" `
	CreditAccount *string             `bson:"creditAccount" json:"creditAccount,omitempty" `
	Created_at    time.Time           `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"date"`
	Updated_at    time.Time           `bson:"updatedAt,omitempty" json:"-" validate:"date"`
	Deleted_at    time.Time           `bson:"deletedAt,omitempty" json:"-" validate:"date"`
	Trans_Type    TxnType             `bson:"txnType" json:"txnType" validate:"required"`
}
