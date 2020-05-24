package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/chainreaction/datastore"
	"github.com/chainreaction/models"
	"github.com/gin-gonic/gin"
)

// JoinExistingGame provides wndpoint to join already created game
func JoinExistingGame(c *gin.Context) {
	var ret gin.H
	roomName := strings.ToLower(c.Param("name"))
	if roomName == "" {
		ret = gin.H{"Error": "Please provide a game instance id"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}
	gInstance, exists := datastore.GetGameInstance(roomName)
	if !exists {
		ret = gin.H{"Error": "No such active game instance found"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusNotFound, ret)
		return
	}
	if gInstance.GetCurrentActivePlayers() == gInstance.PlayersCount {
		ret = gin.H{"Error": "Game is already full."}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusForbidden, ret)
		return
	}
	username := strings.ToLower(c.Query("username"))
	if username == "" {
		ret = gin.H{"Error": "username cannot be empty."}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}
	color := c.Query("color")
	if color == "" {
		ret = gin.H{"Error": "color cannot be empty"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}
	_, pExists := gInstance.GetPlayerByUsername(username)
	if pExists {
		ret = gin.H{"Error": "Username `" + username + "` is already selected by someone else"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusForbidden, ret)
		return
	}
	if !gInstance.CheckIfValidColor(color) {
		ret = gin.H{"Error": "Color `" + color + "` is not a valid color"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusForbidden, ret)
		return
	}
	if gInstance.CheckIfColorSelected(color) {
		ret = gin.H{"Error": "Color `" + color + "` is already selected by someone else"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusForbidden, ret)
		return
	}

	gInstance.AvailableColors[color] = false

	p := models.Player{}
	p.UserName = username
	p.Color = color
	gInstance.AllPlayers = append(gInstance.AllPlayers, p)

	ret = gin.H{"Success": "You have joined the game mothafucka", "game_instance": gInstance,
		"user": gin.H{"username": username, "color": color}}

	log.Println(ret)
	c.JSON(http.StatusOK, ret)
}
