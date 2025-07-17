package api

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nick0323/K8sVision/model"
)

var loginFailMap = make(map[string]struct {
	Count    int
	LastFail time.Time
})

var (
	maxLoginFail = 5
	lockDuration = 10 * time.Minute
)

func init() {
	if v := os.Getenv("LOGIN_MAX_FAIL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxLoginFail = n
		}
	}
	if v := os.Getenv("LOGIN_LOCK_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			lockDuration = time.Duration(n) * time.Minute
		}
	}
}

// LoginHandler 登录接口
// @Summary 用户登录
// @Description 登录获取 JWT Token，连续失败5次10分钟内禁止尝试
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body model.LoginRequest true "登录参数"
// @Success 200 {object} map[string]string "{token: JWT}"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 401 {object} map[string]string "用户名或密码错误"
// @Failure 429 {object} map[string]string "登录失败次数过多"
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request parameters"})
		return
	}
	username := os.Getenv("LOGIN_USERNAME")
	password := os.Getenv("LOGIN_PASSWORD")
	key := req.Username + "|" + c.ClientIP()
	failInfo := loginFailMap[key]
	// 超过锁定时间自动重置失败次数
	if time.Since(failInfo.LastFail) >= lockDuration {
		failInfo.Count = 0
	}
	if failInfo.Count >= maxLoginFail {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "Too many failed login attempts, please try again in 10 minutes"})
		return
	}
	if req.Username == username && req.Password == password {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": req.Username,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString(jwtSecret)
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
		delete(loginFailMap, key)
		return
	}
	// 登录失败
	loginFailMap[key] = struct {
		Count    int
		LastFail time.Time
	}{failInfo.Count + 1, time.Now()}
	c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
}
