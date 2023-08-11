package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGoogleAuth(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/google", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := HandleGoogleAuth

	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusSeeOther, rec.Code)
}
