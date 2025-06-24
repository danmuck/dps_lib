package users

import (
	"context"
	"fmt"

	"github.com/danmuck/dps_lib/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func verifyUsername(username string) bool {
	err := TMP_STORAGE.Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(nil)
	return err == mongo.ErrNoDocuments
}

func verifyEmail(email string) bool {
	err := TMP_STORAGE.Collection("users").FindOne(context.Background(), bson.M{"email": email}).Decode(nil)
	return err == mongo.ErrNoDocuments
}

func verifyNew(username, email string) bool {
	filter := bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": email},
		},
	}

	var existingUser User
	err := TMP_STORAGE.Collection("users").FindOne(context.Background(), filter).Decode(&existingUser)
	return err == mongo.ErrNoDocuments
}

func LoginUser(username, password, secret string) (*User, error) {
	user, err := GetUser(username)
	if err != nil {
		logs.Err("failed to get user %s: %v", username, err)
		return nil, fmt.Errorf("failed to get user %s: %v", username, err)
	}
	logs.Debug("user found: %s", username)
	// validate password against stored hash
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		logs.Err("invalid password for user: %s", username)
		return nil, fmt.Errorf("invalid password for user: %s", username)
	}

	return user, nil
}
