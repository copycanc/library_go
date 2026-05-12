package server

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"

	"log"
	"net/http"
	domainErrors "rest_library/internal/domain/errors"
	"rest_library/internal/domain/models"
)

type Storage interface {
	GetBooksList() ([]models.Book, error)
	SaveBook(book models.Book) error
	RegisterUser(user models.User) error
	GetUser(user models.UserLogin) (models.User, error)
	PrintInfoBook(string) (models.Book, error)
}

type LibraryAPI struct {
	httpSrv *http.Server
	stor    Storage
}

func NewLibraryAPI(stor Storage, addr string, port int) *LibraryAPI {
	httpSrv := http.Server{
		Addr: addr + ":" + strconv.Itoa(port),
	}

	lAPI := LibraryAPI{
		httpSrv: &httpSrv,
		stor:    stor,
	}

	lAPI.configRouters()
	return &lAPI
}

func (lAPI *LibraryAPI) Run() error {
	log.Printf("Library service started on %s ...", lAPI.httpSrv.Addr)
	err := lAPI.httpSrv.ListenAndServe()
	return err
}

func (lAPI *LibraryAPI) configRouters() {
	router := gin.Default()
	user := router.Group("/users")
	{
		user.POST("/register", lAPI.register) // Регистрация пользователя
		user.POST("/login", lAPI.login)       // Авторизация пользователя
		user.PUT("/update/:userID")           // Обновление информации о пользователе
		user.DELETE("/delete/:userID")        // Удаление пользователя
		user.GET("/:userID")                  // Возвращает пользователя по его ID
	}
	books := router.Group("/books")
	{
		books.POST("/create", lAPI.newBook)  // Добавление книги
		books.GET("/list", lAPI.booksList)   // Получить список всех книг
		books.GET("/:bookID", lAPI.bookInfo) // Получение информации о книге
		books.PUT("/update/:bookID")         // Обновление информации о книге
		books.DELETE("/delete/:bookID")      // Удаление книги
	}

	lAPI.httpSrv.Handler = router
}

func (lAPI *LibraryAPI) booksList(ctx *gin.Context) {
	books, err := lAPI.stor.GetBooksList()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (lAPI *LibraryAPI) newBook(ctx *gin.Context) {
	var book models.Book
	if err := ctx.ShouldBindBodyWithJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	lAPI.stor.SaveBook(book)
	ctx.JSON(http.StatusCreated, book)
}

func (lAPI *LibraryAPI) register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	valid := validator.New()
	if err := valid.Struct(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.Password = string(hash)
	if err := lAPI.stor.RegisterUser(user); err != nil {
		if errors.Is(err, domainErrors.ErrUserExist) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.String(http.StatusCreated, "Пользователь зарегистрирован")
}

func (lAPI *LibraryAPI) login(ctx *gin.Context) {
	var user models.UserLogin
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	valid := validator.New()
	if err := valid.Struct(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	dbUser, err := lAPI.stor.GetUser(user)
	if err != nil {
		if errors.Is(err, domainErrors.ErrInvalidCreds) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, dbUser)
}

func (lAPI *LibraryAPI) bookInfo(ctx *gin.Context) {
	bookID := ctx.Param("bookID")
	book, err := lAPI.stor.PrintInfoBook(bookID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrBookNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, book)
}
