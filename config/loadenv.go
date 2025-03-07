package config

import "os"

func GetMongoUri() string {
	return os.Getenv("MONGO_URI")
}
