package emails 

import (
	"log"
	"fmt"
	"errors"
	"strings"
	"os"
	"github.com/LovePelmeni/OnlineStore/EmailService/emails/proto/grpcControllers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"

	"strconv"
	mail "github.com/xhit/go-simple-mail/v2"
	"time"
)

var (
	DebugLogger *log.Logger 
	ErrorLogger *log.Logger 
	WarnLogger *log.Logger 
	InfoLogger *log.Logger 
)

var (
	AcceptBackgroundImage = EmailBackgroundImage{file: nil}
	RejectBackgroundImage = EmailBackgroundImage{file: nil}


	DefaultEmailNotificationImageRoot = "/file.png"
)



// grpc Server Controllers 


type grpcEmailServer struct {} // represents generated Interface by Golang. 
var server grpcEmailServer 


func CreategRPCServer(){
	// Logic Of Setting up gRPC Server, and replacing interface.
	port := os.Getenv("GRPC_APPLICATION_PORT")
	listener, error := net.Listen("tcp", fmt.Sprintf(":%s", port)) 
	if error != nil {panic("Server failed to create a listener..")}
	server := grpc.NewServer()
	grpcControllers.RegisterEmailSenderServer(server, &server)
	reflection.Register(server) 
	if error := server.Serve(listener); error != nil {
		panic(fmt.Sprintf(
		"Failed to Start gRPC Server Error Occurred. %s", error))
	}
}


// Sends Default Email Message to the client with Optional background Image Specified..

func (this *grpcEmailServer) SendEmail(
RequestEmailParams *grpcControllers.DefaultEmailParams) (*grpcControllers.EmailResponse){

	sender := EmailSender{
		customerEmail: RequestEmailParams.CustomerEmail,
	     message: RequestEmailParams.EmailMessage,
	}

	sended, error := sender.sendEmailNotification()
	if error != nil {sended = false}
	DebugLogger.Println(fmt.Sprintf("Email has been sended to %s",
    sender.customerEmail))
	return &grpcControllers.EmailResponse{Delivered: sended}
}



// Sends Order Email Message to the Client with specific background depends on status specified..
func (this *grpcEmailServer) SendOrderEmail(
RequestOrderEmailParams *grpcControllers.OrderEmailParams) (*grpcControllers.EmailResponse){

	sender := EmailSender{
		customerEmail: RequestOrderEmailParams.CustomerEmail,
		message: RequestOrderEmailParams.Message,
	}
	switch RequestOrderEmailParams.Status {

		case grpcControllers.OrderStatus_ACCEPTED:
			sended, error := sender.sendAcceptEmail()
			if error != nil {sended = false}
			return &grpcControllers.EmailResponse{Delivered: sended}
		
		case grpcControllers.OrderStatus_REJECTED:
			sended, error := sender.sendRejectEmail()
			if error != nil {sended = false}
			return &grpcControllers.EmailResponse{Delivered: sended}
		
		default:
			return &grpcControllers.EmailResponse{Delivered: true}
	}
}


// Emails API


type EmailBackgroundImage struct {
	file []byte
}

func (this *EmailBackgroundImage) validateFile() error {
	file := os.File{}
	for _, extension := range []string{"png", "jpeg", "jpg"}{
		if valid := strings.Split(file.Name(), ".")[1];
	    valid == extension {return nil}}
	return errors.New("Invalid File Type")
}


type EmailSender struct {

	customerEmail string 
	message string 
}



var (
	AllStates = []string{"Delivered", "Canceled", "On-The-Way"}

	// email properties
	EmailHTMLBody = ``
)



// Creates Default SMTP Client...
func createSMTPClient() (*mail.SMTPClient, error){
	// creates SMTP Client for managing emails.
	port, error :=  strconv.Atoi(os.Getenv("SMTP_PORT"))
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
	return smtpClient, nil}
	return nil, errors.New("Failed to create SMTP Client.")
}


// Sends Email Notification using mail golang SDK 
func (this *EmailSender) sendEmailNotification(BackgroundImage ...EmailBackgroundImage) (bool, error){
	// some logic of sending email...
	var ValidatedImage = mail.File{}
	FileExtension := "file-extension" // need to parse the file extension..
	client, error := createSMTPClient()
	if error != nil{
	return false, nil}

	if notNone := len(BackgroundImage); notNone != 0 {
		if error := BackgroundImage[0].validateFile(); error == nil {
			ValidatedImage = mail.File{
				Data: BackgroundImage[0].file,
				Name: fmt.Sprintf("%s.%s", time.Now().String(), FileExtension),
			}

		}else{DebugLogger.Println("Invalid Background Image. Skipping...")}
	}

	EmailMessage := mail.NewMSG()
	EmailMessage.AddTo(this.customerEmail).SetSubject(this.message)
	EmailMessage.SetBody(mail.TextHTML, EmailHTMLBody)
	EmailMessage.Attach(&ValidatedImage)
	sended_error := EmailMessage.Send(client)

	if sended_error != nil {
	return false, sended_error}
	return true, nil
}


func (this *EmailSender) sendDefaultEmail(backgroundImage ...EmailBackgroundImage) (bool, error){
	sended, error := this.sendEmailNotification() 
	return sended, error
}
// Method Is used for sending Email Notification to the customer Email, that the order has been rejected.
// Prepares the message and calls `NotifyOrder` method that sends email.
func (this *EmailSender) sendRejectEmail(BackgroundImage ...EmailBackgroundImage) (bool, error){

	sended, error := this.sendEmailNotification()
	if sended != true || error != nil {return false, errors.New(
	"Failed To Send Reject Email Notification.")} else {return true, nil}
}


// Method Is used for sending Email Notification to the customer Email, that the order has been Accepted.
func (this *EmailSender) sendAcceptEmail(BackgroundImage ...EmailBackgroundImage) (bool, error){

	sended, error := this.sendEmailNotification()
	if sended != true || error != nil {return false, errors.New(
	"Failed To Send Accept Email Notification.")} else {return true, nil}
}









