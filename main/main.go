package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/LovePelmeni/OnlineStore/EmailService/emails"
)

var (
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
)

var (
	port = os.Getenv("APPLICATION_PORT")
)

func main() {
	// Creating Default HTTP Router for the application.. 

	router := gin.Default()

	// Cors Policy Goes There...

	AllowedOrigins := []string{} 
	AllowedMethods := []string{} 
	AllowedHeaders := []string{}

	
	router.Use(cors.New(
		cors.Config{
			AllowOrigins: AllowedOrigins,
			AllowMethods: AllowedMethods,
			AllowHeaders: AllowedHeaders,
			AllowCredentials: true,
		},
	))

	// HTTP EndPoints Goes There.

	router.GET("/healthcheck/", func(context *gin.Context) {
		context.JSON(http.StatusOK, nil)
	})
	error := router.Run(fmt.Sprintf(":%s", port))

	// Running gRPC Server. 
	go emails.CreategRPCServer()

	if errors.Is(error, http.ErrServerClosed) {
		ErrorLogger.Println("Failed to Start Server.")
	}
	router.Run(fmt.Sprintf(":%s", port))
}