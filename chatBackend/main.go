package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/gorilla/websocket"
)

var db *gorm.DB
var upgrader = websocket.Upgrader{};

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Email string `gorm:"uniqueIndex"`
	Password string
	Chats []Chat `gorm:"foreignKey:UserId"`
}

type Chat struct {
	gorm.Model
	UserId uint
	User User
	Message string
	RoomId uint	
	Room Room
}

type Room struct {
	gorm.Model
	Chats []Chat `gorm:"foreignKey:RoomId"`
}

func migrate() {
	db.AutoMigrate(&User{}, &Chat{}, &Room{})	
}

func connectDb () error {
	connectionStr := os.Getenv("connectionStr")
	dbClient, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{})
	if(err != nil){
		fmt.Println("error while connecting to db: ", err)
		return err
	}

	db = dbClient
	return nil
}

func Signup (c *gin.Context) {
	var jsonRes struct{
		Email string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := c.BindJSON(&jsonRes);
	if(err != nil){
		fmt.Println("error while parsing body at singup: ", err);
		c.IndentedJSON(500, gin.H{
			"message": "error while parsing body at singup",
			"error": err,
			"success": false,
		})
		return;
	}

	var newUser User;
	newUser.Email = jsonRes.Email;
	newUser.Username = jsonRes.Username;
	newUser.Password = jsonRes.Password;

	db.Create(&newUser)
	fmt.Println("successful signup");
	c.IndentedJSON(500, gin.H{
		"message": "successful signup",
		"success": true,
		"user": struct {
			ID uint
			Username string
			Email string
		}{
			ID: newUser.Model.ID,
			Username: newUser.Username,
			Email: newUser.Email,
		},
	})
}

func GetUser (c *gin.Context) {
	username := c.Params.ByName("username")

	var userByUsername User;
	db.First(&userByUsername, &User{Username: username})

	if(userByUsername.Username != username){
		fmt.Println("no user by username: ", username);
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": fmt.Sprintf("user by username %s not found", username),
		})
		return;
	}

	c.JSON(http.StatusNotFound, gin.H{
		"success": true,
		"user": struct {
			Id uint
			Username string
		}{
			Id: userByUsername.Model.ID,
			Username: userByUsername.Username,
		},
	})
}

var connections []*websocket.Conn

func socket(ctx *gin.Context){
	w := ctx.Writer
	r := ctx.Request

	c, err := upgrader.Upgrade(w, r, nil);
	if(err != nil){
		ctx.JSON(500, gin.H{"message":"failed to upgrade", "success":false})
		return
	}
	connections = append(connections, c);
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage();
		if(err != nil){
			fmt.Println("message retrieval err: ", err);
			ctx.JSON(500, gin.H{"message":"message retrieval err:", "success":false})
			break;
		}

		for _,conn := range connections {
			writeErr := conn.WriteMessage(mt, message);
			if(writeErr != nil){
				fmt.Println("message not send: ", err);
				ctx.JSON(500, gin.H{"message":"message not send", "success":false})
				break;
			}
		}
	}
}

func main () {

	err := godotenv.Load();
	if(err != nil){
		fmt.Println("error while loading env");
		return;
	}
	
	dbErr := connectDb()
	if(dbErr != nil){
		return
	}
	migrate()

	r := gin.Default();

	r.POST("/signup", Signup)
	r.GET("/user/:username", GetUser)
	r.GET("/socket", socket)

	r.Run("localhost:8080")
}