package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	// "github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type character struct {
	ID string `json:"id"`
	NAME string `json:"name"`	
	AGE int32 `json:"age"`
	HP float32 `json:"hp"`
}

var db *sql.DB

func connectDb() {

	connectionStr := os.Getenv("connectionStr");
	dbl, err := sql.Open("postgres", connectionStr);

	if(err != nil){
		fmt.Println("error while connection to db...");
	}

	db = dbl;
	pingErr := db.Ping()
	if(err != nil){
		fmt.Println("err while connecting to db ", pingErr);
	}

	fmt.Println("db connected!")
}
type User struct {
	Id string `db:"id"`
	Username string `db:"username"`
	CreatedAt string `db:"created_at"`
}

func getUsers(ctx *gin.Context) {
	var users []User
	result, err := db.Query("SELECT id, username, created_at from USERS")

	if err != nil {
		fmt.Println("error while fetching, ", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch users"})
		return
	}

	defer result.Close()

	for result.Next() {
		var user User
		if err := result.Scan(&user.Id, &user.Username, &user.CreatedAt); err != nil {
			fmt.Println("get error while scanning")
			
			// The fix is here:
			// You must convert the error to a string before sending it in JSON.
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		// fmt.Println(user)
		users = append(users, user)
	}

	fmt.Println(users)
	ctx.IndentedJSON(http.StatusOK, users)
}


// func getUsers() ([]User, error) {
// 	var users []User

// 	name := "DevAggarwal";
// 	rows, err := db.Query("select id, username, created_at from users")
//     if err != nil {
//         return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
//     }

// 	for rows.Next(){
// 		var user User 
//         if err := rows.Scan(&user.id, &user.username, &user.createdAt); err != nil {
//             return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
//         }
// 		users = append(users, user);
// 	}

// 	return users, nil

// }

func main () {

	err := godotenv.Load();
	if(err != nil){
		fmt.Print("error while loading .env: ", err);
		return;
	}

	connectDb();

	// users, err := getUsers();
	// if err != nil {
	// 	fmt.Println(err);
	// 	return;
	// }
	// fmt.Println(users);

	router := gin.Default();

	router.GET("/users", getUsers)
	router.GET("/getCharacters", getCharacters);
	router.POST("/addCharacter", addCharacter);
	router.GET("/getCharacter/:id", getCharacterById)

	router.Run("localhost:8080")
}

func getCharacterById (ctx *gin.Context){
	id := ctx.Params.ByName("id");
	fmt.Println(id);
	gameChar := gameChars[id];	

	fmt.Println(gameChar);
	if(gameChar.ID == ""){
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "No character found"});
		return;
	}
	ctx.IndentedJSON(http.StatusOK, gameChar);
}

func addCharacter (ctx *gin.Context) {
	var newCharacter character

	if err := ctx.BindJSON(&newCharacter); err != nil {
		fmt.Println("error while binding, ", err);
		return
	}

	gameChars[newCharacter.ID] = newCharacter;
	fmt.Println(gameChars);

	ctx.IndentedJSON(http.StatusOK, gameChars);
	
}

func getCharacters (ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, gameChars);
}

var gameChars = map[string]character{
	"1" : {
		ID: "1",
		NAME: "Alice",
		AGE: 20,
		HP: 1000,
	},
	"2": {
		ID: "2",
		NAME: "Bob",
		AGE: 22,
		HP: 1400,
	},
	"3": {
		ID: "3",
		NAME: "Oscar",
		AGE: 22,
		HP: 1300,
	},
}