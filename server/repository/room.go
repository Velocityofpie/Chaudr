package repository

import (
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const sqliteDbName = "gorm.db"

type Member struct {
	gorm.Model
	Username string
	RoomID   uint
}

type Room struct {
	gorm.Model
	Members []Member
}

type RoomRepository interface {
	CreateRoom(room Room) (Room, error)
	AddMemberToRoom(room Room, member Member) (Room, error)
	DeleteMemberFromRoom(room Room, member Member) (Room, error)
}

type roomRepository struct {
	db *gorm.DB
}

func (r roomRepository) CreateRoom(room Room) (Room, error) {
	tx := r.db.Create(&room)
	if tx.Error != nil {
		return Room{}, errors.Wrap(tx.Error, "could not create room")
	}
	return room, nil
}

func (r roomRepository) AddMemberToRoom(room Room, member Member) (Room, error) {
	room.Members = append(room.Members, member)
	tx := r.db.Save(&room)
	if tx.Error != nil {
		return Room{}, errors.Wrapf(tx.Error, "could not add member %s to room with id %d", member.Username, room.ID)
	}
	return room, nil
}

func (r roomRepository) DeleteMemberFromRoom(room Room, member Member) (Room, error) {
	tx := r.db.Save(&room)
	if tx.Error != nil {
		return Room{}, errors.Wrapf(tx.Error, "could not add member %s to room with id %d", member.Username, room.ID)
	}
	return room, nil
}

func NewSqliteRoomRepository() (RoomRepository, error) {
	db, err := gorm.Open(sqlite.Open(sqliteDbName), &gorm.Config{})
	if err := db.AutoMigrate(&Room{}); err != nil {
		return nil, errors.Wrap(err, "could automigrate room")
	}
	if err := db.AutoMigrate(&Member{}); err != nil {
		return nil, errors.Wrap(err, "could automigrate member")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "could not open sqlite database %s", sqliteDbName)
	}
	return &roomRepository{
		db: db,
	}, nil
}
