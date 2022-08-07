package gmopg

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestNewMember(t *testing.T) {
	id := "test-1"
	name := "Test 1"
	member := NewMember(id, name)
	if id != member.Id {
		t.Errorf("got %s, wanted %s", member.Id, id)
	}
	if name != member.Name {
		t.Errorf("got %s, wanted %s", member.Name, name)
	}
}

func TestSave(t *testing.T) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	id := uuid.String()
	name := fmt.Sprintf("Test %s", id)
	member := NewMember(id, name)
	ctx := context.Background()
	if err = member.Save(ctx); err != nil {
		panic(err)
	}
}

func TestFind(t *testing.T) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	id := uuid.String()
	name := fmt.Sprintf("Test %s", id)
	member := NewMember(id, name)
	ctx := context.Background()
	if err = member.Save(ctx); err != nil {
		panic(err)
	}
	ctx = context.Background()
	member, err = Find(ctx, id)
	if err != nil {
		t.Errorf("Find returns error, %v", err)
	}
	if member.Id != id {
		t.Errorf("got %s, wanted %s", member.Id, id)
	}
	if member.Name != name {
		t.Errorf("got %s, wanted %s", member.Name, name)
	}
	if member.DeleteFlag != 0 {
		t.Errorf("got %d, wanted %d", member.DeleteFlag, 0)
	}
}
