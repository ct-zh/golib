package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const SECRET string = "fb697a356a4363170061f0ae9da5828b"

func main() {
	r := gin.New()

	// 登陆接口， 获取token
	r.POST("login", login)

	// 某个info接口，用于验证token是否正确
	r.POST("info", mwAuthCheck, info)

	r.Run(":12999")
}

// check中间件
func mwAuthCheck(c *gin.Context) {
	bearerLen := len("Bearer ")

	// 通过query的_t或者header的authorization bearer读取token
	token, ok := c.GetQuery("_t")
	if !ok {
		htoken := c.GetHeader("Authorization")
		if len(htoken) < bearerLen {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{
				"msg": "header Authorization has not Bearer token",
			})
			return
		}
		token = strings.TrimSpace(htoken[bearerLen:])
	}

	uid, err := parseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{
			"msg": err.Error(),
		})
	}
	log.Printf("auth check success, uid=%s", uid)

	c.Set("uid", uid)
	c.Next()
}

// 生成token
func generateToken(uid string, ex time.Duration) (string, error) {
	expireTime := time.Now().Add(ex)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        uid,
		//Issuer: "github.com/libragen/felix",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, stdClaims)
	return token.SignedString([]byte(SECRET))
}

// 解析token
func parseToken(tokenStr string) (string, error) {
	claims := jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET), nil
	})
	if err != nil {
		return "", err
	}

	return claims.Id, nil
}

// 登陆接口
func login(ctx *gin.Context) {
	uid := ctx.PostForm("uid")
	log.Printf("uid: %+v ", uid)
	if len(uid) == 0 {
		ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{
			"msg": "uid invalid",
		})
		return
	}

	id, err := strconv.Atoi(uid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{
			"msg": err.Error(),
		})
		return
	}

	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{
			"msg": "uid invalid",
		})
		return
	}

	token, err := generateToken(uid, time.Minute)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{
			"msg": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": token,
	})
}

func info(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": ctx.Value("uid"),
	})
}
