package inceptor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error,omitempty"`
}

type HealthResponse struct {
	BaseResponse
}

type RandResponse struct {
	BaseResponse
	Number uint64 `json:"number"`
}

func success() BaseResponse {
	return BaseResponse{true, nil}
}

func failure(e error) BaseResponse {
	return BaseResponse{false, e}
}

func Server() *gin.Engine {
	r := gin.Default()
	r.GET("/rand", func(c *gin.Context) {
		n, err := Uint64()
		if err == nil {
			c.JSON(http.StatusOK, RandResponse{success(), n})
		} else {
			//TODO: don't pass this error through?
			c.JSON(http.StatusOK, failure(err))
		}
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthResponse{success()})
	})
	return r
}
