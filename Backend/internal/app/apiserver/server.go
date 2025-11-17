package apiserver

import (
	"backed-teacher/internal/app/model"
	"backed-teacher/internal/app/store"
	"encoding/json"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
)

type server struct {
	router    *gin.Engine
	logger    *logrus.Logger
	store     store.Store
	jwtSecret string
}

func newServer(store store.Store, jwtSecret string) *server {
	s := &server{
		router:    gin.Default(),
		logger:    logrus.New(),
		store:     store,
		jwtSecret: jwtSecret,
	}
	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	// Настройка CORS
	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.router.POST("/login-users", s.handleUsersCreate())
	s.router.POST("/sessions", s.handleSessionsCreate())
	
	protected := s.router.Group("/")
	protected.Use(s.authMiddleware())
	{
		protected.GET("/profile", s.handleProfile())
	}
}

func (s *server) handleUsersCreate() gin.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(c *gin.Context) {
		req := &request{}
		if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		u := model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(&u); err != nil {
			s.error(c, http.StatusUnprocessableEntity, err)
			return
		}

		token, err := GenerateToken(u.ID, s.jwtSecret)
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		s.respond(c, map[string]string{"token": token}, http.StatusCreated)

	}
}

func (s *server) handleSessionsCreate() gin.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(c *gin.Context) {
		req := &request{}
		if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}
		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(c, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		token, err := GenerateToken(u.ID, s.jwtSecret)
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		s.respond(c, map[string]string{"token": token}, http.StatusOK)
	}
}

func (s *server) error(c *gin.Context, code int, err error) {
	s.respond(c, map[string]string{"error": err.Error()}, code)
}

func (s *server) respond(c *gin.Context, data interface{}, code int) {
	response := gin.H{"message": "ok"}
	if data != nil {
		response["data"] = data
	}
	c.JSON(code, response)
}

func (s *server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			s.error(c, http.StatusUnauthorized, errors.New("authorization header required"))
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			s.error(c, http.StatusUnauthorized, errors.New("invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := authHeader[len(bearerPrefix):]
		claims, err := ValidateToken(tokenString, s.jwtSecret)
		if err != nil {
			s.error(c, http.StatusUnauthorized, errors.New("invalid token"))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func (s *server) handleProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			s.error(c, http.StatusUnauthorized, errors.New("user not authenticated"))
			return
		}

		u, err := s.store.User().FindByID(userID.(int))
		if err != nil {
			s.error(c, http.StatusNotFound, err)
			return
		}

		u.Sanitize()
		s.respond(c, u, http.StatusOK)
	}
}
