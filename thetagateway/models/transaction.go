package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ThetaTransaction struct {
	TxnID       *primitive.ObjectID `bson:"_id" json:"txnId" validate:"required" `
	Amount      *float64            `bson:"amount" json:"amount" validate:"required"`
	Email       *string             `bson:"email" json:"email" validate:"email"`
	Description *string             `bson:"description,omitempty" json:"description,omitempty"`
	Created_at  *time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"date"`
	Updated_at  *time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty" validate:"date"`
	Deleted_at  *time.Time          `bson:"deletedAt,omitempty" json:"deletedAt,omitempty" validate:"date"`
}
