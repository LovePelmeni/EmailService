package mongo_test

import (
	"os"
	"testing"
	"time"

	"fmt"

	"github.com/LovePelmeni/EmailService/mongo_controllers/"
	"github.com/LovePelmeni/OnlineStore/EmailService/mongo_controllers"
	"github.com/fossoreslp/go-uuid-v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type MongoConsumerSuite struct {
	suite.Suite
	Controller            *gomock.Controller
	MongoConnection       *mongo_controllers.MongoDatabase
	Document              *mongo_controllers.Document
	MockedMongoDBConsumer interface{} // mocked Interface.
}

type MongoDocumentIndexKey uuid.UUID

func (this *MongoConsumerSuite) TestRunMongoConsumerSuite(t *testing.T) {
	suite.Run(t, new(MongoConsumerSuite))
}

func (this *MongoConsumerSuite) SetupTest() {

	TestHost := os.Getenv("TEST_MONGODB_HOST")
	TestPort := os.Getenv("TEST_MONGODB_PORT")
	TestUser := os.Getenv("TEST_MONGODB_USER")
	TestPassword := os.Getenv("TEST_MONGODB_PASSWORD")
	TestDbName := os.Getenv("TEST_MONGODB_DATABASE")

	customerEmail := "some_customer_email@gmail.com"
	emailMessage := "Some Email Message."

	this.Controller = gomock.NewController(this.T())
	generatedUuid, error := uuid.NewString()
	generatedUuid += fmt.Sprintf(
		"%s", time.Now().String())

	if error != nil && len(generatedUuid) == 0 {
		generatedUuid = fmt.Sprintf("%s", time.Now().String())
	}

	this.Document = *mongo_controllers.EmailDocument{
		Uuid:          generatedUuid,
		EmailReceiver: customerEmail,
		EmailMessage:  emailMessage,
		CreatedAt:     time.Now().Date,
	}
	this.MockedMongoDBConsumer = []interface{}{}
	this.MongoConnection = mongo_controllers.MongoDatabase{}
}

func (this *MongoConsumerSuite) TestSaveDocument(t *testing.T) {
	this.MockedMongoDBConsumer.EXPECT().SaveDocument(
		gomock.Eq(this.Document)).Return(true).Times(1)
}

func (this *MongoConsumerSuite) TestDeleteDocument(t *testing.T) {
	this.MockedMongoDBConsumer.EXPECT().DeleteDocument(
		gomock.Eq(this.Document.Uuid)).Return(true).Times(1)
}
