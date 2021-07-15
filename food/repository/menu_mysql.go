package repository
//
//import (
//	"context"
//	"fmt"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//	"gorm.io/gorm"
//	"io"
//	"log"
//	"pft-balance/food/domain"
//	"pft-balance/food/foodpb"
//	"pft-balance/food/utils"
//)
//
//type menuServerMySQL struct {
//	db *gorm.DB
//}
//
//func NewMenuServerMySQL(db *gorm.DB) domain.MenuRepository {
//	return &menuServerMySQL{
//		db: db,
//	}
//}
//
//func (sm *menuServerMySQL) CreateMenu(_ context.Context, req *foodpb.CreateMenuRequest) (*foodpb.CreateMenuResponse, error) {
//	fmt.Println("CreateMenu()")
//
//	data := utils.MenuPbToData(req.GetMenu())
//
//	res := sm.db.Create(data)
//	if res.Error != nil {
//		return nil, status.Errorf(
//			codes.Internal,
//			fmt.Sprintf("Internal error: %v", res.Error))
//	}
//
//	return &foodpb.CreateMenuResponse{}, nil
//}
//
//func (sm *menuServerMySQL) ReadMenu(_ context.Context, req *foodpb.ReadMenuRequest) (*foodpb.ReadMenuResponse, error) {
//	fmt.Println("ReadMenu()")
//	id := req.GetMenuId()
//
//	var data domain.Menu
//
//	res := sm.db.First(&data, "id = ?", id)
//	if res.Error != nil {
//		return nil, status.Errorf(
//			codes.NotFound,
//			fmt.Sprintf("cannnot find menu with specified ID: %v", res.Error),
//		)
//	}
//
//	return &foodpb.ReadMenuResponse{
//		Menu: utils.DataToMenuPb(&data),
//	}, nil
//}
//
//func (sm *menuServerMySQL) UpdateMenu(_ context.Context, req *foodpb.UpdateMenuRequest) (*foodpb.UpdateMenuResponse, error) {
//	fmt.Println("UpdateMenu()")
//	menu := req.GetMenu()
//
//	res := sm.db.First(&domain.Menu{}, "id = ?", menu.GetId())
//	if res.Error != nil {
//		return nil, status.Errorf(
//			codes.NotFound,
//			fmt.Sprintf("cannnot find menu with specified ID: %v", res.Error),
//		)
//	}
//
//	updatedData := utils.MenuPbToData(menu)
//
//	res = sm.db.Save(updatedData)
//	if res.Error != nil {
//		return nil, status.Errorf(
//			codes.NotFound,
//			fmt.Sprintf("unexpected error in UpdateMenu(): %v", res.Error),
//		)
//	}
//
//	var data domain.Menu
//
//	res = sm.db.First(&data, "id = ?", menu.GetId())
//	if res.Error != nil {
//		return nil, status.Errorf(
//			codes.NotFound,
//			fmt.Sprintf("unexpected error in UpdateMenu(): %v", res.Error),
//		)
//	}
//
//	return &foodpb.UpdateMenuResponse{
//		Menu: utils.DataToMenuPb(&data),
//	}, nil
//}
//
//func (sm *menuServerMySQL) DeleteMenu(_ context.Context, req *foodpb.DeleteMenuRequest) (*foodpb.DeleteMenuResponse, error) {
//	fmt.Println("DeleteMenu()")
//	id := req.GetMenuId()
//
//	res := sm.db.Delete(&domain.Menu{}, "id = ?", id)
//	if res.Error != nil {
//		return nil, status.Errorf(
//			codes.Internal,
//			fmt.Sprintf("unexpected error in DeleteMenu(): %v", res.Error),
//		)
//	}
//
//	return &foodpb.DeleteMenuResponse{}, nil
//}
//
//func (sm *menuServerMySQL) ListMenu(stream foodpb.MenuService_ListMenuServer) error {
//	fmt.Println("ListMenu()")
//
//	for {
//		req, err := stream.Recv()
//		if err == io.EOF {
//			return nil
//		}
//		if err != nil {
//			return status.Errorf(
//				codes.Internal,
//				fmt.Sprintf("unexpected error in ListMenu(): %v", err),
//			)
//		}
//
//		var data domain.Menu
//		res := sm.db.First(&data, "id = ?", req.GetMenuId())
//		if res.Error != nil {
//			return status.Errorf(
//				codes.Internal,
//				fmt.Sprintf("unexpected error in ListMenu(): %v\n", res.Error),
//			)
//		}
//
//		err = stream.Send(&foodpb.ListMenuResponse{Menu: utils.DataToMenuPb(&data)})
//		if err != nil {
//			return status.Errorf(
//				codes.Internal,
//				fmt.Sprintf("unexpected error in ListMenu(): %v\n", res.Error),
//			)
//		}
//	}
//}
//
//func (sm *menuServerMySQL) ListAllMenus(_ *foodpb.ListAllMenusRequest, stream foodpb.MenuService_ListAllMenusServer) error {
//	id := 0
//	for {
//		var data domain.Menu
//
//		res := sm.db.First(&data, "id = ?", id)
//		if res.Error != nil {
//			return nil
//		}
//		err := stream.Send(&foodpb.ListAllMenusResponse{Menu: utils.DataToMenuPb(&data)})
//		if err != nil {
//			log.Fatalln(err)
//		}
//		id++
//	}
//}
