package domain

import (
	"context"
	"pft-balance/food/foodpb"
)

type Food struct {
	ID      string `gorm:"primaryKey;index"`
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

//type Menu struct {
//	ID      string  `bson:"_id,omitempty"`
//	Name    string  `bson:"name"`
//	Protein float64 `bson:"protein"`
//	Fat     float64 `bson:"fat"`
//	Carbs   float64 `bson:"carbs"`
//}
//
//type MenuRepository interface {
//	CreateMenu(ctx context.Context, req *foodpb.CreateMenuRequest) (*foodpb.CreateMenuResponse, error)
//	ReadMenu(ctx context.Context, req *foodpb.ReadMenuRequest) (*foodpb.ReadMenuResponse, error)
//	UpdateMenu(ctx context.Context, req *foodpb.UpdateMenuRequest) (*foodpb.UpdateMenuResponse, error)
//	DeleteMenu(ctx context.Context, req *foodpb.DeleteMenuRequest) (*foodpb.DeleteMenuResponse, error)
//	ListMenu(stream foodpb.MenuService_ListMenuServer) error
//	ListAllMenus(req *foodpb.ListAllMenusRequest, stream foodpb.MenuService_ListAllMenusServer) error
//}
//
//type MenuUsecase interface {
//	CreateMenu(ctx context.Context, req *foodpb.CreateMenuRequest) (*foodpb.CreateMenuResponse, error)
//	ReadMenu(ctx context.Context, req *foodpb.ReadMenuRequest) (*foodpb.ReadMenuResponse, error)
//	UpdateMenu(ctx context.Context, req *foodpb.UpdateMenuRequest) (*foodpb.UpdateMenuResponse, error)
//	DeleteMenu(ctx context.Context, req *foodpb.DeleteMenuRequest) (*foodpb.DeleteMenuResponse, error)
//	ListMenu(stream foodpb.MenuService_ListMenuServer) error
//	ListAllMenus(req *foodpb.ListAllMenusRequest, stream foodpb.MenuService_ListAllMenusServer) error
//}
