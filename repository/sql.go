package repository

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const sqliteDbName = "gorm.db"

type memberSqlModel struct {
	gorm.Model
	Username       string
	RoomSqlModelID uint
}

type roomSqlModel struct {
	gorm.Model
	MemberSqlModels []memberSqlModel
}

type sqlRoomRepository struct {
	db *gorm.DB
}

var _ RoomRepository = (*sqlRoomRepository)(nil)

func (s sqlRoomRepository) CreateRoom(room Room) (Room, error) {
	var model roomSqlModel
	for _, member := range room.Members {
		model.MemberSqlModels = append(model.MemberSqlModels, memberSqlModel{
			Username:       member,
			RoomSqlModelID: room.Id,
		})
	}
	tx := s.db.Create(&model)
	if tx.Error != nil {
		return Room{}, errors.Wrap(tx.Error, "could not create room")
	}
	room.Id = model.ID
	room.Members = nil
	for _, m := range model.MemberSqlModels {
		room.Members = append(room.Members, m.Username)
	}
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
	association := s.db.Model(&model).Association("MemberSqlModels")
	if association.Error != nil {
		return room, errors.Wrapf(association.Error, "could not fetch association of room %d", model.ID)
	}
	if txErr := association.Append(&memberSqlModel{Username: member}); txErr != nil {
		return room, errors.Wrapf(txErr, "could not add member %s to room with id %d", member, room.Id)
	}
	model, err = s.getRoomModel(roomSqlModel{
		Model: gorm.Model{
			ID: room.Id,
		},
	})
	if err != nil {
		return room, errors.Wrapf(err, "could not get room %d", room.Id)
	}
	room.Id = model.ID
	room.Members = nil
	for _, m := range model.MemberSqlModels {
		room.Members = append(room.Members, m.Username)
	}
	return room, nil
}

// TODO: when there is no members left in the room, delete the room
func (s sqlRoomRepository) DeleteMemberFromRoom(room Room, member string) (Room, error) {
	model, err := s.getRoomModel(roomSqlModel{
		Model: gorm.Model{
			ID: room.Id,
		},
	})
	if err != nil {
		return room, err
	}
	var toBeDeleteMember memberSqlModel
	var found bool
	for _, m := range model.MemberSqlModels {
		if m.Username == member {
			toBeDeleteMember = m
			found = true
			break
		}
	}
	if !found {
		return room, errors.New(fmt.Sprintf("member %s is not part of room %d", member, room.Id))
	}
	association := s.db.Model(&model).Association("MemberSqlModels")
	if association.Error != nil {
		return room, errors.Wrapf(association.Error, "could not fetch association of room %d", model.ID)
	}
	s.db.Delete(&toBeDeleteMember, toBeDeleteMember.ID)
	model, err = s.getRoomModel(roomSqlModel{
		Model: gorm.Model{
			ID: room.Id,
		},
	})
	if err != nil {
		return room, errors.Wrapf(err, "could not get room %d", room.Id)
	}

	if len(model.MemberSqlModels) == 0 {
		// delete room if no members left
		s.db.Delete(&model)
	}

	room.Id = model.ID
	room.Members = nil
	for _, m := range model.MemberSqlModels {
		room.Members = append(room.Members, m.Username)
	}
	return room, nil
}

func (s sqlRoomRepository) getRoomModel(room roomSqlModel) (roomSqlModel, error) {
	tx := s.db.Model(&room).First(&room, room.ID)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return room, UnknownRoomId
		}
		return room, errors.Wrap(tx.Error, "could not fetch model of room")
	}
	association := s.db.Model(&room).Association("MemberSqlModels")
	if association.Error != nil {
		return room, errors.Wrap(association.Error, "could not fetch model of room")
	}
	var members []memberSqlModel
	if err := association.Find(&members); err != nil {
		return room, errors.New(fmt.Sprintf("could not fetch members of room %d: %v", room.ID, err))
	}
	room.MemberSqlModels = members
	return room, nil
}

func (s sqlRoomRepository) GetRoom(room Room) (Room, error) {
	model, err := s.getRoomModel(roomSqlModel{
		Model: gorm.Model{
			ID: room.Id,
		},
	})
	if err != nil {
		return room, errors.Wrapf(err, "failed to get room %d", room.Id)
	}
	var r Room
	r.Id = model.ID
	r.Members = nil
	for _, m := range model.MemberSqlModels {
		r.Members = append(r.Members, m.Username)
	}
	return r, nil
}

func NewSqliteRoomRepository() (RoomRepository, error) {
	db, err := gorm.Open(sqlite.Open(sqliteDbName), &gorm.Config{})
	if err := db.AutoMigrate(&roomSqlModel{}); err != nil {
		return nil, errors.Wrap(err, "could not automigrate room table")
	}
	if err := db.AutoMigrate(&memberSqlModel{}); err != nil {
		return nil, errors.Wrap(err, "could automigrate member sql table")
	}

	if err != nil {
		return nil, errors.Wrapf(err, "could not open sqlite database %s", sqliteDbName)
	}
	return &sqlRoomRepository{
		db: db,
	}, nil
}
