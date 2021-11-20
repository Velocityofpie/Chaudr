package repository

import (
	"fmt"
	"testing"
)

func TestCreateRoom(t *testing.T) {
	repo, err := NewSqliteRoomRepository()
	if err != nil {
		t.Fatalf("could not create room repository")
	}

	room := roomSqlModel{}
	createdRoom, err := repo.CreateRoom(room)
	if err != nil {
		t.Errorf("could not create room: %v", err)
	}

	if createdRoom.CreatedAt == room.CreatedAt {
		t.Errorf("created_at room time did not change")
	}
	if createdRoom.UpdatedAt == room.UpdatedAt {
		t.Errorf("updated_at room time did not change")
	}
	if createdRoom.ID == room.ID {
		t.Errorf("room id did not change")
	}

	fmt.Println(createdRoom)
}

func TestAddMemberToRoom(t *testing.T) {
	repo, err := NewSqliteRoomRepository()
	if err != nil {
		t.Fatalf("could not create room repository")
	}

	room := roomSqlModel{}
	createdRoom, err := repo.CreateRoom(room)
	if err != nil {
		t.Errorf("could not create room: %v", err)
	}

	if createdRoom.CreatedAt == room.CreatedAt {
		t.Errorf("created_at room time did not change")
	}
	if createdRoom.UpdatedAt == room.UpdatedAt {
		t.Errorf("updated_at room time did not change")
	}
	if createdRoom.ID == room.ID {
		t.Errorf("room id did not change")
	}

	fmt.Println(createdRoom)
}
