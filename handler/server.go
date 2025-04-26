package handler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sidharth-125/token_heap/consts"
	"github.com/sidharth-125/token_heap/entity"
)

type Handler struct {
	pool       *entity.TokenPool
	expiry     int64
	tokenMap   map[string]*entity.Token
	echoLogger echo.Logger
	mutex      *sync.RWMutex
}

func NewHandler(pool *entity.TokenPool, expiry int64, tokenMap map[string]*entity.Token, logger echo.Logger, m *sync.RWMutex) Handler {
	return Handler{pool: pool, expiry: expiry, tokenMap: tokenMap, echoLogger: logger, mutex: m}
}

var reg = regexp.MustCompile(consts.EmailRegex)

func SetupServer() {
	tokenPool := SetupTokenPool()
	tokenMap := make(map[string]*entity.Token, 0)
	loggerObj := echo.New().Logger
	loggerObj.SetLevel(log.DEBUG)
	mutex := new(sync.RWMutex)
	handlerObj := NewHandler(&tokenPool, consts.Expiry, tokenMap, loggerObj, mutex)

	echoObj := echo.New()
	echoObj.POST("/login", handlerObj.UserLogin)
	echoObj.GET("/token", handlerObj.GetQueuedTokens)
	echoObj.GET("/logout", handlerObj.Logout)

	go func() {
		err := echoObj.Start(consts.Port)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to start server %v", err))
		}
	}()

	startTimer := time.NewTicker(time.Second * 5)
	defer startTimer.Stop()

	go func() {
		for range startTimer.C {
			TokenManager(handlerObj.pool, handlerObj.mutex, handlerObj.tokenMap)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	echoObj.Shutdown(context.TODO())

}

func (h *Handler) UserLogin(c echo.Context) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var (
		user     entity.User
		response entity.Token
	)

	err := c.Bind(&user)
	if err != nil {
		h.echoLogger.Error("decoder error", err)
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	err = validateUser(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	response.Expiry = time.Now().UTC().Add(time.Minute * time.Duration(h.expiry))
	response.Id = uuid.New().String()
	response.UserEmail = user.Email

	h.tokenMap[response.Id] = &response
	h.pool.Push(response)
	c.Response().Header().Set(consts.HeaderKey, response.Id)
	h.echoLogger.Info("user logged in successfully")
	return c.JSON(http.StatusOK, "user logged in successfully")
}

func validateUser(u entity.User) error {
	resp := reg.MatchString(u.Email)
	if !resp {
		return errors.New("invalid email")
	}

	if u.Password == "" {
		return errors.New("password cannot be empty")
	}

	return nil
}

func (h *Handler) GetQueuedTokens(c echo.Context) error {
	h.mutex.RLock()
	defer h.mutex.Unlock()

	if len(h.tokenMap) == 0 {
		return c.JSON(http.StatusOK, "no queued tokens")
	}

	var response []entity.Token
	for _, data := range h.tokenMap {
		response = append(response, *data)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) Logout(c echo.Context) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	key := c.Request().Header.Get(consts.HeaderKey)
	if key == "" {
		return c.JSON(http.StatusBadRequest, "invalid token")
	}

	_, ok := h.tokenMap[key]
	if !ok {
		return c.JSON(http.StatusBadRequest, "invalid token")
	}

	delete(h.tokenMap, key)
	h.echoLogger.Info("user deleted successfully")
	return c.JSON(http.StatusOK, "successfully logged out")
}
