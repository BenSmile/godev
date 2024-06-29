package db

import (
	"context"
	"fmt"

	"github.com/bensmile/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollection = "users"
)

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper

	GetUserByID(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	UpdateUser(context.Context, string, *types.UpdateUserParams) (*types.User, error)
	// UpdateUser2(context.Context, string, bson.M) (*types.User, error)
	DeleteUser(context.Context, string) error
}

type MongoUserStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	// validate the correctness of the ID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	// validate the correctness of the ID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {

	_, err := s.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return nil, fmt.Errorf("email already used")
	}
	res, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, id string, params *types.UpdateUserParams) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{
		{Key: "_id", Value: oid},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "firstName", Value: params.FirstName},
			{Key: "lastName", Value: params.LastName},
			{Key: "email", Value: params.Email},
		}},
	}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *MongoUserStore) UpdateUser2(ctx context.Context, id string, values bson.M) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{
		{Key: "_id", Value: oid},
	}
	update := bson.D{
		{Key: "$set", Value: values},
	}

	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var users []*types.User

	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- Dropping users collection ---")
	return s.collection.Drop(ctx)
}

func NewMongoUserStore(client *mongo.Client, dbName string) *MongoUserStore {
	return &MongoUserStore{
		client:     client,
		collection: client.Database(dbName).Collection(userCollection),
	}
}
