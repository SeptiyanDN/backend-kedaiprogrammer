package main

import (
	"fmt"
	"kedaiprogrammer/authorization"
	"kedaiprogrammer/businesses"
	"kedaiprogrammer/categories"
	"kedaiprogrammer/core"
	"kedaiprogrammer/handler"
	"kedaiprogrammer/helper"
	"kedaiprogrammer/helpers"
	"kedaiprogrammer/kedaihelpers"
	"kedaiprogrammer/users"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func main() {
	//---- READ CONFIG JSON ----
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetConfigName("app.conf")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	initGorm, err := core.InitGorm()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	initGorm.AutoMigrate(categories.Category{}, businesses.Business{}, users.User{})
	dbs := core.DBConnect()
	defer dbs.Dbx.Close()

	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:           []string{"https://septiyan.my.id"},
		AllowMethods:           []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:           []string{"Access-Control-Allow-Origin", "Authorization", "Content-Type"},
		ExposeHeaders:          []string{"Content-Length"},
		AllowCredentials:       true,
		AllowAllOrigins:        true,
		AllowWildcard:          true,
		AllowBrowserExtensions: true,
		MaxAge:                 12 * time.Hour,
	}))

	router.Use(gin.Recovery())
	Routing(router, dbs, initGorm)

	tmphttpreadheadertimeout, _ := time.ParseDuration(viper.GetString("server.readheadertimeout") + "s")
	tmphttpreadtimeout, _ := time.ParseDuration(viper.GetString("server.readtimeout") + "s")
	tmphttpwritetimeout, _ := time.ParseDuration(viper.GetString("server.writetimeout") + "s")
	tmphttpidletimeout, _ := time.ParseDuration(viper.GetString("server.idletimeout") + "s")

	s := &http.Server{
		Addr:              ":" + viper.GetString("server.port"),
		Handler:           router,
		ReadHeaderTimeout: tmphttpreadheadertimeout,
		ReadTimeout:       tmphttpreadtimeout,
		WriteTimeout:      tmphttpwritetimeout,
		IdleTimeout:       tmphttpidletimeout,
		//MaxHeaderBytes:    1 << 20,
	}
	fmt.Println("🚀 Server running on port:", viper.GetString("server.port"))
	s.ListenAndServe()

}
func Routing(router *gin.Engine, dbs kedaihelpers.DBStruct, initGorm *gorm.DB) {
	time.Local = time.UTC

	router.Static("/logo-path", viper.GetString("upload_path.logo"))
	router.Any("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "OK",
			"Message": "Welcome to " + viper.GetString("appName"),
		})
	})

	// repository
	userRepository := users.NewRepository(initGorm)
	businessRepository := businesses.NewRepository(initGorm)
	categoryRepository := categories.NewRepository(initGorm, dbs)

	// services
	userServices := users.NewServices(userRepository)
	businessServices := businesses.NewServices(businessRepository)
	categoryServices := categories.NewServices(categoryRepository)
	authServices := authorization.NewServices()

	// handler
	userHandler := handler.NewUserHandler(userServices, authServices)
	businessHandler := handler.NewBusinessHandler(businessServices)
	categoryHandler := handler.NewCategoryHandler(categoryServices)

	versioning := router.Group("/api/v1")

	authRouter := versioning.Group("auth")
	{
		authRouter.POST("login", userHandler.Login)
		authRouter.POST("/register", userHandler.RegisterUser)

	}
	businessRouter := versioning.Group("business")
	{
		businessRouter.Use(authMiddleware(authServices, userServices))
		businessRouter.POST("/", businessHandler.SaveBusiness)
		businessRouter.GET("/list", businessHandler.GetAllBusiness)
	}
	categoryRouter := versioning.Group("categories")
	{
		categoryRouter.Use(authMiddleware(authServices, userServices))
		categoryRouter.POST("/", categoryHandler.SaveCategory)
		categoryRouter.GET("/list", categoryHandler.GetAllCategory)
		categoryRouter.GET("/:id", categoryHandler.GetDetailCategory)
	}
}

func authMiddleware(authServices authorization.Services, userServices users.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized Access", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}
		token, err := authServices.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized Access", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claims, ok := token.Claims.(*authorization.JWTClaim)
		if !ok || !token.Valid {
			response := helpers.APIResponse("Unauthorized Access", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if claims.ExpiresAt < time.Now().Local().Unix() {
			response := helpers.APIResponse("Token Expired", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		c.Set("current", claims)
		c.Next()

	}
}
