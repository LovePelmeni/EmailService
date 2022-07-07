package test_emails

import (
	"testing"
	mock_emails "github.com/LovePelmeni/OnlineStore/EmailService/mocks/emails"
	"github.com/LovePelmeni/OnlineStore/EmailService/emails"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EmailSenderSuite struct {
	suite.Suite
	*require.Assertions
	Controller                  *gomock.Controller
	EmailMessage                string
	EmailReceiver               string
	MockedEmailSenderController *mock_emails.MockEmailSenderInterface 
}

func (this *EmailSenderSuite) SetupTest() {
	this.Controller = gomock.NewController(this.T())
	this.EmailMessage = "Hello, this is test Email Message."
	this.EmailReceiver = "some_email@gmail.com"
	this.MockedEmailSenderController = mock_emails.NewMockEmailSenderInterface(
	this.Controller)
}

func TestRunEmailSenderSuite(t *testing.T) {
	suite.Run(t, new(EmailSenderSuite))
}

func (this *EmailSenderSuite) TearDownTest(t *testing.T) {
	this.Controller.Finish()
}

func (this *EmailSenderSuite) TestEmailSend(t *testing.T) {

	defer this.Controller.Finish()
	EmailMessage := "Test EmailMessage"
	EmailReceiver := "testemail@gmail.com"

	EmailSender := emails.EmailSender{
		CustomerEmail: EmailReceiver,
		Message:       EmailMessage,
	}

	mockedResponse := this.MockedEmailSenderController.EXPECT(
	).SendEmail().Times(1).Return(true, nil)

	response, error := EmailSender.SendEmailNotification()

	assert.Equal(t, mockedResponse.String, response)
	if notError := assert.Equal(t, error, error) &&
		assert.Equal(t, error, nil); notError != true {
		assert.Errorf(t, error,
		"Error Should Equals to None, got %s", error)
	}

	// Mocked Order Accept Assertions.
	mockedOrderResponse := this.MockedEmailSenderController.EXPECT(
	).SendOrderEmail().Times(1).Return(true, nil)

	AcceptResponse, error := EmailSender.SendAcceptEmail() 
	assert.Equal(t, mockedOrderResponse, AcceptResponse)
	assert.Equal(t, error, nil)
	
	// Mocked Order Reject Assertions.

	mockedOrderResponse2 := this.MockedEmailSenderController.EXPECT(
	).SendOrderEmail().Times(1).Return(true, nil)
	RejectResponse, error := EmailSender.SendRejectEmail()
	assert.Equal(t, RejectResponse, mockedOrderResponse2)
	assert.Equal(t, error, nil)
}



