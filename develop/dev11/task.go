package main

import (
	"encoding/json"
	"example/api/model"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "GET" {
		t := time.Now()
		year, month, day := t.Date()
		from := time.Time(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
		to := time.Time(time.Date(year, month, day+1, 0, 0, 0, 0, time.UTC))
		res, err := eventMap.GetGroupFromTo(from, to)
		fmt.Println(from, to)
		if err != nil {
			w.WriteHeader(500)
			v, _ := json.Marshal(model.ErrorReuslt{Err: err.Error()})
			fmt.Fprint(w, string(v))
			return
		}

		fmt.Fprint(w, string(res))
		return
	}
	w.WriteHeader(404)
}
func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "GET" {
		t := time.Now()
		year, month, day := t.Date()
		from := time.Time(time.Date(year, month, day-int(t.Weekday())+1, 0, 0, 0, 0, time.UTC))
		to := time.Time(time.Date(year, month, day+7-int(t.Weekday()), 0, 0, 0, 0, time.UTC))
		res, err := eventMap.GetGroupFromTo(from, to)
		fmt.Println(from, to)
		if err != nil {
			w.WriteHeader(500)
			v, _ := json.Marshal(model.ErrorReuslt{Err: err.Error()})
			fmt.Fprint(w, string(v))
			return
		}

		fmt.Fprint(w, string(res))
		return
	}
	w.WriteHeader(404)
}
func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "GET" {
		t := time.Now()
		year, month, _ := t.Date()
		from := time.Time(time.Date(year, month, 0, 0, 0, 0, 0, time.UTC))
		to := time.Time(time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC))
		res, err := eventMap.GetGroupFromTo(from, to)
		fmt.Println(from, to)
		if err != nil {
			w.WriteHeader(500)
			v, _ := json.Marshal(model.ErrorReuslt{Err: err.Error()})
			fmt.Fprint(w, string(v))
			return
		}

		fmt.Fprint(w, string(res))
		return
	}
	w.WriteHeader(404)
}

func createEventHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "POST" {
		var obj model.Event
		err := json.NewDecoder(r.Body).Decode(&obj)
		if err != nil {
			w.WriteHeader(400)
			v, _ := json.Marshal(model.ErrorReuslt{Err: err.Error()})
			fmt.Fprint(w, string(v))
			return
		}
		ret, err := eventMap.Add(obj)
		if err != nil {
			w.WriteHeader(500)
			v, _ := json.Marshal(model.ErrorReuslt{Err: err.Error()})
			fmt.Fprint(w, string(v))
			return
		}
		fmt.Fprint(w, string(ret))
		return
	}
	w.WriteHeader(400)
}

func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "POST" {
		var obj model.Event
		err := json.NewDecoder(r.Body).Decode(&obj)
		if err != nil {
			w.WriteHeader(400)
			v, _ := json.Marshal(model.ErrorReuslt{Err: err.Error()})
			fmt.Fprint(w, string(v))
			return
		}
		ret, err := eventMap.Update(obj)
		if err != nil {
			w.WriteHeader(400)
			v, _ := json.Marshal(model.ErrorReuslt{Err: err.Error()})
			fmt.Fprint(w, string(v))
			return
		}
		fmt.Fprint(w, string(ret))
		return
	}
	w.WriteHeader(404)
}

func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "POST" {
		var obj model.Event
		err := json.NewDecoder(r.Body).Decode(&obj)
		if err != nil {
			w.WriteHeader(400)
			v, _ := json.Marshal(model.ErrorReuslt{Err: err.Error()})
			fmt.Fprint(w, string(v))
			return
		}
		err = eventMap.Delete(obj)
		if err != nil {
			w.WriteHeader(400)
			v, _ := json.Marshal(model.ErrorReuslt{Err: err.Error()})
			fmt.Fprint(w, string(v))
			return
		}
		fmt.Fprint(w, `{"result":"deleted"}`)
		return
	}
	w.WriteHeader(404)
}

func showAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "GET" {
		val, _ := json.Marshal(eventMap)
		fmt.Fprint(w, string(val))
		return
	}
	w.WriteHeader(404)
}

var eventMap model.Db

var requests *log.Logger

func main() {

	server := http.Server{}

	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&server)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	file.Close()

	eventMap.Storage = make(map[int]model.Event)
	eventMap.Index = 1

	logFile, err := os.OpenFile("requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer logFile.Close()
	requests = log.New(logFile, "request", 0)

	http.HandleFunc("/events_for_day", eventsForDayHandler)
	http.HandleFunc("/events_for_week", eventsForWeekHandler)
	http.HandleFunc("/events_for_month", eventsForMonthHandler)
	http.HandleFunc("/create_event", createEventHandler)
	http.HandleFunc("/update_event", updateEventHandler)
	http.HandleFunc("/delete_event", deleteEventHandler)

	http.HandleFunc("/show_all_events", showAllEventsHandler)

	log.Fatal(server.ListenAndServe())
}
