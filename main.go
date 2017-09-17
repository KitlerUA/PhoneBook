package main

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session *mgo.Session

type Person struct {
	ID        string   `json:"id"`
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Address   *Address `json:"address"`
}

type Address struct {
	City  string
	State string
}

var people []Person

func main() {

	//people = append(people, Person{"1", "Volodymyr", "Kit", &Address{"Sincity", "None"}})
	//people = append(people, Person{"2", "Vika", "Vika", &Address{"Sincity", "NoneToo"}})
	ConnectToDB()
	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(&item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	collection := session.DB("test").C("people")
	err := collection.Insert(&person)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(people)
}
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
	}
	collection := session.DB("test").C("people")
	collection.Remove(params["id"])
	json.NewEncoder(w).Encode(people)
}

func ConnectToDB() {
	var err error
	session, err = mgo.Dial("mongodb://localhost:27017")

	if err != nil {
		log.Fatal(err)
	}
	//defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	err = session.DB("test").C("people").Find(bson.M{}).All(&people)
	if err != nil {
		log.Fatal(err)
	}
	/*c := session.DB("test").C("people")
	err = c.Insert(&people[1])
	if err != nil {
		log.Fatal(err)
	}*/
	/*result := Person{}
	err = c.Find(bson.M{"firstname": "Volodymyr"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}*/

	//fmt.Println("Lastname:", result.Lastname)
}
