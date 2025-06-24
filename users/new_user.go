package users

import (
	"context"
	"fmt"
	"time"

	"github.com/danmuck/dps_lib/auth"
	"github.com/danmuck/dps_lib/logs"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// make sure the user does not already exist
// then create a new user with the given username, email, and password
// and add them to the database
func NewUser(username, email, password string) (*User, error) {
	if exists := verifyUsername(username); exists {
		return nil, fmt.Errorf("username %s already in use", username)
	}
	if exists := verifyEmail(email); exists {
		return nil, fmt.Errorf("email %s already in use", email)
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	roles := []string{"user"}
	if username == "admin" || username == "dirtpig" || username == "danmuck" {
		logs.Dev("assigning admin role to user: %s", username)
		roles = append(roles, "admin")
	}
	user := &User{
		ID:           primitive.NewObjectID(),
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Roles:        roles,
		Bio:          "Welcome to my office!",
		AvatarURL:    "",
		Token:        "", // will be set after signing
		CreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:    primitive.NewDateTimeFromTime(time.Now()),
	}
	logs.Dev("creating user: %s", user.Username)

	if _, err := TMP_STORAGE.Collection("users").InsertOne(context.Background(), user); err != nil {
		logs.Err("failed to insert user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logs.Dev("user created successfully: %+v", user)
	return user, nil
}
