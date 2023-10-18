package services

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

var roomNameMap sync.Map
var roomIdMap sync.Map

type Room struct {
	RoomId    string `json:"room_id"`
	Name      string `json:"name"`
	Owner     *User  `json:"owner"`
	sessions  map[string]*Session
	CreatedAt time.Time `json:"created_at"`

	//操作member时加锁"
	lock *sync.RWMutex
}

func (r *Room) JoinSession(session *Session) ResponsiveError {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.sessions[session.UserId]; ok {
		return userAlreadyJoined
	} else {
		r.sessions[session.UserId] = session
		return nil
	}
}

func (r *Room) RemoveSession(session *Session) ResponsiveError {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.sessions[session.UserId]; !ok {
		return userNotInRoom
	} else {
		// 调用session的注销函数
		_ = session.ShutDown()
		delete(r.sessions, session.UserId)
		return nil
	}
}

func (r *Room) clearSession() ResponsiveError {
	r.lock.Lock()
	defer r.lock.Unlock()

	var errs = make([]error, 0)

	wg := sync.WaitGroup{}
	wg.Add(len(r.sessions))
	for _, session := range r.sessions {
		s := session
		go func() {
			if err := s.ShutDown(); err != nil {
				errs = append(errs, err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	if len(errs) != 0 {
		return WrapBaseError(errs[len(errs)-1])
	} else {
		return nil
	}
}

func (r *Room) Broadcast(message *Message) (err ResponsiveError) {
	if message == nil {
		return nil
	}
	r.lock.RLock()
	defer r.lock.RUnlock()

	var errs = make([]error, 0)

	wg := sync.WaitGroup{}
	wg.Add(len(r.sessions))
	for _, session := range r.sessions {
		s := session
		go func() {
			if err := s.Send(message); err != nil {
				errs = append(errs, err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	if len(errs) != 0 {
		return WrapBaseError(errs[len(errs)-1])
	} else {
		return nil
	}
}

func CreateRoom(roomName string, owner *User) (*Room, ResponsiveError) {
	room := &Room{
		RoomId:    uuid.NewString(),
		Name:      roomName,
		Owner:     owner,
		sessions:  make(map[string]*Session),
		CreatedAt: time.Now(),
		lock:      &sync.RWMutex{},
	}

	if _, loaded := roomNameMap.LoadOrStore(roomName, room); loaded {
		// loaded表示房间名已存在，则返回报错
		return nil, roomConflict
	}

	roomIdMap.Store(room.RoomId, room)
	return room, nil
}

func GetRoom(roomId string) (*Room, ResponsiveError) {
	room, ok := roomIdMap.Load(roomId)
	if ok {
		return room.(*Room), nil
	} else {
		return nil, roomNotExist

	}
}

func GetRooms() []*Room {
	var roomsSlice = make([]*Room, 0)
	roomNameMap.Range(func(key, value any) bool {
		roomsSlice = append(roomsSlice, value.(*Room))
		return true
	})

	return roomsSlice
}

func DeleteRoom(roomId string, userId string) ResponsiveError {
	roomAny, ok := roomIdMap.Load(roomId)
	room := roomAny.(*Room)
	if ok {
		if userId != room.Owner.UserId {
			return PermissionNotAllowd
		}
		userIdMap.Delete(roomId)
		userNameMap.Delete(room.Name)
		return room.clearSession()
	} else {
		return roomNotExist
	}
}

func GetRoomUsers(roomId string) ([]*User, ResponsiveError) {
	roomAny, ok := roomIdMap.Load(roomId)
	if !ok {
		return nil, roomNotExist
	}
	room := roomAny.(*Room)

	room.lock.RLock()
	defer room.lock.RUnlock()

	var users = make([]*User, 0)
	for _, session := range room.sessions {
		users = append(users, session.User)
	}

	return users, nil
}
