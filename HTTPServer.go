package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "nbs2004"
	dbname   = "postgres"

	hostRemoto     = "172.18.50.67"
	portRemoto     = 5432
	userRemoto     = "postgres"
	passwordRemoto = "postgres"
	dbnameRemoto   = "blogNoOrm"
)

type Id struct {
	Id int
}

type Titolo struct {
	Titolo string
}

type Channel struct {
	Channel string
}

type Search struct {
	Name    string
	Channel string
}

type Response struct {
	Response string
}

type CheckResponse struct {
	Response bool
}

type ResponseArray struct {
	Response []Post
}

type ResponseFile struct {
	Response File
}

type ResponseArrayDB struct {
	Response []PostDB
}

type ResponseArrayFileDB struct {
	Response []File
}

type Post struct {
	Id     int
	Titolo string
	Corpo  string
}

type PostDB struct {
	Id     int
	Corpo  string
	Titolo string
}

type File struct {
	Id       int
	Name     string
	File     []byte
	MimeType string
	Channel  string
}

type UpdateFile struct {
	Id   int
	Name string
}

type UpdatePost struct {
	Id     int
	Titolo string
	Corpo  string
}

type Commento struct {
	Proprietario string
	Commento     string
	Idblog       int
}

type UpdateEventCalendar struct {
	Id        int
	TypeEevnt string
	Event     string
}

type EventCalendar struct {
	Id        int
	Canale    string
	Data      string
	TypeEevnt string
	Event     string
}

type DataInizioFineEventCalendar struct {
	DataInizio string
	DataFine   string
}

type ResponseArrayEventDB struct {
	Response []EventCalendar
}

type DataEventCalendar struct {
	Id   int
	Data string
}

type Questions struct {
	Id              int
	Tipo            string
	Descrizione     string
	Risposta_a      string
	Risposta_b      string
	Risposta_c      string
	Risposta_d      string
	Risposta_esatta string
	Titolo          string
}

type ResponseArrayQuestions struct {
	Response []Questions
}

type RispostaUtente struct {
	Id       int
	Risposta string
}

type Survey struct {
	Titolo      string
	Descrizione string
}

type NewSurvey struct {
	Titolo        string
	Descrizione   string
	VecchioTitolo string
}

type ResponseArrayDataEventCalendar struct {
	Response []DataEventCalendar
}

