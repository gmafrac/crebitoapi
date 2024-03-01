package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {

	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file %v", err)
		}
	}
}

func GetDBUrl() string {

	db_user := os.Getenv("DATABASE_USER")
	db_password := os.Getenv("DATABASE_PASSWORD")
	db_name := os.Getenv("DATABASE_NAME")
	db_host := os.Getenv("DATABASE_HOST")
	db_port := os.Getenv("DATABASE_PORT")
	db_engine := os.Getenv("DATABASE_ENGINE")

	return db_engine + "://" + db_user + ":" + db_password + "@" + db_host + ":" + db_port + "/" + db_name
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
