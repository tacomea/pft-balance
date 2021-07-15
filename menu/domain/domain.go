package domain

import (
	"context"
	"pft-balance/menu/foodpb"
)

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
	ListMenu(stream foodpb.MenuService_ListMenuServer) error
	ListAllMenus(req *foodpb.ListAllMenusRequest, stream foodpb.MenuService_ListAllMenusServer) error
}

