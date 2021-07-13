package usecase

import (
	"context"
	"pft-balance/food/domain"
	"pft-balance/food/foodpb"
)

type foodUsecase struct {
	foodRepo domain.FoodRepository
}

func NewFoodUsecase(fr domain.FoodRepository) domain.FoodUsecase {
	return &foodUsecase{
		foodRepo: fr,
	}
}

func (u *foodUsecase) CreateFood(ctx context.Context, req *foodpb.CreateFoodRequest) (*foodpb.CreateFoodResponse, error) {
	return u.foodRepo.CreateFood(ctx, req)
}

func (u *foodUsecase) ReadFood(ctx context.Context, req *foodpb.ReadFoodRequest) (*foodpb.ReadFoodResponse, error) {
	return u.foodRepo.ReadFood(ctx, req)
}

func (u *foodUsecase) UpdateFood(ctx context.Context, req *foodpb.UpdateFoodRequest) (*foodpb.UpdateFoodResponse, error) {
	return u.foodRepo.UpdateFood(ctx, req)
}

func (u *foodUsecase) DeleteFood(ctx context.Context, req *foodpb.DeleteFoodRequest) (*foodpb.DeleteFoodResponse, error) {
	return u.foodRepo.DeleteFood(ctx, req)
}

func (u *foodUsecase) ListFoods(stream foodpb.FoodService_ListFoodsServer) error {
	return u.foodRepo.ListFoods(stream)
}
