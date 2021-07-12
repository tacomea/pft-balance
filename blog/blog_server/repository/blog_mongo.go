package repository

import (
	"blog_server/blogpb"
	"blog_server/domain"
	"blog_server/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2/bson"
)

func NewBlogItemServerMongo(c *mongo.Collection) domain.BlogItemRepository {
	return &blogItemServerMongo{
		collection: c,
	}
}

type blogItemServerMongo struct {
	collection *mongo.Collection
}

func (sm *blogItemServerMongo) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("CreateBlog()")
	blog := req.GetBlog()

	data := domain.BlogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	res, err := sm.collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err))
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannot convert to OID"))
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.AuthorId,
			Title:    blog.Title,
			Content:  blog.Content,
		},
	}, nil
}

func (sm *blogItemServerMongo) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("ReadBlog()")
	blogId := req.GetBlogId()

	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	// create an empty struct
	data := &domain.BlogItem{}

	filter := bson.M{"_id": oid}

	res := sm.collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannnot find blog with specified ID: %v", err),
		)
	}

	return &blogpb.ReadBlogResponse{
		Blog: utils.DataToBlogPb(data),
	}, nil
}

func (sm *blogItemServerMongo) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("UpdateBlog()")
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("cannot parse ID"),
		)
	}

	// create an empty struct
	data := &domain.BlogItem{}
	filter := bson.M{"_id": oid}

	res := sm.collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannot find blog with specified iD: %v", err),
		)
	}

	// update
	data.AuthorID = blog.GetAuthorId()
	data.Title = blog.GetTitle()
	data.Content = blog.GetContent()

	_, err = sm.collection.ReplaceOne(context.Background(), filter, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannot update object in mongo db: %v", err),
		)
	}

	return &blogpb.UpdateBlogResponse{
		Blog: utils.DataToBlogPb(data),
	}, nil
}

func (sm *blogItemServerMongo) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	fmt.Println("DeleteBlog()")
	oid, err := primitive.ObjectIDFromHex(req.GetBlogId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("cannot parse ID"),
		)
	}

	// create an empty struct
	filter := bson.M{"_id": oid}
	res, err := sm.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannnot delete blogin mongo db: %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannnot find blog with specified ID: %v", err),
		)
	}

	return &blogpb.DeleteBlogResponse{
		BlogId: req.GetBlogId(),
	}, nil

}

func (sm *blogItemServerMongo) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	fmt.Println("ListBlog()")
	cur, err := sm.collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("unknown internal error: %v", err),
		)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := &domain.BlogItem{}
		if err := cur.Decode(data); err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("error while decoding data from mongo db: %v", err),
			)
		}
		stream.Send(&blogpb.ListBlogResponse{Blog: utils.DataToBlogPb(data)})
	}
	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("unknown internal error: %v", err),
		)
	}

	return nil
}
