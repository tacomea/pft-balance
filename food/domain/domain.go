package domain

import (
	"context"
	"pft-balance/food/foodpb"
)

type Food struct {
	ID      uint32 `gorm:"primaryKey;autoIncrement;index"`
	Name    string
	Protein float64
	Fat     float64
	Carbs   float64
}

type FoodRepository interface {
	CreateFood(ctx context.Context, req *foodpb.CreateFoodRequest) (*foodpb.CreateFoodResponse, error)
	ReadFood(ctx context.Context, req *foodpb.ReadFoodRequest) (*foodpb.ReadFoodResponse, error)
	UpdateFood(ctx context.Context, req *foodpb.UpdateFoodRequest) (*foodpb.UpdateFoodResponse, error)
	DeleteFood(ctx context.Context, req *foodpb.DeleteFoodRequest) (*foodpb.DeleteFoodResponse, error)
	ListFoods(stream foodpb.FoodService_ListFoodsServer) error
}

type FoodUsecase interface {
	CreateFood(ctx context.Context, req *foodpb.CreateFoodRequest) (*foodpb.CreateFoodResponse, error)
	ReadFood(ctx context.Context, req *foodpb.ReadFoodRequest) (*foodpb.ReadFoodResponse, error)
	UpdateFood(ctx context.Context, req *foodpb.UpdateFoodRequest) (*foodpb.UpdateFoodResponse, error)
	DeleteFood(ctx context.Context, req *foodpb.DeleteFoodRequest) (*foodpb.DeleteFoodResponse, error)
	ListFoods(stream foodpb.FoodService_ListFoodsServer) error
}
