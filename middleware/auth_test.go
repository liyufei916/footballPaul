package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/liyufei916/footballPaul/config"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func makeTestToken(secret string, userID uint, expired bool) string {
	exp := time.Now().Add(time.Hour)
	if expired {
		exp = time.Now().Add(-time.Hour)
	}
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"exp":     exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{Secret: "test-secret"},
	}
	router := gin.New()
	router.Use(AuthMiddleware(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{Secret: "test-secret"},
	}
	router := gin.New()
	router.Use(AuthMiddleware(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{Secret: "test-secret"},
	}
	router := gin.New()
	router.Use(AuthMiddleware(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	secret := "test-secret"
	cfg := &config.Config{
		JWT: config.JWTConfig{Secret: secret},
	}
	router := gin.New()
	router.Use(AuthMiddleware(cfg))

	var capturedUserID uint
	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if exists {
			capturedUserID = userID.(uint)
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	token := makeTestToken(secret, 42, false)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, uint(42), capturedUserID)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	secret := "test-secret"
	cfg := &config.Config{
		JWT: config.JWTConfig{Secret: secret},
	}
	router := gin.New()
	router.Use(AuthMiddleware(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	token := makeTestToken(secret, 42, true)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
