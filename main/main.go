package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LovePelmeni/OnlineStore/EmailService/emails"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
)

var (
	port = os.Getenv("APPLICATION_PORT")

	OrderApplicationHost   = os.Getenv("ORDER_APPLICATION_HOST")
	OrderApplicationPort   = os.Getenv("ORDER_APPLICATION_PORT")
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

var (
	grpcPort = os.Getenv("GRPC_SERVER_PORT")
	grpcHost = os.Getenv("GRPC_SERVER_HOST")
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
			AllowOrigins:     AllowedOrigins,
			AllowMethods:     AllowedMethods,
			AllowHeaders:     AllowedHeaders,
			AllowCredentials: true,
		},
	))

	go emails.CreategRPCServer()

	// Allowed Proxies...
	router.SetTrustedProxies([]string{
		fmt.Sprintf("%s:%s", NGINX_PROXY_HOST, NGINX_PROXY_PORT)})

	router.GET("/healthcheck/", func(context *gin.Context) {
		context.JSON(http.StatusOK, nil)
	})
	router.Run(fmt.Sprintf(":%s", port))
}
