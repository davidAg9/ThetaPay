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
	Number string
}

type VisaTransaction struct {
	CardNo     string
	Cvv        string
	ExpiryDate time.Time
}
type ThetaTransaction struct {
	TxnID         *primitive.ObjectID `bson:"_id" json:"txnId" validate:"required" `
	Amount        *float64            `bson:"amount" json:"amount" validate:"required"`
	Email         *string             `bson:"email" json:"email" validate:"email"`
	Description   *string             `bson:"description,omitempty" json:"description,omitempty"`
	AcountId      *string             `bson:"accountId" json:"accountID" validate:"required"`
	CreditAccount []string            `bson:"creditAccount" json:"creditAccount" validate:"required"`
	Created_at    *time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"date"`
	Updated_at    *time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty" validate:"date"`
	Deleted_at    *time.Time          `bson:"deletedAt,omitempty" json:"deletedAt,omitempty" validate:"date"`
}
