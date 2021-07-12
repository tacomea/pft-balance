#mysql -uroot -ppassword --local-infile food_data -e "LOAD DATA LOCAL INFILE '/docker-entrypoint-initdb.d/food_ja.csv' INTO TABLE foods FIELDS TERMINATED BY ','LINES TERMINATED BY '\n'"
