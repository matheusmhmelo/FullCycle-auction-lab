package auction

import (
	"context"
	"errors"
	"fullcycle-auction_go/internal/entity/auction_entity"
	mock_auction "fullcycle-auction_go/internal/infra/database/auction/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestAuctionRepository_CreateAuction(t *testing.T) {
	auctionId := uuid.New().String()
	auction := &auction_entity.Auction{
		Id:          auctionId,
		ProductName: "TestProduct",
		Category:    "TestCategory",
		Description: "TestDescription",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	tests := []struct {
		name           string
		interval       time.Duration
		assertMock     func(T *testing.T, update chan<- bson.M) *mock_auction.MockMongoCollection
		expectedUpdate bson.M
		expectError    bool
	}{
		{
			name:     "success",
			interval: time.Millisecond,
			assertMock: func(t *testing.T, updateGot chan<- bson.M) *mock_auction.MockMongoCollection {
				ctrl := gomock.NewController(t)
				m := mock_auction.NewMockMongoCollection(ctrl)

				m.EXPECT().
					InsertOne(gomock.Any(), gomock.Any()).
					Return(nil, nil)

				m.EXPECT().
					UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
						updateGot <- update.(bson.M)
						return nil, nil
					})

				return m
			},
			expectError:    false,
			expectedUpdate: bson.M{"$set": bson.M{"status": auction_entity.Completed}},
		},
		{
			name:     "error to create auction",
			interval: time.Millisecond,
			assertMock: func(t *testing.T, updateGot chan<- bson.M) *mock_auction.MockMongoCollection {
				ctrl := gomock.NewController(t)
				m := mock_auction.NewMockMongoCollection(ctrl)

				m.EXPECT().
					InsertOne(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error to create auction"))

				return m
			},
			expectError: true,
		},
		{
			name:     "error to close auction",
			interval: time.Millisecond,
			assertMock: func(t *testing.T, updateGot chan<- bson.M) *mock_auction.MockMongoCollection {
				ctrl := gomock.NewController(t)
				m := mock_auction.NewMockMongoCollection(ctrl)

				m.EXPECT().
					InsertOne(gomock.Any(), gomock.Any()).
					Return(nil, nil)

				m.EXPECT().
					UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
						updateGot <- update.(bson.M)
						return nil, errors.New("error to close auction")
					})

				return m
			},
			expectError:    false,
			expectedUpdate: bson.M{"$set": bson.M{"status": auction_entity.Completed}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateChan := make(chan bson.M)
			db := tt.assertMock(t, updateChan)
			r := AuctionRepository{
				Collection: db,
				interval:   tt.interval,
			}

			err := r.CreateAuction(context.Background(), auction)
			if tt.expectError {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)

				select {
				case update := <-updateChan:
					require.Equal(t, tt.expectedUpdate, update)
				case <-time.After(time.Second):
					t.Fail()
				}
			}
		})
	}
}
