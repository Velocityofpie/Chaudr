package repository

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestCreateRepository(t *testing.T) {
	t.Cleanup(func() {
		os.Remove(sqliteDbName)
	})
	repo, err := NewSqliteRoomRepository()
	if err != nil {
		t.Fatalf("could not create room repository")
	}
	fmt.Println(repo)
}

func TestCreateRoom(t *testing.T) {
	t.Cleanup(func() {
		os.Remove(sqliteDbName)
	})
	repo, err := NewSqliteRoomRepository()
	if err != nil {
		t.Fatalf("could not create room repository")
	}

	room := Room{
		Members: []string{"freeguy"},
	}
	createdRoom, err := repo.CreateRoom(room)
	if err != nil {
		t.Errorf("could not create room: %v", err)
	}

	if createdRoom.Id == room.Id {
		t.Errorf("created room does not have id")
	}

	if len(createdRoom.Members) != len(room.Members) {
		t.Errorf("created room does not the same members. want: %v, got %v", room.Members, createdRoom.Members)
	}

	fmt.Println(createdRoom)
}

func TestAddMemberToRoom(t *testing.T) {
	t.Cleanup(func() {
		os.Remove(sqliteDbName)
	})
	repo, err := NewSqliteRoomRepository()
	if err != nil {
		t.Fatalf("could not create room repository")
	}

	room := Room{
		Members: []string{"freeguy"},
	}
	createdRoom, err := repo.CreateRoom(room)
	if err != nil {
		t.Errorf("could not create room: %v", err)
	}

	if createdRoom.Id == room.Id {
		t.Errorf("created room does not have id")
	}

	if len(createdRoom.Members) != len(room.Members) {
		t.Errorf("created room does not the same members. want: %v, got %v", room.Members, createdRoom.Members)
	}

	newMember := "dodowater"

	updatedRoom, err := repo.AddMemberToRoom(createdRoom, newMember)
	if err != nil {
		t.Fatalf("could not add member to room: %s", err)
	}

	if updatedRoom.Id != createdRoom.Id {
		t.Errorf("updated room does not have same id as created room. want: %v, got %v", createdRoom.Id, updatedRoom.Id)
	}

	if len(createdRoom.Members)+1 != len(updatedRoom.Members) {
		t.Errorf("updated room has unexpected number of members. prev: %v, now: %v", updatedRoom.Members, createdRoom.Members)
	}

	for _, um := range updatedRoom.Members {
		for _, cm := range createdRoom.Members {
			if cm != um && um != newMember {
				t.Errorf("unknown member %s", um)
			}
		}
	}

	if updatedRoom.Id != createdRoom.Id {
		t.Errorf("updated room does not have same id as created room. want: %v, got %v", createdRoom.Id, updatedRoom.Id)
	}

	fmt.Println(updatedRoom)
}

func TestDeleteMemberInRoom(t *testing.T) {
	t.Cleanup(func() {
		os.Remove(sqliteDbName)
	})
	repo, err := NewSqliteRoomRepository()
	if err != nil {
		t.Fatalf("could not create room repository")
	}

	room := Room{
		Members: []string{"freeguy", "dodowater"},
	}
	createdRoom, err := repo.CreateRoom(room)
	if err != nil {
		t.Errorf("could not create room: %v", err)
	}

	if createdRoom.Id == room.Id {
		t.Errorf("created room does not have id")
	}

	if len(createdRoom.Members) != len(room.Members) {
		t.Errorf("created room does not the same members. want: %v, got %v", room.Members, createdRoom.Members)
	}

	deletedMember := "dodowater"

	updatedRoom, err := repo.DeleteMemberFromRoom(createdRoom, deletedMember)
	if err != nil {
		t.Fatalf("could not add member to room: %s", err)
	}

	if updatedRoom.Id != createdRoom.Id {
		t.Errorf("updated room does not have same id as created room. want: %v, got %v", createdRoom.Id, updatedRoom.Id)
	}

	if len(createdRoom.Members)-1 != len(updatedRoom.Members) {
		t.Errorf("created room does not the same members. want: %v, got %v", createdRoom.Members, updatedRoom.Members)
	}

	for _, cm := range createdRoom.Members {
		for _, um := range updatedRoom.Members {
			if cm != um && cm != deletedMember {
				t.Errorf("unknown member %s", um)
			}
		}
	}

	if updatedRoom.Id != createdRoom.Id {
		t.Errorf("updated room does not have same id as created room. want: %v, got %v", createdRoom.Id, updatedRoom.Id)
	}

	fmt.Println(createdRoom)
	fmt.Println(updatedRoom)
}

