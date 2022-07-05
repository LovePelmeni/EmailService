package mongo_test

import (
	"os"
	"testing"
	"time"

	"github.com/LovePelmeni/EmailService/mongo_controllers/"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/fossoreslp/go-uuid-v4"
	"fmt"
)

type MongoConsumerSuite struct {
	suite.Suite
	Controller            *gomock.Controller
	MongoConnection       interface{}
	Document              *mongo_controllers.Document
	MockedMongoDBConsumer interface{}
}

type MongoDocumentIndexKey uuid.UUID 


func (this *MongoConsumerSuite) SetupTest() {

	Host := os.Getenv("TEST_MONGODB_HOST")
	Port := os.Getenv("TEST_MONGODB_PORT")
	User := os.Getenv("TEST_MONGODB_USER")
	Password := os.Getenv("TEST_MONGODB_PASSWORD")
	DbName := os.Getenv("TEST_MONGODB_DATABASE")
	
	customerEmail := "some_customer_email@gmail.com"
	emailMessage := "Some Email Message."

	this.Controller = gomock.NewController(this.T())
	generatedUuid, error := uuid.NewString()
 	generatedUuid += fmt.Sprintf(
	"%s", time.Now().String())

	if error != nil && len(generatedUuid) == 0 {
	generatedUuid = fmt.Sprintf("%s", time.Now().String())}


	this.Document = *mongo_controllers.EmailDocument{
		Uuid: generatedUuid,
		EmailReceiver: customerEmail,
		EmailMessage: emailMessage,
		CreatedAt: time.Now().Date,
	}
	this.MockedMongoDBConsumer = []interface{}{}
	this.MongoConnection = []interface{}{}
}

func (this *MongoConsumerSuite) TestSaveDocument(t *testing.T) {
	this.MockedMongoDBConsumer.EXPECT().SaveDocument(
	gomock.Eq(this.Document)).Return(true).Times(1)
}

func (this *MongoConsumerSuite) TestDeleteDocument(t *testing.T) {
	this.MockedMongoDBConsumer.EXPECT().DeleteDocument(
	gomock.Eq(this.Document.Uuid)).Return(true).Times(1)
}
