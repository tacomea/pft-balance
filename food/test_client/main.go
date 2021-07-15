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

	cc, err := grpc.Dial("data_server:50051", opts)
	if err != nil {
		log.Fatalf("could not conntect: %v\n", err)
	}
	defer cc.Close()

	c := foodpb.NewMenuServiceClient(cc)

	// create food
	createMenu(c)
	createMenu(c)
	createMenu(c)

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
		Name: "test",
		Protein: 10.10,
		Fat: 10.10,
		Carbs: 10.10,
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
		Name: "changed",
		Protein: 20.20,
		Fat: 20.20,
		Carbs: 20.20,
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
