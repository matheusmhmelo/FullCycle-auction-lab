package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCollection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
}

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection MongoCollection
	interval   time.Duration
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
		interval:   getAuctionInterval(),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	// go function to watch auction and close after duration.
	go ar.handleAuctionClose(ctx, auctionEntity)

	return nil
}

func (ar *AuctionRepository) handleAuctionClose(ctx context.Context, auc *auction_entity.Auction) {
	logger.Info(fmt.Sprintf("starting watch auction %s", auc.Id))

	select {
	case <-time.After(ar.interval):
		// set auction status as completed
		filter := bson.M{"_id": auc.Id}
		update := bson.M{"status": auction_entity.Completed}

		_, err := ar.Collection.UpdateOne(ctx, filter, bson.M{"$set": update})
		if err != nil {
			logger.Error("error updating auction", err)
		}
		logger.Info(fmt.Sprintf("auction %s closed successfully", auc.Id))
		return
	case <-ctx.Done():
		logger.Info(fmt.Sprintf("context canceled, auction %s cannot be closed", auc.Id))
		return
	}
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}

	return duration
}
