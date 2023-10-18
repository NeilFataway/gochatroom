package services

import (
	"github.com/google/uuid"
	"sync"
)

type User struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar"` //Avatar url
}

var userIdMap sync.Map
var userNameMap sync.Map

func GetUser(userId string) (*User, ResponsiveError) {
	user, ok := userIdMap.Load(userId)
	if ok {
		return user.(*User), nil
	} else {
		return nil, UserNotExists
	}
}

func GetUsers() []*User {
	var usersSlice = make([]*User, 0)
	userIdMap.Range(func(key, value any) bool {
		usersSlice = append(usersSlice, value.(*User))
		return true
	})

	return usersSlice
}

func CreateUser(userName string, avatar string) (*User, ResponsiveError) {
	user := &User{
		UserName: userName,
		Avatar:   avatar,
		UserId:   uuid.NewString(),
	}

	if _, loaded := userNameMap.LoadOrStore(userName, user); loaded {
		// loaded表示用户已存在，则返回报错
		return nil, UserAlreadyExists
	}

	userIdMap.Store(user.UserId, user)
	return user, nil
}

func RemoveUser(userId string) ResponsiveError {
	if _, loaded := userIdMap.LoadAndDelete(userId); loaded {
		return nil
	} else {
		return UserNotExists
	}
}
