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
)

type myData struct {
	Name string
}

type response struct {
	Response string
}

type responseArray struct {
	Response []blog
}

type responseArrayDB struct {
	Response []blogDB
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
	Name string
	File []byte
}

func requestTime(w http.ResponseWriter, r *http.Request) {
	u := myData{Name: "Pippo"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	fmt.Println(b)
	fmt.Fprint(w, b)
}

func requestSay(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	var u myData
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Println(u.Name)
	resp := response{Response: u.Name}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
	/*
		if u.Name != "" {
			fmt.Fprintf(w, "Hello %s!", u.Name)
		} else {
			fmt.Fprintf(w, "Hello ... you.")
		}*/
}

/*
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var t myData
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	log.Println(t.Name)
	if t.Name != "" {
		fmt.Fprintf(w, "Hello %s!", t.Name)
	} else {
		fmt.Fprintf(w, "Hello ... you.")
	}

}
*/

/*
func getBlog (w http.ResponseWriter, r *http.Request){
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	var blogs = []blog{}
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	blogs = append(blogs, blog{Id:1, Name:"Pippo", Surname:"Pluto"}, blog{Id:2, Name:"Pippo", Surname:"Pluto"})
	fmt.Println(blogs)
	resp := responseArray{Response: blogs}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	fmt.Println(b)
	fmt.Fprint(w, b)
}
*/

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
	insert, err := db.Query("INSERT INTO blog(titolo, corpo) VALUES ( $1, $2 )", u.Name, u.Surname)

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

	var u Image

	err1 := json.NewDecoder(r.Body).Decode(&u)
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(u)
	fmt.Println("U", u.Name, u.File)
	insert, err := db.Query("INSERT INTO mediarepo(namemedia, file) VALUES ($1, $2)", u.Name, u.File)

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

	/*fmt.Println(r.FormValue("image"))
	fmt.Println(r.FormValue("objArr"))
	fmt.Println(r.Form.Get("objArr"))

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
	fmt.Println([]byte(r.FormValue("objArr")))
	var byte = []byte(r.FormValue("objArr"))

	//body, err := ioutil.ReadAll(bytes.NewReader(byte))
	if err != nil {
		panic(err)
	}
	//var decoded []interface{}
	//err = json.Unmarshal(body, &decoded)
	//fmt.Print("Decoded", decoded)
	var u Image

	err1 := json.NewDecoder(bytes.NewReader(byte)).Decode(&u)
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(u)
	fmt.Println("U", u.Name)
	/*fmt.Println(u.Id, u.Name, u.File)*/
	//body, err := ioutil.ReadAll(r.Body)
	//log.Println(string(body))
	/*insert, err := db.Query("INSERT INTO mediarepo(namemedia, file) VALUES ('ProvaApi', $1)",r.FormValue("image"))

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
	//fmt.Fprint(w, b)*/

}

func main() {
	fmt.Println("Starting server on port :3001")

	http.HandleFunc("/say", requestSay)
	http.HandleFunc("/blog", saveBlog)
	http.HandleFunc("/getBlog", getBlog)
	http.HandleFunc("/saveImage", saveImage)

	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
