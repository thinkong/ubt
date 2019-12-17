package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
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
