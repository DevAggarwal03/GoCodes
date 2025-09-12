package main

import (
	"fmt"
	"net/http"
	"os"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

var dbClient *gorm.DB

func connectDb (connectionStr string) {
	db, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{});
	if(err != nil){
		fmt.Print("error while connection db: ", err);
		return;
	}
	dbClient = db;
}

type User struct {
	gorm.Model
    Username  string    `json:"username" gorm:"unique"`
    Email     string    `json:"email" gorm:"unique"`
    Password  string    `json:"password"`
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&User{});	
}

func main () {
	err := godotenv.Load();
	if(err != nil){
		fmt.Print("error while loading .env: ", err);
		return
	}
	connectionStr := os.Getenv("connectionStr")

	connectDb(connectionStr)	

	migrate(dbClient)

	router := gin.Default()

	router.POST("/signUp", createProfile)
	router.POST("/signIn", logIn)
	
	router.Run("localhost:8080")
}

type creatProfileBody struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func createProfile (c *gin.Context) {

	var credentials creatProfileBody
	
	if err := c.BindJSON(&credentials); err != nil {
		fmt.Println("error while binding json in create profile: ", err);
		c.IndentedJSON(http.StatusBadGateway, gin.H{"error": err});
		return;
	}

	fmt.Println(credentials);
	var newUser User;
	newUser.Email = credentials.Email;
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 16);
	if(err != nil){
		fmt.Println("error hasing the fun: ", err);
		c.IndentedJSON(500, gin.H{"error": err, "message":"err while hashing password", "success":false})
		return;
	}
	newUser.Password = string(hashedBytes);
	newUser.Username = credentials.Username;

	dbClient.Create(&newUser);

	c.IndentedJSON(http.StatusOK, gin.H{"message": "user created!", "success": true, "userId": newUser.Model.ID})
}

type logInBody struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func logIn (c *gin.Context) { 
	var logInCred logInBody

	err := c.BindJSON(&logInCred);
	if err != nil {
		fmt.Println("error occ", err);
		c.IndentedJSON(500, gin.H{"success": false, "message":"error occ", "error": err})
		return
	}
	
	fmt.Println(logInCred.Email);
	var user User
	// dbClient.First(&user);
	dbClient.First(&user, User{Email: logInCred.Email});

	hashErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(logInCred.Password));
	if(hashErr != nil){
		fmt.Println("error occ", err);
		c.IndentedJSON(403, gin.H{"success": false, "message":"error occ", "error": err})
		return
	}

	// gen token

	c.IndentedJSON(http.StatusOK, gin.H{
		"success": true,
		"message": "successfully logged in",
		"user": struct {
			Email string
			Username string
			Id uint 
		}{
			Email: user.Email,
			Username: user.Username,
			Id: user.Model.ID,
		},
	})
	fmt.Println(user);
}