package services

import (
	"errors"
	"sync"
	"time"
)

var rooms sync.Map

// 常见错误
var roomConflict = errors.New("房间已存在，命名冲突")
var roomNotExist = errors.New("房间不存在")
var userAlreadyJoined = errors.New("冲突啦，用户已在聊天室中")
var userNotInRoom = errors.New("用户不在聊天室里")

type Room struct {
	name      string
	owner     User
	members   map[string]struct{}
	createdAt time.Time

	//操作member时加锁"
	lock *sync.RWMutex
}

func (r *Room) JoinMember(member string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.members[member]; ok {
		return userAlreadyJoined
	} else {
		r.members[member] = struct{}{}
		return nil
	}
}

func (r *Room) RemoveMember(member string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.members[member]; !ok {
		return userNotInRoom
	} else {
		delete(r.members, member)
		return nil
	}
}

func IsChatRoomExists(roomName string) bool {
	_, ok := rooms.Load(roomName)
	if ok {
		return true
	} else {
		return false
	}
}

func CreateRoom(roomName string, owner User) (*Room, error) {
	room := &Room{
		name:      roomName,
		owner:     owner,
		members:   make(map[string]struct{}),
		createdAt: time.Now(),
		lock:      &sync.RWMutex{},
	}

	if _, loaded := rooms.LoadOrStore(roomName, room); loaded {
		// loaded表示房间已存在，则返回nil
		return nil, roomConflict
	}

	return room, nil
}

func JoinRoom(roomName string, userName string) error {
	// 加入聊天室
	room, ok := rooms.Load(roomName)
	if !ok {
		return roomNotExist
	}

	return room.(*Room).JoinMember(userName)
}

func ExitRoom(roomName string, userName string) error {
	// 退出聊天室
	room, ok := rooms.Load(roomName)
	if !ok {
		return roomNotExist
	}

	return room.(*Room).RemoveMember(userName)
}
