package handler

import (
	"context"
	"net/http"
	"os"

	"RIP_lab1/internal/api"
	"RIP_lab1/internal/api/repository"
	"RIP_lab1/internal/models"
	"RIP_lab1/internal/pkg/auth"
	"RIP_lab1/internal/pkg/hash"
	"RIP_lab1/internal/pkg/minio"
	"RIP_lab1/internal/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	logger *logrus.Entry

	minio minio.Client
	redis redis.Client
	repo  api.Repo

	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
}

func initConfig(vp *viper.Viper) error {
	vp.AddConfigPath("/home/alik/GolandProjects/RIP_lab1/config")
	vp.SetConfigName("config")

	return vp.ReadInConfig()
}

func NewHandler(logger *logrus.Logger) *Handler {
	vp := viper.New()
	if err := initConfig(vp); err != nil {
		logger.Fatalf("error initializing configs: %s", err.Error())
	}

	repo, err := repository.NewRepo(logger, vp)
	if err != nil {
		logger.Error(err)
	}

	minioConfig := minio.InitConfig(vp)

	minioClient, err := minio.NewMinioClient(context.Background(), minioConfig, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	redisConfig := redis.InitRedisConfig(vp, logger)

	redisClient, err := redis.NewRedisClient(context.Background(), redisConfig, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	tokenManager, err := auth.NewManager(os.Getenv("TOKEN_SECRET"))
	if err != nil {
		logger.Fatalln(err)
	}

	return &Handler{
		repo:         repo,
		minio:        minioClient,
		logger:       logger.WithField("component", "handler"),
		redis:        redisClient,
		hasher:       hash.NewSHA256Hasher(os.Getenv("SALT")),
		tokenManager: tokenManager,
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	h.logger.Println("Server start up")

	r := gin.Default()
	r.Use(CORSMiddleware())

	r.LoadHTMLGlob("templates/*")

	r.Static("/image", "./resources")
	r.Static("/style", "./style")
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", h.Ping)

	// полезные нагрузки
	r.GET("/payloads", h.WithAuthCheck([]models.Role{}), h.GetPayloadList)
	r.GET("/payloads/:id", h.GetCardPayloadById)
	r.POST("/payloads", h.WithAuthCheck([]models.Role{models.Admin}), h.CreateNewPayload)
	r.PUT("/payloads/:id", h.WithAuthCheck([]models.Role{models.Admin}), h.ChangePayload)
	r.DELETE("/payloads/:id", h.WithAuthCheck([]models.Role{models.Admin}), h.DeletePayloadById)

	// удалить после перехода на фронт
	// r.POST("/payload/:id", h.WithAuthCheck([]models.Role{models.Admin}), h.DeletePayloadById)

	// полёты ракет-носителей
	r.GET("/rocket_flights", h.WithAuthCheck([]models.Role{models.Admin, models.Client}), h.GetRocketFlightList)
	r.GET("/rocket_flights/:id", h.WithAuthCheck([]models.Role{models.Client, models.Admin}), h.GetRocketFlightById)
	r.PUT("/rocket_flights", h.WithAuthCheck([]models.Role{models.Client, models.Admin}), h.ChangeRocketFlight)
	r.PUT("/rocket_flights/form", h.WithAuthCheck([]models.Role{models.Client, models.Admin}), h.FormRocketFlight)
	r.PUT("/rocket_flights/:id/response", h.WithAuthCheck([]models.Role{models.Admin}), h.ResponceRocketFlight)
	r.DELETE("/rocket_flights", h.WithAuthCheck([]models.Role{models.Client}), h.DeleteRocketFlight)
	// async service
	r.PUT("rocket_flights/finish_calculating", h.FinishCalculating)

	// формирование информации о будущем полёте через полезные нагрузки
	r.POST("/payloads/rocket_flight", h.WithAuthCheck([]models.Role{models.Client, models.Admin}), h.AddPayloadToFlight)

	// m-m
	r.DELETE("/flights_payloads/payload/:id", h.WithAuthCheck([]models.Role{models.Client}), h.DeletePayloadFromFlight)
	r.PUT("/flights_payloads/payload/:id/count/:count", h.WithAuthCheck([]models.Role{models.Client}), h.ChangeCountFlightsPayload)

	//user
	r.POST("/sign_in", h.SignIn)
	r.POST("/sign_up", h.SignUp)
	r.POST("/logout", h.Logout)
	r.PUT("/profile", h.ChangeProfile)

	// для фронта на будущее
	r.GET("/check-auth", h.WithAuthCheck([]models.Role{models.Client, models.Admin}), h.CheckAuth)

	return r
}

// Ping godoc
// @Summary      Show hello text
// @Description  very very friendly response
// @Tags         Tests
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       /ping [get]
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "pong",
		})
}
