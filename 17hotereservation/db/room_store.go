package db

import (
	"context"

	"github.com/bensmile/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoRoomStore(client *mongo.Client, dbName string) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(dbName).Collection("rooms"),
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = resp.InsertedID.(primitive.ObjectID)

	//update the hotel with this room id

	return room, nil
}
