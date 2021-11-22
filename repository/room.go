package repository

import "errors"

type Room struct {
	Id      uint
	Members []string
}

var UnknownRoomId = errors.New("unknown room id")

type RoomRepository interface {
	CreateRoom(room Room) (Room, error)
	AddMemberToRoom(room Room, member string) (Room, error)
	DeleteMemberFromRoom(room Room, member string) (Room, error)
	GetRoom(room Room) (Room, error)
}
