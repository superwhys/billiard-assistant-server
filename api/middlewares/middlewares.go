package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	
	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/token"
	"gitlab.hoven.com/billiard/billiard-assistant-server/server"
)

const (
	tokenContextPrefix = "billiard:token"
	tokenHeaderKey     = "X-BILLIARD-Token"
)

type UserToken struct {
	TokenId  string
	Uid      int
	Username string
	// WechatId is openid
	WechatId         string
	WechatUnionId    string
	WechatSessionKey string
}

func NewUserToken(uid int, wechatId string, username string) *UserToken {
	return &UserToken{
		TokenId:  uuid.New().String(),
		Uid:      uid,
		Username: username,
		WechatId: wechatId,
	}
}

func (t *UserToken) GetKey() string {
	if t.TokenId == "" {
		t.TokenId = uuid.New().String()
	}
	return t.TokenId
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
			
			err := m.manager.Read(tokenStr, nt)
			if errors.Is(err, redis.ErrNil) {
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
		
		if err := m.manager.Save(afterProcessToken); err != nil {
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
				pgin.ErrorRet(http.StatusUnauthorized, "token required"),
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
	
	return userToken.Uid, nil
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
	
	return m.server.UserSrv.GetUserById(c, userToken.Uid)
}

func (m *BilliardMiddleware) AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := m.CurrentUser(c)
		if err != nil || !user.IsAdmin() {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				pgin.ErrorRet(http.StatusUnauthorized, "admin required"),
			)
			return
		}
	}
}
