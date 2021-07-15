package utils

import (
	"pft-balance/menu/domain"
	"pft-balance/menu/foodpb"
)

func MenuPbToData(menu *foodpb.Menu) *domain.Menu {
	return &domain.Menu{
		Name:    menu.GetName(),
		Protein: menu.GetProtein(),
		Fat:     menu.GetFat(),
		Carbs:   menu.GetCarbs(),
	}
}

func DataToMenuPb(data *domain.Menu) *foodpb.Menu {
	return &foodpb.Menu{
		Id:      data.ID,
		Name:    data.Name,
		Protein: data.Protein,
		Fat:     data.Fat,
		Carbs:   data.Carbs,
	}
}
