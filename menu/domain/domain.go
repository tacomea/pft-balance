package domain

import (
	"context"
	"pft-balance/menu/foodpb"
)

type Menu struct {
	ID      string  `bson:"_id,omitempty"`
	Name    string  `bson:"name"`
	Protein float64 `bson:"protein"`
	Fat     float64 `bson:"fat"`
	Carbs   float64 `bson:"carbs"`
}

type MenuRepository interface {
	CreateMenu(ctx context.Context, req *foodpb.CreateMenuRequest) (*foodpb.CreateMenuResponse, error)
	ReadMenu(ctx context.Context, req *foodpb.ReadMenuRequest) (*foodpb.ReadMenuResponse, error)
	UpdateMenu(ctx context.Context, req *foodpb.UpdateMenuRequest) (*foodpb.UpdateMenuResponse, error)
	DeleteMenu(ctx context.Context, req *foodpb.DeleteMenuRequest) (*foodpb.DeleteMenuResponse, error)
	ListMenu(stream foodpb.MenuService_ListMenuServer) error
	ListAllMenus(req *foodpb.ListAllMenusRequest, stream foodpb.MenuService_ListAllMenusServer) error
}
