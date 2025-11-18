package apiserver

import (
	"backed-teacher/internal/app/model"
	"backed-teacher/internal/app/store"
	"backed-teacher/internal/mailer"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errInvalidResetToken        = errors.New("invalid reset token")
	errExpiredResetToken        = errors.New("reset token expired")
	errInvalidVerificationToken = errors.New("invalid verification token")
	errExpiredVerificationToken = errors.New("verification token expired")
	errEmailAlreadyVerified     = errors.New("email already verified")
)

const (
	userTokenTTL           = 24 * time.Hour
	adminTokenTTL         = 7 * 24 * time.Hour
	resetTokenTTL          = 30 * time.Minute
	verificationTokenTTL   = 24 * time.Hour
)

type server struct {
	router            *gin.Engine
	logger            *logrus.Logger
	store             store.Store
	jwtSecret         string
	mailer            *mailer.Mailer
	passwordResetURL  string
	verificationURL   string
	uploadDir         string
}

func newServer(store store.Store, jwtSecret string, mailer *mailer.Mailer, passwordResetURL string, verificationURL string, uploadDir string) *server {
	s := &server{
		router:            gin.Default(),
		logger:            logrus.New(),
		store:             store,
		jwtSecret:         jwtSecret,
		mailer:            mailer,
		passwordResetURL:  passwordResetURL,
		verificationURL:   verificationURL,
		uploadDir:         uploadDir,
	}
	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(func(c *gin.Context) {
		s.logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("Incoming request")
		c.Next()
	})

	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.router.Static("/uploads", s.uploadDir)

	s.router.POST("/login-users", s.handleUsersCreate())
	s.router.POST("/sessions", s.handleSessionsCreate())
	s.router.POST("/verify-email", s.handleVerifyEmail())
	s.router.POST("/resend-verification", s.handleResendVerification())
	s.router.GET("/pages/*page", s.handlePublicPageBlocks())
	s.router.GET("/posts", s.handlePublicPostsList())
	s.router.GET("/posts/:id", s.handlePublicPostDetail())
	s.router.GET("/posts/:id/comments", s.handlePublicPostComments())
	s.router.GET("/users/:id", s.handlePublicUserGet())

	admin := s.router.Group("/admin")
	{
		admin.POST("/login", s.handleAdminLogin())
		admin.POST("/forgot-password", s.handleAdminForgotPassword())
		admin.POST("/reset-password", s.handleAdminResetPassword())
		protectedAdmin := admin.Group("")
		protectedAdmin.Use(s.adminAuthMiddleware())
		{
			protectedAdmin.GET("/me", s.handleAdminMe())
			protectedAdmin.GET("/dashboard", s.handleAdminDashboard())
			protectedAdmin.GET("/blocks", s.handleAdminBlocksList())
			protectedAdmin.POST("/blocks", s.handleAdminBlockCreate())
			protectedAdmin.PATCH("/blocks/:id", s.handleAdminBlockUpdate())
			protectedAdmin.DELETE("/blocks/:id", s.handleAdminBlockDelete())
			protectedAdmin.PATCH("/blocks/reorder", s.handleAdminBlockReorder())

			protectedAdmin.GET("/posts", s.handleAdminPostsList())
			protectedAdmin.POST("/posts", s.handleAdminPostCreate())
			protectedAdmin.GET("/posts/:id", s.handleAdminPostDetail())
			protectedAdmin.PATCH("/posts/:id", s.handleAdminPostUpdate())
			protectedAdmin.DELETE("/posts/:id", s.handleAdminPostDelete())

			protectedAdmin.GET("/gallery", s.handleAdminGalleryList())
			protectedAdmin.POST("/gallery", s.handleAdminGalleryCreate())
			protectedAdmin.PATCH("/gallery/:id", s.handleAdminGalleryUpdate())
			protectedAdmin.DELETE("/gallery/:id", s.handleAdminGalleryDelete())

			protectedAdmin.GET("/slider", s.handleAdminSliderList())
			protectedAdmin.POST("/slider", s.handleAdminSliderCreate())
			protectedAdmin.PATCH("/slider/:id", s.handleAdminSliderUpdate())
			protectedAdmin.DELETE("/slider/:id", s.handleAdminSliderDelete())

			protectedAdmin.POST("/upload", s.handleAdminUpload())
			protectedAdmin.PATCH("/slider/reorder", s.handleAdminSliderReorder())

			protectedAdmin.GET("/comments", s.handleAdminCommentsList())
			protectedAdmin.POST("/comments/:id/reply", s.handleAdminCommentReply())
			protectedAdmin.DELETE("/comments/:id", s.handleAdminCommentDelete())

			protectedAdmin.GET("/users", s.handleAdminUsersList())
			protectedAdmin.GET("/users/:id", s.handleAdminUserGet())
			protectedAdmin.PATCH("/users/:id/ban", s.handleAdminUserBan())
		}
	}

	protected := s.router.Group("/")
	protected.Use(s.authMiddleware())
	{
		protected.GET("/profile", s.handleProfile())
		protected.POST("/posts/:id/comments", s.handleUserCommentCreate())
		protected.POST("/comments/:id/reply", s.handleUserCommentReply())
		protected.DELETE("/comments/:id", s.handleUserCommentDelete())
	}

	s.router.NoRoute(func(c *gin.Context) {
		s.error(c, http.StatusNotFound, errors.New("page not found"))
	})
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

		// Генерируем токен верификации
		verificationToken, err := generateResetToken()
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		// Сохраняем токен в БД
		if err := s.store.User().SaveVerificationToken(u.ID, verificationToken, time.Now().Add(verificationTokenTTL)); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		// Отправляем письмо с верификацией
		if s.mailer != nil {
			verificationLink := s.buildVerificationLink(verificationToken)
			if err := s.mailer.SendVerificationEmail(u.Email, verificationLink); err != nil {
				s.logger.WithError(err).Error("Failed to send verification email")
				// Не возвращаем ошибку, чтобы пользователь мог зарегистрироваться
				// Но логируем проблему
			} else {
				s.logger.WithField("email", u.Email).Info("Verification email sent")
			}
		} else {
			s.logger.Warn("Mailer not configured - verification email not sent. Please configure SMTP settings.")
		}

		token, err := GenerateToken(u.ID, TokenScopeUser, s.jwtSecret, userTokenTTL)
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		s.respond(c, map[string]interface{}{
			"token":          token,
			"email_verified": false,
			"message":        "Регистрация успешна. Пожалуйста, проверьте вашу почту для подтверждения email адреса.",
		}, http.StatusCreated)

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

		token, err := GenerateToken(u.ID, TokenScopeUser, s.jwtSecret, userTokenTTL)
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		s.respond(c, map[string]interface{}{
			"token":          token,
			"email_verified": u.EmailVerified,
		}, http.StatusOK)
	}
}

