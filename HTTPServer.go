package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type myData struct {
	Name string
}

type response struct {
	Response string
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
func main() {
	fmt.Println("Starting server on port :3001")

	http.HandleFunc("/time", requestTime)
	http.HandleFunc("/say", requestSay)

	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
