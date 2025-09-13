package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func createRoom (c *gin.Context) {
	var jsonBody struct {
		Username string `json:"username"`
	}
	err := c.BindJSON(&jsonBody);
	if(err != nil){
		fmt.Println(jsonBody)
		c.JSON(500, gin.H{"error": err, "sucess":false})
		return;
	}
	fmt.Println("hi af", jsonBody)
	var newRoom Room
	db.Create(&newRoom)
	c.JSON(http.StatusOK, gin.H{"success": true, "roomCrated": newRoom})
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
	c.IndentedJSON(200, gin.H{
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

type userRoom struct {
	Conn *websocket.Conn
	Id uint
	RoomId uint
}

var activeUsers = struct {
	sync.RWMutex
	m map[uint]*userRoom
}{
	m: make(map[uint]*userRoom),
}

var activeRooms = struct {
	sync.RWMutex
	m map[uint][]*userRoom
}{
	m: make(map[uint][]*userRoom),
} 

type socketMsg struct {
	Type string `json:"type"`
	Username string `json:"username"`
	Msg string `json:"msg"` 
	RoomId uint `json:"roomId"`
}

func joinRoom(userId uint, roomId uint, c *websocket.Conn){
	activeUsers.RWMutex.Lock();
	activeRooms.RWMutex.Lock();

	userRoomInfo := activeUsers.m[userId]
	if(userRoomInfo.RoomId == roomId){
		fmt.Println("already in room")
		userRoomInfo.Conn = c;
		for _, users := range activeRooms.m[roomId] {
			if(users.Id == userId){
				users.Conn = c;
			}
		}
	}else{
		userRoomInfo.RoomId = roomId;	
		activeRooms.m[roomId] = append(activeRooms.m[roomId], userRoomInfo);
	}	

	activeUsers.RWMutex.Unlock();
	activeRooms.RWMutex.Unlock();
}

func broadcast(sktMsg *socketMsg, roomId uint, c *websocket.Conn, messageType int) {
	var connection []*websocket.Conn
	activeRooms.RWMutex.RLock();
	for _, userRoomInfo := range activeRooms.m[roomId] {
		if userRoomInfo.Conn != c  {
			connection = append(connection, userRoomInfo.Conn)
		}
	}
	activeRooms.RWMutex.RUnlock();

	for _, conn := range connection {
		err := conn.WriteMessage(messageType, []byte(sktMsg.Msg))
		if(err != nil){
			fmt.Println("err while broadcasting")
		}
	}
}

func socket(ctx *gin.Context){
	w := ctx.Writer
	r := ctx.Request

	username := r.URL.Query().Get("username") 
	if(username == ""){
		ctx.JSON(403, gin.H{"message":"no user found", "success":false})
		return
	}

	var connUser User
	dbError := db.Model(&User{}).
		Select("id", "username").
		Where("username = ?", username).
		First(&connUser).Error

	if(dbError != nil){
		ctx.JSON(403, gin.H{"message":"no user found", "success":false})
		return;
	}


	c, err := upgrader.Upgrade(w, r, nil);
	if(err != nil){
		ctx.JSON(500, gin.H{"message":"failed to upgrade", "success":false})
		return
	}

	activeUsers.RWMutex.Lock();

	user, ok := activeUsers.m[connUser.Model.ID];
	if(!ok){
		activeUsers.m[connUser.Model.ID] = &userRoom{
			Conn: c,
			Id: connUser.Model.ID,
			RoomId: 0,
		};
	}else{
		user.RoomId = 0
		user.Conn = c
	}

	activeUsers.RWMutex.Unlock();

	defer func () {
		var deletedUser userRoom
		activeUsers.Lock()
		deletedUser = *activeUsers.m[connUser.Model.ID]
		delete(activeUsers.m, connUser.Model.ID)
		activeUsers.Unlock()
		activeRooms.Lock()

		roomUsers, ok := activeRooms.m[deletedUser.RoomId];
		if(ok){
			currUsers := make([]*userRoom, 0, len(roomUsers) - 1)
			for _, user := range roomUsers {
				if(user.Id != deletedUser.Id){
					currUsers = append(currUsers, user)	
				}
			}

			activeRooms.m[deletedUser.RoomId] = currUsers;
		}

		activeRooms.Unlock()
	}()

	for {
		mt, message, err := c.ReadMessage();
		if(err != nil){
			fmt.Println("message retrieval err: ", err);
			ctx.JSON(500, gin.H{"message":"message retrieval err:", "success":false})
			break;
		}

		var socketMsgJson socketMsg;
		jsonerr := json.Unmarshal(message, &socketMsgJson)
		if jsonerr != nil {
			log.Printf("Error: Failed to unmarshal JSON: %v", jsonerr)

			log.Printf("Raw message that failed to parse: %s", message) 
			 
		}

		fmt.Println("socketMsgJson", socketMsgJson);

		if socketMsgJson.Type == "join" {
			go joinRoom(connUser.Model.ID, socketMsgJson.RoomId, c)
		}else{
			dbChat := &Chat{
				UserId: connUser.Model.ID,
				Message: socketMsgJson.Msg,
				RoomId: socketMsgJson.RoomId,
			}	

			dbRes := db.Create(&dbChat)
			if dbRes.Error != nil {
				fmt.Println("error while writing the chat to db");
				break;
			}
			go broadcast(&socketMsgJson, socketMsgJson.RoomId, c, mt)
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
	r.POST("/createRoom", createRoom)
	r.Run("localhost:8080")
}