type ResponseArraySurveys struct {
	Response []Survey
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	results, err := db.Query("SELECT * FROM blog")
	var postArrayDB = []PostDB{}
	for results.Next() {
		var tag PostDB
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.Id, &tag.Titolo, &tag.Corpo)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		postArrayDB = append(postArrayDB, tag)
	}
	fmt.Println(postArrayDB)
	resp := ResponseArrayDB{Response: postArrayDB}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var u Id
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err1 := json.NewDecoder(r.Body).Decode(&u)
	if err1 != nil {
		http.Error(w, err1.Error(), 400)
		return
	}
	fmt.Println(u.Id)
	delete, err := db.Query("DELETE FROM blog where id = $1", u.Id)
	if err != nil {
		panic(err.Error())
		resp := Response{Response: "errore"}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(resp)
		fmt.Println(b)
		fmt.Fprint(w, b)
	}
	defer delete.Close()
	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func savePost(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var u Post
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err1 := json.NewDecoder(r.Body).Decode(&u)
	if err1 != nil {
		http.Error(w, err1.Error(), 400)
		return
	}
	fmt.Println(u.Titolo, u.Corpo)
	insert, err := db.Query("INSERT INTO blog(titolo, corpo) VALUES ( $1, $2 ) ", u.Titolo, u.Corpo)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func editPost(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var u UpdatePost
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err1 := json.NewDecoder(r.Body).Decode(&u)
	if err1 != nil {
		http.Error(w, err1.Error(), 400)
		return
	}
	fmt.Println(u.Id, u.Titolo, u.Corpo)
	update, err := db.Query("UPDATE blog SET titolo=$1, corpo=$2 WHERE id = $3", u.Titolo, u.Corpo, u.Id)
	if err != nil {
		panic(err.Error())
		resp := Response{Response: "errore"}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(resp)
		fmt.Println(b)
		fmt.Fprint(w, b)
	}
	defer update.Close()
	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func addComment(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var com Commento

	err1 := json.NewDecoder(r.Body).Decode(&com)
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	if err1 != nil {
		fmt.Print("Sono Qui")
		panic(err1)
	}
	fmt.Println("U", com.Proprietario, com.Commento, com.Idblog)

	insert, err := db.Query("INSERT INTO commenti(proprietario, commento, idblog) VALUES ($1, $2, $3)", com.Proprietario, com.Commento, com.Idblog)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func saveFile(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var u File

	err1 := json.NewDecoder(r.Body).Decode(&u)
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	if err1 != nil {
		fmt.Print("Sono Qui")
		panic(err1)
	}
	fmt.Println("U", u.Name, u.MimeType, u.Channel)
	//fmt.Println(db.Query("INSERT INTO mediarepo(namemedia, file, mimetype, channel) VALUES ($1, $2, $3, $4)", u.Name, u.File, u.MimeType, u.Channel))
	insert, err := db.Query("INSERT INTO mediarepo(namemedia, file, mimetype, channel) VALUES ($1, $2, $3, $4)", u.Name, u.File, u.MimeType, u.Channel)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func getFileByChannel(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var channel Channel
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err1 := json.NewDecoder(r.Body).Decode(&channel)
	if err1 != nil {
		fmt.Print("Sono Qui")
		panic(err1)
	}
	fmt.Println(channel)

	results, err := db.Query("SELECT id, namemedia, mimetype, channel FROM mediarepo where channel = $1", channel.Channel)
	var files = []File{}
	for results.Next() {
		var tag File
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.Id, &tag.Name, &tag.MimeType, &tag.Channel)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		files = append(files, tag)
	}
	fmt.Println(files)
	resp := ResponseArrayFileDB{Response: files}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

/*
func getAllImage(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	results, err := db.Query("SELECT * FROM mediarepo LIMIT 3")
	var images = []Image{}
	for results.Next() {
		var tag Image
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.Id, &tag.Name, &tag.File, &tag.MimeType, &tag.Channel)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		images = append(images, tag)
	}
	fmt.Println(images)
	resp := responseArrayImageDB{Response: images}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}*/

func getFileById(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	var id Id
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err1 := json.NewDecoder(r.Body).Decode(&id)
	if err1 != nil {
		fmt.Print("Sono Qui")
		panic(err1)
	}
	fmt.Println(id)
	results, err := db.Query("SELECT * FROM mediarepo WHERE id = $1", id.Id)
	var tag File
	for results.Next() {
		err = results.Scan(&tag.Id, &tag.Name, &tag.File, &tag.MimeType, &tag.Channel)
	}
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println(tag)
	resp := ResponseFile{Response: tag}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var u Id
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err1 := json.NewDecoder(r.Body).Decode(&u)
	if err1 != nil {
		http.Error(w, err1.Error(), 400)
		return
	}
	fmt.Println(u.Id)
	delete, err := db.Query("DELETE FROM mediarepo where id = $1", u.Id)
	if err != nil {
		panic(err.Error())
		resp := Response{Response: "errore"}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(resp)
		fmt.Println(b)
		fmt.Fprint(w, b)
	}
	defer delete.Close()
	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func editFile(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var u UpdateFile
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err1 := json.NewDecoder(r.Body).Decode(&u)
	if err1 != nil {
		http.Error(w, err1.Error(), 400)
		return
	}
	fmt.Println(u.Id, u.Name)
	update, err := db.Query("UPDATE mediarepo SET namemedia=$1 WHERE id = $2", u.Name, u.Id)
	if err != nil {
		panic(err.Error())
		resp := Response{Response: "errore"}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(resp)
		fmt.Println(b)
		fmt.Fprint(w, b)
	}
	defer update.Close()
	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func searchFileByChannel(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var search Search
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err1 := json.NewDecoder(r.Body).Decode(&search)
	if err1 != nil {
		fmt.Print("Sono Qui")
		panic(err1)
	}
	fmt.Println(search)

	results, err := db.Query("SELECT id, namemedia, mimetype, channel FROM mediarepo where namemedia LIKE '%' || $1 || '%' AND channel=$2", search.Name, search.Channel)
	var files = []File{}
	for results.Next() {
		var tag File
		// for each row, scan the result into our tag composite object
		fmt.Print("Sono Qui")

		err = results.Scan(&tag.Id, &tag.Name, &tag.MimeType, &tag.Channel)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		files = append(files, tag)
	}
	fmt.Println(files)
	resp := ResponseArrayFileDB{Response: files}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

/*API per il calendario*/
func saveCalendarEvent(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var ev EventCalendar

	err1 := json.NewDecoder(r.Body).Decode(&ev)
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	if err1 != nil {
		fmt.Print("Sono Qui")
		panic(err1)
	}
	fmt.Println("U", ev.Canale, ev.Data, ev.TypeEevnt, ev.Event)

	insert, err := db.Query("INSERT INTO calendario(canale, data, tipo, descrizione) VALUES ($1, $2, $3, $4)", ev.Canale, ev.Data, ev.TypeEevnt, ev.Event)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func getAEventsDay(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var ev EventCalendar

	err1 := json.NewDecoder(r.Body).Decode(&ev)
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	if err1 != nil {
		fmt.Print("Sono Qui")
		panic(err1)
	}
	fmt.Println("U", ev.Canale, ev.Data, ev.TypeEevnt, ev.Event)

	results, err := db.Query("SELECT id, canale, data, tipo, descrizione FROM calendario WHERE data =$1", ev.Data)

	var events = []EventCalendar{}

	for results.Next() {
		var tag EventCalendar

		err = results.Scan(&tag.Id, &tag.Canale, &tag.Data, &tag.TypeEevnt, &tag.Event)
		events = append(events, tag)

	}
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println(events)
	resp := ResponseArrayEventDB{Response: events}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)

}

func getAllAEvents(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var dataInizioFine DataInizioFineEventCalendar

	err1 := json.NewDecoder(r.Body).Decode(&dataInizioFine)
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	if err1 != nil {
		fmt.Print("Sono Qui")
		panic(err1)
	}
	fmt.Println("U", dataInizioFine.DataInizio, dataInizioFine.DataFine)

	results, err := db.Query("SELECT id, data FROM calendario WHERE data >= $1 AND data <= $2", dataInizioFine.DataInizio, dataInizioFine.DataFine)

	var events = []DataEventCalendar{}

	for results.Next() {
		var tag DataEventCalendar

		err = results.Scan(&tag.Id, &tag.Data)
		events = append(events, tag)

	}
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println(events)
	resp := ResponseArrayDataEventCalendar{Response: events}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)

}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var u Id
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err1 := json.NewDecoder(r.Body).Decode(&u)
	if err1 != nil {
		http.Error(w, err1.Error(), 400)
		return
	}
	fmt.Println(u.Id)
	delete, err := db.Query("DELETE FROM calendario where id = $1", u.Id)
	if err != nil {
		panic(err.Error())
		resp := Response{Response: "errore"}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(resp)
		fmt.Println(b)
		fmt.Fprint(w, b)
	}
	defer delete.Close()
	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func editEvent(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var u UpdateEventCalendar
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err1 := json.NewDecoder(r.Body).Decode(&u)
	if err1 != nil {
		http.Error(w, err1.Error(), 400)
		return
	}
	fmt.Println(u.Id, u.TypeEevnt, u.Event)
	update, err := db.Query("UPDATE calendario SET tipo=$1, descrizione=$2 WHERE id = $3", u.TypeEevnt, u.Event, u.Id)
	if err != nil {
		panic(err.Error())
		resp := Response{Response: "errore"}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(resp)
		fmt.Println(b)
		fmt.Fprint(w, b)
	}
	defer update.Close()
	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

/*API per il Self Assessment*/
func saveQuestions(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var qu Questions

	err1 := json.NewDecoder(r.Body).Decode(&qu)
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	if err1 != nil {
		fmt.Print("Sono Qui")
		panic(err1)
	}
	fmt.Println("U", qu.Id, qu.Tipo, qu.Descrizione, qu.Risposta_a, qu.Risposta_b, qu.Risposta_c, qu.Risposta_d, qu.Risposta_esatta)

	insert, err := db.Query("INSERT INTO domande(tipo, descrizione, risposta_a, risposta_b, risposta_c, risposta_d ,risposta_esatta, titolo ) VALUES ($1, $2,$3,$4,$5,$6,$7, $8)", qu.Tipo, qu.Descrizione, qu.Risposta_a, qu.Risposta_b, qu.Risposta_c, qu.Risposta_d, qu.Risposta_esatta, qu.Titolo)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()

	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func saveSurvey(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var su Survey

	err1 := json.NewDecoder(r.Body).Decode(&su)
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	if err1 != nil {
		fmt.Print("Sono Qui")
		panic(err1)
	}
	fmt.Println("U", su.Titolo, su.Descrizione)

	insert, err := db.Query("INSERT INTO questionari(titolo, descrizione) VALUES ($1, $2)", su.Titolo, su.Descrizione)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()

	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func getAllSurveys(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	results, err := db.Query("SELECT titolo, descrizione FROM questionari ")

	var surveys = []Survey{}

	for results.Next() {
		var tag Survey

		err = results.Scan(&tag.Titolo, &tag.Descrizione)
		surveys = append(surveys, tag)

	}
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println(surveys)
	resp := ResponseArraySurveys{Response: surveys}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func deleteSurvey(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var t Titolo
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err1 := json.NewDecoder(r.Body).Decode(&t)
	if err1 != nil {
		http.Error(w, err1.Error(), 400)
		return
	}
	fmt.Println(t.Titolo)
	delete, err := db.Query("DELETE FROM questionari where titolo = $1", t.Titolo)
	if err != nil {
		panic(err.Error())
		resp := Response{Response: "errore"}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(resp)
		fmt.Println(b)
		fmt.Fprint(w, b)
	}
	defer delete.Close()
	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func editSurvey(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var su NewSurvey
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err1 := json.NewDecoder(r.Body).Decode(&su)
	if err1 != nil {
		http.Error(w, err1.Error(), 400)
		return
	}
	fmt.Println(su.Titolo, su.Descrizione, su.VecchioTitolo)
	update, err := db.Query("UPDATE questionari SET titolo=$1, descrizione=$2 WHERE titolo = $3", su.Titolo, su.Descrizione, su.VecchioTitolo)
	if err != nil {
		panic(err.Error())
		resp := Response{Response: "errore"}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(resp)
		fmt.Println(b)
		fmt.Fprint(w, b)
	}
	defer update.Close()
	resp := Response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func getSurvey(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var qu Questions

	err1 := json.NewDecoder(r.Body).Decode(&qu)
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	if err1 != nil {
		fmt.Print("Sono Qui")
		panic(err1)
	}
	fmt.Println("U", qu.Titolo, qu.Descrizione)

	results, err := db.Query("SELECT * FROM domande WHERE titolo =$1", qu.Titolo)

	var questions = []Questions{}

	for results.Next() {
		var tag Questions

		err = results.Scan(&tag.Id, &tag.Tipo, &tag.Descrizione, &tag.Risposta_a, &tag.Risposta_b, &tag.Risposta_c, &tag.Risposta_d, &tag.Risposta_esatta, &tag.Titolo)
		questions = append(questions, tag)

	}
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println(questions)
	resp := ResponseArrayQuestions{Response: questions}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)

}

func checkResponse(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRemoto, portRemoto, userRemoto, passwordRemoto, dbnameRemoto)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var ri RispostaUtente

	err1 := json.NewDecoder(r.Body).Decode(&ri)
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	if err1 != nil {
		fmt.Print("Sono Qui")
		panic(err1)
	}
	fmt.Println("U", ri.Id, ri.Risposta)

	results, err := db.Query("SELECT risposta_esatta FROM domande WHERE id=$1", ri.Id)

	var rispostaEsatta string

	for results.Next() {
		err = results.Scan(&rispostaEsatta)
	}

	var correzzione bool = false
	if ri.Risposta == rispostaEsatta {
		correzzione = true
	} else {
		correzzione = false
	}

	if err != nil {
		panic(err.Error())
		resp := Response{Response: "errore"}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(resp)
		fmt.Println(b)
		fmt.Fprint(w, b)
	}
	defer results.Close()
	resp := CheckResponse{Response: correzzione}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func main() {
	fmt.Println("Starting server on port :3002")

	http.HandleFunc("/savePost", savePost)
	http.HandleFunc("/getAllPosts", getAllPosts)
	http.HandleFunc("/getImage", getFileByChannel)
	http.HandleFunc("/deletePost", deletePost)
	//http.HandleFunc("/getAllImage", getAllImage)
	http.HandleFunc("/saveFile", saveFile)
	http.HandleFunc("/getFileById", getFileById)
	http.HandleFunc("/deleteFile", deleteFile)
	http.HandleFunc("/editFile", editFile)
	http.HandleFunc("/editPost", editPost)
	http.HandleFunc("/addComment", addComment)
	http.HandleFunc("/searchFileByChannel", searchFileByChannel)
	http.HandleFunc("/saveCalendarEvent", saveCalendarEvent)
	http.HandleFunc("/getAEventsDay", getAEventsDay)
	http.HandleFunc("/getAllAEvents", getAllAEvents)
	http.HandleFunc("/deleteEvent", deleteEvent)
	http.HandleFunc("/editEvent", editEvent)
	http.HandleFunc("/saveQuestions", saveQuestions)
	http.HandleFunc("/saveSurvey", saveSurvey)
	http.HandleFunc("/getAllSurveys", getAllSurveys)
	http.HandleFunc("/deleteSurvey", deleteSurvey)
	http.HandleFunc("/editSurvey", editSurvey)
	http.HandleFunc("/getSurveyQuestions", getSurvey)
	http.HandleFunc("/checkResponse", checkResponse)

	err := http.ListenAndServe(":3002", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
