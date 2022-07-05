package emails 

import (
	"log"
	"fmt"
	"errors"
	"github.com/LovePelmeni/OnlineStore/EmailService/emails/exceptions"
	"strings"
	"encoding/json"
	"os"
	"github.com/LovePelmeni/OnlineStore/EmailService/emails/proto/grpcControllers"
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
)



// grpc Server Controllers 


type grpcEmailServer struct {} // represents generated Interface by Golang. 



// Sends Default Email Message to the client with Optional background Image Specified..

func (this *grpcEmailServer) SendEmail(RequestEmailParams ) {

}


// Sends Order Email Message to the Client with specific background depends on status specified..
func (this *grpcEmailServer) SendOrderEmail(RequestOrderEmailParams) {

}








// Emails API


type EmailBackgroundImage struct {
	file byte
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

func (this *EmailSender) sendEmail(backgroundImage ...EmailBackgroundImage) (bool, error) {}


func (this *EmailSender) sendAcceptEmail(customerEmail string, message string) (bool, error) {}


func (this *EmailSender) sendRejectEmail(customerEmail string, message string) (bool, error) {}



