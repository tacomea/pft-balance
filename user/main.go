package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"pft-balance/user/foodpb"
	"strconv"
	"time"
)

var myFitnessPalConfig = &oauth2.Config{
	ClientID:     "this is test",
	ClientSecret: "this is test",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://api.myfitnesspal.com/oauth2/auth",
		TokenURL: "https://api.myfitnesspal.com/oauth2/token",
	},
	RedirectURL: "http://localhost:8080/oauth/mfp/receive",
	Scopes:      []string{"diary"},
}

// key is state from oauth login, value is expiration time
var stateDb = map[string]time.Time{}

var userId string

var oauthClient *http.Client

var tpl *template.Template

type serviceClient struct {
	mc foodpb.MenuServiceClient
}

type Menu struct {
	ID      string
	Name    string
	Protein float64
	Fat     float64
	Carbs   float64
}

type Diary struct {
	Items []struct {
		Type                string `json:"type"`
		Date                string `json:"date"`
		DiaryMeal           string `json:"diary_meal"`
		NutritionalContents struct {
			Carbohydrates float64 `json:"carbohydrates"`
			Fat           float64 `json:"fat"`
			Protein       int     `json:"protein"`
		} `json:"nutritional_contents"`
	} `json:"items"`
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*html"))
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	//gRPC
	cc, err := grpc.Dial("localhost:50050", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not conntect: %v\n", err)
	}
	defer func(cc *grpc.ClientConn) {
		err := cc.Close()
		if err != nil {
			log.Println(err)
		}
	}(cc)

	mc := foodpb.NewMenuServiceClient(cc)

	c := serviceClient{mc: mc}

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/order", c.orderHandler)
	mux.HandleFunc("/oauth/mfp/add", c.oauthMFPAddHandler)
	mux.HandleFunc("/oauth/mfp/check", c.oauthMFPCheckHandler)
	mux.HandleFunc("/oauth/mfp/grant", oauthMFPGrantHandler)
	mux.HandleFunc("/oauth/mfp/receive", oauthMFPReceiveHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	html := `
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Document</title>
</head>
<body>
<p>%s</p>
<p>Connect My Fitness Pal??</p>
<form method="POST" action="/oauth/mfp/grant">
    <input type="submit" value="CONNECT">
</form>
<div>
	<a href="/order">注文する<a/>
</div>
</body>
</html>`

	_, err := fmt.Fprintf(w, html, msg)
	if err != nil {
		log.Println(err)
	}
}

func (c *serviceClient) orderHandler(w http.ResponseWriter, _ *http.Request) {
	stream, err := c.mc.ListAllMenus(context.Background(), &foodpb.ListAllMenusRequest{})
	if err != nil {
		log.Println(err)
	}

	var menus []Menu
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}
		menu := Menu{
			ID:      res.GetMenu().GetId(),
			Name:    res.GetMenu().GetName(),
			Protein: res.GetMenu().GetProtein(),
			Fat:     res.GetMenu().GetFat(),
			Carbs:   res.GetMenu().GetCarbs(),
		}
		menus = append(menus, menu)
	}

	err = tpl.ExecuteTemplate(w, "order.html", menus)
	if err != nil {
		log.Println(err)
	}
}

func (c *serviceClient) oauthMFPCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	err := tpl.ExecuteTemplate(w, "mfp_check.html", id)
	if err != nil {
		log.Println(err)
	}
}

func (c *serviceClient) oauthMFPAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userId := r.FormValue("user_id")
	menuId := r.FormValue("menu_id")

	res, err := c.mc.ReadMenu(context.Background(), &foodpb.ReadMenuRequest{MenuId: menuId})
	if err != nil {
		log.Println(err)
	}
	protein := res.GetMenu().GetProtein()
	protein := res.GetMenu().GetProtein()
	protein := res.GetMenu().GetProtein()

	var diary Diary
	diary.Items[0].Type = "diary_meal"
	diary.Items[0].Date = time.Now().Format("2006-1-2")
	// todo
	diary.Items[0].DiaryMeal = r.FormValue("diary_meal")
	diary.Items[0].NutritionalContents.Protein = protein
	diary.Items[0].NutritionalContents.Fat =
	diary.Items[0].NutritionalContents.Carbohydrates =

	req, err := http.NewRequest("POST", "https://api.myfitnesspal.com/diary", bytes.NewBuffer())

	res, err := oauthClient.Do(req)
	if err != nil {
		msg := url.QueryEscape(err.Error())
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		msg := url.QueryEscape(strconv.Itoa(res.StatusCode))
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}
	defer res.Body.Close()
	msg := url.QueryEscape("data added to your my fitness pal account!")
	http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
}

func oauthMFPGrantHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	state := uuid.NewString()
	stateDb[state] = time.Now().Add(time.Hour)

	http.Redirect(w, r, myFitnessPalConfig.AuthCodeURL(state), http.StatusSeeOther)
}

func oauthMFPReceiveHandler(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")
	if state == "" || code == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if time.Now().After(stateDb[state]) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	t, err := myFitnessPalConfig.Exchange(r.Context(), code)
	if err != nil {
		msg := url.QueryEscape(err.Error())
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}
	userId = t.Extra("user_id").(string)

	tokenSource := myFitnessPalConfig.TokenSource(r.Context(), t)
	oauthClient = oauth2.NewClient(r.Context(), tokenSource)

	msg := url.QueryEscape("from oauth2")
	http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
}
