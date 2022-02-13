package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role int

const (
	Admin Role = iota
	SuperAdmin
)

type User struct {
	ID          *primitive.ObjectID `bson:"_id"`
	FullName    *string             `bson:"fullName" json:"fullName"  validate:"required, min=3, max=150"`
	PhoneNumber *string             `bson:"phoneNumber" json:"phoneNumber" validate:"phonenumber,required"`
	Role        *Role               `bson:"role" json:"role" validate:"required, max=1"`
	Password    *string             `bson:"password,omitempty" json:"password,omitempty" validate:"min=6"`
	Created_at  *time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"date,required"`
	Updated_at  *time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty" validate:"date,required"`
	Deleted_at  *time.Time          `bson:"deletedAt,omitempty" json:"deletedAt,omitempty" validate:"date"`
	Token       *string             `bson:"token" json:"token"`
}
