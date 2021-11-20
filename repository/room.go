package repository

type Room struct {
	Id      uint
	Members []string
}

type RoomRepository interface {
	CreateRoom(room Room) (Room, error)
	AddMemberToRoom(room Room, member string) (Room, error)
	DeleteMemberFromRoom(room Room, member string) (Room, error)
	GetRoom(room Room) (Room, error)
}
