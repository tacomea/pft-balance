# pft-balance 

### what is this?
this is a pft-balance calculator.  
(Haven't you ever thought about keeping record of your pft balance? And also how it would be a pain?)  
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
    - Server for users to connect their tracking app and order (track pft-balance)
    - `localhost:8080`

### Usage
1.`docker-compose up` from pft-balance directory  
- docker-compose (also, MySQL) currently doesn't work. (I'm still looking up about how to connect different containers... )
- This uses MySQL

2.`make` from the inside of each directory (food, menu, store, user)  
- make sure to run food & menu first
- This uses MongoDB

To use MyFitnessPal API, you need to make `oauthOn` true and insert client ID & Client Secret.
