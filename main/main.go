package main 

import (
	"net/http"
	"fmt"
	"os"
	"log"
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	ErrorLogger *log.Logger 
	DebugLogger *log.Logger 
	InfoLogger *log.Logger 
)

var (
	port = os.Getenv("APPLICATION_PORT")
)

func main(){
	router := gin.Default()
	router.GET("/healthcheck/", func(context *gin.Context){
	context.JSON(http.StatusOK, nil)})
	error := router.Run(fmt.Sprintf(":%s", port))
	if errors.Is(error, http.ErrServerClosed){
	ErrorLogger.Println("Failed to Start Server.")}
}


