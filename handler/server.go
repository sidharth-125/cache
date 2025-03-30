package handler

import (
	"errors"
	"fmt"

	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sidharth-125/token_heap/entity"
)

type Handler struct {
	pool       *entity.TokenPool
	expiry     int
	tokenMap   map[string]*entity.Token
	echoLogger echo.Logger
}

func NewHandler(pool *entity.TokenPool, expiry int, tokenMap map[string]*entity.Token, logger echo.Logger) Handler {
	return Handler{pool: pool, expiry: expiry, tokenMap: tokenMap, echoLogger: logger}
}

func SetupServer() {
	tokenPool := SetupTokenPool()
	tokenMap := make(map[string]*entity.Token, 0)
	loggerObj := echo.New().Logger
	loggerObj.SetLevel(log.DEBUG)
	handlerObj := NewHandler(&tokenPool, expiry, tokenMap, loggerObj)

	echoObj := echo.New()
	echoObj.POST("/login", handlerObj.UserLogin)

	err := echoObj.Start(port)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to start server %v", err))
	}

}

const (
	port   string = ":8080"
	expiry int    = 60
)

func (h Handler) UserLogin(c echo.Context) error {
	var (
		user     entity.User
		response entity.Token
	)

	err := c.Bind(&user)
	if err != nil {
		h.echoLogger.Error("decoder error", err)
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	err = ValidateUser(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	response.Expiry = time.Now().Unix() + int64(expiry)
	response.Id = uuid.New().String()
	response.UserEmail = user.Email

	h.pool.Push(response)
	h.echoLogger.Info("user logged in successfully", response)
	return c.JSON(http.StatusOK, response)
}

func ValidateUser(u entity.User) error {
	reg := regexp.MustCompile(`^\w+[-.+]?\w+@[\w]+([-.]\w+)*\.[a-zA-Z]{2,}$`)
	resp := reg.MatchString(u.Email)
	if !resp {
		return errors.New("invalid email")
	}

	if u.Password == "" {
		return errors.New("password cannot be empty")
	}

	return nil
}

// func (h Handler) GetQueuedTokens(c echo.Context) error {
// 	if h.pool.Len() == 0 {
// 		return c.JSON(http.StatusOK, "no items present")
// 	}

// }
