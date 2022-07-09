package mongo_test

import (
	"os"
	"testing"
	"time"

	mock_mongo "github.com/LovePelmeni/OnlineStore/EmailService/mocks/mongo"
	"github.com/LovePelmeni/OnlineStore/EmailService/mongo_controllers"
	"github.com/fossoreslp/go-uuid-v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoConsumerSuite struct {
	suite.Suite
	*require.Assertions
	Controller            *gomock.Controller
	MongoConnection       *mongo_controllers.MongoDatabase
	Document              *mongo_controllers.EmailDocument
	MockedMongoDBConsumer *mock_mongo.MockMongoDatabaseInterface // mocked Interface.
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
	this.Document = &mongo_controllers.EmailDocument{
		Uuid:          primitive.NewObjectID(),
		EmailReceiver: customerEmail,
		Message:       emailMessage,
		CreatedAt:     time.Now(),
	}
	this.MockedMongoDBConsumer = mock_mongo.NewMockMongoDatabaseInterface(this.Controller)
	this.MongoConnection = &mongo_controllers.MongoDatabase{
		Host: TestHost, Port: TestPort, DbName: TestDbName,
		User: TestUser, Password: TestPassword,
	}
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(MongoConsumerSuite))
}

func (this *MongoConsumerSuite) TearDownTest() {
	this.Controller.Finish()
}

func (this *MongoConsumerSuite) TestSaveDocument() {
	TestDocument := mongo_controllers.EmailDocument{
		Uuid:          primitive.NewObjectID(),
		EmailReceiver: "some_email@gmail.com",
		Message:       "Test Email Message",
		CreatedAt:     time.Now(),
	}

	InterfaceResponse := this.MockedMongoDBConsumer.EXPECT().SaveDocument(&TestDocument).Return(true, nil).Times(1)

	StructResponse, error := this.MongoConnection.SaveDocument(&TestDocument)
	assert.Equal(this.T(), InterfaceResponse, StructResponse, "Failed To Compare Responses.")
	assert.Equal(this.T(), error, nil, "Error Should equal to None")
}

func (this *MongoConsumerSuite) TestUpdateDocument() {

	TestDocument := mongo_controllers.EmailDocument{
		Uuid:          primitive.NewObjectID(),
		EmailReceiver: "some_email@gmail.com",
		Message:       "Test Email Message",
		CreatedAt:     time.Now(),
	}

	saved, error := this.MongoConnection.SaveDocument(&TestDocument)
	if saved && error != nil {
		assert.Errorf(this.T(), error, "Mongo Test Database Is Not Running...")
	}
	UpdatedDocumentData := map[string]string{}

	InterfaceResponse := this.MockedMongoDBConsumer.EXPECT().UpdateDocument(TestDocument.Uuid.String(),
		UpdatedDocumentData).Return(true, nil).Times(1)

	StructResponse, error := this.MongoConnection.UpdateDocument(
		TestDocument.Uuid.String(), UpdatedDocumentData)

	assert.Equal(this.T(), StructResponse, InterfaceResponse)
	assert.Equal(this.T(), error, nil)
}

func (this *MongoConsumerSuite) TestDeleteDocument() {

	TestDocument := mongo_controllers.EmailDocument{
		Uuid:          primitive.NewObjectID(),
		EmailReceiver: "some_email@gmail.com",
		Message:       "Test Email Message",
		CreatedAt:     time.Now(),
	}

	saved, error := this.MongoConnection.SaveDocument(&TestDocument)
	if error != nil && saved != true {
		assert.Errorf(
			this.T(), error, "FAILED TO SAVE DOCUMENT TO THE TEST DB. CHECK IF ITS RUNNING.")
	}

	this.MockedMongoDBConsumer.EXPECT().DeleteDocument(TestDocument.Uuid.String()).Return(true).Times(1)
	this.MongoConnection.DeleteDocument(
		TestDocument.Uuid.String())

	var receivedObj *mongo_controllers.EmailDocument
	ReceivedDeletedObject, error := this.MongoConnection.GetDocument(
	TestDocument.Uuid.String())


	ReceivedDeletedObject.Decode(receivedObj)
	assert.Equal(this.T(), receivedObj, nil, "Object Should be None, Because It Has been Deleted.")
}

func (this *MongoConsumerSuite) TestGetDocument() {
	newObj := mongo_controllers.EmailDocument{
		Uuid: primitive.NewObjectID(),
		Message: "",
		EmailReceiver: "",
	}
	var receivedObject *mongo_controllers.EmailDocument
	this.MockedMongoDBConsumer.EXPECT().GetDocument(
	this.Document.Uuid.String()).Return(&newObj, nil).Times(1)

	StructResponse, error := this.MongoConnection.GetDocument(newObj.Uuid.String())
	assert.Equal(this.T(), error, nil)
	DecodedObjError := StructResponse.Decode(&receivedObject)
	if DecodedObjError != nil {assert.Errorf(this.T(), error, "Failed to  Decode Response Object.")}
	assert.Equal(this.T(), receivedObject.Uuid, newObj.Uuid)
}


func (this *MongoConsumerSuite) TestGetQuerysetDocument() {
	newObj := mongo_controllers.EmailDocument{
		Uuid: primitive.NewObjectID(),
		Message: "Test Message",
		EmailReceiver: "testemail@gmail.com",
	}
	var receivedObject *mongo_controllers.EmailDocument


	this.MockedMongoDBConsumer.EXPECT().GetDocumentList().Return(&newObj, nil).Times(1)

	StructResponse, error := this.MongoConnection.GetDocumentList()
	assert.Equal(this.T(), error, nil)
	DecodedObjError := StructResponse.Decode(&receivedObject)
	if DecodedObjError != nil {assert.Error(this.T(), error, "Failed to  Decode Response Object.")}
	assert.Equal(this.T(), receivedObject.Uuid, newObj.Uuid)

}

