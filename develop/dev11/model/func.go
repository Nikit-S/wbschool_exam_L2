package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (t ISOtime) String() string {
	return time.Time(t).String()
}
func (c *ISOtime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value) //parse time
	if err != nil {
		return err
	}
	*c = ISOtime(t) //set result using the pointer
	return nil
}

func (c ISOtime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format("2006-01-02") + `"`), nil
}

func (m *Db) Add(event Event) ([]byte, error) {
	res := Result{}
	id := m.Index
	m.Index++
	event.Id = id
	m.Storage[id] = event
	res.Res = event
	res.Res.Id = id
	return json.Marshal(res)
}

func (m *Db) Update(event Event) ([]byte, error) {

	res := Result{}
	_, ok := m.Storage[event.Id]
	if !ok {
		return []byte{}, fmt.Errorf("no such id found")
	}
	m.Storage[event.Id] = event
	res.Res = event
	return json.Marshal(res)
}

func (m *Db) Get(id int) ([]byte, error) {

	res := Result{}
	v, ok := m.Storage[id]
	if !ok {
		return []byte{}, fmt.Errorf("no such id found")
	}
	res.Res = v
	return json.Marshal(res)
}

func (m *Db) GetGroupFromTo(b time.Time, a time.Time) ([]byte, error) {
	ret := AggrResult{}
	tempRes := Result{}
	for key := range m.Storage {
		temp := time.Time(m.Storage[key].Date)
		if (b.Before(temp) && a.After(temp)) || b == temp || a == temp {
			tempRes.Res = m.Storage[key]
			ret.Res = append(ret.Res, tempRes.Res)
		}
	}
	return json.Marshal(ret)
}

func (m *Db) Delete(event Event) error {
	_, ok := m.Storage[event.Id]
	if !ok {
		return fmt.Errorf("no such id found")
	}
	delete(m.Storage, event.Id)
	return nil
}
