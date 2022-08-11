package gmopg

import (
	"context"
	"testing"
)

func TestCreateCard(t *testing.T) {
	ctx := context.Background()
	member := &Member{}
	if err := member.Create(ctx); err != nil {
		t.Errorf("member.Create(ctx) returns error: %v", err)
	}
	cardInput := &CardInput{No: "4111111111111111", Expire: "202212", SecurityCode: "1234", Holder: "Test"}
	_, err := CreateCard(ctx, member.Id, cardInput)
	if err != nil {
		t.Errorf("CreateCard(ctx, member.Id, cardInput) returns error: %v", err)
	}
}
