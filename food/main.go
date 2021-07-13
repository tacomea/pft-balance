package main

import (
	"fmt"
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

func connect() *gorm.DB {
	connection, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	connection.AutoMigrate(&domain.Food{})

	return connection
}

//func initDb(db *gorm.DB) {
//	file, err := os.Open("csv/food_en.csv")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer file.Close()
//
//	reader := csv.NewReader(file)
//	var line []string
//
//	i := 0
//	for {
//		fmt.Println(i)
//		i++
//		line, err = reader.Read()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			log.Fatalln(err)
//		}
//
//		protein, err := strconv.ParseFloat(line[2], 64)
//		if err != nil {
//			log.Fatal(err)
//		}
//		fat, err := strconv.ParseFloat(line[3], 64)
//		if err != nil {
//			log.Fatal(err)
//		}
//		carbs, err := strconv.ParseFloat(line[4], 64)
//		if err != nil {
//			log.Fatal(err)
//		}
//		res := db.Create(domain.Food{
//			Name: line[1],
//			Protein: protein,
//			Fat: fat,
//			Carbs: carbs,
//		})
//		if res.Error != nil {
//			log.Fatal(res.Error)
//		}
//	}
//}

func main() {
	// MySQL
	db := connect()
	fmt.Println("initializing DB")
	//initDb(db)

	fmt.Println("Food Database Server Started!!")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	sm := repository.NewFoodServerMySQL(db)
	//fu := usecase.NewFoodUsecase(sm)
	mm := repository.NewMenuServerMySQL(db)

	foodServer := grpc.NewServer()
	menuServer := grpc.NewServer()

	foodpb.RegisterFoodServiceServer(foodServer, sm)
	foodpb.RegisterMenuServiceServer(menuServer, mm)

	// Register reflection service on gRPC server
	//reflection.Register(s)

	go func() {
		if err := foodServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve : %v", err)
		}
		if err := menuServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve : %v", err)
		}
	}()

	// wait for control C to stop
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// block until a signal is received
	<-ch
	fmt.Println("stopping the server")
	foodServer.Stop()
	fmt.Println("Closing the lister")
	lis.Close()
	fmt.Println("End of program")
}
