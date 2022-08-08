package gmopg

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/lucsky/cuid"
)

func TestCreate1(t *testing.T) {
	id := cuid.New()
	name := fmt.Sprintf("Test %s", id)
	member := &Member{Id: id, Name: name}
	ctx := context.Background()
	if err := member.Create(ctx); err != nil {
		t.Errorf("Save returns error, %v", err)
	}
}

func TestCreate2(t *testing.T) {
	member := &Member{}
	ctx := context.Background()
	if err := member.Create(ctx); err != nil {
		t.Errorf("Save returns error, %v", err)
	}
	if len(member.Id) == 0 {
		t.Errorf("got %s, wanted %s", member.Id, "some cuid")
	}
}

func TestFind(t *testing.T) {
	id := cuid.New()
	name := fmt.Sprintf("Test %s", id)
	member := &Member{Id: id, Name: name}
	ctx := context.Background()
	if err := member.Create(ctx); err != nil {
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
	member := &Member{Id: id, Name: name}
	ctx := context.Background()
	if err := member.Create(ctx); err != nil {
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
	member := &Member{Id: id, Name: name}
	ctx := context.Background()
	if err := member.Create(ctx); err != nil {
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
	merrs, ok := err.(*multierror.Error)
	if ok == false {
		t.Errorf("Find does not return multierror.Error")
	}
	if len(merrs.Errors) != 1 {
		t.Errorf("got %d, wanted %d", len(merrs.Errors), 1)
	}
	merr := merrs.Errors[0]
	if merr.Error() != "E01390002" {
		t.Errorf("got %s, wanted %s", merr.Error(), "E01390002")
	}
}
