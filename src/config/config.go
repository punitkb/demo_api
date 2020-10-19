package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/go-redis/redis"
	"os"	
	"strings"
)


//init mysql connection
func InitDb() (*gorm.DB, error) {
	dbConfig := viper.Get("db").(map[string]interface{})
	dbConnectionString := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig["user"].(string),
		dbConfig["password"].(string),
		dbConfig["protocol"].(string),
		dbConfig["host"].(string),
		dbConfig["port"].(string),
		dbConfig["name"].(string),
	)
	return gorm.Open(dbConfig["dialect"].(string), dbConnectionString)
}

//init redis connectinon
func InitRedisdb() (*redis.Client, error) {
	redisConfig := viper.Get("redis").(map[string]interface{})
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisConfig["host"].(string) +":"+redisConfig["port"].(string),
		Password: "",
		DB: 0,
	})
    _, err := redisClient.Ping().Result()
    if err != nil {
            return nil, err
    }
    return redisClient, nil
}

//
func GetRealRootDirectory() string {
	reverse := func(pathParts []string) []string {

		for i, j := 0, len(pathParts)-1; i < j; i, j = i+1, j-1 {
			pathParts[i], pathParts[j] = pathParts[j], pathParts[i]
		}
		return pathParts
	}

	potentialWorkingDirectory, _ := os.Getwd()
	reversedPathParts := reverse(strings.Split(potentialWorkingDirectory, "/"))

	realWorkingDirectoryPathParts := make([]string, 0)
	shouldAppend := false
	for _, pathPart := range reversedPathParts {
		if pathPart == "sezzle_api" {
			shouldAppend = true
		}
		if shouldAppend && len(pathPart) > 0 {
			realWorkingDirectoryPathParts = append(realWorkingDirectoryPathParts, pathPart)
		}
	}
	return "/" + strings.Join(reverse(realWorkingDirectoryPathParts), "/")

}
