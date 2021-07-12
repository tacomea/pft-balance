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
