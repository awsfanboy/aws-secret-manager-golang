package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"aws-secret-manager-test/Logger"
	"aws-secret-manager-test/Models"
	"github.com/joho/godotenv"
	"log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	_ "github.com/lib/pq"
	"os"
)

var DB *sql.DB

func init() {
	var err error
	err = godotenv.Load() 
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	var err error

	databaseAuth := getDatabaseAuth()
	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		databaseAuth.Host, databaseAuth.Port, databaseAuth.UserName, databaseAuth.Password, os.Getenv("DB_NAME"))

	DB, err = sql.Open("postgres", psql)
	if err != nil {
		Logger.AddLogger(Logger.ERROR, "Database driver error")
		panic(err)

	}
	if err = DB.Ping(); err != nil {
		Logger.AddLogger(Logger.ERROR, "Database parameters error")
		panic(err)
	}
	Logger.AddLogger(Logger.INFO, "Connected to Database")
}

func getDatabaseAuth() Models.DatabaseAuth {
	secretName := os.Getenv("AWS_SECRET_NAME")
	region := os.Getenv("AWS_REGION")

	svc := secretsmanager.New(session.New(&aws.Config {
		Region: &region,
	}))

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(input)
	var databaseAuth = Models.DatabaseAuth{}

	if err == nil {
		var secretString, decodedBinarySecret string

		if result.SecretString != nil {
			secretString = *result.SecretString
			json.Unmarshal([]byte(secretString) , &databaseAuth)
		} else {
			decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
			len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
			if err != nil {
				fmt.Println("Base64 Decode Error:", err)
			}
			decodedBinarySecret = string(decodedBinarySecretBytes[:len])
			json.Unmarshal([]byte(decodedBinarySecret) , &databaseAuth)
		}
	}
	return databaseAuth
}
