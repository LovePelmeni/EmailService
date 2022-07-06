package test_emails

import (
	"testing"
	"github.com/LovePelmeni/OnlineStore/EmailService/mocks/mock_emails"
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
	MockedEmailSenderController *mock_emails.NewMockEmailSenderInterface
}

func (this *EmailSenderSuite) SetupTest() {
	this.Controller = gomock.NewController(this.T())
	this.EmailMessage = "Hello, this is test Email Message."
	this.EmailReceiver = "some_email@gmail.com"
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
		Message:       EmailMessage}

	mocked_response, error := this.MockedEmailSenderController.EXPECT().SendEmail(
		gomock.Eq([]string{this.EmailReceiver, this.EmailMessage}),
	).Return(true, nil).Times(1).NoError(error)

	response, error := EmailSender.SendEmailNotification()

	assert.Equal(mocked_response, response)
	if notError := assert.Equal(error, error) &&
		assert.Equal(error, nil); notError != true {
		assert.Errorf(
			"Error Should Equals to None, got %s", error)
	}
}