func (s *server) handleAdminLogin() gin.HandlerFunc {
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

		admin, err := s.store.Admin().FindByEmail(req.Email)
		if err != nil || !admin.ComparePassword(req.Password) {
			s.error(c, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		token, err := GenerateToken(admin.ID, TokenScopeAdmin, s.jwtSecret, adminTokenTTL)
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		s.respond(c, map[string]string{"token": token}, http.StatusOK)
	}
}

func (s *server) handleAdminForgotPassword() gin.HandlerFunc {
	type request struct {
		Email string `json:"email"`
	}

	return func(c *gin.Context) {
		req := &request{}
		if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		admin, err := s.store.Admin().FindByEmail(req.Email)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.logger.WithField("email", req.Email).Warn("password reset requested for unknown admin")
				s.respond(c, map[string]string{"status": "if the email exists, instructions were sent"}, http.StatusOK)
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		token, err := generateResetToken()
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		if err := s.store.Admin().SaveResetToken(admin.ID, token, time.Now().Add(resetTokenTTL)); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		if s.mailer != nil {
			if err := s.mailer.SendResetPasswordEmail(admin.Email, s.buildResetLink(token)); err != nil {
				s.error(c, http.StatusInternalServerError, err)
				return
			}
		} else {
			s.logger.Warn("Mailer not configured - reset password email not sent")
			s.error(c, http.StatusInternalServerError, errors.New("email service not configured"))
			return
		}

		s.respond(c, map[string]string{"status": "if the email exists, instructions were sent"}, http.StatusOK)
	}
}

func (s *server) handleAdminResetPassword() gin.HandlerFunc {
	type request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	return func(c *gin.Context) {
		req := &request{}
		if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		if req.Token == "" || req.NewPassword == "" {
			s.error(c, http.StatusBadRequest, errors.New("token and new_password are required"))
			return
		}

		admin, err := s.store.Admin().FindByResetToken(req.Token)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusBadRequest, errInvalidResetToken)
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		if admin.ResetTokenExpires == nil || time.Now().After(*admin.ResetTokenExpires) {
			s.error(c, http.StatusBadRequest, errExpiredResetToken)
			return
		}

		if err := admin.SetPassword(req.NewPassword); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Admin().UpdatePassword(admin.ID, admin.PasswordHash); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		if err := s.store.Admin().ClearResetToken(admin.ID); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		s.respond(c, map[string]string{"status": "password updated"}, http.StatusOK)
	}
}

