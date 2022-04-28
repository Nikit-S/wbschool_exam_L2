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
//метод упрощенного каста ошибки в джейсон нужного формата
func returnErr(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	v, _ := json.Marshal(model.ErrorReuslt{Err: err})
	fmt.Fprint(w, string(v))
}

//func returnRes(w http.ResponseWriter, ret interface{}) {
//	v, _ := json.Marshal(model.Result{Res: ret.(model.Event)})
//	fmt.Fprint(w, string(v))
//}

// метод каста объекта типа Ивент из полученных данных

func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "GET" {
		var obj model.Event
		err := obj.GenerateByUrl(r)
		if err != nil {
			returnErr(w, 500, err.Error())
			return
		}
		year, month, day := (time.Time(obj.Date)).Date()
		//диапозон в с нынешнегодня по следующий минус одна наносекунда
		from := time.Time(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
		to := time.Time(time.Date(year, month, day+1, 0, 0, 0, -1, time.UTC))
		res, err := eventMap.GetGroupFromTo(from, to, obj)
		//fmt.Println(from, to)
		if err != nil {
			returnErr(w, 500, err.Error())
			return
		}
		fmt.Fprint(w, string(res))
		return
	}
	returnErr(w, 400, "Wrong Method")
}
func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "GET" {
		var obj model.Event
		err := obj.GenerateByUrl(r)
		if err != nil {
			returnErr(w, 500, err.Error())
			return
		}
		year, month, day := (time.Time(obj.Date)).Date()
		move := []int{6, 0, 1, 2, 3, 4, 5}
		from := time.Time(time.Date(year, month, day-move[time.Time(obj.Date).Weekday()], 0, 0, 0, 0, time.UTC))
		to := time.Time(time.Date(year, month, day+7-move[time.Time(obj.Date).Weekday()], 0, 0, 0, -1, time.UTC))
		res, err := eventMap.GetGroupFromTo(from, to, obj)
		//fmt.Println(from, to, move[t.Weekday()])
		if err != nil {
			returnErr(w, 500, err.Error())
			return
		}

		fmt.Fprint(w, string(res))
		return
	}
	returnErr(w, 400, "Wrong Method")
}
func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "GET" {
		var obj model.Event
		err := obj.GenerateByUrl(r)
		if err != nil {
			returnErr(w, 500, err.Error())
			return
		}
		year, month, _ := (time.Time(obj.Date)).Date()
		from := time.Time(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC))
		to := time.Time(time.Date(year, month+1, 1, 0, 0, 0, -1, time.UTC))
		res, err := eventMap.GetGroupFromTo(from, to, obj)
		fmt.Println(from, to)
		if err != nil {
			returnErr(w, 500, err.Error())
			return
		}

		fmt.Fprint(w, string(res))
		return
	}
	returnErr(w, 400, "Wrong Method")
}

func validateCreateForm(r *http.Request) (int, error) {
	err := r.ParseForm()
	if err != nil {
		return 400, err
	}
	if _, ok := r.Form["user_id"]; !ok {
		return 400, fmt.Errorf("no user id")
	}
	if _, ok := r.Form["name"]; !ok {
		return 400, fmt.Errorf("no name")
	}
	if _, ok := r.Form["date"]; !ok {
		return 400, fmt.Errorf("no date")
	}
	return 200, nil
}

func createEventHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "POST" {
		var obj model.Event
		status, err := validateCreateForm(r)
		if err != nil {
			returnErr(w, status, err.Error())
			return
		}
		err = obj.GenerateByForm(r)
		if err != nil {
			returnErr(w, 400, err.Error())
			return
		}
		ret, err := eventMap.Add(obj)
		if err != nil {
			returnErr(w, 500, err.Error())
			return
		}
		fmt.Fprint(w, string(ret))
		return
	}
	returnErr(w, 400, "Wrong Method")
}

func validateUpdateForm(r *http.Request) (int, error) {
	var isSmth bool

	err := r.ParseForm()
	if err != nil {
		return 400, err
	}
	if _, ok := r.Form["id"]; !ok {
		return 400, fmt.Errorf("No id provided")
	}
	if _, ok := r.Form["user_id"]; ok {
		isSmth = true
	}
	if _, ok := r.Form["name"]; ok {
		isSmth = true
	}
	if _, ok := r.Form["date"]; ok {
		isSmth = true
	}
	if isSmth {
		return 200, nil
	}
	return 400, fmt.Errorf("No args provided")
}

func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "POST" {
		var obj model.Event
		status, err := validateUpdateForm(r)
		if err != nil {
			returnErr(w, status, err.Error())
			return
		}
		err = obj.GenerateByForm(r)
		if err != nil {
			returnErr(w, 400, err.Error())
			return
		}
		ret, err := eventMap.Update(obj, r.Form)
		if err != nil {
			returnErr(w, 400, err.Error())
			return
		}
		fmt.Fprint(w, string(ret))
		return
	}
	returnErr(w, 400, "Wrong Method")
}

func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "POST" {
		var obj model.Event
		err := r.ParseForm()
		if err != nil {
			returnErr(w, 400, err.Error())
			return
		}
		if _, ok := r.Form["id"]; !ok {
			returnErr(w, 400, "no id provided")
			return
		}
		err = obj.GenerateByForm(r)
		if err != nil {
			returnErr(w, 400, err.Error())
			return
		}
		err = eventMap.Delete(obj)
		if err != nil {
			returnErr(w, 400, err.Error())
			return
		}
		fmt.Fprint(w, `{"result":"deleted"}`)
		return
	}

	returnErr(w, 400, "Wrong Method")
}

func showAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	requests.Println(*r)
	if r.Method == "GET" {
		val, _ := json.Marshal(eventMap)
		fmt.Fprint(w, string(val))
		return
	}
	returnErr(w, 400, "Wrong Method")
}

//хранит данные как кеш\БД
var eventMap model.Db

//отдельный логгер для запросов вк ачестве глобальной
var requests *log.Logger

func main() {

	server := http.Server{}

	//используем файл конфигурации, там только порт
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&server)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	file.Close()

	//инициализируем кеш с индексом
	eventMap.Storage = make(map[int]model.Event)
	eventMap.Index = 1

	//(пере)открываем файл с логированием и вешаем запись на логирование запросов
	logFile, err := os.OpenFile("requests.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer logFile.Close()
	requests = log.New(logFile, "request", 0)

	//по методу на ендпойнт
	http.HandleFunc("/events_for_day", eventsForDayHandler)
	http.HandleFunc("/events_for_week", eventsForWeekHandler)
	http.HandleFunc("/events_for_month", eventsForMonthHandler)
	http.HandleFunc("/create_event", createEventHandler)
	http.HandleFunc("/update_event", updateEventHandler)
	http.HandleFunc("/delete_event", deleteEventHandler)

	http.HandleFunc("/show_all_events", showAllEventsHandler)

	log.Fatal(server.ListenAndServe())
}
