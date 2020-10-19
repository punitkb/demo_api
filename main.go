package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"flag"
	"sezzle_api/src/setup_wizard"
	"sezzle_api/src/api"
	"os"
)


func main() {

	fmt.Println("Sezzle setting up ")
	fmt.Println("---------------------")
    configPath := flag.String("config", "", "Configuration file location")
    process := flag.String("process", "", "Process type")
    flag.Parse()

    viper.SetConfigFile(*configPath)
    err := viper.ReadInConfig() // Find and read the config file
    if err != nil { // Handle errors reading the config file
            panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }

	
	switch *process {
		case "setup":
			setup_wizard.RunWizard()
		case "run_api":
			api.RunServer()
		default:
			fmt.Println("Please choose options from 'setup' or 'run_api' , eg. => go run main.go -config=config.json -process=setup")
			os.Exit(1)
	}
}