func (s *server) handleAdminMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID, exists := c.Get("adminID")
		if !exists {
			s.error(c, http.StatusUnauthorized, errors.New("admin not authenticated"))
			return
		}

		admin, err := s.store.Admin().FindByID(adminID.(int))
		if err != nil {
			s.error(c, http.StatusNotFound, err)
			return
		}

		admin.Sanitize()
		s.respond(c, admin, http.StatusOK)
	}
}

func (s *server) handleAdminDashboard() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats := map[string]interface{}{
			"message":      "Добро пожаловать в админскую панель",
			"generated_at": time.Now(),
		}
		blocks, _ := s.store.Block().List(store.BlockFilter{})
		posts, _ := s.store.Post().List(store.PostFilter{IncludeDrafts: true})
		users, _ := s.store.User().List()
		stats["blocks_count"] = len(blocks)
		stats["posts_count"] = len(posts)
		stats["users_count"] = len(users)

		s.respond(c, stats, http.StatusOK)
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
		claims, ok := s.extractClaimsFromHeader(c)
		if !ok {
			return
		}

		if claims.Scope != TokenScopeUser {
			s.error(c, http.StatusUnauthorized, errors.New("invalid token scope"))
			c.Abort()
			return
		}

		c.Set("userID", claims.SubjectID)
		c.Next()
	}
}

func (s *server) adminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.logger.WithFields(logrus.Fields{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		}).Info("Admin auth middleware called")
		
		claims, ok := s.extractClaimsFromHeader(c)
		if !ok {
			s.logger.Warn("Failed to extract claims from header")
			return
		}

		if claims.Scope != TokenScopeAdmin {
			s.logger.Warn("Invalid token scope")
			s.error(c, http.StatusUnauthorized, errors.New("invalid token scope"))
			c.Abort()
			return
		}

		s.logger.Info("Admin authenticated successfully")
		c.Set("adminID", claims.SubjectID)
		c.Next()
	}
}

func (s *server) extractClaimsFromHeader(c *gin.Context) (*Claims, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		s.error(c, http.StatusUnauthorized, errors.New("authorization header required"))
		c.Abort()
		return nil, false
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || !strings.HasPrefix(strings.ToLower(authHeader), strings.ToLower(bearerPrefix)) {
		s.error(c, http.StatusUnauthorized, errors.New("invalid authorization header format"))
		c.Abort()
		return nil, false
	}

	tokenString := strings.TrimSpace(authHeader[len(bearerPrefix):])
	if tokenString == "" {
		s.error(c, http.StatusUnauthorized, errors.New("token is required"))
		c.Abort()
		return nil, false
	}

	claims, err := ValidateToken(tokenString, s.jwtSecret)
	if err != nil {
		s.error(c, http.StatusUnauthorized, errors.New("invalid token"))
		c.Abort()
		return nil, false
	}

	return claims, true
}

func generateResetToken() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

func (s *server) buildResetLink(token string) string {
	if strings.Contains(s.passwordResetURL, "%s") {
		return fmt.Sprintf(s.passwordResetURL, token)
	}
	return fmt.Sprintf("%s%s", s.passwordResetURL, token)
}

