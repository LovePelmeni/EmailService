package mongo_controllers

import (
	"sync"
	"time"
	"fmt"
	"context"
	exceptions "github.com/LovePelmeni/OnlineStore/EmailService/mongo_controllers/exceptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"errors"
	"github.com/fossoreslp/go-uuid-v4"
	"os"
	"log"
	"reflect"
)


var (
	DebugLogger *log.Logger 
	ErrorLogger *log.Logger 
	InfoLogger *log.Logger 
	WarnLogger *log.Logger 
)

var (
	EmailCollectionName = os.Getenv("EMAIL_COLLECTION_NAME")
)
type MongoDatabase struct {
	mutex sync.RWMutex 
	User string 
	Password string 
	Host string 
	Port string 
	DbName string 
}

type EmailDocument struct {

	mutex sync.RWMutex 
	Uuid primitive.ObjectID `bson:"_id"` // Unique Idenitifier for Every Document..
	EmailReceiver string `bson:"receiverEmail"`
	Message string `bson:"Message"`
	CreatedAt time.Time  `bson:"CreatedAt`
}

func InitializeMongoDatabase(){ 

	MongoDatabase := MongoDatabase{}
	Connection, error := MongoDatabase.Connect()

	if error != nil {
	if errors.Is(error, mongo.ErrWrongClient) || 
	errors.Is(error, mongo.ErrClientDisconnected){
	panic("Failed to Connect To MongoDB. Check if its running!")}
	}

	CreatedError := Connection.Database(fmt.Sprintf(
	"%s", MongoDatabase.DbName)).CreateCollection(
		context.Background(),
		EmailCollectionName,
    )
	if CreatedError != nil {panic("Failed To Initialize MongoDB.")}
}

func (this *MongoDatabase) Connect() (*mongo.Client, error){

	RequestContext, error := context.WithTimeout(
	context.Background(), 10 * time.Second)
	if error != nil {panic("Timeout Context for Mongo Connection, Creation Failed..")}
	mongoURL := fmt.Sprintf("mongodb://%s:%s/%s:%s/%s", this.User, this.Password, 
	this.Host, this.Port, this.DbName)

	defer error()
	Client, error_ := mongo.Connect(RequestContext, options.Client().ApplyURI(mongoURL))
	if errors.Is(error_, mongo.ErrClientDisconnected) || errors.Is(error, mongo.ServerError){
	return nil, errors.New("Failed To Connect, Server Error")}

	defer func(){if error := Client.Disconnect(RequestContext); error != nil {
	panic("Failed To Disconnect.")}
	}() // defering disconnection..
	return Client, nil 
}

func GenerateMongoDocumentIndexUuid() (string, error){
	generatedUuid, error := uuid.NewString()
	if error != nil {return "", error}
	generatedUuid = generatedUuid + fmt.Sprintf(
	"%s", time.Now().String())
	return generatedUuid, nil 
}


// going to be a goroutine...

func (this *MongoDatabase) saveDocument(document *EmailDocument) (bool, error) {

	Session, Exception := this.Connect()
	if errors.Is(Exception, mongo.ErrWrongClient) ||
	errors.Is(Exception, mongo.ErrWrongClient){panic("Invalid Client Credentials.")}
	if Exception != nil {return false, exceptions.ConnectionFailed()}

	Collection := Session.Database(this.DbName).Collection(
	fmt.Sprintf("%s", EmailCollectionName))

	document.mutex.Lock()
	inserted, error := Collection.InsertOne(context.Background(), document)
	DebugLogger.Println("Inserted.")
	_ = inserted 
	
	defer document.mutex.Unlock() 

	if errors.Is(error, mongo.ErrNilValue) || 
	errors.Is(error, mongo.ErrNilDocument) {return false, errors.New("Document Is Empty.")}
	return true, nil
}

func (this *MongoDatabase) updateDocument(documentUuid string,
UpdatedData ...map[string]string) (bool, error) {

	if none := reflect.TypeOf(UpdatedData).NumField(); none == 0 {
		return false, exceptions.ConnectionFailed()
	}
	Connection, error := this.Connect()
	if errors.Is(error, mongo.ErrWrongClient,) || errors.Is(error, mongo.ErrClientDisconnected){
		return false, exceptions.InvalidMongoClientError()
	}

	Collection := Connection.Database(
	this.DbName).Collection(EmailCollectionName)

	this.mutex.Lock()
	updated, error := Collection.UpdateOne(context.Background(),
	map[string]string{"_id": documentUuid}, UpdatedData)
	_ = updated 


	defer this.mutex.Unlock()
	if error != nil {
		InfoLogger.Println(fmt.Sprintf("Failed To Update Object: %s", error))
	 	return false, exceptions.OperationFailed("Update", error)
	}
	return true, nil 
}

func (this *MongoDatabase) deleteDocument(DocumentUuid string) (bool, error) {
	connection, error := this.Connect()
	Collection := connection.Database(this.DbName).Collection(EmailCollectionName)

	if error != nil {return false, exceptions.OperationFailed("Delete", error)}
	document, error := Collection.DeleteOne(context.Background(),
	map[string]string{"_id": DocumentUuid})
	_ = document 


	if error != nil {return true, nil}else{return false, 
	exceptions.OperationFailed("Delete", error)}
}

var document EmailDocument

func (this *MongoDatabase) getDocument(DocumentUuid string) (*mongo.SingleResult, error){

	connection, error := this.Connect() 
	if error != nil {return nil, exceptions.ConnectionFailed()}

	Collection := connection.Database(this.DbName).Collection(EmailCollectionName)
	Document := Collection.FindOne(context.Background(),
	map[string]string{"_id": DocumentUuid})

	if errors.Is(Document.Err(), mongo.ErrNilDocument) ||
    errors.Is(Document.Err(), mongo.ErrNilValue){
	DebugLogger.Println("No such Email Documents with Uuid:",
    DocumentUuid)}
	return Document, nil 
}

func (this *MongoDatabase) getDocumentList() (*mongo.Cursor, error) {
	connection, error := this.Connect() 
	if error != nil {return nil, nil}

	Collection := connection.Database(this.DbName).Collection(EmailCollectionName)
	Documents, error := Collection.Find(context.Background(), nil)

	if errors.Is(Documents.Err(), mongo.ErrNilCursor) ||
    errors.Is(Documents.Err(), mongo.ErrNilValue){
	DebugLogger.Println("No Documents Were Found.")}
	return Documents, nil 
}