package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"pft-balance/store/domain"
	"pft-balance/store/foodpb"
	"strconv"
)

var tpl *template.Template

type serviceClient struct {
	fc foodpb.FoodServiceClient
	mc foodpb.MenuServiceClient
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*html"))
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// gRPC
	cc1, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not conntect: %v\n", err)
	}
	defer func(cc *grpc.ClientConn) {
		err := cc.Close()
		if err != nil {
			log.Println(err)
		}
	}(cc1)

	cc2, err := grpc.Dial("localhost:50050", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not conntect: %v\n", err)
	}
	defer func(cc *grpc.ClientConn) {
		err := cc.Close()
		if err != nil {
			log.Println(err)
		}
	}(cc2)

	fc := foodpb.NewFoodServiceClient(cc1)
	mc := foodpb.NewMenuServiceClient(cc2)

	c := serviceClient{fc: fc, mc: mc}

	// routing
	r := mux.NewRouter()
	r.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))

	// GET
	r.HandleFunc("/", c.indexHandler)
	r.HandleFunc("/show", c.showHandler)
	r.HandleFunc("/add", c.addHandler)
	r.HandleFunc("/edit", c.editHandler)

	// POST
	r.HandleFunc("/create", c.createHandler)
	r.HandleFunc("/update", c.updateHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	fmt.Println("Starting Store Server...")
	log.Fatalln(http.ListenAndServe(":"+port, r))
}

func (c *serviceClient) indexHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")

	err := tpl.ExecuteTemplate(w, "index.html", msg)
	if err != nil {
		log.Println(err)
	}
}

func (c *serviceClient) showHandler(w http.ResponseWriter, _ *http.Request) {
	stream, err := c.mc.ListAllMenus(context.Background(), &foodpb.ListAllMenusRequest{})
	if err != nil {
		log.Println(err)
	}

	var menus []domain.Menu
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}
		menu := domain.Menu{
			ID:      res.GetMenu().GetId(),
			Name:    res.GetMenu().GetName(),
			Protein: res.GetMenu().GetProtein(),
			Fat:     res.GetMenu().GetFat(),
			Carbs:   res.GetMenu().GetCarbs(),
		}
		menus = append(menus, menu)
	}

	err = tpl.ExecuteTemplate(w, "show.html", menus)
	if err != nil {
		log.Println(err)
	}
}

func (c *serviceClient) addHandler(w http.ResponseWriter, _ *http.Request) {
	err := tpl.ExecuteTemplate(w, "add.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func (c *serviceClient) editHandler(w http.ResponseWriter, _ *http.Request) {
	err := tpl.ExecuteTemplate(w, "edit.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func (c *serviceClient) createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")

	protein, fat, carbs := calcNutri(r, c.fc)
	_, err := c.mc.CreateMenu(context.Background(), &foodpb.CreateMenuRequest{Menu: &foodpb.Menu{
		Name:    name,
		Protein: protein,
		Fat:     fat,
		Carbs:   carbs,
	}})
	if err != nil {
		log.Println(err)
	}

	msg := url.QueryEscape("New Menu Created!")
	http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
}

func (c *serviceClient) updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	id := r.FormValue("id")

	protein, fat, carbs := calcNutri(r, c.fc)

	_, err := c.mc.UpdateMenu(context.Background(), &foodpb.UpdateMenuRequest{Menu: &foodpb.Menu{
		Id:      id,
		Name:    name,
		Protein: protein,
		Fat:     fat,
		Carbs:   carbs,
	}})
	if err != nil {
		log.Println(err)
	}

	msg := url.QueryEscape("New Menu Created!")
	http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
}

func calcNutri(r *http.Request, fc foodpb.FoodServiceClient) (float64, float64, float64) {
	amount1, err := strconv.ParseFloat(r.FormValue("amount1"), 32)
	if err != nil {
		log.Println(err)
	}
	amount2, err := strconv.ParseFloat(r.FormValue("amount1"), 32)
	if err != nil {
		log.Println(err)
	}
	amount3, err := strconv.ParseFloat(r.FormValue("amount1"), 32)
	if err != nil {
		log.Println(err)
	}

	res1, err := fc.ReadFood(context.Background(), &foodpb.ReadFoodRequest{FoodId: r.FormValue("id1")})
	if err != nil {
		log.Println(err)
	}
	res2, err := fc.ReadFood(context.Background(), &foodpb.ReadFoodRequest{FoodId: r.FormValue("id2")})
	if err != nil {
		log.Println(err)
	}
	res3, err := fc.ReadFood(context.Background(), &foodpb.ReadFoodRequest{FoodId: r.FormValue("id3")})
	if err != nil {
		log.Println(err)
	}

	protein := res1.GetFood().GetProtein() * amount1 / 100
	fat := res1.GetFood().GetFat() * amount1 / 100
	carbs := res1.GetFood().GetCarbs() * amount1 / 100

	protein += res2.GetFood().GetProtein() * amount2 / 100
	fat += res2.GetFood().GetFat() * amount2 / 100
	carbs += res2.GetFood().GetCarbs() * amount2 / 100

	protein += res3.GetFood().GetProtein() * amount3 / 100
	fat += res3.GetFood().GetFat() * amount3 / 100
	carbs += res3.GetFood().GetCarbs() * amount3 / 100

	return protein, fat, carbs
}
