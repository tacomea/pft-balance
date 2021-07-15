package utils

import (
	"pft-balance/food/domain"
	"pft-balance/food/foodpb"
)

func FoodPbToData(food *foodpb.Food) *domain.Food {
	return &domain.Food{
		Name:    food.GetName(),
		Protein: food.GetProtein(),
		Fat:     food.GetFat(),
		Carbs:   food.GetCarbs(),
	}
}

func DataToFoodPb(data *domain.Food) *foodpb.Food {
	return &foodpb.Food{
		Id:      data.ID,
		Name:    data.Name,
		Protein: data.Protein,
		Fat:     data.Fat,
		Carbs:   data.Carbs,
	}
}

//func MenuPbToData(menu *foodpb.Menu) *domain.Menu {
//	return &domain.Menu{
//		Name:    menu.GetName(),
//		Protein: menu.GetProtein(),
//		Fat:     menu.GetFat(),
//		Carbs:   menu.GetCarbs(),
//	}
//}
//
//func DataToMenuPb(data *domain.Menu) *foodpb.Menu {
//	return &foodpb.Menu{
//		Id:      data.ID,
//		Name:    data.Name,
//		Protein: data.Protein,
//		Fat:     data.Fat,
//		Carbs:   data.Carbs,
//	}
//}