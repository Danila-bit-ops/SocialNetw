package api

import (
	"log"
	"net/http"
	"newsSite/internal/model"
	"newsSite/internal/service"
	"newsSite/pkg"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type api struct {
	srv *service.Service
}

func InitApi(srv *service.Service) *api {
	return &api{
		srv: srv,
	}
}

func (a *api) LoadIndexHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (a *api) InitHandlers(r *gin.Engine) {
	cmdDir, err := filepath.Abs(filepath.Dir("./cmd"))
	if err != nil {
		log.Fatalf("failed to get cmd dir: %v", err)
	}
	assetsDir := filepath.Join(cmdDir, "..", "assets")
	r.LoadHTMLGlob(filepath.Join(assetsDir, "public/*.html"))
	r.GET("/", a.LoadIndexHTML)
	r.Static("/assets/src", filepath.Join(assetsDir, "src"))

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/register", a.RegisterUser)
		apiGroup.POST("/login", a.LoginUser)
		// apiGroup.GET("/news", a.GetAllPosts)

		// Защищённые маршруты
		protected := apiGroup.Group("/")
		protected.Use(pkg.AuthMiddleware())
		{
			protected.GET("/news", a.GetAllPosts)
			protected.POST("/posts", a.AddPost)
			// Добавьте другие защищённые маршруты здесь
		}
	}
}

func (a *api) InitRouter() *gin.Engine {
	router := gin.Default()
	a.InitHandlers(router)
	return router
}

func (a *api) GetAllPosts(c *gin.Context) {
	posts, err := a.srv.GetAllPosts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (a *api) RegisterUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании пароля"})
		return
	}
	user.Password = string(hashedPassword)

	if err := a.srv.RegisterUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно зарегистрирован"})
}

func (a *api) LoginUser(c *gin.Context) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.srv.GetUserByEmail(c.Request.Context(), input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}

	// Сравнение хешированного пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}

	// Генерация JWT
	token, err := pkg.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *api) AddPost(c *gin.Context) {
	var post model.Posts
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEmail := c.GetString("user_email")
	user, err := a.srv.GetUserByEmail(c.Request.Context(), userEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных пользователя"})
		return
	}
	post.AuthorID = user.ID
	post.Author = user.Username
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	if err := a.srv.AddPost(c.Request.Context(), post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Пост успешно добавлен"})
}
