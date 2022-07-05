package mongo_controllers

import (
	"sync"
	"time"
	"github.com/"
	"fmt"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"errors"
	"github.com/fossoreslp/go-uuid-v4"
)

type MongoDatabaseConnection struct {}

type MongoDatabase struct {
	User string 
	Password string 
	Host string 
	Port string 
	DbName string 
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
	return Client, nil 
}

type EmailDocument struct {
	sync.RWMutex 
	EmailReceiver string `bson:"receiverEmail"`
	Message string `bson:"Message"`
	CreatedAt time.Time  `bson:"CreatedAt`
}

func GenerateMongoDocumentIndexUuid() (string, error){
	generatedUuid, error := uuid.NewString()
	if error != nil {return "", error}
	generatedUuid = generatedUuid.String() + fmt.Sprintf(
	"%s", time.Now().String())
	return generatedUuid, nil 
}

func (this *MongoDatabase) saveDocument() (bool, error) {

}

func (this *MongoDatabase) updateDocument() (bool, error) {

}

func (this *MongoDatabase) deleteDocument() (bool, error) {
	
}