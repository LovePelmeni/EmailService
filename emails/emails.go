package emails

import (
	"errors"
	"fmt"

	"log"
	"net"
	"os"

	"strings"
	"context"

	"strconv"
	"time"

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


// Default Themes for Accepted and Rejected Images.
var (
	AcceptBackgroundImage = mail.File{Data: []byte(""), Name: "AcceptedOrderEmail.png"}
	RejectBackgroundImage = mail.File{Data: []byte(""), Name: "RejectedOrderEmail.png"}
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
	port := os.Getenv("GRPC_APPLICATION_PORT")
	listener, error := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if error != nil {
		panic("Server failed to create a listener..")
	}
	server := grpc.NewServer()
	grpcControllers.RegisterEmailSenderServer(server, &serverInterface)
	reflection.Register(server)

	defer listener.Close()
	if error := server.Serve(listener); error != nil {
		panic(fmt.Sprintf(
			"Failed to Start gRPC Server Error Occurred. %s", error))
	}
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

	go func(customerEmail string, Message string){
		// Saving Email to Mongo database Asynchronously...
		mongoDatabase := mongo_controllers.MongoDatabase{
			User: os.Getenv("MONGO_DATABASE_USER"),
			Password: os.Getenv("MONGO_DATABASE_PASSWORD"),
			Host: os.Getenv("MONGO_DATABASE_HOST"),
			Port: os.Getenv("MONGO_DATABASE_PORT"),
		}
		Document := mongo_controllers.EmailDocument{
			Uuid: primitive.NewObjectID(),
			Message: Message,
			EmailReceiver: customerEmail,
		}
		response, error := mongoDatabase.SaveDocument(&Document)
		if response && error == nil {DebugLogger.Println(fmt.Sprintf("Document has been saved"))}else{
			DebugLogger.Println(fmt.Sprintf("Failed TO Save Email Document, Exception %s", error))
		}
	}(sender.CustomerEmail, sender.Message)


	DebugLogger.Println(fmt.Sprintf("Email has been sended to %s",
		sender.CustomerEmail))
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
		sended, error := sender.sendAcceptEmail()
		if error != nil {
			sended = false
		}
		return &grpcControllers.EmailResponse{Delivered: sended}, nil

	case grpcControllers.OrderStatus_REJECTED:
		sended, error := sender.sendRejectEmail()
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

func (this *EmailBackgroundImage) validateFile() error {
	file := os.File{}
	for _, extension := range []string{"png", "jpeg", "jpg"} {
		if valid := strings.Split(file.Name(), ".")[1]; valid == extension {
			return nil
		}
	}
	return errors.New("Invalid File Type")
}

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
	port, error := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if error != nil {
		return nil, error
	}

	client := mail.NewSMTPClient()
	client.Encryption = mail.EncryptionSTARTTLS
	client.Username = os.Getenv("SMTP_EMAIL")
	client.Password = os.Getenv("SMTP_PASSWORD")
	client.Port = port
	client.Host = os.Getenv("SMTP_HOST")
	client.ConnectTimeout = 10 * time.Second
	client.SendTimeout = 10 * time.Second
	smtpClient, error := client.Connect()

	if error == nil && smtpClient != nil {
		return smtpClient, nil
	}
	return nil, errors.New("Failed to create SMTP Client.")
}

// Sends Email Notification using mail golang SDK
func (this *EmailSender) SendEmailNotification(BackgroundImage ...EmailBackgroundImage) (bool, error) {
	// some logic of sending email...
	var ValidatedImage = mail.File{} // default theme for the email is going to be Empty file.
	FileExtension := "file-extension" // need to parse the file extension..
	client, error := createSMTPClient()
	
	if error != nil {return false, nil}

	if notNone := len(BackgroundImage); notNone != 0 {
		if error := BackgroundImage[0].validateFile(); error == nil {
			ValidatedImage = mail.File{
				Data: BackgroundImage[0].file,
				Name: fmt.Sprintf("%s.%s", time.Now().String(), FileExtension),
			}
		} else {
			DebugLogger.Println("Invalid Background Image. Skipping...")
		}
	}

	EmailMessage := mail.NewMSG()
	EmailMessage.AddTo(this.CustomerEmail).SetSubject(this.Message)
	EmailMessage.SetBody(mail.TextHTML, EmailHTMLBody)
	EmailMessage.Attach(&ValidatedImage)
	sended_error := EmailMessage.Send(client)

	switch sended_error.(Error) {

	case nil:
		go func(customerEmail string, Message string){
			// Saving Email to Mongo database Asynchronously...
			mongoDatabase := mongo_controllers.MongoDatabase{
				User: os.Getenv("MONGO_DATABASE_USER"),
				Password: os.Getenv("MONGO_DATABASE_PASSWORD"),
				Host: os.Getenv("MONGO_DATABASE_HOST"),
				Port: os.Getenv("MONGO_DATABASE_PORT"),
			}
			Document := mongo_controllers.EmailDocument{
				Uuid: primitive.NewObjectID(),
				Message: Message,
				EmailReceiver: customerEmail,
			}
			response, error := mongoDatabase.SaveDocument(&Document)
			if response && error == nil {DebugLogger.Println(fmt.Sprintf("Document has been saved"))}else{
				DebugLogger.Println(fmt.Sprintf("Failed TO Save Email Document, Exception %s", error))
			}
		}(this.CustomerEmail, this.Message)
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
func (this *EmailSender) sendRejectEmail(BackgroundImage ...EmailBackgroundImage) (bool, error) {

	sended, error := this.SendEmailNotification()
	if sended != true || error != nil {
		return false, errors.New(
			"Failed To Send Reject Email Notification.")
	} else {
		return true, nil
	}
}

// Method Is used for sending Email Notification to the customer Email, that the order has been Accepted.
func (this *EmailSender) sendAcceptEmail(BackgroundImage ...EmailBackgroundImage) (bool, error) {

	sended, error := this.SendEmailNotification()
	if sended != true || error != nil {
		return false, errors.New(
			"Failed To Send Accept Email Notification.")
	} else {
		return true, nil
	}
}



