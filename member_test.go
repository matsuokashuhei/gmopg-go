package gmopg

import (
	"context"
	"fmt"
	"testing"

	"github.com/lucsky/cuid"
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
	id := cuid.New()
	name := fmt.Sprintf("Test %s", id)
	member := NewMember(id, name)
	ctx := context.Background()
	if err := member.Save(ctx); err != nil {
		t.Errorf("Save returns error, %v", err)
	}
}

func TestFind(t *testing.T) {
	id := cuid.New()
	name := fmt.Sprintf("Test %s", id)
	member := NewMember(id, name)
	ctx := context.Background()
	if err := member.Save(ctx); err != nil {
		t.Fatalf("Save returns error, %v", err)
	}
	ctx = context.Background()
	member, err := Find(ctx, id)
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

func TestUpdate(t *testing.T) {
	id := cuid.New()
	name := fmt.Sprintf("Test %s", id)
	member := NewMember(id, name)
	ctx := context.Background()
	if err := member.Save(ctx); err != nil {
		t.Fatalf("Save returns error, %v", err)
	}
	ctx = context.Background()
	name = "New Name"
	member.Name = name
	if err := member.Update(ctx); err != nil {
		t.Errorf("Update returns error, %v", err)
	}
	if member.Name != name {
		t.Errorf("got %s, wanted %s", member.Name, name)
	}
}

func TestDelete(t *testing.T) {
	id := cuid.New()
	name := fmt.Sprintf("Test %s", id)
	member := NewMember(id, name)
	ctx := context.Background()
	if err := member.Save(ctx); err != nil {
		t.Fatalf("Save returns error, %v", err)
	}
	ctx = context.Background()
	if err := member.Delete(ctx); err != nil {
		t.Errorf("Delete returns error, %v", err)
	}
	_, err := Find(ctx, id)
	if err == nil {
		t.Errorf("Find does not return error")
	}
	if err.Error() != "E01390002" {
		t.Errorf("got %s, wanted %s", err.Error(), "E01390002")
	}
}
