package helper

import (
	"errors"
	"net/http"
)

const UserIDKey = "user_id"

func GetUserID(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}
