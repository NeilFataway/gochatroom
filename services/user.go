package services

import (
	"errors"
	"sync"
)

var UserAlreadyExists = errors.New("用户名已存在")
var UserAvatarNotExists = errors.New("用户头像地址不存在")

type User struct {
	userName string
	avatar   string //avatar url
}

var users sync.Map

func GetUser(userName string) {

}
