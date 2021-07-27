package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"strconv"
)

type Person struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	City      string `json:"city"`
}

var err error
var db *gorm.DB

func init() {
	db, err = gorm.Open("mysql", "root:dg123456@(192.168.11.143)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}

	db.LogMode(true) //
	db.SingularTable(true)
	db.AutoMigrate(Person{})
}

func main() {
	defer db.Close()

	r := gin.Default()
	r.GET("/persons", GetPersons)
	r.GET("/person/:id", GetPerson)
	r.POST("/person", CreatePerson)
	r.PUT("/person/:id", UpdatePerson)
	r.DELETE("/person/:id", DeletePerson)
	r.Run(":8080")
}

func DeletePerson(context *gin.Context) {
	idString := context.Param("id")
	id, _ := strconv.Atoi(idString)
	if id > 0 {
		var person Person
		if err := db.Find(&person, id).Error; err != nil {
			if rows := db.Delete(&person).RowsAffected; rows > 0 {
				context.JSON(http.StatusOK, gin.H{
					"message": "deleted success",
				})
			} else {
				context.JSON(http.StatusOK, gin.H{
					"message": "deleted failed",
				})
			}
		} else {
			context.JSON(http.StatusOK, gin.H{
				"message": "data not exists",
			})
		}
	} else {
		context.JSON(http.StatusOK, gin.H{
			"message": "params error",
		})
	}
}

func GetPerson(c *gin.Context) {
	var person Person
	id := c.Param("id")
	if err := db.Find(&person).Where("id=?", id).Error; err != nil {
		c.JSON(http.StatusOK, err)
	} else {
		c.JSON(http.StatusOK, person)
	}
}

func GetPersons(c *gin.Context) {
	var persons []Person
	if err := db.Find(&persons).Error; err != nil {
		c.JSON(http.StatusOK, err)
	} else {
		c.JSON(http.StatusOK, persons)
	}
}

func CreatePerson(c *gin.Context) {
	var person Person
	c.BindJSON(&person)
	db.Create(&person)
	c.JSON(200, person)
}

func UpdatePerson(c *gin.Context) {
	var person Person
	var id = c.Param("id")
	if err := db.Find(&person, id).Error; err != nil {
		c.JSON(http.StatusOK, nil)
	} else {
		c.BindJSON(&person)
		db.Save(&person)
		c.JSON(http.StatusOK, person)
	}
}
