package gmopg

import (
	"context"
	"fmt"
	"matsuokashuhei/gmopg-go/transaction"
	"testing"

	"github.com/lucsky/cuid"
)

func TestCreateTransaction(t *testing.T) {
	ctx := context.Background()
	member := &Member{Id: cuid.New(), Name: fmt.Sprintf("Test-%s", cuid.New())}
	if err := member.Create(ctx); err != nil {
		t.Fatalf("Save returns error, %v", err)
	}
	card, err := CreateCard(ctx, member.Id, member.Name, "4111111111111111", "2212", "1234")
	if err != nil {
		t.Errorf("CreateCard returns error: %v", err)
	}
	t1, err := CreateTransaction(
		ctx,
		member.Id,
		card.Seq,
		// fmt.Sprintf("order-%s", cuid.New()),
		cuid.New(),
		transaction.AUTH,
		1000,
		100,
	)
	if err != nil {
		t.Errorf("CreateTransaction returns error: %v", err)
	}
	if t1.JobCd != transaction.AUTH {
		t.Errorf("got %v, wanted %v", t1.JobCd, transaction.AUTH)
	}
	if t1.Amount != 1000 {
		t.Errorf("got %d, wanted %d", t1.Amount, 1000)
	}
	if t1.Tax != 100 {
		t.Errorf("got %d, wanted %d", t1.Tax, 100)
	}
}

func TestFindTransaction(t *testing.T) {
	ctx := context.Background()
	member := &Member{Id: cuid.New(), Name: fmt.Sprintf("Test-%s", cuid.New())}
	if err := member.Create(ctx); err != nil {
		t.Fatalf("Save returns error, %v", err)
	}
	card, err := CreateCard(ctx, member.Id, member.Name, "4111111111111111", "2212", "1234")
	if err != nil {
		t.Errorf("CreateCard returns error: %v", err)
	}
	t1, err := CreateTransaction(
		ctx,
		member.Id,
		card.Seq,
		// fmt.Sprintf("order-%s", cuid.New()),
		cuid.New(),
		transaction.AUTH,
		1000,
		100,
	)
	if err != nil {
		t.Errorf("CreateTransaction returns error: %v", err)
	}
	t2, err := FindTransaction(ctx, t1.OrderId)
	if err != nil {
		t.Errorf("FindTransaction returns error: %v", err)
	}
	if t2.Amount != 1000 {
		t.Errorf("got %d, wanted %d", t2.Amount, 1000)
	}
	if t2.Tax != 100 {
		t.Errorf("got %d, wanted %d", t2.Tax, 100)
	}
}