func (s *server) buildVerificationLink(token string) string {
	if strings.Contains(s.verificationURL, "%s") {
		return fmt.Sprintf(s.verificationURL, token)
	}
	return fmt.Sprintf("%s%s", s.verificationURL, token)
}

func (s *server) handleVerifyEmail() gin.HandlerFunc {
	type request struct {
		Token string `json:"token"`
	}

	return func(c *gin.Context) {
		req := &request{}
		if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		if req.Token == "" {
			s.error(c, http.StatusBadRequest, errors.New("token is required"))
			return
		}

		user, err := s.store.User().FindByVerificationToken(req.Token)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusBadRequest, errInvalidVerificationToken)
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		if user.EmailVerified {
			s.error(c, http.StatusBadRequest, errEmailAlreadyVerified)
			return
		}

		if user.EmailVerificationTokenExpires == nil || time.Now().After(*user.EmailVerificationTokenExpires) {
			s.error(c, http.StatusBadRequest, errExpiredVerificationToken)
			return
		}

		if err := s.store.User().VerifyEmail(user.ID); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		s.respond(c, map[string]string{"status": "email verified"}, http.StatusOK)
	}
}

func (s *server) handleResendVerification() gin.HandlerFunc {
	type request struct {
		Email string `json:"email"`
	}

	return func(c *gin.Context) {
		req := &request{}
		if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		user, err := s.store.User().FindByEmail(req.Email)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				// Не раскрываем, что пользователь не существует
				s.respond(c, map[string]string{"status": "if the email exists, verification email was sent"}, http.StatusOK)
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		if user.EmailVerified {
			s.error(c, http.StatusBadRequest, errEmailAlreadyVerified)
			return
		}

		// Генерируем новый токен верификации
		verificationToken, err := generateResetToken()
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		// Сохраняем токен в БД
		if err := s.store.User().SaveVerificationToken(user.ID, verificationToken, time.Now().Add(verificationTokenTTL)); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		// Отправляем письмо с верификацией
		if s.mailer != nil {
			verificationLink := s.buildVerificationLink(verificationToken)
			if err := s.mailer.SendVerificationEmail(user.Email, verificationLink); err != nil {
				s.error(c, http.StatusInternalServerError, err)
				return
			}
		} else {
			s.logger.Warn("Mailer not configured - verification email not sent")
			s.error(c, http.StatusInternalServerError, errors.New("email service not configured"))
			return
		}

		s.respond(c, map[string]string{"status": "if the email exists, verification email was sent"}, http.StatusOK)
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

func (s *server) handlePublicPageBlocks() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.Param("page")
		if strings.HasPrefix(page, "/") {
			page = page[1:]
		}
		if page == "" || page == "home" {
			page = "/"
		} else {
			page = "/" + page
		}
		blocks, err := s.store.Block().List(store.BlockFilter{Page: page})
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, blocks, http.StatusOK)
	}
}

func (s *server) handlePublicPostsList() gin.HandlerFunc {
	return func(c *gin.Context) {
		posts, err := s.store.Post().List(store.PostFilter{IncludeDrafts: false})
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, posts, http.StatusOK)
	}
}

func (s *server) handlePublicPostDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		post, err := s.store.Post().FindByID(id)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusNotFound, err)
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		if !post.IsPublished {
			s.error(c, http.StatusNotFound, errors.New("post not found"))
			return
		}
		s.respond(c, post, http.StatusOK)
	}
}

func (s *server) handlePublicPostComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("id")
		comments, err := s.store.Comment().ListByPost(postID)
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, comments, http.StatusOK)
	}
}

func (s *server) handleUserCommentCreate() gin.HandlerFunc {
	type request struct {
		Text string `json:"text" binding:"required"`
	}
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		postID := c.Param("id")
		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		if _, err := s.store.Post().FindByID(postID); err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusNotFound, errors.New("post not found"))
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		comment := &model.Comment{
			PostID:  postID,
			UserID:  userID,
			Text:    req.Text,
			IsAdmin: false,
		}
		if err := s.store.Comment().Create(comment); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, comment, http.StatusCreated)
	}
}

