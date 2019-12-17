package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

var MongoSession *mgo.Session

func main() {
	var err error
	if MongoSession, err = mgo.Dial("localhost"); err != nil {
		panic(err)
	}
	defer MongoSession.Close()

	MongoSession.SetMode(mgo.Monotonic, true)
	// setup mongodb..
	collection := MongoSession.DB("test").C("users")
	// Login ID SHOULD be unique.. so create index here?
	// if err = collection.DropIndex("Login"); err != nil {
	// 	fmt.Println(err)
	// 	//panic("cannot drop index")
	// }
	if err = collection.EnsureIndex(mgo.Index{
		Key:        []string{"login"},
		Unique:     true,
		DropDups:   true,
		Background: false,
		Sparse:     false,
	}); err != nil {
		fmt.Println(err)
		panic("cannot ensure mongodb index")
	}
	// initialize mongodb session store
	sessionCollection := MongoSession.DB("test").C("sessions")
	if err = sessionCollection.EnsureIndex(mgo.Index{
		Key:         []string{"sessioncreatedate"},
		Unique:      false,
		ExpireAfter: 3600,
	}); err != nil {
		fmt.Println(err)
		panic("cannot ensure mongodb session index")
	}
	fmt.Println("Starting Server at...")

	r := gin.Default()
	//http.ListenAndServe(":8081", r)
	api := r.Group("/api")
	{
		auth := api.Group("/authentication")
		{

			auth.POST("/register", registerEndpoint)
			auth.POST("/login", loginEndpoint)
			auth.GET("/login", loginEndpoint)
		}
		api.GET("/configurations", confEndpoint)
	}
	//r.Get("/api")
	r.RunTLS(":8081", "./certs/localhost.pem", "./certs/localhost-key.pem")
	//http.ListenAndServe(":8081", r)
}
