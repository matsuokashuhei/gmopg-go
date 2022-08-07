package gmopg

import (
	"context"
	"net/url"
	"strconv"
)

type Member struct {
	Id         string
	Name       string
	DeleteFlag int
}

func NewMember(id string, name string) *Member {
	return &Member{Id: id, Name: name}
}

func Find(ctx context.Context, id string) (*Member, error) {
	values := url.Values{}
	values.Set("MemberID", id)
	res, err := SearchMember.Call(&values)
	if err != nil {
		return nil, err
	}
	m := &Member{}
	m.parse(res)
	return m, nil
}

func (m *Member) Save(ctx context.Context) error {
	values := url.Values{}
	values.Set("MemberID", m.Id)
	values.Set("MemberName", m.Name)
	_, err := SaveMember.Call(&values)
	if err != nil {
		return err
	}
	return nil
}

func (m *Member) Update(ctx context.Context) error {
	values := url.Values{}
	values.Set("MemberID", m.Id)
	values.Set("MemberName", m.Name)
	_, err := UpdateMember.Call(&values)
	if err != nil {
		return err
	}
	return nil
}

func (m *Member) Delete(ctx context.Context) error {
	values := url.Values{}
	values.Set("MemberID", m.Id)
	_, err := DeleteMember.Call(&values)
	if err != nil {
		return err
	}
	return nil
}

func (m *Member) RegisterCard(ctx context.Context, token string) (*Card, error) {
	values := url.Values{}
	values.Set("MemberID", m.Id)
	values.Set("Token", token)
	result, err := SaveCard.Call(&values)
	if err != nil {
		return nil, err
	}
	card := &Card{}
	card.Parse(result)
	return card, nil
}

func (m *Member) parse(body map[string]string) error {
	m.Id = body["MemberID"]
	m.Name = body["MemberName"]
	var err error
	if m.DeleteFlag, err = strconv.Atoi(body["DeleteFlag"]); err != nil {
		return err
	}
	return nil
}
