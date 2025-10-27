package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"elogika.vsb.cz/backend/docs"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/modules/activities"
	"elogika.vsb.cz/backend/modules/auth"
	authCrons "elogika.vsb.cz/backend/modules/auth/crons"
	"elogika.vsb.cz/backend/modules/auth/helpers"
	"elogika.vsb.cz/backend/modules/auth/middlewares"
	"elogika.vsb.cz/backend/modules/categories"
	"elogika.vsb.cz/backend/modules/chapters"
	"elogika.vsb.cz/backend/modules/classes"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/course_item_terms"
	"elogika.vsb.cz/backend/modules/course_items"
	"elogika.vsb.cz/backend/modules/courses"
	"elogika.vsb.cz/backend/modules/files"
	"elogika.vsb.cz/backend/modules/print"
	printCrons "elogika.vsb.cz/backend/modules/print/crons"
	"elogika.vsb.cz/backend/modules/questions"
	"elogika.vsb.cz/backend/modules/recognizer"
	"elogika.vsb.cz/backend/modules/templates"
	"elogika.vsb.cz/backend/modules/tests"
	testCrons "elogika.vsb.cz/backend/modules/tests/crons"
	"elogika.vsb.cz/backend/modules/users"
	"elogika.vsb.cz/backend/utils/certstore"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB(true)
	if initializers.GlobalAppConfig.ACCESS_TOKEN_REVOKE_SYNC {
		helpers.StartRevokedTokenSync(initializers.DB, 10*time.Minute)
	}
}

func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error however you want
				fmt.Printf("panic recovered: %s\n%s\n", err, debug.Stack())

				// Replace default 500 response with custom JSON
				jsonObj := &common.ErrorResponse{
					Code:    500,
					Message: "Thread panicked",
				}
				c.Abort()
				c.JSON(jsonObj.Code, jsonObj)
			}
		}()
		c.Next()
	}
}

// @title           eLogika public API
// @version         1.0
// @description     Public api for e-learning system eLogika developer at VŠB-TUO

// @contact.name   Daniel Makovský
// @contact.url    https://www.makovsky.dev/contact
// @contact.email  daniel@makovsky.dev

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9000

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @supportedSubmitMethods []
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	if initializers.GlobalAppConfig.GIN_RELEASE_MODE {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(CustomRecovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8000", "https://elogika.vsb.cz", "http://elogika.vsb.cz"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "x-AS-ROLE"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	docs.SwaggerInfo.BasePath = ""

	authCrons.DeleteExpiredExpirations()
	testCrons.ExpireReadyTests()
	testCrons.FinishActiveTests()
	printCrons.ClearTempDir()

	c := cron.New()
	c.AddFunc("0 0 0 0 0", func() {
		log.Println("Running job: DeleteExpiredExpirations", time.Now())
		go authCrons.DeleteExpiredExpirations()
		log.Println("Running job: ClearTempDir", time.Now())
		go printCrons.ClearTempDir()
	})
	c.AddFunc("* * * * 5", func() {
		log.Println("Running job: ExpireReadyTests", time.Now())
		go testCrons.ExpireReadyTests()
	})
	c.AddFunc("* * * * 5", func() {
		log.Println("Running job: FinishActiveTests", time.Now())
		go testCrons.FinishActiveTests()
	})
	c.Start()

	v2api := r.Group("/api/v2")
	{
		public := v2api.Group("")
		{
			auth.RegisterRoutesUnauth(public)
			files.RegisterRoutesUnauth(public)
		}
		private := v2api.Group("", middlewares.AuthMiddleware())
		{
			questions.RegisterRoutes(private)
			auth.RegisterRoutes(private)
			chapters.RegisterRoutes(private)
			files.RegisterRoutes(private)
			categories.RegisterRoutes(private)
			courses.RegisterRoutes(private)
			users.RegisterRoutes(private)
			templates.RegisterRoutes(private)
			course_items.RegisterRoutes(private)
			course_item_terms.RegisterRoutes(private)
			tests.RegisterRoutes(private)
			print.RegisterRoutes(private)
			classes.RegisterRoutes(private)
			activities.RegisterRoutes(private)
			recognizer.RegisterRoutes(private)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		err := common.ErrorResponse{
			Code:    404,
			Message: "Requested endpoint not found",
		}
		c.JSON(err.Code, err)

	})

	r.NoMethod(func(c *gin.Context) {
		err := common.ErrorResponse{
			Code:    405,
			Message: "Requested method not found",
		}
		c.JSON(err.Code, err)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if initializers.GlobalAppConfig.PROTOCOL == "https" {
		addr := ":" + strconv.Itoa(int(initializers.GlobalAppConfig.PORT))
		cert, err := certstore.LoadCertCommon(initializers.GlobalAppConfig.CERTPATH, initializers.GlobalAppConfig.CERTPASS)
		if err != nil {
			panic(err.Error())
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{*cert},
		}

		server := &http.Server{
			Addr:      addr,
			Handler:   r,
			TLSConfig: tlsConfig,
		}

		log.Println("Listening and serving HTTPS on " + addr + "\n")
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	} else {
		r.Run()
	}
}
