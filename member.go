package gmopg

import (
	"context"
	"net/url"
	"strconv"

	"matsuokashuhei/gmopg-go/api"

	"github.com/lucsky/cuid"
)

type Member struct {
	Id         string
	Name       string
	DeleteFlag int
}

func FindMember(ctx context.Context, id string) (*Member, error) {
	values := url.Values{}
	values.Set("MemberID", id)
	result, err := api.SearchMember.Call(&values)
	if err != nil {
		return nil, err
	}
	m := &Member{}
	m.parse(result[0])
	return m, nil
}

func CreateMember(ctx context.Context, id string, name string) (*Member, error) {
	member := &Member{Id: id, Name: name}
	err := member.Create(ctx)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (m *Member) Create(ctx context.Context) error {
	if len(m.Id) == 0 {
		m.Id = cuid.New()
	}
	values := url.Values{
		"MemberID":   {m.Id},
		"MemberName": {m.Name},
	}
	_, err := api.SaveMember.Call(&values)
	if err != nil {
		return err
	}
	return nil
}

func (m *Member) Update(ctx context.Context) error {
	values := url.Values{
		"MemberID":   {m.Id},
		"MemberName": {m.Name},
	}
	_, err := api.UpdateMember.Call(&values)
	if err != nil {
		return err
	}
	return nil
}

func (m *Member) Delete(ctx context.Context) error {
	values := url.Values{"MemberID": {m.Id}}
	_, err := api.DeleteMember.Call(&values)
	if err != nil {
		return err
	}
	return nil
}

func (m *Member) CreateCard(ctx context.Context, number string, expiryDate string, securityCode string) (*Card, error) {
	// card, err := CreateCard(ctx, m, number, expiryDate, securityCode)
	// if err != nil {
	// 	return nil, err
	// }
	// FindCard(ctx, m.Id, card.Seq)
	return CreateCard(ctx, m.Id, m.Name, number, expiryDate, securityCode)
}

// func (m *Member) RegisterCard(ctx context.Context, cardInput *CardInput) (*Card, error) {
// 	card, err := CreateCard(ctx, m.Id, cardInput)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return card, nil
// }

func (m *Member) parse(body map[string]*string) error {
	m.Id = *body["MemberID"]
	m.Name = *body["MemberName"]
	deleteFlag, err := strconv.Atoi(*body["DeleteFlag"])
	if err != nil {
		return err
	}
	m.DeleteFlag = deleteFlag
	return nil
}
