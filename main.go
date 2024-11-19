package main

import (
	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/pgorm"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/predis"
	"gitlab.hoven.com/billiard/billiard-assistant-server/api"
	"gitlab.hoven.com/billiard/billiard-assistant-server/models"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/oss/minio"
	"gitlab.hoven.com/billiard/billiard-assistant-server/server"

	consulpuzzle "github.com/go-puzzles/puzzles/cores/puzzles/consul-puzzle"
	httppuzzle "github.com/go-puzzles/puzzles/cores/puzzles/http-puzzle"
)

var (
	port          = pflags.Int("port", 29920, "Server run port")
	srvConfigFlag = pflags.Struct("conf", (*models.Config)(nil), "server config")
	redisConfFlag = pflags.Struct("redisAuth", (*predis.RedisConf)(nil), "redis auth config")
	mysqlConfFlag = pflags.Struct("mysqlAuth", (*pgorm.MysqlConfig)(nil), "mysql auth config")
	minioConfFlag = pflags.Struct("minioAuth", (*models.MinioConfig)(nil), "minio auth config")
)

func main() {
	pflags.Parse()
	srvConfig, redisConf, mysqlConf, minioConf := models.ParseConfig(
		srvConfigFlag,
		redisConfFlag,
		mysqlConfFlag,
		minioConfFlag,
	)

	minioClient := minio.NewMinioOss(srvConfig.UserApi, minioConf)
	redisClient := predis.NewRedisClient(redisConf.DialRedisPool())
	plog.PanicError(pgorm.RegisterSqlModelWithConf(mysqlConf, dal.AllTables()...))
	plog.PanicError(pgorm.AutoMigrate(mysqlConf))

	db := pgorm.GetDbByConf(mysqlConf)

	billiardSrv := server.NewBilliardServer(srvConfig, db, redisClient, minioClient)
	engine := api.SetupRouter(srvConfig, redisClient, billiardSrv)
	srv := cores.NewPuzzleCore(
		cores.WithService(pflags.GetServiceName()),
		consulpuzzle.WithConsulRegister(),
		httppuzzle.WithCoreHttpCORS(),
		httppuzzle.WithCoreHttpPuzzle("/api", engine),
	)

	plog.PanicError(cores.Start(srv, port()))
}
