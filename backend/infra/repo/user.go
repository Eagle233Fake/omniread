package repo

import (
	"context"
	"time"

	"github.com/Boyuan-IT-Club/go-kit/logs"
	"github.com/Eagle233Fake/omniread/backend/infra/model"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserRepoSet = wire.NewSet(NewUserRepo)

type UserRepo struct {
	coll *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	r := &UserRepo{
		coll: db.Collection("users"),
	}
	r.ensureIndexes()
	return r
}

func (r *UserRepo) ensureIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Username unique
	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		logs.Errorf("Failed to create username index: %v", err)
	}

	// Email unique (sparse because it's optional)
	_, err = r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetSparse(true),
	})
	if err != nil {
		logs.Errorf("Failed to create email index: %v", err)
	}

	// Phone unique (sparse)
	_, err = r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "phone", Value: 1}},
		Options: options.Index().SetUnique(true).SetSparse(true),
	})
	if err != nil {
		logs.Errorf("Failed to create phone index: %v", err)
	}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	if user.Status == "" {
		user.Status = "active"
	}
	_, err := r.coll.InsertOne(ctx, user)
	return err
}

func (r *UserRepo) FindOne(ctx context.Context, filter interface{}) (*model.User, error) {
	var user model.User
	err := r.coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.FindOne(ctx, bson.M{"username": username})
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.FindOne(ctx, bson.M{"email": email})
}

func (r *UserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return r.FindOne(ctx, bson.M{"phone": phone})
}

func (r *UserRepo) UpdateLastLogin(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"lastLoginAt": now,
		},
	})
	return err
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": user,
	}
	_, err := r.coll.UpdateOne(ctx, filter, update)
	return err
}

func (r *UserRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	return r.FindOne(ctx, bson.M{"_id": id})
}