func (s *server) handleUserCommentReply() gin.HandlerFunc {
	type request struct {
		Text string `json:"text" binding:"required"`
	}
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		parentID := c.Param("id")
		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		parent, err := s.store.Comment().FindByID(parentID)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusNotFound, errors.New("comment not found"))
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		comment := &model.Comment{
			PostID:  parent.PostID,
			UserID:  userID,
			Text:    req.Text,
			IsAdmin: false,
			ReplyTo: &parent.ID,
		}
		if err := s.store.Comment().Create(comment); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, comment, http.StatusCreated)
	}
}

func (s *server) handleUserCommentDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		commentID := c.Param("id")
		userID := c.GetInt("userID")

		comment, err := s.store.Comment().FindByID(commentID)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusNotFound, errors.New("comment not found"))
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		if comment.UserID != userID {
			s.error(c, http.StatusForbidden, errors.New("cannot delete foreign comment"))
			return
		}

		if err := s.store.Comment().Delete(commentID); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, map[string]string{"status": "deleted"}, http.StatusOK)
	}
}

func (s *server) handleAdminBlocksList() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.Query("page")
		blocks, err := s.store.Block().List(store.BlockFilter{Page: page})
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, blocks, http.StatusOK)
	}
}

func (s *server) handleAdminBlockCreate() gin.HandlerFunc {
	type request struct {
		Page         string          `json:"page" binding:"required"`
		Pages        []string        `json:"pages"`
		Type         model.BlockType `json:"type" binding:"required"`
		Content      json.RawMessage `json:"content"`
		Order        int             `json:"order"`
		DisplayOrder int             `json:"display_order"`
	}
	return func(c *gin.Context) {
		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		pages, err := json.Marshal(req.Pages)
		if err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		order := req.DisplayOrder
		if order == 0 && req.Order != 0 {
			order = req.Order
		}

		block := &model.Block{
			Page:    req.Page,
			Pages:   pages,
			Type:    req.Type,
			Content: req.Content,
			Order:   order,
		}
		if err := s.store.Block().Create(block); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, block, http.StatusCreated)
	}
}

func (s *server) handleAdminBlockUpdate() gin.HandlerFunc {
	type request struct {
		Page         string          `json:"page"`
		Pages        []string        `json:"pages"`
		Type         model.BlockType `json:"type"`
		Content      json.RawMessage `json:"content"`
		Order        *int            `json:"order"`
		DisplayOrder *int            `json:"display_order"`
	}
	return func(c *gin.Context) {
		id := c.Param("id")
		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		existing, err := s.store.Block().FindByID(id)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusNotFound, errors.New("block not found"))
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		if req.Page != "" {
			existing.Page = req.Page
		}
		if req.Pages != nil {
			if pages, err := json.Marshal(req.Pages); err == nil {
				existing.Pages = pages
			}
		}
		if req.Type != "" {
			existing.Type = req.Type
		}
		if req.Content != nil {
			existing.Content = req.Content
		}
		if req.DisplayOrder != nil {
			existing.Order = *req.DisplayOrder
		} else if req.Order != nil {
			existing.Order = *req.Order
		}

		if err := s.store.Block().Update(existing); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, existing, http.StatusOK)
	}
}

func (s *server) handleAdminBlockDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := s.store.Block().Delete(id); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, map[string]string{"status": "deleted"}, http.StatusOK)
	}
}

func (s *server) handleAdminBlockReorder() gin.HandlerFunc {
	type request struct {
		Items []store.BlockOrderUpdate `json:"items" binding:"required,dive"`
	}
	return func(c *gin.Context) {
		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}
		if err := s.store.Block().UpdateOrders(req.Items); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, map[string]string{"status": "reordered"}, http.StatusOK)
	}
}

func (s *server) handleAdminPostsList() gin.HandlerFunc {
	return func(c *gin.Context) {
		includeDrafts := c.Query("includeDrafts") == "true"
		posts, err := s.store.Post().List(store.PostFilter{IncludeDrafts: includeDrafts})
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, posts, http.StatusOK)
	}
}

func (s *server) handleAdminPostDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		post, err := s.store.Post().FindByID(id)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusNotFound, err)
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, post, http.StatusOK)
	}
}

