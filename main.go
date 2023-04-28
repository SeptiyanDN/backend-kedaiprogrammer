package main

import (
	"fmt"
	"kedaiprogrammer/authorization"
	"kedaiprogrammer/core"
	"kedaiprogrammer/handler"
	"kedaiprogrammer/helper"
	"kedaiprogrammer/helpers"
	"kedaiprogrammer/kedaihelpers"
	"kedaiprogrammer/master/articles"
	"kedaiprogrammer/master/businesses"
	"kedaiprogrammer/master/categories"
	"kedaiprogrammer/master/services"
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
	tmphttpreadheadertimeout, _ := time.ParseDuration(viper.GetString("server.readheadertimeout") + "s")
	tmphttpreadtimeout, _ := time.ParseDuration(viper.GetString("server.readtimeout") + "s")
	tmphttpwritetimeout, _ := time.ParseDuration(viper.GetString("server.writetimeout") + "s")
	tmphttpidletimeout, _ := time.ParseDuration(viper.GetString("server.idletimeout") + "s")
	initGorm, err := core.InitGorm()
	router := gin.New()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	dbs := core.DBConnect()
	defer dbs.Dbx.Close()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://edukasi.kedaiprogrammer.com", "https://cms-edu.kedaiprogrammer.com", "http://localhost:3501", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Access-Control-Allow-Origin", "Authorization", "Content-Type", "x-requested-with"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		MaxAge:           12 * time.Hour,
		AllowCredentials: true,
	}))

	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Next()
	})
	S3 := core.ConnectAws()
	router.Use(func(c *gin.Context) {
		c.Set("S3", S3)
		c.Next()
	})
	Routing(router, dbs, initGorm)

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
	serviceRepository := services.NewRepository(initGorm, dbs)
	categoryRepository := categories.NewRepository(initGorm, dbs)
	articleRepository := articles.NewRepository(initGorm, dbs)

	// services
	userServices := users.NewServices(userRepository)
	businessServices := businesses.NewServices(businessRepository)
	serviceServices := services.NewServices(serviceRepository)
	categoryServices := categories.NewServices(categoryRepository)
	articleServices := articles.NewServices(articleRepository)
	authServices := authorization.NewServices()

	// handler
	userHandler := handler.NewUserHandler(userServices, authServices)
	businessHandler := handler.NewBusinessHandler(businessServices)
	serviceHandler := handler.NewServiceHandler(serviceServices)
	categoryHandler := handler.NewCategoryHandler(categoryServices)
	articleHandler := handler.NewArticleHandler(articleServices)

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
	serviceRouter := versioning.Group("services")
	{
		serviceRouter.Use(authMiddleware(authServices, userServices))
		serviceRouter.POST("/", serviceHandler.SaveService)
		serviceRouter.GET("/list", serviceHandler.GetAllServices)
		serviceRouter.GET("/:id", serviceHandler.GetDetailService)
	}
	categoryRouter := versioning.Group("categories")
	{
		categoryRouter.POST("/", categoryHandler.SaveCategory)
		categoryRouter.GET("/list", categoryHandler.GetAllCategory)
		categoryRouter.Use(authMiddleware(authServices, userServices))
		categoryRouter.GET("/:id", categoryHandler.GetDetailCategory)
	}
	articleRouter := versioning.Group("articles")
	{
		articleRouter.GET("/list", articleHandler.GetAll)
		articleRouter.GET("/:article_id", articleHandler.GetDetailArticle)
		articleRouter.Use(authMiddleware(authServices, userServices))
		articleRouter.POST("/", articleHandler.CreateData)
	}
	domainRouter := versioning.Group("domain")
	{
		domainRouter.GET("/account/balance", func(ctx *gin.Context) { handler.GetBalanceAccount(ctx) })
		domainRouter.GET("/price", func(ctx *gin.Context) { handler.GetPriceDomain(ctx) })
		domainRouter.GET("/list-all", func(ctx *gin.Context) { handler.GetAllDomainsHandler(ctx) })
		domainRouter.GET("/avaibility/:keyword", func(ctx *gin.Context) { handler.GetAvailabiltyDomain(ctx) })
		domainRouter.GET("/detail/:domain", func(ctx *gin.Context) { handler.GetDetailManageDomain(ctx) })
		domainRouter.GET("/new/:keyword", func(ctx *gin.Context) { handler.NewPrice(ctx) })

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
