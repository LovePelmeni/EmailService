package emails

import (
	"errors"
	"fmt"

	"log"
	"net"
	"os"

	"sync"
	"context"
	"net/http"
	"strings"

	"strconv"
	"time"
	"io/ioutil"

	"github.com/LovePelmeni/OnlineStore/EmailService/emails/proto/grpcControllers"
	"github.com/LovePelmeni/OnlineStore/EmailService/mongo_controllers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	mail "github.com/xhit/go-simple-mail/v2"
)

var (
	DebugLogger *log.Logger
	ErrorLogger *log.Logger
	WarnLogger  *log.Logger
	InfoLogger  *log.Logger
)

func init() {
	LogFile, error := os.OpenFile("emails.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if error != nil {
		fmt.Print("FAILED TO SET UP LOGGING IN EMAILS.GO")
	}
	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(LogFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Default Themes for Accepted and Rejected Images.
var (
	AcceptBackgroundImage = mail.File{Data: []byte(""), Name: "AcceptedOrderEmail.png"}
	RejectBackgroundImage = mail.File{Data: []byte(""), Name: "RejectedOrderEmail.png"}
)

// GRPC Server Credentials

var (
	grpcHost = os.Getenv("GRPC_SERVER_HOST")
	grpcPort = os.Getenv("GRPC_SERVER_PORT")
)

type Error error

// grpc Server Controllers

type grpcEmailServer struct {
	grpcControllers.UnimplementedEmailSenderServer
}

// represents generated Interface by Golang.
var serverInterface grpcEmailServer

func CreategRPCServer() {
	// Logic Of Setting up gRPC Server, and replacing interface.
	listener, error := net.Listen("tcp", fmt.Sprintf("%s:%s", grpcHost, grpcPort))
	if error != nil {
		panic("Server failed to create a listener..")
	}
	server := grpc.NewServer()
	grpcControllers.RegisterEmailSenderServer(server, &serverInterface)
	reflection.Register(server)

	if error := server.Serve(listener); error != nil {
		panic(fmt.Sprintf(
			"Failed to Start gRPC Server Error Occurred. %s", error))
	}
	defer listener.Close()
}

// Sends Default Email Message to the client with Optional background Image Specified..

func (this *grpcEmailServer) SendEmail(context context.Context,
	RequestEmailParams *grpcControllers.DefaultEmailParams) (*grpcControllers.EmailResponse, error) {

	sender := EmailSender{
		CustomerEmail: RequestEmailParams.CustomerEmail,
		Message:       RequestEmailParams.EmailMessage,
	}

	sended, error := sender.SendEmailNotification()
	if error != nil {
		sended = false
	}
	DebugLogger.Println(fmt.Sprintf("Email has been sended to ..."))
	return &grpcControllers.EmailResponse{Delivered: sended}, nil
}

// Sends Order Email Message to the Client with specific background depends on status specified..
func (this *grpcEmailServer) SendOrderEmail(context context.Context,
	RequestOrderEmailParams *grpcControllers.OrderEmailParams) (*grpcControllers.EmailResponse, error) {

	sender := EmailSender{
		CustomerEmail: RequestOrderEmailParams.CustomerEmail,
		Message:       RequestOrderEmailParams.Message,
	}
	switch RequestOrderEmailParams.Status {

	case grpcControllers.OrderStatus_ACCEPTED:
		sended, error := sender.SendAcceptEmail()
		if error != nil {
			sended = false
		}
		return &grpcControllers.EmailResponse{Delivered: sended}, nil

	case grpcControllers.OrderStatus_REJECTED:
		sended, error := sender.SendRejectEmail()
		if error != nil {
			sended = false
		}
		return &grpcControllers.EmailResponse{Delivered: sended}, nil

	default:
		return &grpcControllers.EmailResponse{Delivered: true}, nil

	}
}

// Emails API

type EmailBackgroundImage struct {
	file []byte 
}

func (this *EmailBackgroundImage) ToMailFile() (mail.File){
	mailFile := mail.File{
		Data: this.file, 
		Name: "File",
	}
	return mailFile 
}

//go:generate mockgen -destination=mocks/emails.go --build_flags=--mod=mod . EmailSenderInterface

type EmailSenderInterface interface {

	// Interface that represents Email Sender.
	// The Implementation should have attributes:
	// 1. CustomerEmail
	// 2. Message

	CustomerEmail() string
	Message() string
	SendEmail() (bool, error)
	SendOrderEmail() (bool, error)
}

type EmailSender struct {
	CustomerEmail string
	Message       string
}

var (
	AllStates = []string{"Delivered", "Canceled", "On-The-Way"}

	// email properties
	EmailHTMLBody = ``
)

// Creates Default SMTP Client...
func createSMTPClient() (*mail.SMTPClient, error) {
	// creates SMTP Client for managing emails.
	port, error := strconv.Atoi(os.Getenv("SMTP_SERVER_PORT"))
	if error != nil {
		return nil, error
	}

	client := mail.NewSMTPClient()
	client.Encryption = mail.EncryptionSTARTTLS
	client.Username = os.Getenv("SMTP_SERVER_EMAIL")
	client.Password = os.Getenv("SMTP_SERVER_PASSWORD")
	client.Port = port
	client.Host = os.Getenv("SMTP_SERVER_HOST")
	client.ConnectTimeout = 10 * time.Second
	client.SendTimeout = 10 * time.Second
	smtpClient, error := client.Connect()

	if error != nil {
		panic(error)
	}
	return smtpClient, nil
}


// Sends Email Notification using mail golang SDK
func (this *EmailSender) SendEmailNotification(
BackgroundImage ...EmailBackgroundImage) (bool, error) {

	
	client, error := createSMTPClient()

	if error != nil {
		panic(error)
	}

	if notNone := len(BackgroundImage); notNone != 0 {
		FileExtension := strings.Split(http.DetectContentType(BackgroundImage[0].file), "/")[1]
		ValidatedImage := mail.File{
			Data: BackgroundImage[0].file,
			Name: fmt.Sprintf("%s.%s", time.Now().String(), FileExtension),
		}
		_ = ValidatedImage 
		
	} else {
		DebugLogger.Println("Invalid Background Image. Skipping...")
	}

	DebugLogger.Println(fmt.Sprintf("Sended Notification to customer: %s",
    this.CustomerEmail))

	EmailMessage := mail.NewMSG()
	EmailMessage.AddTo(this.CustomerEmail).SetSubject(this.Message)
	EmailMessage.SetBody(mail.TextHTML, EmailHTMLBody)
	// EmailMessage.Attach(&ValidatedImage)
	sended_error := EmailMessage.Send(client)
	if sended_error != nil {
		panic(sended_error)
	}

	switch sended_error.(Error) {
	case nil:

		Group := sync.WaitGroup{}
		go func(group sync.WaitGroup, customerEmail string, Message string) {
			// Saving Email to Mongo database Asynchronously...

			mongoDatabase := mongo_controllers.MongoDatabase{

				User:     os.Getenv("MONGO_DATABASE_USER"),
				Password: os.Getenv("MONGO_DATABASE_PASSWORD"),
				Host:     os.Getenv("MONGO_DATABASE_HOST"),
				Port:     os.Getenv("MONGO_DATABASE_PORT"),
			}

			Document := mongo_controllers.EmailDocument{
				Uuid:          primitive.NewObjectID(),
				Message:       Message,
				EmailReceiver: customerEmail,
			}

			response, error := mongoDatabase.SaveDocument(&Document)


			if response && error == nil {
				DebugLogger.Println(fmt.Sprintf("Document has been saved"))
			} else {
				DebugLogger.Println(fmt.Sprintf("Failed TO Save Email Document, Exception %s", error))
			}

			Group.Done()
		}(Group, this.CustomerEmail, this.Message)

		Group.Wait()
		return true, nil

	default:
		return false, errors.New("Failed To Send Notification.")
	}
}

func (this *EmailSender) sendDefaultEmail(backgroundImage ...EmailBackgroundImage) (bool, error) {
	sended, error := this.SendEmailNotification()
	return sended, error
}

// Method Is used for sending Email Notification to the customer Email, that the order has been rejected.
// Prepares the message and calls `NotifyOrder` method that sends email.
func (this *EmailSender) SendRejectEmail() (bool, error) {

	FileByteData, ReadError := ioutil.ReadFile(os.Getenv("REJECT_ORDER_EMAIL_BACKGROUND_IMAGE_PATH")) // parsing reject email schema.
	if ReadError != nil {ErrorLogger.Println("Failed to Parse Reject Order Email File Path.")}
	BackgroundImage := EmailBackgroundImage{
		file: FileByteData,
	}
	sended, error := this.SendEmailNotification(BackgroundImage)
	if sended != true || error != nil {
		return false, errors.New(
			"Failed To Send Reject Email Notification.")
	} else {
		return true, nil
	}
}

// Method Is used for sending Email Notification to the customer Email, that the order has been Accepted.
func (this *EmailSender) SendAcceptEmail() (bool, error) {

	fileByteData, ReadError := ioutil.ReadFile(os.Getenv("ACCEPT_ORDER_EMAIL_BACKGROUND_IMAGE_PATH")) // parsing accept email schema
	if ReadError != nil {ErrorLogger.Println("Failed to Parse Accept Email File Path.")}
	backgroundImage := EmailBackgroundImage{
		file: fileByteData,
	}
	sended, error := this.SendEmailNotification(backgroundImage)
	if sended != true || error != nil {
		return false, errors.New(
			"Failed To Send Accept Email Notification.")
	} else {
		return true, nil
	}
}
