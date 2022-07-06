package mongo_test

import (
	"os"
	"testing"
	"time"

	"fmt"

	"github.com/LovePelmeni/OnlineStore/EmailService/mocks/mock_mongo"
	"github.com/LovePelmeni/OnlineStore/EmailService/mongo_controllers"
	"github.com/fossoreslp/go-uuid-v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/stretchr/testify/assert"
)


type MongoConsumerSuite struct {
	suite.Suite
	*require.Assertions 
	Controller            *gomock.Controller
	MongoConnection       *mongo_controllers.MongoDatabase
	Document              *mongo_controllers.EmailDocument
	MockedMongoDBConsumer *mock_mongo.NewMockMongoDatabaseInterface // mocked Interface.
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

	this.Document = &mongo_controllers.EmailDocument{
		Uuid:          primitive.NewObjectID(),
		EmailReceiver: customerEmail,
		Message:  emailMessage,
		CreatedAt:     time.Now(),
	}
	this.MockedMongoDBConsumer = *mock_mongo.NewMockMongoDatabaseInterface{}
	this.MongoConnection = &mongo_controllers.MongoDatabase{
		Host: TestHost, Port: TestPort, DbName: TestDbName,
	    User: TestUser, Password: TestPassword,
	}
}

func (this *MongoConsumerSuite) TearDownTest(t *testing.T) {
	this.Controller.Finish()
}

func (this *MongoConsumerSuite) TestSaveDocument(t *testing.T) {
	TestDocument := mongo_controllers.EmailDocument{
		Uuid: primitive.NewObjectID(), 
		EmailReceiver: "some_email@gmail.com",
		Message: "Test Email Message",
		CreatedAt: time.Now(),
	}

	InterfaceResponse, error_ := this.MockedMongoDBConsumer.EXPECT().SaveDocument(
	gomock.Eq(this.Document)).Return(true, nil).Times(1)
	
	StructResponse, error := this.MongoConnection.SaveDocument(&TestDocument)
	
	assert.Equal(t, InterfaceResponse, StructResponse)
	assert.Equal(t, error, error_)
	assert.Equal(t, error_, nil)
}

func (this *MongoConsumerSuite) TestUpdateDocument(t *testing.T){

	TestDocument := mongo_controllers.EmailDocument{
		Uuid: primitive.NewObjectID(), 
		EmailReceiver: "some_email@gmail.com",
		Message: "Test Email Message",
		CreatedAt: time.Now(),
	}

	saved, error := this.MongoConnection.SaveDocument(&TestDocument)
	if saved && error != nil {assert.Errorf(t, error, "Mongo Test Database Is Not Running...")}
	UpdatedDocumentData := map[string]string{}

	InterfaceResponse, error_ := this.MockedMongoDBConsumer.EXPECT().UpdateDocument(
	gomock.Eq(UpdatedDocumentData)).Return(true, nil).Times(1)

	StructResponse, error := this.MongoConnection.UpdateDocument(
	TestDocument.Uuid.String(), UpdatedDocumentData) 

	assert.Equal(t, StructResponse, InterfaceResponse)
	assert.Equal(t, error, error_)
	assert.Equal(t, error_, nil)
}	

func (this *MongoConsumerSuite) TestDeleteDocument(t *testing.T) {
	this.MockedMongoDBConsumer.EXPECT().DeleteDocument(
		gomock.Eq(this.Document.Uuid)).Return(true).Times(1)
}
