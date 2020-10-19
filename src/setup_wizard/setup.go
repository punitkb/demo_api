package setup_wizard

import (
	"fmt"
	"sezzle_api/src/config"
	"sezzle_api/src/repository"

)


func RunWizard() {
	fmt.Println("Sezzle initial setup wizard")
	fmt.Println("---------------------")

	//initialize DB connection
	db, err := config.InitDb()
	defer db.Close()
	if err != nil {
		panic(err)
	}
	//intialize redis connection
	redisDb,err := config.InitRedisdb()

	connectionRepo := repository.NewRepository(db,redisDb)
    if err := connectionRepo.Init(); err != nil {
        panic(err)
    }

    // create tables 
	err  = connectionRepo.CreateTables(db)
	if err != nil {
		panic(err)
	}
}



