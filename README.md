# demo_api

## requirement( ideal) (you can run on lower setup also)
 - go version go1.13.4
 - mysql  Ver 8.0.20-0ubuntu0.19.10.1
 - redis-cli 5.0.5
 
## run
 - create MySql database ( :=> create database myDb) where myDb is database name
    - if you want to change host, database name, port and password please change it in config.json file
 
 - to setup tables in db =>  go run main.go -config=config.json -process=setup 
 
 - to run api  =>  go run main.go -config=config.json -process=run_api 
    - server will listen on port 9001 you can modigy port number in server.go file

