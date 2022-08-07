package gmopg

import "strconv"

type Card struct {
	Seq     int
	No      string
	Forward string
}

func NewCard(seq int, no, string, forward string) *Card {
	return &Card{Seq: seq, No: no, Forward: forward}
}

func (c *Card) Parse(body map[string]string) error {
	var err error
	if c.Seq, err = strconv.Atoi(body["CardSeq"]); err != nil {
		return err
	}
	c.No = body["CardNo"]
	c.Forward = body["Forward"]
	return nil
}
