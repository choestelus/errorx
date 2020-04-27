package ginerrorx

import (
	"errors"
	"net/http"

	"github.com/choestelus/errorx"
	"github.com/gin-gonic/gin"
)

// ErrorExtractor unwinds and returns wrapped errors as stack trace
func ErrorExtractor() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		if !c.Writer.Written() && c.Writer.Status() == 200 {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}

		c.Abort()
		c.PureJSON(
			c.Writer.Status(),
			gin.H{"errors": unwind(err.Err)},
		)
	}
}

func unwind(e error) (errorMessages []map[string]interface{}) {
	if e == nil {
		return nil
	}
	switch err := e.(type) {
	case errorx.E:
		errMessage := map[string]interface{}{
			"code":    err.Code,
			"message": err.Error(),
		}
		if wrappedE := err.Unwrap(); wrappedE == nil {
			return []map[string]interface{}{errMessage}
		}
		return append(unwind(err.Unwrap()), errMessage)
	case *errorx.E:
		errMessage := map[string]interface{}{
			"code":    err.Code,
			"message": err.Error(),
		}
		if wrappedE := err.Unwrap(); wrappedE == nil {
			return []map[string]interface{}{errMessage}
		}
		return append(unwind(err.Unwrap()), errMessage)
	case error:
		errMessage := map[string]interface{}{
			"code":    "std go error",
			"message": err.Error(),
		}
		return append(unwind(errors.Unwrap(err)), errMessage)
	default:
		return unwind(errors.Unwrap(e))
	}
}
