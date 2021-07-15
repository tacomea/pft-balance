package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"pft-balance/menu/domain"
	"pft-balance/menu/foodpb"
	"pft-balance/menu/utils"
)

func NewMenuServerMongo(c *mongo.Collection) domain.MenuRepository {
	return &menuServerMongo{
		collection: c,
	}
}

type menuServerMongo struct {
	collection *mongo.Collection
}

func (sm *menuServerMongo) CreateMenu(_ context.Context, req *foodpb.CreateMenuRequest) (*foodpb.CreateMenuResponse, error) {
	fmt.Println("CreateMenu()")
	data := utils.MenuPbToData(req.GetMenu())

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

	return &foodpb.CreateMenuResponse{}, nil
}

func (sm *menuServerMongo) ReadMenu(_ context.Context, req *foodpb.ReadMenuRequest) (*foodpb.ReadMenuResponse, error) {
	fmt.Println("ReadMenu()")
	id := req.GetMenuId()

	//oid, err := primitive.ObjectIDFromHex(foodId)
	//if err != nil {
	//	return nil, status.Errorf(
	//		codes.InvalidArgument,
	//		fmt.Sprintf("Cannot parse ID"),
	//	)
	//}

	// create an empty struct
	data := &domain.Menu{}

	filter := bson.M{"id": id}

	res := sm.collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannnot find food with specified ID: %v", err),
		)
	}

	return &foodpb.ReadMenuResponse{
		Menu: utils.DataToMenuPb(data),
	}, nil
}

func (sm *menuServerMongo) UpdateMenu(_ context.Context, req *foodpb.UpdateMenuRequest) (*foodpb.UpdateMenuResponse, error) {
	fmt.Println("UpdateMenu()")
	food := req.GetMenu()
	//oid, err := primitive.ObjectIDFromHex(food.GetId())
	//if err != nil {
	//	return nil, status.Errorf(
	//		codes.InvalidArgument,
	//		fmt.Sprintf("cannot parse ID"),
	//	)
	//}

	// create an empty struct
	data := &domain.Menu{}
	filter := bson.M{"id": food.GetId()}

	res := sm.collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannot find food with specified iD: %v", err),
		)
	}

	data = utils.MenuPbToData(food)

	_, err := sm.collection.ReplaceOne(context.Background(), filter, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannot update object in mongo db: %v", err),
		)
	}

	return &foodpb.UpdateMenuResponse{
		Menu: utils.DataToMenuPb(data),
	}, nil
}

func (sm *menuServerMongo) DeleteMenu(_ context.Context, req *foodpb.DeleteMenuRequest) (*foodpb.DeleteMenuResponse, error) {
	fmt.Println("DeleteMenu()")
	id := req.GetMenuId()
	//oid, err := primitive.ObjectIDFromHex(req.GetMenuId())
	//if err != nil {
	//	return nil, status.Errorf(
	//		codes.InvalidArgument,
	//		fmt.Sprintf("cannot parse ID"),
	//	)
	//}

	// create an empty struct
	filter := bson.M{"id": id}
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

	return &foodpb.DeleteMenuResponse{}, nil

}

func (sm *menuServerMongo) ListMenu(stream foodpb.MenuService_ListMenuServer) error {
	fmt.Println("ListMenu()")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("unexpected error in ListMenu(): %v", err),
			)
		}

		data := &domain.Menu{}

		filter := bson.M{"id": req.GetMenuId()}

		res := sm.collection.FindOne(context.Background(), filter)
		if err := res.Decode(data); err != nil {
			return status.Errorf(
				codes.NotFound,
				fmt.Sprintf("cannnot find food with specified ID: %v", err),
			)
		}

		err = stream.Send(&foodpb.ListMenuResponse{Menu: utils.DataToMenuPb(data)})
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("unexpected error in ListMenu(): %v\n", err),
			)
		}
	}
}

func (sm *menuServerMongo) ListAllMenus(_ *foodpb.ListAllMenusRequest, stream foodpb.MenuService_ListAllMenusServer) error {
	fmt.Println("ListAllMenus()")
	cur, err := sm.collection.Find(context.Background(), nil)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("unknown internal error: %v", err),
		)
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Println(err)
		}
	}(cur, context.Background())

	for cur.Next(context.Background()) {
		data := &domain.Menu{}
		if err := cur.Decode(data); err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("error while decoding data from mongo db: %v", err),
			)
		}
		err := stream.Send(&foodpb.ListAllMenusResponse{Menu: utils.DataToMenuPb(data)})
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("unknown internal error: %v", err),
			)
		}
	}
	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("unknown internal error: %v", err),
		)
	}

	return nil
}