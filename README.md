# pfc-balance-calculator

### what is this?
Haven't you ever thought about keeping record of your pfc balance? And also how it would be a pain in the ass?
This lets you do that easily when you use food delivery service by allowing stores to put the amount of the ingredients in the first place.  
Supported Tracker App: `None` (I put the code to use MyFitnessPal, but I need to submit a form to use their API... 😢)

- blog
    - Just my test file. Please ignore this.
    
- food
    - Server to CRUD Food Nutrition Database
    
- menu
    - Server to CRUD Menu Database
    
- store
    - Server for stores to add their menu
    - `localhost:80`
    
- user
    - Server for users to connect their tracking app and order (& track pfc)
    - `localhost:8080`

### How to run
#### Method1
1. run mongodb locally
2. for the first time, comment out initMongo(food/main.go line:89) to initialize Food Nutrition DB
3. `make`
4. access `localhost:80` to create a menu and then, access `localhost:8080` to order & track data (data will show up in the terminal running user/main,go)
  
#### Method2 

1. change db from Mongo to MySQL in food/main.go & menu/main.go
2. `docker-compose up` in pfc-balance directory 

In either way, to use MyFitnessPal API, you need to make `oauthOn` true in user/main.go line39 and insert client ID & Client Secret in user/main.go line 22/23.
