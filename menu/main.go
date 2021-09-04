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
	"pft-balance/menu/domain"
	"pft-balance/menu/foodpb"
	"pft-balance/menu/repository"
)

//var (
//	schema         = "%s:%s@tcp(mysql:3306)/%s?charset=utf8&parseTime=True&loc=Local"
//	username       = os.Getenv("MYSQL_USER")
//	password       = os.Getenv("MYSQL_PASSWORD")
//	userDbName     = os.Getenv("MYSQL_DATABASE")
//	dataSourceName = fmt.Sprintf(schema, username, password, userDbName)
//)

var (
	schema			= "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
	dbHost			= os.Getenv("DB_HOST")
	dbPort			= os.Getenv("DB_PORT")
	username		= os.Getenv("MYSQL_USER")
	password		= os.Getenv("MYSQL_PASSWORD")
	userDbName		= os.Getenv("MYSQL_DATABASE")
	dataSourceName	= fmt.Sprintf(schema, username, dbHost, dbPort, password, userDbName)
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func connectMySQL() *gorm.DB {
	connection, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	connection.AutoMigrate(&domain.Menu{})

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
	fmt.Println("Menu RPC Starting...")

	lis, err := net.Listen("tcp", "0.0.0.0:50050")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// MySQL
	db := connectMySQL()
	mm := repository.NewMenuServerMySQL(db)

	// Mongo
	//client := connectMongo()
	//colMenu := client.Database("food_db").Collection("menu")
	//mm := repository.NewMenuServerMongo(colMenu)

	// Initializing DB
	//initDb(colFood)

	menuServer := grpc.NewServer()

	foodpb.RegisterMenuServiceServer(menuServer, mm)

	// Register reflection service on gRPC server
	//reflection.Register(s)

	go func() {
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
	//foodServer.Stop()
	menuServer.Stop()
	fmt.Println("Closing the lister")
	lis.Close()
	fmt.Println("closing the mongodb connection")
	//client.Disconnect(context.TODO())

	fmt.Println("End of program")
}
