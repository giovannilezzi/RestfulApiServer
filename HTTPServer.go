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
	Id string
}

type Channel struct {
	Channel string
}

type Response struct {
	Response string
}

type ResponseArray struct {
	Response []Post
}

type ResponseFile struct {
	Response Image
}

type ResponseArrayDB struct {
	Response []PostDB
}

type ResponseArrayImageDB struct {
	Response []Image
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

type Image struct {
	Id       int
	Name     string
	File     []byte
	MimeType string
	Channel  string
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
	fmt.Println(u.Id, u.Titolo, u.Corpo)
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
	insert, err := db.Query("INSERT INTO blog(titolo, corpo) VALUES ( $1, $2 ) LIMIT 3", u.Titolo, u.Corpo)

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

func saveImage(w http.ResponseWriter, r *http.Request) {
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

func getImageByChannel(w http.ResponseWriter, r *http.Request) {
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
	resp := ResponseArrayImageDB{Response: images}
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

func getImageById(w http.ResponseWriter, r *http.Request) {
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
	var tag Image
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

func main() {
	fmt.Println("Starting server on port :3002")

	http.HandleFunc("/savePost", savePost)
	http.HandleFunc("/getAllPosts", getAllPosts)
	http.HandleFunc("/getImage", getImageByChannel)
	http.HandleFunc("/deletePost", deletePost)
	//http.HandleFunc("/getAllImage", getAllImage)
	http.HandleFunc("/saveImage", saveImage)
	http.HandleFunc("/getImageById", getImageById)

	err := http.ListenAndServe(":3002", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
