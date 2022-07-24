package test_emails

import (
	"testing"

	"github.com/LovePelmeni/OnlineStore/EmailService/emails"
	mock_emails "github.com/LovePelmeni/OnlineStore/EmailService/mocks/emails"
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

func (this *EmailSenderSuite) TearDownTest() {
	this.Controller.Finish()
}

func TestRunEmailSenderSuite(t *testing.T) {
	suite.Run(t, new(EmailSenderSuite))
}

func (this *EmailSenderSuite) TestEmailSend() {

	EmailMessage := "Test EmailMessage"
	EmailReceiver := "testemail@gmail.com"

	EmailSender := emails.EmailSender{
		CustomerEmail: EmailReceiver,
		Message:       EmailMessage,
	}

	mockedResponse := this.MockedEmailSenderController.EXPECT().SendEmail().Times(1).Return(true, nil)
	response, error := EmailSender.SendEmailNotification(emails.EmailBackgroundImage{})

	assert.Equal(this.T(), mockedResponse.String, response)
	if notError := assert.Equal(this.T(), error, error) &&
		assert.Equal(this.T(), error, nil); notError != true {
		assert.Errorf(this.T(), error,
			"Error Should Equals to None, got %s", error)
	}

	// Mocked Order Accept Assertions.
	mockedOrderResponse := this.MockedEmailSenderController.EXPECT().SendOrderEmail().Times(1).Return(true, nil)

	AcceptResponse, error := EmailSender.SendAcceptEmail()
	assert.Equal(this.T(), mockedOrderResponse, AcceptResponse)
	assert.Equal(this.T(), error, nil)

	// Mocked Order Reject Assertions.

	mockedOrderResponse2 := this.MockedEmailSenderController.EXPECT().SendOrderEmail().Times(1).Return(true, nil)
	RejectResponse, error := EmailSender.SendRejectEmail()
	assert.Equal(this.T(), RejectResponse, mockedOrderResponse2)
	assert.Equal(this.T(), error, nil)

	defer this.Controller.Finish()
}
