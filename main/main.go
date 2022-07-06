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

	OrderApplicationHost = os.Getenv("ORDER_APPLICATION_HOST")
	OrderApplicationPort = os.Getenv("ORDER_APPLICATION_PORT")
	ProductApplicationHost = os.Getenv("PRODUCT_APPLICATION_HOST")
	ProductApplicationPort = os.Getenv("PRODUCT_APPLICATION_PORT")

	ApplicationHost = os.Getenv("APPLICATION_HOST")
	ApplicationPort = os.Getenv("APPLICATION_PORT")
)

// Nginx Proxy configuration... 

var (
	NGINX_PROXY_HOST = os.Getenv("NGINX_PROXY_HOST")
	NGINX_PROXY_PORT = os.Getenv("NGINX_PROXY_PORT")
)

func main() {
	// Creating Default HTTP Router for the application.. 

	router := gin.Default()

	// Cors Policy Goes There...

	AllowedOrigins := []string{
		fmt.Sprintf("http://%s:%s", OrderApplicationHost, OrderApplicationPort),
    	fmt.Sprintf("http://%s:%s", ProductApplicationHost, ProductApplicationPort), 
		fmt.Sprintf("http://%s:%s", ApplicationHost, ApplicationPort),
	} 
	AllowedMethods := []string{"GET"} 
	AllowedHeaders := []string{"*"}

	router.Use(cors.New(
		cors.Config{
			AllowOrigins: AllowedOrigins,
			AllowMethods: AllowedMethods,
			AllowHeaders: AllowedHeaders,
			AllowCredentials: true,
		},
	))

	// Allowed Proxies...
	router.SetTrustedProxies([]string{
	fmt.Sprintf("%s:%s", NGINX_PROXY_HOST, NGINX_PROXY_PORT)})


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