package repo

import (
	"context"
	"time"

	"github.com/Eagle233Fake/omniread/backend/infra/model"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ReadingRepoSet = wire.NewSet(NewBookRepo, NewReadingProgressRepo, NewReadingSessionRepo)

// BookRepo
type BookRepo struct {
	coll *mongo.Collection
}

func NewBookRepo(db *mongo.Database) *BookRepo {
	return &BookRepo{
		coll: db.Collection("books"),
	}
}

func (r *BookRepo) Create(ctx context.Context, book *model.Book) error {
	if book.ID.IsZero() {
		book.ID = primitive.NewObjectID()
	}
	now := time.Now()
	book.CreatedAt = now
	book.UpdatedAt = now
	_, err := r.coll.InsertOne(ctx, book)
	return err
}

func (r *BookRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Book, error) {
	var book model.Book
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&book)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepo) List(ctx context.Context, limit int64, offset int64) ([]*model.Book, error) {
	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"created_at": -1})
	cursor, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	var books []*model.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}
	return books, nil
}

// ReadingProgressRepo
type ReadingProgressRepo struct {
	coll *mongo.Collection
}

func NewReadingProgressRepo(db *mongo.Database) *ReadingProgressRepo {
	return &ReadingProgressRepo{
		coll: db.Collection("reading_progress"),
	}
}

func (r *ReadingProgressRepo) Save(ctx context.Context, progress *model.ReadingProgress) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{
		"user_id": progress.UserID,
		"book_id": progress.BookID,
	}
	update := bson.M{
		"$set": bson.M{
			"progress":    progress.Progress,
			"current_loc": progress.CurrentLoc,
			"status":      progress.Status,
			"updated_at":  time.Now(),
		},
	}
	_, err := r.coll.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *ReadingProgressRepo) Find(ctx context.Context, userID, bookID primitive.ObjectID) (*model.ReadingProgress, error) {
	var p model.ReadingProgress
	err := r.coll.FindOne(ctx, bson.M{"user_id": userID, "book_id": bookID}).Decode(&p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *ReadingProgressRepo) ListByUser(ctx context.Context, userID primitive.ObjectID) ([]*model.ReadingProgress, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	var list []*model.ReadingProgress
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// ReadingSessionRepo
type ReadingSessionRepo struct {
	coll *mongo.Collection
}

func NewReadingSessionRepo(db *mongo.Database) *ReadingSessionRepo {
	return &ReadingSessionRepo{
		coll: db.Collection("reading_sessions"),
	}
}

func (r *ReadingSessionRepo) Create(ctx context.Context, session *model.ReadingSession) error {
	if session.ID.IsZero() {
		session.ID = primitive.NewObjectID()
	}
	_, err := r.coll.InsertOne(ctx, session)
	return err
}

func (r *ReadingSessionRepo) FindByUserAndDateRange(ctx context.Context, userID primitive.ObjectID, start, end time.Time) ([]*model.ReadingSession, error) {
	filter := bson.M{
		"user_id": userID,
		"start_time": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var list []*model.ReadingSession
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}
