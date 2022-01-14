package main

import (
	"bwa_golang/auth"
	"bwa_golang/handler"
	"bwa_golang/helper"
	"bwa_golang/image"
	"bwa_golang/project"
	"bwa_golang/transaction"
	"bwa_golang/user"
	webHandler "bwa_golang/web/handler"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OptionMessage(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
}
func main() {
	dsn := "root:root@tcp(127.0.0.1:8889)/bwa_golang?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	// REPOSITORY START
	userRepository := user.NewRepository(db)
	projectRepository := project.NewRepository(db)
	imageRepository := image.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)
	// REPOSITORY END

	// SERVICE START
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	projectService := project.NewService(projectRepository)
	imageService := image.NewService(imageRepository)
	transactionService := transaction.NewService(transactionRepository, projectRepository)
	// SERVICE END

	// HANDLER START
	userHandler := handler.NewUserHandler(userService, authService)
	projectHandler := handler.NewProjectHandler(projectService)
	imageHandler := handler.NewImageHandler(imageService, projectService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	userWebHandler := webHandler.NewUserHandler(userService)
	projectWebHandler := webHandler.NewProjectHandler(projectService, userService, imageService)
	transactionWebHandler := webHandler.NewTransactionHandler(transactionService)
	sessionHandler := webHandler.NewSessionHandler(userService)
	// HANDLER END

	// DEFAULT API
	router := gin.Default()

	router.Static("/project_image", "./project_image")
	router.Static("/images", "./images")
	router.Static("/assets", "./web/assets")
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Authorization,content-type"}
	config.AllowOriginFunc = func(origin string) bool {
		return origin == "http://localhost:3000"
	}
	router.Use(cors.New(config))

	cookieStore := cookie.NewStore([]byte("bwa_secret_key_s4j4"))

	router.Use(sessions.Sessions("adminId", cookieStore))

	router.HTMLRender = loadTemplates("./web/templates")

	api := router.Group("/api/v1")
	// ROUTES USER START
	api.POST("/user", userHandler.RegisterUser)
	api.POST("/session", userHandler.LoginUser)
	api.POST("/email_checkers", userHandler.CheckEmail)
	api.POST("/photo", authMiddleware(userService, authService), userHandler.UploadImage)
	api.GET("/users/fetch", authMiddleware(userService, authService), userHandler.FetchUser)
	// api.OPTIONS("/users/fetch", OptionMessage)
	// ROUTES USER END

	// ROUTES PROJECT START
	api.GET("/projects", projectHandler.GetProject)
	api.GET("/project/:id", projectHandler.DetailProject)
	api.POST("/project", authMiddleware(userService, authService), projectHandler.CreateProject)
	api.PUT("/project/:id", authMiddleware(userService, authService), projectHandler.UpdateProject)
	// ROUTES PROJECT END

	// ROUTES IMAGES START
	api.POST("/image/:id", authMiddleware(userService, authService), imageHandler.UploadImage)
	// ROUTES IMAGES END

	// ROUTES TRANSACTION START
	api.GET("/projects/:id/transactions", authMiddleware(userService, authService), transactionHandler.GetTransactionByProjectId)
	api.GET("/transactions", authMiddleware(userService, authService), transactionHandler.GetTransactionByUserId)
	api.POST("/transactions", authMiddleware(userService, authService), transactionHandler.CreateTransaction)
	api.POST("/payment", transactionHandler.ResponsePayment)
	// ROUTES TRANSACTION END

	// WEB ADMIN START

	// USER MENU START
	router.GET("/admin_user", authAdminMiddleware(), userWebHandler.Index)
	router.GET("/admin_user/input", authAdminMiddleware(), userWebHandler.New)
	router.GET("/admin_user/edit/:id", authAdminMiddleware(), userWebHandler.Edit)
	router.GET("/admin_user/delete/:id", authAdminMiddleware(), userWebHandler.Delete)
	router.POST("/admin_user", authAdminMiddleware(), userWebHandler.Create)
	router.POST("/admin_user/edit/:id", authAdminMiddleware(), userWebHandler.Update)
	// USER MENU END

	// PROJECT MENU START
	router.GET("/admin_project", authAdminMiddleware(), projectWebHandler.Index)
	router.GET("/admin_project/input", authAdminMiddleware(), projectWebHandler.New)
	router.GET("/admin_project/image/:id", authAdminMiddleware(), projectWebHandler.Image)
	router.GET("/admin_project/edit/:id", authAdminMiddleware(), projectWebHandler.Edit)
	router.GET("/admin_project/show/:id", authAdminMiddleware(), projectWebHandler.Show)
	router.POST("/admin_project", authAdminMiddleware(), projectWebHandler.Create)
	router.POST("/admin_project/image/:id", authAdminMiddleware(), projectWebHandler.CreateImage)
	router.POST("/admin_project/edit/:id", authAdminMiddleware(), projectWebHandler.Update)
	// PROJECT MENU END

	// TRANSACTION MENU START
	router.GET("/admin_transaction", authAdminMiddleware(), transactionWebHandler.Index)
	// TRANSACTION MENU END

	router.GET("/login_admin", sessionHandler.Login)
	router.GET("/logout_admin", sessionHandler.Logout)
	router.POST("/login_admin", sessionHandler.Sign_in)
	// WEB ADMIN END
	router.Run()

	fmt.Println("Success")
}
func authMiddleware(u user.Service, a auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", "error", http.StatusUnauthorized, nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		var tokenString string

		arrayToken := strings.Split(authHeader, " ")

		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		valid, err := a.ValidateToken(tokenString)

		if err != nil || !valid.Valid {
			response := helper.APIResponse("Unauthorized", "error", http.StatusUnauthorized, nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := valid.Claims.(jwt.MapClaims)

		if !ok {
			response := helper.APIResponse("Unauthorized", "error", http.StatusUnauthorized, nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		id := int(claim["id"].(float64))
		userData, err := u.GetUserById(id)

		if err != nil || userData.Id == 0 {
			response := helper.APIResponse("Unauthorized", "error", http.StatusUnauthorized, nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", userData)
	}
}
func authAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("adminId") == nil {
			c.Redirect(http.StatusFound, "/login_admin")
			return
		}
	}
}
func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")

	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")

	if err != nil {
		panic(err.Error())
	}

	for _, v := range includes {
		layoutsCopy := make([]string, len(layouts))
		copy(layoutsCopy, layouts)
		files := append(layoutsCopy, v)
		r.AddFromFiles(filepath.Base(v), files...)
	}
	return r
}
