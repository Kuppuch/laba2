package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	Name string
}

var nFlag = flag.String("db", "", "database machine ip")

func main() {
	flag.Parse()
	matched, _ := regexp.MatchString(`\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}`, *nFlag)
	if !matched {
		fmt.Println("ip is bad. Try use flag db. Example --db=1.1.1.1")
		return
	}
	fmt.Println(matched) // true
	err := dbConnect()
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})
	r.LoadHTMLGlob("pages/*.html")

	r.GET("/", getAll)
	r.POST("/", create)

	r.Run()
}

func dbConnect() error {
	dsn := fmt.Sprintf("phpmyadmin:phpmyadminpass@tcp(%v:3306)/phpmyadmin?charset=utf8mb4&parseTime=True&loc=Local", *nFlag)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	err = DB.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	return nil
}

func getAll(c *gin.Context) {
	var users []User
	DB.Find(&users)
	c.HTML(http.StatusOK, "index.html", gin.H{"users": users})
}

func create(c *gin.Context) {
	var user User
	raw, _ := c.GetRawData()
	_ = json.Unmarshal(raw, &user)
	DB.Create(&User{
		Name: user.Name,
	})
	var users []User
	DB.Find(&users)
	c.HTML(http.StatusOK, "index.html", gin.H{"users": users})
}
