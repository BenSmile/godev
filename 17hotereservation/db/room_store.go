package db

import (
	"context"
	"fmt"

	"github.com/bensmile/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
	GetRoomById(context.Context, string) (*types.Room, error)
	UpdateRoom(context.Context, string, bson.M) (*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection
	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(DBNAME).Collection("rooms"),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {

	resp, err := s.collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var rooms []*types.Room
	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = resp.InsertedID.(primitive.ObjectID)

	//update the hotel with this room id

	filter := bson.M{"_id": room.HotelID}
	updateRomm := bson.M{"$push": bson.M{"rooms": room.ID}}

	if err := s.HotelStore.UpdateHotel(ctx, filter, updateRomm); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *MongoRoomStore) GetRoomById(ctx context.Context, id string) (*types.Room, error) {

	roomOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid room id")
	}

	var room types.Room
	if err := s.collection.FindOne(ctx, bson.M{"_id": roomOID}).Decode(&room); err != nil {
		return nil, err
	}

	return &room, nil
}

func (s *MongoRoomStore) UpdateRoom(ctx context.Context, id string, values bson.M) (*types.Room, error) {
	// Convert the string ID to an ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %w", err)
	}

	// Build the filter and update documents
	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: values}}

	// Update the room in the collection
	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update room: %w", err)
	}

	// If no documents were updated, return an error
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("room with ID %s not found", id)
	}

	// Fetch the updated room
	var updatedRoom types.Room
	err = s.collection.FindOne(ctx, filter).Decode(&updatedRoom)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated room: %w", err)
	}

	return &updatedRoom, nil
}
