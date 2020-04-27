package ginerrorx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/choestelus/errorx"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestUnwind(t *testing.T) {
	var e error = errorx.Wrap(fmt.Errorf("init"), "1st wrap message")
	e = errorx.Wrap(e, "2nd wrap message")
	e = errorx.Wrap(e, "3rd wrap message")
	e = fmt.Errorf("4th wrap message: %w", e)
	e = errorx.Wrap(e, "5th wrap message")

	err, ok := e.(*errorx.E)
	require.True(t, ok)

	errorMessages := unwind(err)

	require.Len(t, errorMessages, 6)
}

func TestErrorExtractorMiddleware(t *testing.T) {
	var e error = errorx.Wrap(fmt.Errorf("init"), "1st wrap message")
	e = errorx.Wrap(e, "2nd wrap message")
	e = errorx.Wrap(e, "3rd wrap message")
	e = fmt.Errorf("4th wrap message: %w", e)
	e = errorx.Wrap(e, "5th wrap message")

	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()

	c, r := gin.CreateTestContext(rec)
	r.Use(ErrorExtractor())
	r.GET("/test-stacktrace", func(c *gin.Context) {
		c.Error(e)
	})

	req, reqErr := http.NewRequest(http.MethodGet, "/test-stacktrace", nil)
	require.NoError(t, reqErr)
	c.Request = req

	r.ServeHTTP(rec, c.Request)
	body := rec.Body.String()
	require.NotEmpty(t, body)

	encodedJSONErr, err := json.Marshal(gin.H{"errors": unwind(e)})
	require.NoError(t, err)
	require.JSONEq(t, string(encodedJSONErr), body)

	require.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestErrorExtractorMiddlewareHTTPStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()

	c, r := gin.CreateTestContext(rec)
	r.Use(ErrorExtractor())
	r.GET("/test-stacktrace", func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
		c.Error(errors.New("something went wrong"))
	})

	req, reqErr := http.NewRequest(http.MethodGet, "/test-stacktrace", nil)
	require.NoError(t, reqErr)
	c.Request = req

	r.ServeHTTP(rec, c.Request)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
