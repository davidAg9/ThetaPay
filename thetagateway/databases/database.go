package databases

import (
	"context"

	"github.com/davidAg9/thetagateway/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ThetaDatabaseProtocol interface {
	CreateTransaction(ctx *context.Context, txn *models.ThetaTransaction) (bool, error)
	UpdateTransaction(txn *models.ThetaTransaction) *models.ThetaTransaction
}

type ThetaDatabase struct {
	*mongo.Database
}

func ConnnectDatabase(ctx context.Context, opts *options.ClientOptions) (*mongo.Client, error) {

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (thetadb *ThetaDatabase) CreateCollection(collectionName string) *mongo.Collection {
	return thetadb.Collection(collectionName)

}
