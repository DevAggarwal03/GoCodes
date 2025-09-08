package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type character struct {
	ID string `json:"id"`
	NAME string `json:"name"`	
	AGE int32 `json:"age"`
	HP float32 `json:"hp"`
}

func main () {
	router := gin.Default();

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