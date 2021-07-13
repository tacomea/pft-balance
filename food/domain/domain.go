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

type Menu struct {
	ID      uint32 `gorm:"primaryKey;autoIncrement;index"`
	Name    string
	Protein float64
	Fat     float64
	Carbs   float64
}

type MenuRepository interface {
	CreateMenu(ctx context.Context, req *foodpb.CreateMenuRequest) (*foodpb.CreateMenuResponse, error)
	ReadMenu(ctx context.Context, req *foodpb.ReadMenuRequest) (*foodpb.ReadMenuResponse, error)
	UpdateMenu(ctx context.Context, req *foodpb.UpdateMenuRequest) (*foodpb.UpdateMenuResponse, error)
	DeleteMenu(ctx context.Context, req *foodpb.DeleteMenuRequest) (*foodpb.DeleteMenuResponse, error)
	ListMenus(stream foodpb.MenuService_ListMenusServer) error
}

type MenuUsecase interface {
	CreateMenu(ctx context.Context, req *foodpb.CreateMenuRequest) (*foodpb.CreateMenuResponse, error)
	ReadMenu(ctx context.Context, req *foodpb.ReadMenuRequest) (*foodpb.ReadMenuResponse, error)
	UpdateMenu(ctx context.Context, req *foodpb.UpdateMenuRequest) (*foodpb.UpdateMenuResponse, error)
	DeleteMenu(ctx context.Context, req *foodpb.DeleteMenuRequest) (*foodpb.DeleteMenuResponse, error)
	ListMenus(stream foodpb.MenuService_ListMenusServer) error
}
