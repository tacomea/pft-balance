package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"os/signal"
	"pft-balance/food/domain"
	"pft-balance/food/foodpb"
	"pft-balance/food/repository"
)

var (
	schema         = "%s:%s@tcp(mysql:3306)/%s?charset=utf8&parseTime=True&loc=Local"
	username       = os.Getenv("MYSQL_USER")
	password       = os.Getenv("MYSQL_PASSWORD")
	userDbName     = os.Getenv("MYSQL_DATABASE")
	dataSourceName = fmt.Sprintf(schema, username, password, userDbName)
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func connectMySQL() *gorm.DB {
	connection, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	connection.AutoMigrate(&domain.Food{})

	return connection
}

func connectMongo() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func main() {

	fmt.Println("Food Database Server Starting...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// MySQL
	//db := connectMySQL()
	//sm := repository.NewFoodServerMySQL(db)
	//mm := repository.NewMenuServerMySQL(db)

	// Mongo
	client := connectMongo()

	// Food
	colFood := client.Database("food_db").Collection("food")
	sm := repository.NewFoodServerMongo(colFood)
	foodServer := grpc.NewServer()
	foodpb.RegisterFoodServiceServer(foodServer, sm)

	// Menu
	//colMenu := client.Database("food_db").Collection("menu")
	//mm := repository.NewMenuServerMongo(colMenu)
	//menuServer := grpc.NewServer()
	//foodpb.RegisterMenuServiceServer(menuServer, mm)

	// Initializing DB
	//initDb(colFood)

	// Register reflection service on gRPC server
	//reflection.Register(s)

	go func() {
		if err := foodServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve : %v", err)
		}

	}()

	// wait for control C to stop
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// block until a signal is received
	<-ch
	fmt.Println("stopping the server")
	//foodServer.Stop()
	foodServer.Stop()
	fmt.Println("Closing the lister")
	lis.Close()
	fmt.Println("closing the mongodb connection")
	client.Disconnect(context.TODO())
	fmt.Println("End of program")
}

//func initDb(collection *mongo.Collection) {
//	file, err := os.Open("csv/food_en.csv")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer file.Close()
//
//	reader := csv.NewReader(file)
//	var line []string
//
//	for {
//		line, err = reader.Read()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			log.Fatalln(err)
//		}
//
//		id, err := strconv.ParseUint(line[0], 10, 32)
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Println(id)
//		protein, err := strconv.ParseFloat(line[2], 64)
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Println(line[3])
//		fat, err := strconv.ParseFloat(line[3], 64)
//		if err != nil {
//			log.Fatal(err)
//		}
//		carbs, err := strconv.ParseFloat(line[4], 64)
//		if err != nil {
//			log.Fatal(err)
//		}
//		// My SQL
//		//res := db.Create(domain.Food{
//		//	Name: line[1],
//		//	Protein: protein,
//		//	Fat: fat,
//		//	Carbs: carbs,
//		//})
//		//if res.Error != nil {
//		//	log.Fatal(res.Error)
//		//}
//
//		// Mongo
//		_, err = collection.InsertOne(context.Background(), domain.Food{
//			Name:    line[1],
//			ID: uint32(id),
//			Protein: protein,
//			Fat:     fat,
//			Carbs:   carbs,
//		})
//		if err != nil {
//			log.Fatal(err)
//		}
//
//
//	}
//}
