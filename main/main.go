package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LovePelmeni/OnlineStore/EmailService/emails"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/LovePelmeni/OnlineStore/EmailService/emails/proto/grpcControllers"
	"google.golang.org/grpc"
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

func CreateClient() (grpcControllers.EmailSenderClient, error) {
	connection, error := grpc.Dial(fmt.Sprintf("%s:%s", grpcHost, grpcPort), grpc.WithInsecure())
	if error != nil {
		panic(error)
	}
	client := grpcControllers.NewEmailSenderClient(connection)
	return client, nil
}

type RequestHTTPMessage struct {
	message string
}

var reqMessage RequestHTTPMessage

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

	// HTTP EndPoints Goes There.

	// GRPCLient via HTTP Endpoints...

	router.POST("send/default/email/", func(context *gin.Context) {
		customerEmail := context.Query("customerEmail")
		client, error := CreateClient()
		message := context.Query("Message")
		if error != nil {
			panic(error)
		}
		grpcRequestParams := grpcControllers.DefaultEmailParams{
			CustomerEmail: customerEmail,
			EmailMessage:  message,
		}
		response, error := client.SendEmail(context, &grpcRequestParams)
		_ = response
		context.JSON(http.StatusOK, gin.H{"Delivered": response.Delivered})
	})

	router.POST("send/order/accept/email/", func(context *gin.Context) {

		customerEmail := context.Query("customerEmail")
		client, error := CreateClient()
		if error != nil {
			panic("Failed To Create Client.")
		}
		message := context.Query("Message")

		orderEmailParams := grpcControllers.OrderEmailParams{

			Status:        grpcControllers.OrderStatus_ACCEPTED,
			CustomerEmail: customerEmail,
			Message:       message,
		}
		response, error := client.SendOrderEmail(
			context, &orderEmailParams)

		context.JSON(http.StatusOK,
			gin.H{"Delivered": response.Delivered})
	})

	router.POST("send/order/reject/email/:customerEmail", func(context *gin.Context) {

		customerEmail := context.Query("customerEmail")
		client, error := CreateClient()
		if error != nil {
			panic("Failed To Create Client.")
		}
		message := context.Query("Message")
		orderEmailParams := grpcControllers.OrderEmailParams{

			Status:        grpcControllers.OrderStatus_REJECTED,
			CustomerEmail: customerEmail,
			Message:       message,
		}
		response, error := client.SendOrderEmail(
			context, &orderEmailParams)
		context.JSON(http.StatusOK,
			gin.H{"Delivered": response.Delivered})
	})

	router.GET("/healthcheck/", func(context *gin.Context) {
		context.JSON(http.StatusOK, nil)
	})
	error := router.Run(fmt.Sprintf(":%s", port))

	// Running gRPC Server.

	if errors.Is(error, http.ErrServerClosed) {
		ErrorLogger.Println("Failed to Start Server.")
	}
	router.Run(fmt.Sprintf(":%s", port))
}
