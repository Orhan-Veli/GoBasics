package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type MyService interface {
	Get()
	GetAll()
	Create()
	Update()
	Delete()
}

type Person struct {
	Id   int
	Name string
	Age  int
}

func main() {
	router := gin.Default()

	router.GET("/getall", GetAll)
	router.GET("/", HealtCheck)
	router.GET("/getbyid/:id", GetById)
	router.Run(":5005")
}

func GetById(c *gin.Context) {
	id := c.Param("id")
	db, _ := sql.Open("sqlite3", "./SqliteDeneme")
	Persons, _ := db.Query("SELECT * FROM Kisiler k WHERE Id = ?", id)
	var kisi Person
	for Persons.Next() {
		Persons.Scan(&kisi.Id, &kisi.Name, &kisi.Age)
		Persons.Close()
	}
	db.Close()
	uj, _ := json.Marshal(kisi)
	c.String(http.StatusOK, string(uj))
}

func GetAll(c *gin.Context) {
	db, _ := sql.Open("sqlite3", "./SqliteDeneme")
	Persons, _ := db.Query("SELECT * FROM kisiler")
	var kisiler []Person

	for Persons.Next() {
		var kisi Person
		var Id int
		var Name string
		var Age int
		Persons.Scan(&Id, &Name, &Age)
		kisi.Id = Id
		kisi.Name = Name
		kisi.Age = Age
		fmt.Println(kisi)
		kisiler = append(kisiler, kisi)
	}
	Persons.Close()
	db.Close()
	uj, _ := json.Marshal(kisiler)
	c.String(http.StatusOK, string(uj))
}

func HealtCheck(c *gin.Context) {
	c.String(http.StatusOK, "Deneme deneme")
}
