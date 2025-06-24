package users

import (
	"context"
	"fmt"

	"github.com/danmuck/dps_lib/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUser(username string) (*User, error) {
	// retrieve the raw map from storage
	var user *User
	err := TMP_STORAGE.Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user %s not found", username)
		}
		logs.Err("failed to find user %s: %v", username, err)
		return nil, fmt.Errorf("failed to find user %s: %v", username, err)
	}
	logs.Debug("got user %s", user.String())
	return user, nil
}
