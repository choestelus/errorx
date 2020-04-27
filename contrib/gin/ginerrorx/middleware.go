package ginerrorx

import (
	"errors"
	"net/http"

	"github.com/choestelus/errorx"
	"github.com/gin-gonic/gin"
)

func initHTTPStatusCode(code int) int {
	if code == 0 {
		return http.StatusInternalServerError
	}
	return code
}

// ErrorExtractor unwinds and returns wrapped errors as stack trace
func ErrorExtractor() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		switch ae := err.Err.(type) {
		case errorx.E:
			code := initHTTPStatusCode(ae.HTTPStatusCode)
			c.JSON(code, unwind(&ae))
		case *errorx.E:
			code := initHTTPStatusCode(ae.HTTPStatusCode)
			c.JSON(code, unwind(ae))
		default:
			return
		}

	}
}

func unwind(e *errorx.E) []map[string]interface{} {
	errMessages := []map[string]interface{}{}
	for {
		errMessage := map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
		}
		errMessages = append(errMessages, errMessage)

		if wrappedE := errors.Unwrap(e); wrappedE == nil {
			break
		}
	}
	return errMessages
}
