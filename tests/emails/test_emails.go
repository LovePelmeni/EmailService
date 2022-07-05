package test_emails

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EmailSenderSuite struct {
	suite.Suite
	*require.Assertions
	Controller                  *gomock.Controller
	EmailMessage                string
	EmailReceiver               string
	MockedEmailSenderController interface{}
}

func (this *EmailSenderSuite) SetupTest() {
	this.Controller = gomock.NewController(this.T())
	this.EmailMessage = "Hello, this is test Email Message."
	this.EmailReceiver = "some_email@gmail.com"
}
func TestRunEmailSenderSuite(t *testing.T) {
	suite.Run(t, new(EmailSenderSuite))
}

func (this *EmailSenderSuite) TestEmailSend(t *testing.T) {
	defer this.Controller.Finish()
	this.MockedEmailSenderController.EXPECT().SendEmail(
		gomock.Eq([]string{this.EmailReceiver, this.EmailMessage}),
	).Return(true, nil).Times(1)
}
