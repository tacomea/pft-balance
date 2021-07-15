package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2/bson"
	"io"
	"pft-balance/food/domain"
	"pft-balance/food/foodpb"
	"pft-balance/food/utils"
)

func NewFoodServerMongo(c *mongo.Collection) domain.FoodRepository {
	return &foodServerMongo{
		collection: c,
	}
}

type foodServerMongo struct {
	collection *mongo.Collection
}

func (sm *foodServerMongo) CreateFood(_ context.Context, req *foodpb.CreateFoodRequest) (*foodpb.CreateFoodResponse, error) {
	fmt.Println("CreateFood()")
	data := utils.FoodPbToData(req.GetFood())

	_, err := sm.collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err))
	}

	//oid, ok := res.InsertedID.(primitive.ObjectID)
	//if !ok {
	//	return nil, status.Errorf(
	//		codes.Internal,
	//		fmt.Sprintf("cannot convert to OID"))
	//}

	return &foodpb.CreateFoodResponse{}, nil
}

func (sm *foodServerMongo) ReadFood(_ context.Context, req *foodpb.ReadFoodRequest) (*foodpb.ReadFoodResponse, error) {
	fmt.Println("ReadFood()")
	id := req.GetFoodId()

	//oid, err := primitive.ObjectIDFromHex(foodId)
	//if err != nil {
	//	return nil, status.Errorf(
	//		codes.InvalidArgument,
	//		fmt.Sprintf("Cannot parse ID"),
	//	)
	//}

	// create an empty struct
	data := &domain.Food{}

	filter := bson.M{"_id": id}

	res := sm.collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannnot find food with specified ID: %v", err),
		)
	}

	return &foodpb.ReadFoodResponse{
		Food: utils.DataToFoodPb(data),
	}, nil
}

func (sm *foodServerMongo) UpdateFood(_ context.Context, req *foodpb.UpdateFoodRequest) (*foodpb.UpdateFoodResponse, error) {
	fmt.Println("UpdateFood()")
	food := req.GetFood()
	//oid, err := primitive.ObjectIDFromHex(food.GetId())
	//if err != nil {
	//	return nil, status.Errorf(
	//		codes.InvalidArgument,
	//		fmt.Sprintf("cannot parse ID"),
	//	)
	//}

	// create an empty struct
	data := &domain.Food{}
	filter := bson.M{"_id": food.GetId()}

	res := sm.collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannot find food with specified iD: %v", err),
		)
	}

	data = utils.FoodPbToData(food)

	_, err := sm.collection.ReplaceOne(context.Background(), filter, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannot update object in mongo db: %v", err),
		)
	}

	return &foodpb.UpdateFoodResponse{
		Food: utils.DataToFoodPb(data),
	}, nil
}

func (sm *foodServerMongo) DeleteFood(_ context.Context, req *foodpb.DeleteFoodRequest) (*foodpb.DeleteFoodResponse, error) {
	fmt.Println("DeleteFood()")
	id := req.GetFoodId()
	//oid, err := primitive.ObjectIDFromHex(req.GetFoodId())
	//if err != nil {
	//	return nil, status.Errorf(
	//		codes.InvalidArgument,
	//		fmt.Sprintf("cannot parse ID"),
	//	)
	//}

	// create an empty struct
	filter := bson.M{"_id": id}
	res, err := sm.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannnot delete foodin mongo db: %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannnot find food with specified ID: %v", err),
		)
	}

	return &foodpb.DeleteFoodResponse{}, nil

}

func (sm *foodServerMongo) ListFoods(stream foodpb.FoodService_ListFoodsServer) error {
	fmt.Println("ListFood()")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("unexpected error in ListFood(): %v", err),
			)
		}

		data := &domain.Food{}

		filter := bson.M{"_id": req.GetFoodId()}

		res := sm.collection.FindOne(context.Background(), filter)
		if err := res.Decode(data); err != nil {
			return status.Errorf(
				codes.NotFound,
				fmt.Sprintf("cannnot find food with specified ID: %v", err),
			)
		}

		err = stream.Send(&foodpb.ListFoodResponse{Food: utils.DataToFoodPb(data)})
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("unexpected error in ListFood(): %v\n", err),
			)
		}
	}
}