func (s *server) handleAdminPostCreate() gin.HandlerFunc {
	type request struct {
		Title          string   `json:"title" binding:"required"`
		Content        string   `json:"content" binding:"required"`
		Images         []string `json:"images"`
		Videos         []string `json:"videos"`
		Pages          []string `json:"pages"`
		IsPublished    bool     `json:"is_published"`
		Alignment      string   `json:"alignment"`
		TitlePosition  string   `json:"title_position"`
		ContentPosition string  `json:"content_position"`
	}
	return func(c *gin.Context) {
		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		images, err := json.Marshal(req.Images)
		if err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		videos, err := json.Marshal(req.Videos)
		if err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		pages, err := json.Marshal(req.Pages)
		if err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		post := &model.Post{
			Title:          req.Title,
			Content:        req.Content,
			Images:         images,
			Videos:         videos,
			Pages:          pages,
			IsPublished:    req.IsPublished,
			Alignment:      req.Alignment,
			TitlePosition:  req.TitlePosition,
			ContentPosition: req.ContentPosition,
		}
		if err := s.store.Post().Create(post); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, post, http.StatusCreated)
	}
}

func (s *server) handleAdminPostUpdate() gin.HandlerFunc {
	type request struct {
		Title          *string  `json:"title"`
		Content        *string  `json:"content"`
		Images         []string `json:"images"`
		Videos         []string `json:"videos"`
		Pages          []string `json:"pages"`
		IsPublished    *bool    `json:"is_published"`
		Alignment      *string  `json:"alignment"`
		TitlePosition  *string  `json:"title_position"`
		ContentPosition *string `json:"content_position"`
	}
	return func(c *gin.Context) {
		id := c.Param("id")
		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		post, err := s.store.Post().FindByID(id)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusNotFound, err)
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		if req.Title != nil {
			post.Title = *req.Title
		}
		if req.Content != nil {
			post.Content = *req.Content
		}
		if req.Images != nil {
			if images, err := json.Marshal(req.Images); err == nil {
				post.Images = images
			}
		}
		if req.Videos != nil {
			if videos, err := json.Marshal(req.Videos); err == nil {
				post.Videos = videos
			}
		}
		if req.Pages != nil {
			if pages, err := json.Marshal(req.Pages); err == nil {
				post.Pages = pages
			}
		}
		if req.IsPublished != nil {
			post.IsPublished = *req.IsPublished
		}
		if req.Alignment != nil {
			post.Alignment = *req.Alignment
		}
		if req.TitlePosition != nil {
			post.TitlePosition = *req.TitlePosition
		}
		if req.ContentPosition != nil {
			post.ContentPosition = *req.ContentPosition
		}

		if err := s.store.Post().Update(post); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, post, http.StatusOK)
	}
}

func (s *server) handleAdminPostDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := s.store.Post().Delete(id); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, map[string]string{"status": "deleted"}, http.StatusOK)
	}
}

func (s *server) handleAdminGalleryList() gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := s.store.Gallery().List()
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, items, http.StatusOK)
	}
}

func (s *server) handleAdminGalleryCreate() gin.HandlerFunc {
	type request struct {
		ImageURL    string   `json:"image_url" binding:"required"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Text        string   `json:"text"`
		Pages       []string `json:"pages"`
	}
	return func(c *gin.Context) {
		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		pages, err := json.Marshal(req.Pages)
		if err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		item := &model.GalleryItem{
			ImageURL:    req.ImageURL,
			Title:       req.Title,
			Description: req.Description,
			Text:        req.Text,
			Pages:       pages,
		}
		if err := s.store.Gallery().Create(item); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, item, http.StatusCreated)
	}
}

func (s *server) handleAdminGalleryUpdate() gin.HandlerFunc {
	type request struct {
		ImageURL    *string  `json:"image_url"`
		Title       *string  `json:"title"`
		Description *string  `json:"description"`
		Text        *string  `json:"text"`
		Pages       []string `json:"pages"`
	}
	return func(c *gin.Context) {
		id := c.Param("id")
		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		item, err := s.store.Gallery().FindByID(id)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusNotFound, err)
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		if req.ImageURL != nil {
			item.ImageURL = *req.ImageURL
		}
		if req.Title != nil {
			item.Title = *req.Title
		}
		if req.Description != nil {
			item.Description = *req.Description
		}
		if req.Text != nil {
			item.Text = *req.Text
		}
		if req.Pages != nil {
			if pages, err := json.Marshal(req.Pages); err == nil {
				item.Pages = pages
			}
		}

		if err := s.store.Gallery().Update(item); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, item, http.StatusOK)
	}
}

func (s *server) handleAdminGalleryDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := s.store.Gallery().Delete(id); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, map[string]string{"status": "deleted"}, http.StatusOK)
	}
}

func (s *server) handleAdminSliderList() gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := s.store.Slider().List()
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, items, http.StatusOK)
	}
}

func (s *server) handleAdminSliderCreate() gin.HandlerFunc {
	type request struct {
		ImageURL    string   `json:"image_url" binding:"required"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Pages       []string `json:"pages"`
		Order       int      `json:"order"`
	}
	return func(c *gin.Context) {
		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		pages, err := json.Marshal(req.Pages)
		if err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		item := &model.SliderItem{
			ImageURL:    req.ImageURL,
			Title:       req.Title,
			Description: req.Description,
			Pages:       pages,
			Order:       req.Order,
		}
		if err := s.store.Slider().Create(item); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, item, http.StatusCreated)
	}
}

