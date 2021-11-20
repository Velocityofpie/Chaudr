package repository

import (
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const sqliteDbName = "gorm.db"

type memberSqlModel struct {
	gorm.Model
	Username string
	RoomID   uint
}

type roomSqlModel struct {
	gorm.Model
	Members []memberSqlModel
}

type sqlRoomRepository struct {
	db *gorm.DB
}

var _ RoomRepository = (*sqlRoomRepository)(nil)

func (s sqlRoomRepository) CreateRoom(room Room) (Room, error) {
	var model roomSqlModel
	for _, member := range room.Members {
		model.Members = append(model.Members, memberSqlModel{
			Username: member,
			RoomID:   room.Id,
		})
	}
	tx := s.db.Create(&model)
	if tx.Error != nil {
		return Room{}, errors.Wrap(tx.Error, "could not create room")
	}
	room.Id = model.ID
	return room, nil
}

func (s sqlRoomRepository) AddMemberToRoom(room Room, member string) (Room, error) {
	model, err := s.getRoomModel(roomSqlModel{
		Model: gorm.Model{
			ID: room.Id,
		},
	})
	if err != nil {
		return room, err
	}
	model.Members = append(model.Members, memberSqlModel{Username: member})
	tx := s.db.Save(&model)
	if tx.Error != nil {
		return room, errors.Wrapf(tx.Error, "could not add member %s to room with id %d", member, room.Id)
	}
	room.Id = model.ID
	room.Members = nil
	for _, m := range model.Members {
		room.Members = append(room.Members, m.Username)
	}
	return room, nil
}

func (s sqlRoomRepository) DeleteMemberFromRoom(room Room, member string) (Room, error) {
	model, err := s.getRoomModel(roomSqlModel{
		Model: gorm.Model{
			ID: room.Id,
		},
	})
	if err != nil {
		return room, err
	}
	for i, m := range model.Members {
		if m.Username == member {
			model.Members[i] = model.Members[len(model.Members)-1]
			model.Members = model.Members[:len(model.Members)-1]
		}
	}
	tx := s.db.Save(&model)
	if tx.Error != nil {
		return room, errors.Wrapf(tx.Error, "could not delete member %s to room with id %d", member, room.Id)
	}
	room.Id = model.ID
	for _, m := range model.Members {
		room.Members = append(room.Members, m.Username)
	}
	return room, nil
}

func (s sqlRoomRepository) getRoomModel(room roomSqlModel) (roomSqlModel, error) {
	tx := s.db.First(&room, room.ID)
	if tx.Error != nil {
		return room, errors.Wrap(tx.Error, "could not fetch model of room")
	}
	return room, nil
}

func (s sqlRoomRepository) GetRoom(room Room) (Room, error) {
	model := roomSqlModel{
		Model: gorm.Model{
			ID: room.Id,
		},
	}
	var err error
	model, err = s.getRoomModel(model)
	if err != nil {
		return room, errors.Wrapf(err, "failed to get room %d", room.Id)
	}
	var r Room
	r.Id = model.ID
	for _, m := range model.Members {
		r.Members = append(r.Members, m.Username)
	}
	return r, nil
}

func NewSqliteRoomRepository() (RoomRepository, error) {
	db, err := gorm.Open(sqlite.Open(sqliteDbName), &gorm.Config{})
	if err := db.AutoMigrate(&roomSqlModel{}); err != nil {
		return nil, errors.Wrap(err, "could automigrate room")
	}
	if err := db.AutoMigrate(&memberSqlModel{}); err != nil {
		return nil, errors.Wrap(err, "could automigrate member")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "could not open sqlite database %s", sqliteDbName)
	}
	return &sqlRoomRepository{
		db: db,
	}, nil
}