func TestDeleteMemberInRoomAndReAdd(t *testing.T) {
	t.Cleanup(func() {
		os.Remove(sqliteDbName)
	})
	repo, err := NewSqliteRoomRepository()
	if err != nil {
		t.Fatalf("could not create room repository")
	}

	room := Room{
		Members: []string{"freeguy", "dodowater"},
	}
	createdRoom, err := repo.CreateRoom(room)
	if err != nil {
		t.Errorf("could not create room: %v", err)
	}

	if createdRoom.Id == room.Id {
		t.Errorf("created room does not have id")
	}

	if len(createdRoom.Members) != len(room.Members) {
		t.Errorf("created room does not the same members. want: %v, got %v", room.Members, createdRoom.Members)
	}

	deletedMember := "dodowater"

	updatedRoom, err := repo.DeleteMemberFromRoom(createdRoom, deletedMember)
	if err != nil {
		t.Fatalf("could not delete member in room: %s", err)
	}

	if updatedRoom.Id != createdRoom.Id {
		t.Errorf("updated room does not have same id as created room. want: %v, got %v", createdRoom.Id, updatedRoom.Id)
	}

	if len(createdRoom.Members)-1 != len(updatedRoom.Members) {
		t.Errorf("created room does not the same members. want: %v, got %v", createdRoom.Members, updatedRoom.Members)
	}

	for _, cm := range createdRoom.Members {
		for _, um := range updatedRoom.Members {
			if cm != um && cm != deletedMember {
				t.Errorf("unknown member %s", um)
			}
		}
	}

	if updatedRoom.Id != createdRoom.Id {
		t.Errorf("updated room does not have same id as created room. want: %v, got %v", createdRoom.Id, updatedRoom.Id)
	}

	u2Room, err := repo.AddMemberToRoom(updatedRoom, deletedMember)
	if err != nil {
		t.Fatalf("could not add member to room: %s", err)
	}

	if u2Room.Id != updatedRoom.Id {
		t.Errorf("updated room does not have same id as created room. want: %v, got: %v", createdRoom.Id, updatedRoom.Id)
	}

	if len(updatedRoom.Members)+1 != len(u2Room.Members) {
		t.Errorf("updated room has unexpected number of members. prev: %v, updated: %v", updatedRoom.Members, u2Room.Members)
	}

	for _, um := range updatedRoom.Members {
		for _, u2m := range u2Room.Members {
			if um != u2m && u2m != deletedMember {
				t.Errorf("unknown member %s", um)
			}
		}
	}

	fmt.Println(createdRoom)
	fmt.Println(updatedRoom)
	fmt.Println(u2Room)
}

func TestDeleteAllMembersInRoom(t *testing.T) {
	t.Cleanup(func() {
		os.Remove(sqliteDbName)
	})
	repo, err := NewSqliteRoomRepository()
	if err != nil {
		t.Fatalf("could not create room repository")
	}

	room := Room{
		Members: []string{"freeguy"},
	}
	createdRoom, err := repo.CreateRoom(room)
	if err != nil {
		t.Errorf("could not create room: %v", err)
	}

	if createdRoom.Id == room.Id {
		t.Errorf("created room does not have id")
	}

	if len(createdRoom.Members) != len(room.Members) {
		t.Errorf("created room does not the same members. want: %v, got %v", room.Members, createdRoom.Members)
	}

	updatedRoom, err := repo.DeleteMemberFromRoom(createdRoom, "freeguy")
	if err != nil {
		t.Fatalf("could not delete member in room: %s", err)
	}

	if updatedRoom.Id != createdRoom.Id {
		t.Errorf("updated room does not have same id as created room. want: %v, got %v", createdRoom.Id, updatedRoom.Id)
	}

	if len(createdRoom.Members)-1 != len(updatedRoom.Members) {
		t.Errorf("created room does not the same members. want: %v, got %v", createdRoom.Members, updatedRoom.Members)
	}

	_, err = repo.GetRoom(updatedRoom)
	if err == nil {
		t.Errorf("expected error")
	} else {
		if !errors.Is(err, UnknownRoomId) {
			t.Errorf("unexpected error: %s", err)
		}
	}
}
