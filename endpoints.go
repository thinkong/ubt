package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type UserConfig struct {
	SomeData string
}

type UserAccount struct {
	Login    string
	Password string
	Config   UserConfig
}

type UserSession struct {
	Login             string
	SessionAuthID     bson.ObjectId `bson:"_id,omitempty"`
	SessionCreateDate time.Time
}

func encryptPasswd(p string) string {
	bSlice := []byte(p)
	var hash []byte
	var err error
	if hash, err = bcrypt.GenerateFromPassword(bSlice, bcrypt.DefaultCost); err != nil {
		log.Println(err)
	}
	return string(hash)
}

func registerEndpoint(c *gin.Context) {
	var user UserAccount
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad form or json",
		})
		return
	}
	user.Password = encryptPasswd(user.Password)
	collection := MongoSession.DB("test").C("users")
	if err := collection.Insert(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Already exists",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"register": "ok",
	})
}

func loginEndpoint(c *gin.Context) {
	var user UserAccount
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad form or json",
		})
		return
	}

	collection := MongoSession.DB("test").C("users")
	var foundUser UserAccount
	if err := collection.Find(bson.M{"login": user.Login}).One(&foundUser); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		log.Println(user.Login, err)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "wrong password",
		})
		return
	}
	var sessionData UserSession
	sessionCollection := MongoSession.DB("test").C("sessions")

	sessionData = UserSession{
		Login:             user.Login,
		SessionCreateDate: time.Now(),
		SessionAuthID:     bson.NewObjectId(),
	}
	sessionCollection.Insert(sessionData)

	c.JSON(http.StatusOK, gin.H{
		"authenticationToken": sessionData.SessionAuthID,
	})
}

func confEndpoint(c *gin.Context) {
	authToken := c.GetHeader("Authorization")
	log.Println(authToken)
	authCode := strings.Split(authToken, "Bearer")
	if len(authCode) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}
	authTokenFinal := strings.TrimSpace(authCode[1])
	// get login from db
	var sessionData UserSession
	if err := MongoSession.DB("test").C("sessions").Find(bson.M{"_id": bson.ObjectIdHex(authTokenFinal)}).One(&sessionData); err != nil {
		log.Println("no session found ::", authTokenFinal)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token",
		})
		return
	}
	// get config from database
	var userData UserAccount
	if err := MongoSession.DB("test").C("users").Find(bson.M{"login": sessionData.Login}).One(&userData); err != nil {
		log.Println("no user found")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token with user",
		})
	}
	c.JSON(http.StatusOK, userData.Config)
}
