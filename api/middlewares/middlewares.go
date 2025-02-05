package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/auth"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/token"
	"gitea.hoven.com/billiard/billiard-assistant-server/server"
	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/redis/go-redis/v9"
)

const (
	tokenContextPrefix = "billiard:token"
	tokenHeaderKey     = "X-BILLIARD-Token"
)

type UserToken struct {
	UserId int
	Token  *auth.Token
}

func NewUserLoginToken(uid int, accessToken, refreshToken string) *UserToken {
	return &UserToken{
		UserId: uid,
		Token: &auth.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
}

func (t *UserToken) GetKey() string {
	return t.Token.AccessToken
}

type BilliardMiddleware struct {
	manager *token.Manager
	server  *server.BilliardServer
}

func NewBilliardMiddleware(manager *token.Manager, billiardSrv *server.BilliardServer) *BilliardMiddleware {
	return &BilliardMiddleware{
		manager: manager,
		server:  billiardSrv,
	}
}

func (m *BilliardMiddleware) getTokenContextKey(t token.Token) string {
	return m.getTokenContextKeyByReflect(reflect.TypeOf(t).Elem())
}

func (m *BilliardMiddleware) getTokenContextKeyByReflect(rt reflect.Type) string {
	return fmt.Sprintf("%s:%s", tokenContextPrefix, rt.Name())
}

func (m *BilliardMiddleware) SaveToken(t token.Token, c *gin.Context) {
	c.Set(m.getTokenContextKey(t), t)
}

func (m *BilliardMiddleware) CancelToken(c *gin.Context, t token.Token) error {
	c.Set(m.getTokenContextKey(t), nil)
	return m.manager.Remove(c, t)
}

func (m *BilliardMiddleware) UserLoginStatMiddleware() gin.HandlerFunc {
	return m.headerTokenMiddleware(tokenHeaderKey, &UserToken{})
}

func (m *BilliardMiddleware) headerTokenMiddleware(headerKey string, tokenTmpl token.Token) gin.HandlerFunc {
	t := reflect.TypeOf(tokenTmpl)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		plog.Fatalf("TokenManagerMiddleware: token template should be ptr to struct")
	}

	t = t.Elem()
	tokenContextKey := m.getTokenContextKeyByReflect(t)

	return func(c *gin.Context) {
		tokenStr := c.GetHeader(headerKey)

		var nt token.Token
		if tokenStr != "" {
			nt = reflect.New(t).Interface().(token.Token)

			err := m.manager.Read(c, tokenStr, nt)
			if errors.Is(err, redis.Nil) {
				nt = nil
			} else if err != nil {
				plog.Errorf("token manager read token: %v error: %v", tokenStr, err)
			}
		}

		c.Set(tokenContextKey, nt)
		c.Next()

		afterProcessTokenTmp, exists := c.Get(tokenContextKey)
		if !exists || afterProcessTokenTmp == nil {
			return
		}

		afterProcessToken, ok := afterProcessTokenTmp.(token.Token)
		if !ok {
			plog.Errorf("AfterProcessToken should be a Token object")
			return
		}

		if err := m.manager.Save(c, afterProcessToken); err != nil {
			plog.Errorf("token manager save token: %v error: %v", afterProcessToken, err)
			return
		}
	}
}

func (m *BilliardMiddleware) UserLoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := m.GetLoginToken(c)
		if t == nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				pgin.ErrorRet(http.StatusUnauthorized, "登录过期或未登录"),
			)
			return
		}
	}
}

func (m *BilliardMiddleware) GetLoginToken(c *gin.Context) token.Token {
	val, exists := c.Get(m.getTokenContextKey(&UserToken{}))
	if !exists || val == nil {
		return nil
	}

	return val.(token.Token)
}

func (m *BilliardMiddleware) CurrentUserId(c *gin.Context) (int, error) {
	t := m.GetLoginToken(c)

	if t == nil {
		return -1, errors.New("token required")
	}

	userToken, ok := t.(*UserToken)
	if !ok {
		return -1, errors.New("invalid token")
	}

	return userToken.UserId, nil
}

func (m *BilliardMiddleware) CurrentUser(c *gin.Context) (*user.User, error) {
	t := m.GetLoginToken(c)

	if t == nil {
		return nil, errors.New("token required")
	}

	userToken, ok := t.(*UserToken)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return m.server.UserSrv.GetUserProfile(c, userToken.Token.AccessToken)
}

func (m *BilliardMiddleware) AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := m.CurrentUser(c)
		if user != nil && user.IsAdmin() {
			return
		}

		if err != nil {
			plog.Errorf("get current user error: %v", err)
		}

		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			pgin.ErrorRet(http.StatusUnauthorized, "admin required"),
		)
		return
	}
}
