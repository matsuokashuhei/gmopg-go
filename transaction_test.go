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
	transaction, err := CreateTransaction(
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
	if transaction == nil {
		t.Errorf("CreateTransaction returns nil")
	}
}
