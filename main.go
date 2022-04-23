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
	GetById()
	GetAll()
	Create()
	Update()
	Delete()
}

type Person struct {
	Id   int `json:"id"`
	Name string `json:"Name"`
	Age  int `json:"Age"`
}

type UpdatePersonDto struct{
	UpdateId int `json:"UpdateId"`
	Id int `json:"Id"`
	Name string `json:"Name"`
	Age int `json:"Age"`
}

// func init(){
// 	db, _ := sql.Open("sqlite3:127.0.0.1/9000", "SqliteDeneme")
// 	db.Exec("CREATE Database SqliteDeneme ( Id int, Name nvarchar(250), Age int );")
// 	db.Exec("CREATE TABLE Kisiler ( Id int, Name nvarchar(250), Age int );")
// 	fmt.Println(res);
// 	db.Close()
// }

func main() {
	router := gin.Default()

	router.GET("/getall", GetAll)
	router.GET("/", HealtCheck)
	router.GET("/getbyid/:id", GetById)
	router.POST("/",Create)
	router.DELETE(":id",Delete)
	router.PUT("/",Update)
	router.Run(":5005")
}

func Update(c *gin.Context){
	var update UpdatePersonDto
	
	err := c.BindJSON(&update);
	if err != nil{
		c.String(http.StatusBadRequest, err.Error())
	}

	db, _ := sql.Open("sqlite3", "./SqliteDeneme")
	updated, _ := db.Exec("UPDATE Kisiler SET Id=?, Name=?, Age=? WHERE Id=?",update.Id,update.Name,update.Age,update.UpdateId)
	
	data, _ := updated.RowsAffected()
	ser, _ := json.Marshal(data)
	c.String(http.StatusOK,string(ser))

}

func Delete(c *gin.Context){
	id := c.Param("id");

	fmt.Println(id);

	db, err := sql.Open("sqlite3", "./SqliteDeneme")

	if err != nil{
		c.String(http.StatusInternalServerError, err.Error() + "open");
	}

	result, preErr := db.Prepare("DELETE FROM Kisiler WHERE Id =?")

	if preErr != nil{
		c.String(http.StatusInternalServerError, preErr.Error() + "prep")
	}

	_ , errExec :=result.Exec(id);
	if preErr != nil{
		c.String(http.StatusInternalServerError, errExec.Error()+ "exec")
	}

	c.String(http.StatusNoContent,"executed")
}

func GetById(c *gin.Context) {
	id := c.Param("id")
	db, _ := sql.Open("sqlite3", "./SqliteDeneme")
	Persons, _ := db.Query("SELECT * FROM Kisiler WHERE Id =?", id)
	var kisi Person
	for Persons.Next() {
		Persons.Scan(&kisi.Id, &kisi.Name, &kisi.Age)
		Persons.Close()
	}
	db.Close()
	uj, _ := json.Marshal(kisi)
	c.String(http.StatusOK, string(uj))
}

func Create(c *gin.Context){
	var newPerson Person;
	err := c.BindJSON(&newPerson);
	if err != nil {
		c.String(http.StatusBadRequest, err.Error());
	}
	db, _ := sql.Open("sqlite3", "./SqliteDeneme")
	result, _ := db.Exec("INSERT INTO Kisiler (Id,Name,Age) VALUES(?,?,?)", newPerson.Id,newPerson.Name,newPerson.Age);
	db.Close();
	id, _ := result.LastInsertId();

	uj, _ := json.Marshal(id)
	c.String(http.StatusCreated,string(uj));
	
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
