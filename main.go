package main

import (
	"gitea.hoven.com/billiard/billiard-assistant-server/api"
	"gitea.hoven.com/billiard/billiard-assistant-server/models"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/dal"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/oss/minio"
	"gitea.hoven.com/billiard/billiard-assistant-server/server"
	"gitea.hoven.com/core/auth-core/api/grpc/proto/authenticationpb"
	"gitea.hoven.com/core/auth-core/api/grpc/proto/userpb"
	"gitea.hoven.com/core/auth-core/api/grpc/proto/verifycodepb"
	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/dialer/grpc"
	"github.com/go-puzzles/puzzles/goredis"
	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/pgorm"
	"github.com/go-puzzles/puzzles/plog"

	consulpuzzle "github.com/go-puzzles/puzzles/cores/puzzles/consul-puzzle"
	httppuzzle "github.com/go-puzzles/puzzles/cores/puzzles/http-puzzle"
)

var (
	port          = pflags.Int("port", 29920, "Server run port")
	authCoreSrv   = pflags.String("authCoreSrv", "auth-core", "auth-core server name")
	srvConfigFlag = pflags.Struct("conf", (*models.Config)(nil), "server config")
	redisConfFlag = pflags.Struct("redisAuth", (*goredis.RedisConf)(nil), "redis auth config")
	mysqlConfFlag = pflags.Struct("mysqlAuth", (*pgorm.MysqlConfig)(nil), "mysql auth config")
	minioConfFlag = pflags.Struct("minioAuth", (*minio.MinioConfig)(nil), "minio auth config")
)

func main() {
	pflags.Parse()
	configs := models.ParseConfig(
		srvConfigFlag,
		redisConfFlag,
		mysqlConfFlag,
		minioConfFlag,
	)

	authCoreConn, err := grpc.DialGrpc(authCoreSrv())
	plog.PanicError(err)
	authenticationClient := authenticationpb.NewAuthCoreAuthenticationHandlerClient(authCoreConn)
	userClient := userpb.NewAuthCoreUserHandlerClient(authCoreConn)
	verifycodeClient := verifycodepb.NewAuthCoreVerifyCodeHandlerClient(authCoreConn)

	redisClient := configs.RedisConf.DialRedisClient()
	minioClient := minio.NewMinioOss(configs.SrvConf.BaseApi, configs.MinioConf)
	plog.PanicError(pgorm.RegisterSqlModelWithConf(configs.MysqlConf, dal.AllTables()...))
	plog.PanicError(pgorm.AutoMigrate(configs.MysqlConf))

	db := pgorm.GetDbByConf(configs.MysqlConf)

	billiardSrv := server.NewBilliardServer(
		configs.SrvConf,
		db, redisClient, minioClient,
		authenticationClient, userClient, verifycodeClient,
	)
	engine := api.SetupRouter(configs.SrvConf, redisClient, minioClient, billiardSrv)
	srv := cores.NewPuzzleCore(
		cores.WithService(pflags.GetServiceName()),
		consulpuzzle.WithConsulRegister(),
		httppuzzle.WithCoreHttpCORS(),
		httppuzzle.WithCoreHttpPuzzle("/api", engine),
	)

	plog.PanicError(cores.Start(srv, port()))
}
