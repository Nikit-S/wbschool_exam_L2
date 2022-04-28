package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (event *Event) Validate(val url.Values) error {
	//fmt.Println(val)

	//провериь что всё цифры в юзер айди
	if _, ok := val["user_id"]; ok {
		for _, v := range event.UserId {
			if v < '0' || v > '9' {
				return fmt.Errorf("User id")
			}
		}
	}

	//проверить нащвание
	if _, ok := val["name"]; ok {
		if len(event.Name) == 0 || len(event.Name) > 100 {
			return fmt.Errorf("invalid name")
		}
	}

	//проверка айди
	if _, ok := val["id"]; ok {
		if event.Id <= 0 {
			return fmt.Errorf("invalid id")
		}
	}
	if _, ok := val["date"]; ok {
		if time.Time(event.Date).Equal(time.Time{}) {
			return fmt.Errorf("invalid time")
		}
	}
	fmt.Println(event)
	return nil
}

func (event *Event) GenerateByForm(r *http.Request) error {

	t, _ := time.Parse("2006-01-02", r.FormValue("date"))
	event.Date = ISOtime(t)
	event.UserId = r.FormValue("user_id")
	event.Name = r.FormValue("name")
	event.Id, _ = strconv.Atoi(r.FormValue("id"))
	return event.Validate(r.Form)
}

func (event *Event) GenerateByUrl(r *http.Request) error {
	var obj Event
	obj.UserId = r.URL.Query().Get("user_id")
	obj.Name = r.URL.Query().Get("name")
	t, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		return err
	}
	obj.Date = ISOtime(t)
	return nil
}

func (t ISOtime) String() string {
	return time.Time(t).String()
}
func (c *ISOtime) UnmarshalJSON(b []byte) error {
	fmt.Println(string(b))
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

func (m *Db) Update(event Event, val url.Values) ([]byte, error) {

	res := Result{}
	_, ok := m.Storage[event.Id]
	if !ok {
		return []byte{}, fmt.Errorf("no such id found")
	}
	if _, ok := val["date"]; !ok {
		event.Date = m.Storage[event.Id].Date
	}
	if _, ok := val["user_id"]; !ok {
		event.UserId = m.Storage[event.Id].UserId
	}
	if _, ok := val["name"]; !ok {
		event.Name = m.Storage[event.Id].Name
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

func (m *Db) GetGroupFromTo(b time.Time, a time.Time, obj Event) ([]byte, error) {
	ret := AggrResult{}
	tempRes := Result{}
	for key := range m.Storage {
		temp := time.Time(m.Storage[key].Date)
		if ((b.Before(temp) && a.After(temp)) || b == temp || a == temp) &&
			obj.UserId == m.Storage[key].UserId {
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
