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
	if member.Id != id {
		t.Errorf("got %s, wanted %s", member.Id, id)
	}
	if member.Name != name {
		t.Errorf("got %s, wanted %s", member.Name, name)
	}
}

func TestCreate2(t *testing.T) {
	member := &Member{}
	ctx := context.Background()
	if err := member.Create(ctx); err != nil {
		t.Errorf("Save returns error, %v", err)
	}
	if member.Id == "" {
		t.Errorf("member.Id is empty, %v", member.Id)
	}
	if member.Name != "" {
		t.Errorf("member.Name is not empty, %v", member.Name)
	}
}

func TestFind(t *testing.T) {
	ctx := context.Background()
	m1 := &Member{Id: cuid.New(), Name: fmt.Sprintf("Test-%s", cuid.New())}
	if err := m1.Create(ctx); err != nil {
		t.Fatalf("Save returns error, %v", err)
	}
	m2, err := FindMember(ctx, m1.Id)
	if err != nil {
		t.Errorf("Find returns error, %v", err)
	}
	if m2 == nil {
		t.Errorf("Find returns nil")
	}
	if m2.Id != m1.Id {
		t.Errorf("got %s, wanted %s", m2.Id, m1.Id)
	}
	if m2.Name != m1.Name {
		t.Errorf("got %s, wanted %s", m2.Name, m1.Name)
	}
	if m2.DeleteFlag != 0 {
		t.Errorf("got %d, wanted %d", m2.DeleteFlag, 0)
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	m1 := &Member{Id: cuid.New(), Name: fmt.Sprintf("Test-%s", cuid.New())}
	if err := m1.Create(ctx); err != nil {
		t.Fatalf("Save returns error, %v", err)
	}
	newName := "New name"
	m1.Name = newName
	if err := m1.Update(ctx); err != nil {
		t.Errorf("Update returns error, %v", err)
	}
	if m1.Name != newName {
		t.Errorf("got %s, wanted %s", m1.Name, newName)
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	m1 := &Member{Id: cuid.New(), Name: fmt.Sprintf("Test-%s", cuid.New())}
	if err := m1.Create(ctx); err != nil {
		t.Fatalf("Save returns error, %v", err)
	}
	if err := m1.Delete(ctx); err != nil {
		t.Errorf("Delete returns error, %v", err)
	}
	_, err := FindMember(ctx, m1.Id)
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
