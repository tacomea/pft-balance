package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"pft-balance/food/foodpb"
)

func main() {
	fmt.Println("Creating food request")

	opts := grpc.WithInsecure()

	//cc1, err := grpc.Dial("0.0.0.0:50051", opts)
	//if err != nil {
	//	log.Fatalf("could not conntect: %v\n", err)
	//}
	//defer cc1.Close()

	cc2, err := grpc.Dial("0.0.0.0:50050", opts)
	if err != nil {
		log.Fatalf("could not conntect: %v\n", err)
	}
	defer cc2.Close()

	//c1 := foodpb.NewFoodServiceClient(cc1)
	c2 := foodpb.NewMenuServiceClient(cc2)

	//res, err := c1.ReadFood(context.Background(), &foodpb.ReadFoodRequest{FoodId: "1"})
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//fmt.Println(res.GetFood())

	menu := foodpb.Menu{
		Name:    "test",
		Protein: 1.2,
		Fat:     1.2,
		Carbs:   1.2,
	}
	_, err = c2.CreateMenu(context.Background(), &foodpb.CreateMenuRequest{Menu: &menu})
	if err != nil {
		log.Println(err)
	}

	response, err := c2.ReadMenu(context.Background(), &foodpb.ReadMenuRequest{MenuId: "60eff9760b468d550bc191ad"})
	if err != nil {
		return
	}
	fmt.Println(response.GetMenu())

	res1, err := c2.ListAllMenus(context.Background(), &foodpb.ListAllMenusRequest{})
	if err != nil {
		log.Println(err)
		return
	}
	for {
		msg, err := res1.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}
		log.Println(msg.GetMenu())
	}

	// create food
	//createMenu(c2)
	//createMenu(c2)
	//createMenu(c2)

	// read food
	//readMenu(c, foodId)

	//// update food
	//updateMenu(c, foodId)
	//
	//// delete food
	//deleteMenu(c, foodId)
	//
	//// list foods
	//listMenu(c)

}

func createMenu(c foodpb.MenuServiceClient) {
	food := foodpb.Menu{
		Name:    "test",
		Protein: 10.10,
		Fat:     10.10,
		Carbs:   10.10,
	}
	createMenuRes, err := c.CreateMenu(context.Background(), &foodpb.CreateMenuRequest{Menu: &food})
	if err != nil {
		log.Fatalf("unexpected error : %v\n", err)
	}
	fmt.Println("food has been created: ", createMenuRes)
}

func readMenu(c foodpb.MenuServiceClient, foodId string) {
	_, err := c.ReadMenu(context.Background(), &foodpb.ReadMenuRequest{MenuId: "jeafed"})
	if err != nil {
		fmt.Printf("correct - error happened while reading: %v\n", err)
	}

	readMenuReq := &foodpb.ReadMenuRequest{MenuId: foodId}
	readMenuRes, err := c.ReadMenu(context.Background(), readMenuReq)
	if err != nil {
		fmt.Printf("wrong - error happened while reading: %v\n", err)
	}
	fmt.Printf("food has been read: %v\n", readMenuRes)
}

func updateMenu(c foodpb.MenuServiceClient, foodId string) {
	newMenu := foodpb.Menu{
		Name:    "changed",
		Protein: 20.20,
		Fat:     20.20,
		Carbs:   20.20,
	}
	updateRes, err := c.UpdateMenu(context.Background(), &foodpb.UpdateMenuRequest{Menu: &newMenu})
	if err != nil {
		fmt.Println("wrong - error while updating")
	}
	fmt.Printf("food has been updated: %v\n", updateRes)
}

func deleteMenu(c foodpb.MenuServiceClient, foodId string) {
	deleteRes, err := c.DeleteMenu(context.Background(), &foodpb.DeleteMenuRequest{MenuId: foodId})
	if err != nil {
		fmt.Printf("wrong - error happened while deleting: %v\n", err)
	}
	fmt.Printf("food has been updated: %v\n", deleteRes)
}

func listAllMenus(c foodpb.MenuServiceClient) {
	stream, err := c.ListAllMenus(context.Background(), &foodpb.ListAllMenusRequest{})
	if err != nil {
		log.Fatalf("error in ListAllMenus(): %v \n", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("something happened: %v \n", err)
		}
		fmt.Println("list of foods: ", res.GetMenu())
	}
}
