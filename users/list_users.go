package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func ListUsers() ([]*User, error) {
	var users []*User
	cursor, err := TMP_STORAGE.Collection("users").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
