package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {

	if err := godotenv.Load(); err != nil {

		wd, err := os.Getwd()
		if err != nil {
			log.Printf("get current working directory failed, err=%v", err)
		} else {
			log.Printf("[env] current working directory: %s", wd)
		}
		log.Fatalf("load .env file failed, err=%v", err)
	}
}

func MustHasEnvs(envs ...string) {
	for _, env := range envs {
		if os.Getenv(env) == "" {
			log.Fatalf("env %s is required. please check your .env file", env)
		}
	}
}
