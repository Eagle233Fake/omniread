package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email,omitempty" json:"email"`
	Phone     string             `bson:"phone,omitempty" json:"phone"`
	Password  string             `bson:"password" json:"-"` // Don't return password in JSON
	Gender    string             `bson:"gender,omitempty" json:"gender"`
	Birthdate *time.Time         `bson:"birthdate,omitempty" json:"birthdate"`
	Avatar    string             `bson:"avatar,omitempty" json:"avatar"`
	Bio       string             `bson:"bio,omitempty" json:"bio"`
	Status    string             `bson:"status" json:"status"` // active, inactive
	// Extended fields
	OAuthID         string     `bson:"oauth_id,omitempty" json:"oauth_id,omitempty"`
	Provider        string     `bson:"provider,omitempty" json:"provider,omitempty"`
	EmailVerified   bool       `bson:"email_verified" json:"email_verified"`
	PhoneVerified   bool       `bson:"phone_verified" json:"phone_verified"`
	ResetToken      string     `bson:"reset_token,omitempty" json:"-"`
	ResetTokenExpAt *time.Time `bson:"reset_token_exp_at,omitempty" json:"-"`

	CreatedAt   time.Time  `bson:"createdAt" json:"created_at"`
	UpdatedAt   time.Time  `bson:"updatedAt" json:"updated_at"`
	LastLoginAt *time.Time `bson:"lastLoginAt,omitempty" json:"last_login_at"`
}
