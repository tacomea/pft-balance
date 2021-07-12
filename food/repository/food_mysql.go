package repository

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"io"
	"pft-balance/food/domain"
	"pft-balance/food/foodpb"
	"pft-balance/food/utils"
)

type foodServerMySQL struct {
	db *gorm.DB
}

func NewFoodServerMySQL(db *gorm.DB) domain.FoodRepository {
	return &foodServerMySQL{
		db: db,
	}
}

func (sm *foodServerMySQL) CreateFood(_ context.Context, req *foodpb.CreateFoodRequest) (*foodpb.CreateFoodResponse, error) {
	fmt.Println("CreateFood()")

	data := utils.FoodPbToData(req.GetFood())

	res := sm.db.Create(data)
	if res.Error != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", res.Error))
	}

	return &foodpb.CreateFoodResponse{}, nil
}

func (sm *foodServerMySQL) ReadFood(_ context.Context, req *foodpb.ReadFoodRequest) (*foodpb.ReadFoodResponse, error) {
	fmt.Println("ReadFood()")
	id := req.GetFoodId()

	var data domain.Food

	res := sm.db.First(&data, "id = ?", id)
	if res.Error != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannnot find food with specified ID: %v", res.Error),
		)
	}

	return &foodpb.ReadFoodResponse{
		Food: utils.DataToFoodPb(&data),
	}, nil
}

func (sm *foodServerMySQL) UpdateFood(_ context.Context, req *foodpb.UpdateFoodRequest) (*foodpb.UpdateFoodResponse, error) {
	fmt.Println("UpdateFood()")
	food := req.GetFood()

	res := sm.db.First(&domain.Food{}, "id = ?", food.GetId())
	if res.Error != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannnot find food with specified ID: %v", res.Error),
		)
	}

	updatedData := utils.FoodPbToData(food)

	res = sm.db.Save(updatedData)
	if res.Error != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("unexpected error in UpdateFood(): %v", res.Error),
		)
	}

	var data domain.Food

	res = sm.db.First(&data, "id = ?", food.GetId())
	if res.Error != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("unexpected error in UpdateFood(): %v", res.Error),
		)
	}

	return &foodpb.UpdateFoodResponse{
		Food: utils.DataToFoodPb(&data),
	}, nil
}

func (sm *foodServerMySQL) DeleteFood(_ context.Context, req *foodpb.DeleteFoodRequest) (*foodpb.DeleteFoodResponse, error) {
	fmt.Println("DeleteFood()")
	id := req.GetFoodId()

	res := sm.db.Delete(&domain.Food{}, "id = ?", id)
	if res.Error != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("unexpected error in DeleteFood(): %v", res.Error),
		)
	}

	return &foodpb.DeleteFoodResponse{}, nil
}

func (sm *foodServerMySQL) ListFoods(stream foodpb.FoodService_ListFoodsServer) error {
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

		var data domain.Food
		res := sm.db.First(&data, "id = ?", req.GetFoodId())
		if res.Error != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("unexpected error in ListFood(): %v\n", res.Error),
			)
		}

		err = stream.Send(&foodpb.ListFoodResponse{Food: utils.DataToFoodPb(&data)})
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("unexpected error in ListFood(): %v\n", res.Error),
			)
		}
	}
}
