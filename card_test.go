package gmopg

import (
	"context"
	"testing"
)

func TestCreateCard(t *testing.T) {
	ctx := context.Background()
	member := &Member{Name: "TestCreateCard"}
	if err := member.Create(ctx); err != nil {
		t.Errorf("member.Create(ctx) returns error: %v", err)
	}
	card, err := CreateCard(ctx, member.Id, member.Name, "4111111111111111", "2212", "1234")
	if err != nil {
		t.Errorf("CreateCard returns error: %v", err)
	}
	if card == nil {
		t.Error("card is nil")
	}
	if card.Seq != 0 {
		t.Errorf("wanted %d, got %d", 0, card.Seq)
	}
	if card.HolderName != "TestCreateCard" {
		t.Errorf("wanted %s, got %s", "TestCreateCard", card.HolderName)
	}
	if card.No != "*************111" {
		t.Errorf("wanted %s, got %s", "*************111", card.No)
	}
	if card.Expire != "2212" {
		t.Errorf("wanted %s, got %s", "2212", card.Expire)
	}
	if card.DefaultFlag != 0 {
		t.Errorf("wanted %d, got %d", 0, card.DefaultFlag)
	}
	if card.DeleteFlag != 0 {
		t.Errorf("wanted %d, got %d", 0, card.DeleteFlag)
	}
}

func TestFindCard(t *testing.T) {
	ctx := context.Background()
	member := &Member{Name: "TestFindCard"}
	if err := member.Create(ctx); err != nil {
		t.Errorf("member.Create returns error: %v", err)
	}
	card, err := CreateCard(ctx, member.Id, member.Name, "4111111111111111", "2212", "1234")
	if err != nil {
		t.Errorf("CreateCard returns error: %v", err)
	}
	card, err = FindCard(ctx, member.Id, card.Seq)
	if err != nil {
		t.Errorf("FindCard returns error: %v", err)
	}
	if card == nil {
		t.Error("card is nil")
	}
	if card.Seq != 0 {
		t.Errorf("wanted %d, got %d", 0, card.Seq)
	}
	if card.HolderName != "TestFindCard" {
		t.Errorf("wanted %s, got %s", "TestFindCard", card.HolderName)
	}
	if card.No != "*************111" {
		t.Errorf("wanted %s, got %s", "*************111", card.No)
	}
	if card.Expire != "2212" {
		t.Errorf("wanted %s, got %s", "2212", card.Expire)
	}
	if card.DefaultFlag != 0 {
		t.Errorf("wanted %d, got %d", 0, card.DefaultFlag)
	}
	if card.DeleteFlag != 0 {
		t.Errorf("wanted %d, got %d", 0, card.DeleteFlag)
	}
}

func TestFindCards(t *testing.T) {
	ctx := context.Background()
	member := &Member{Name: "TestFindCards"}
	if err := member.Create(ctx); err != nil {
		t.Errorf("member.Create returns error: %v", err)
	}
	var err error
	_, err = CreateCard(ctx, member.Id, member.Name, "4111111111111111", "2212", "1234")
	if err != nil {
		t.Errorf("CreateCard returns error: %v", err)
	}
	_, err = CreateCard(ctx, member.Id, member.Name, "2111111111111111", "2312", "5678")
	if err != nil {
		t.Errorf("CreateCard returns error: %v", err)
	}
	cards, err := FindCards(ctx, member.Id)
	if err != nil {
		t.Errorf("CreateCard returns error: %v", err)
	}
	if len(cards) != 2 {
		t.Errorf("wanted: 2, got: %d", len(cards))
	}
	var card *Card
	card = cards[0]
	if card.Seq != 0 {
		t.Errorf("wanted %d, got %d", 0, card.Seq)
	}
	if card.HolderName != "TestFindCards" {
		t.Errorf("wanted %s, got %s", "TestFindCards", card.HolderName)
	}
	if card.No != "*************111" {
		t.Errorf("wanted %s, got %s", "*************111", card.No)
	}
	if card.Expire != "2212" {
		t.Errorf("wanted %s, got %s", "2212", card.Expire)
	}
	if card.DefaultFlag != 0 {
		t.Errorf("wanted %d, got %d", 0, card.DefaultFlag)
	}
	if card.DeleteFlag != 0 {
		t.Errorf("wanted %d, got %d", 0, card.DeleteFlag)
	}
	card = cards[1]
	if card.Seq != 1 {
		t.Errorf("wanted %d, got %d", 1, card.Seq)
	}
	if card.HolderName != "TestFindCards" {
		t.Errorf("wanted %s, got %s", "TestFindCards", card.HolderName)
	}
	if card.No != "*************111" {
		t.Errorf("wanted %s, got %s", "*************111", card.No)
	}
	if card.Expire != "2312" {
		t.Errorf("wanted %s, got %s", "2312", card.Expire)
	}
	if card.DefaultFlag != 0 {
		t.Errorf("wanted %d, got %d", 0, card.DefaultFlag)
	}
	if card.DeleteFlag != 0 {
		t.Errorf("wanted %d, got %d", 0, card.DeleteFlag)
	}
}

func TestDeleteCard(t *testing.T) {
	ctx := context.Background()
	member := &Member{Name: "TestFindCards"}
	if err := member.Create(ctx); err != nil {
		t.Errorf("member.Create returns error: %v", err)
	}
	card1, err := CreateCard(ctx, member.Id, member.Name, "4111111111111111", "2212", "1234")
	if err != nil {
		t.Errorf("CreateCard returns error: %v", err)
	}
	card2, err := CreateCard(ctx, member.Id, member.Name, "2111111111111111", "2312", "5678")
	if err != nil {
		t.Errorf("CreateCard returns error: %v", err)
	}
	if err := DeleteCard(ctx, member.Id, card1.Seq); err != nil {
		t.Errorf("DeleteCard returns error: %v", err)
	}
	cards, err := FindCards(ctx, member.Id)
	if err != nil {
		t.Errorf("CreateCard returns error: %v", err)
	}
	if len(cards) != 1 {
		t.Errorf("wanted: 1, got: %d", len(cards))
	}
	card := cards[0]
	if card.Expire != card2.Expire {
		t.Errorf("wanted: %s, got: %s", card2.Expire, card.Expire)
	}
}
