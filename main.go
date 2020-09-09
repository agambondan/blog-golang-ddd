package main

import (
	"Repository-Pattern/infrastructure/auth"
	"Repository-Pattern/infrastructure/persistence"
	"Repository-Pattern/interfaces"
	"Repository-Pattern/interfaces/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	//To load our environmental variables.
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

func main() {
	dbDriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	//redis details
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	services, err := persistence.NewRepositories(dbDriver, user, password, port, host, dbname)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	err = services.AutoMigrate()
	if err != nil {
		panic(err)
	}

	redisService, err := auth.NewRedisDB(redisHost, redisPort, redisPassword)
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()

	users := interfaces.NewUsers(services.User, redisService.Auth, tk)
	authenticate := interfaces.NewAuthenticate(services.User, redisService.Auth, tk)

	r := gin.Default()
	r.Use(middlewares.CORSMiddleware()) //For CORS

	//user routes
	r.POST("/users", users.SaveUser)
	r.GET("/users", users.GetUsers)
	r.GET("/users/:user_id", users.GetUser)

	//authentication routes
	r.POST("/login", authenticate.Login)
	r.POST("/logout", authenticate.Logout)
	r.POST("/refresh", authenticate.Refresh)

	//Starting the application
	appPort := os.Getenv("PORT") //using heroku host
	if appPort == "" {
		appPort = "8888" //localhost
	}
	log.Fatal(r.Run(":" + appPort))
}
