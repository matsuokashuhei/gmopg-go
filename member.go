package gmopg

import (
	"context"
	"io"
	"net/url"
	"regexp"
	"strconv"
)

type Parser[C any] interface {
	Parse(body []byte) error
}

type Member struct {
	Id         string `json:"MemberID"`
	Name       string `json:"MemberName"`
	DeleteFlag int    `json:"DeleteFlag"`
}

func NewMember(id string, name string) *Member {
	return &Member{Id: id, Name: name}
}

func (m *Member) Parse(body []byte) error {
	s := string(body)
	re, _ := regexp.Compile(`MemberID=(?P<MemberID>\w+)&MemberName=(?P<MemberName>.*)&DeleteFlag=(?P<DeleteFlag>[01]{1})`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return nil
	}
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}
	m.Id = result["MemberID"]
	m.Name = result["MemberName"]
	var err error
	if m.DeleteFlag, err = strconv.Atoi(result["DeleteFlag"]); err != nil {
		return err
	}
	return nil
}

func Find(ctx context.Context, id string) (*Member, error) {
	values := url.Values{}
	values.Set("MemberID", id)
	res, err := SearchMember.Call(&values)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	if IsError(body) {
		NewError(body)
	}
	m := &Member{}
	if err := m.Parse(body); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Member) Save(ctx context.Context) error {
	values := url.Values{}
	values.Set("MemberID", m.Id)
	values.Set("MemberName", m.Name)
	res, err := SaveMember.Call(&values)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if IsError(body) {
		return NewError(body)
	}
	return nil
}