func (s *server) handleAdminSliderUpdate() gin.HandlerFunc {
	type request struct {
		ImageURL    *string  `json:"image_url"`
		Title       *string  `json:"title"`
		Description *string  `json:"description"`
		Pages       []string `json:"pages"`
		Order       *int     `json:"order"`
	}
	return func(c *gin.Context) {
		id := c.Param("id")
		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		item, err := s.store.Slider().FindByID(id)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusNotFound, err)
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		if req.ImageURL != nil {
			item.ImageURL = *req.ImageURL
		}
		if req.Title != nil {
			item.Title = *req.Title
		}
		if req.Description != nil {
			item.Description = *req.Description
		}
		if req.Pages != nil {
			if pages, err := json.Marshal(req.Pages); err == nil {
				item.Pages = pages
			}
		}
		if req.Order != nil {
			item.Order = *req.Order
		}

		if err := s.store.Slider().Update(item); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, item, http.StatusOK)
	}
}

func (s *server) handleAdminSliderDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := s.store.Slider().Delete(id); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, map[string]string{"status": "deleted"}, http.StatusOK)
	}
}

func (s *server) handleAdminSliderReorder() gin.HandlerFunc {
	type request struct {
		Items []store.SliderOrderUpdate `json:"items" binding:"required,dive"`
	}
	return func(c *gin.Context) {
		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Slider().UpdateOrders(req.Items); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, map[string]string{"status": "reordered"}, http.StatusOK)
	}
}

func (s *server) handleAdminCommentsList() gin.HandlerFunc {
	return func(c *gin.Context) {
		comments, err := s.store.Comment().ListAll()
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, comments, http.StatusOK)
	}
}

func (s *server) handleAdminCommentReply() gin.HandlerFunc {
	type request struct {
		Text string `json:"text" binding:"required"`
	}
	return func(c *gin.Context) {
		adminID := c.GetInt("adminID")
		parentID := c.Param("id")
		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		parent, err := s.store.Comment().FindByID(parentID)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusNotFound, errors.New("comment not found"))
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		adminUserID, err := s.ensureAdminUser(adminID)
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		comment := &model.Comment{
			PostID:  parent.PostID,
			UserID:  adminUserID,
			Text:    req.Text,
			IsAdmin: true,
			ReplyTo: &parent.ID,
		}
		if err := s.store.Comment().Create(comment); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, comment, http.StatusCreated)
	}
}

func (s *server) handleAdminCommentDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := s.store.Comment().Delete(id); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		s.respond(c, map[string]string{"status": "deleted"}, http.StatusOK)
	}
}

func (s *server) handleAdminUsersList() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := s.store.User().List()
		if err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}
		for _, u := range users {
			u.Sanitize()
		}
		s.respond(c, users, http.StatusOK)
	}
}

