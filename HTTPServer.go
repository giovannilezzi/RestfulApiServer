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

	hostRem     = "172.18.50.67"
	portRem     = 5432
	userRem     = "postgres"
	passwordRem = "postgres"
	dbnameRem   = "blogNoOrm"
)

type myData struct {
	Name string
}

type id struct {
	Id string
}

type channel struct {
	Channel string
}

type response struct {
	Response string
}

type responseArray struct {
	Response []blog
}

type responseFile struct {
	Response Image
}

type responseArrayDB struct {
	Response []blogDB
}

type responseArrayImageDB struct {
	Response []Image
}

type blog struct {
	Id      int
	Name    string
	Surname string
}

type blogDB struct {
	Id     int
	Corpo  string
	Titolo string
}

type Image struct {
	Id       int
	Name     string
	File     []byte
	MimeType string
	Channel  string
}

func getBlog(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
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
	var blogsDB = []blogDB{}
	for results.Next() {
		var tag blogDB
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.Titolo, &tag.Corpo, &tag.Id)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		blogsDB = append(blogsDB, tag)
	}
	fmt.Println(blogsDB)
	resp := responseArrayDB{Response: blogsDB}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func saveBlog(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
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

	var u blog
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err1 := json.NewDecoder(r.Body).Decode(&u)
	if err1 != nil {
		http.Error(w, err1.Error(), 400)
		return
	}
	fmt.Println(u.Id, u.Name, u.Surname)
	insert, err := db.Query("INSERT INTO blog(titolo, corpo) VALUES ( $1, $2 ) LIMIT 3", u.Name, u.Surname)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	resp := response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func saveImage(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRem, portRem, userRem, passwordRem, dbnameRem)
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

	var u Image

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
	fmt.Println(db.Query("INSERT INTO mediarepo(namemedia, file, mimetype, channel) VALUES ($1, $2, $3, $4)", u.Name, u.File, u.MimeType, u.Channel))
	insert, err := db.Query("INSERT INTO mediarepo(namemedia, file, mimetype, channel) VALUES ($1, $2, $3, $4)", u.Name, u.File, u.MimeType, u.Channel)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	resp := response{Response: "ok"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRem, portRem, userRem, passwordRem, dbnameRem)
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

	var channel channel
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
	var images = []Image{}
	for results.Next() {
		var tag Image
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.Id, &tag.Name, &tag.MimeType, &tag.Channel)
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
}

/*
func getAllImage(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRem, portRem, userRem, passwordRem, dbnameRem)
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

func getImageById(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostRem, portRem, userRem, passwordRem, dbnameRem)
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

	var id id
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
	var tag Image
	for results.Next() {
		err = results.Scan(&tag.Id, &tag.Name, &tag.File, &tag.MimeType, &tag.Channel)
	}
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println(tag)
	resp := responseFile{Response: tag}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func main() {
	fmt.Println("Starting server on port :3002")

	http.HandleFunc("/blog", saveBlog)
	http.HandleFunc("/getBlog", getBlog)
	http.HandleFunc("/getImage", getImage)
	//http.HandleFunc("/getAllImage", getAllImage)
	http.HandleFunc("/saveImage", saveImage)
	http.HandleFunc("/getImageById", getImageById)

	err := http.ListenAndServe(":3002", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
