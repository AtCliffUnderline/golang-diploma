package http

import (
	"github.com/AtCliffUnderline/golang-diploma/internal/database"
	"go.uber.org/zap"
	"time"

	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func SetupRouter(h Handler, ur database.UserStorage, l *zap.Logger) *gin.Engine {

	r := gin.Default()

	r.Use(ginzap.Ginzap(l, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(l, true))

	r.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	r.POST("/api/user/register", h.UserRegister)
	r.POST("/api/user/login", h.UserLogin)

	authorized := r.Group("/").Use(AuthMiddleware(ur))

	authorized.POST("/api/user/orders", h.OrderRegister)
	authorized.GET("/api/user/orders", h.OrdersGet)
	authorized.GET("/api/user/balance", h.BalanceGet)
	authorized.POST("/api/user/balance/withdraw", h.WithdrawRegister)
	authorized.GET("/api/user/balance/withdrawals", h.WithdrawalsGet)

	return r
}
