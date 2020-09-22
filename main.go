package main

import (
	"Repository-Pattern/infrastructure/auth"
	"Repository-Pattern/infrastructure/persistence"
	"Repository-Pattern/interfaces"
	"Repository-Pattern/interfaces/fileupload"
	"Repository-Pattern/interfaces/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
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
	err = services.Seeder()
	if err != nil {
		panic(err)
	}
	err = services.AddForeignKey()
	if err != nil {
		panic(err)
	}

	redisService, err := auth.NewRedisDB(redisHost, redisPort, redisPassword)
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()
	fd := fileupload.NewFileUpload()

	users := interfaces.NewUsers(services.User, services.Post, services.Role, redisService.Auth, tk)
	posts := interfaces.NewPost(services.Post, services.User, services.PostLabel, services.PostCategory, fd, redisService.Auth, tk)
	labels := interfaces.NewLabel(services.Label, services.Role, services.User, redisService.Auth, tk)
	categories := interfaces.NewCategory(services.Category, services.Role, services.User, redisService.Auth, tk)
	roles := interfaces.NewRole(services.Role, services.User, redisService.Auth, tk)
	authenticate := interfaces.NewAuthenticate(services.User, redisService.Auth, tk)

	r := gin.Default()
	r.Use(middlewares.CORSMiddleware()) //For CORS

	// home routes
	r.GET("/", home)

	//user routes
	r.POST("/users", middlewares.MaxSizeAllowed(8192000), users.SaveUser)
	r.GET("/users", users.GetUsers)
	r.GET("/users/:user_id", users.GetUser)
	r.GET("/private/users", middlewares.AuthMiddleware(), users.GetPrivateUsers)

	//role routes
	r.POST("/role", middlewares.AuthMiddleware(), roles.SaveRole)
	r.PUT("/role/:role_id", middlewares.AuthMiddleware(), roles.UpdateRole)
	r.GET("/role/:role_id", roles.GetRole)
	r.DELETE("/role/:role_id", middlewares.AuthMiddleware(), roles.DeleteRole)
	r.GET("/roles", roles.GetAllRole)

	//label routes
	r.POST("/label", middlewares.AuthMiddleware(), labels.SaveLabel)
	r.PUT("/label/:label_id", middlewares.AuthMiddleware(), labels.UpdateLabel)
	r.GET("/label/:label_id", labels.GetLabel)
	r.DELETE("/label/:label_id", middlewares.AuthMiddleware(), labels.DeleteLabel)
	r.GET("/labels", labels.GetAllLabel)

	//category routes
	r.POST("/category", middlewares.AuthMiddleware(), categories.SaveCategory)
	r.PUT("/category/:category_id", middlewares.AuthMiddleware(), categories.UpdateCategory)
	r.GET("/category/:category_id", categories.GetCategory)
	r.DELETE("/category/:category_id", middlewares.AuthMiddleware(), categories.DeleteCategory)
	r.GET("/categories", categories.GetAllCategory)

	//post routes
	r.POST("/post", middlewares.AuthMiddleware(), middlewares.MaxSizeAllowed(8192000), posts.SavePost)
	r.PUT("/post/:post_id", middlewares.AuthMiddleware(), middlewares.MaxSizeAllowed(8192000), posts.UpdatePost)
	r.GET("/post/:post_id", middlewares.AuthMiddleware(), posts.GetPostAndCreator)
	r.DELETE("/post/:post_id", middlewares.AuthMiddleware(), posts.DeletePost)
	r.GET("/posts", posts.GetAllPost)

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

func home(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello World From My Office")
}