func (s *server) handleAdminUserGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			s.error(c, http.StatusBadRequest, errors.New("invalid user ID"))
			return
		}

		user, err := s.store.User().FindByID(userID)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusNotFound, errors.New("user not found"))
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		commentsCount, err := s.store.User().GetCommentsCount(userID)
		if err != nil {
			commentsCount = 0
		}

		user.Sanitize()
		response := map[string]interface{}{
			"id":            user.ID,
			"email":         user.Email,
			"is_admin":      user.IsAdmin,
			"banned":        user.Banned,
			"created_at":    user.CreatedAt,
			"updated_at":    user.UpdatedAt,
			"comments_count": commentsCount,
		}

		s.respond(c, response, http.StatusOK)
	}
}

func (s *server) handleAdminUserBan() gin.HandlerFunc {
	type request struct {
		Banned bool `json:"banned"`
	}
	return func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			s.error(c, http.StatusBadRequest, errors.New("invalid user ID"))
			return
		}

		req := request{}
		if err := c.ShouldBindJSON(&req); err != nil {
			s.error(c, http.StatusBadRequest, err)
			return
		}

		if err := s.store.User().UpdateBanned(userID, req.Banned); err != nil {
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		s.respond(c, map[string]interface{}{
			"status": "updated",
			"banned": req.Banned,
		}, http.StatusOK)
	}
}

func (s *server) handlePublicUserGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			s.error(c, http.StatusBadRequest, errors.New("invalid user ID"))
			return
		}

		user, err := s.store.User().FindByID(userID)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(c, http.StatusNotFound, errors.New("user not found"))
				return
			}
			s.error(c, http.StatusInternalServerError, err)
			return
		}

		commentsCount, err := s.store.User().GetCommentsCount(userID)
		if err != nil {
			commentsCount = 0
		}

		user.Sanitize()
		response := map[string]interface{}{
			"id":             user.ID,
			"email":          user.Email,
			"is_admin":       user.IsAdmin,
			"created_at":     user.CreatedAt,
			"updated_at":     user.UpdatedAt,
			"comments_count": commentsCount,
		}
		// Не показываем статус banned для обычных пользователей

		s.respond(c, response, http.StatusOK)
	}
}

func (s *server) ensureAdminUser(adminID int) (int, error) {
	admin, err := s.store.Admin().FindByID(adminID)
	if err != nil {
		return 0, err
	}

	user, err := s.store.User().FindByEmail(admin.Email)
	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			randomPassword, genErr := generateResetToken()
			if genErr != nil {
				return 0, genErr
			}
			newUser := &model.User{
				Email:    admin.Email,
				Password: randomPassword,
			}
			if err := s.store.User().Create(newUser); err != nil {
				return 0, err
			}
			return newUser.ID, nil
		}
		return 0, err
	}

	return user.ID, nil
}

func (s *server) handleAdminUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.logger.Info("handleAdminUpload called")
		file, err := c.FormFile("file")
		if err != nil {
			s.error(c, http.StatusBadRequest, fmt.Errorf("file is required: %w", err))
			return
		}

		contentType := file.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") && !strings.HasPrefix(contentType, "video/") {
			s.error(c, http.StatusBadRequest, errors.New("only image and video files are allowed"))
			return
		}

		if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
			s.error(c, http.StatusInternalServerError, fmt.Errorf("failed to create upload directory: %w", err))
			return
		}
		ext := filepath.Ext(file.Filename)
		randomBytes := make([]byte, 16)
		if _, err := rand.Read(randomBytes); err != nil {
			s.error(c, http.StatusInternalServerError, fmt.Errorf("failed to generate filename: %w", err))
			return
		}
		filename := hex.EncodeToString(randomBytes) + ext
		fullPath := filepath.Join(s.uploadDir, filename)

		src, err := file.Open()
		if err != nil {
			s.error(c, http.StatusInternalServerError, fmt.Errorf("failed to open file: %w", err))
			return
		}
		defer src.Close()

		dst, err := os.Create(fullPath)
		if err != nil {
			s.error(c, http.StatusInternalServerError, fmt.Errorf("failed to create file: %w", err))
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			s.error(c, http.StatusInternalServerError, fmt.Errorf("failed to save file: %w", err))
			return
		}

		url := fmt.Sprintf("/uploads/%s", filename)
		s.respond(c, map[string]string{"url": url}, http.StatusOK)
	}
}
