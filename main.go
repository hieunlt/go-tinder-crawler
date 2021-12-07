package main

import (
	"github.com/joho/godotenv"
	"go-tinder-crawler/infrastructure"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	dbConn, err := infrastructure.NewDBConn(os.Getenv("DB_URI"))
	if err != nil {
		log.Fatal(err)
	}
	defer func(dbConn *infrastructure.DBConn) {
		_ = dbConn.Close()
	}(dbConn)
	for true {
		rand.Seed(time.Now().UnixNano())
		api := infrastructure.NewTinderAPI(os.Getenv("TOKEN"))
		profiles, err := api.GetNearbyProfiles()
		if err != nil {
			log.Fatal(err)
		}
		if err := dbConn.UpsertProfiles(profiles); err != nil {
			log.Fatal(err)
		}
		sleepDuration := time.Duration(rand.Intn(30))
		log.Printf("Sleeping for %02d second(s)\n", sleepDuration)
		time.Sleep(sleepDuration * time.Second)
	}

}
