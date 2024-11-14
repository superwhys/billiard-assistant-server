package main

import (
	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/pgorm"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/predis"
	"github.com/superwhys/snooker-assistant-server/api"
	"github.com/superwhys/snooker-assistant-server/models"
	"github.com/superwhys/snooker-assistant-server/pkg/dal"
	"github.com/superwhys/snooker-assistant-server/pkg/oss/minio"
	"github.com/superwhys/snooker-assistant-server/server"

	httppuzzle "github.com/go-puzzles/puzzles/cores/puzzles/http-puzzle"
)

var (
	port          = pflags.Int("port", 29920, "Server run port")
	saConfigFlag  = pflags.Struct("conf", (*models.SaConfig)(nil), "server config")
	redisConfFlag = pflags.Struct("redisAuth", (*predis.RedisConf)(nil), "redis auth config")
	mysqlConfFlag = pflags.Struct("mysqlAuth", (*pgorm.MysqlConfig)(nil), "mysql auth config")
	minioConfFlag = pflags.Struct("minioAuth", (*models.MinioConfig)(nil), "minio auth config")
)

func main() {
	pflags.Parse(pflags.WithConsulEnable())

	saConfig, redisConf, mysqlConf, minioConf := models.ParseConfig(
		saConfigFlag,
		redisConfFlag,
		mysqlConfFlag,
		minioConfFlag,
	)

	minioClient := minio.NewMinioOss(minioConf)
	redisClient := predis.NewRedisClient(redisConf.DialRedisPool())
	plog.PanicError(pgorm.RegisterSqlModelWithConf(mysqlConf, dal.AllTables()...))
	plog.PanicError(pgorm.AutoMigrate(mysqlConf))

	db := pgorm.GetDbByConf(mysqlConf)

	saServer := server.NewSaServer(saConfig, db, redisClient, minioClient)
	engine := api.SetupRouter(redisClient, saServer)
	srv := cores.NewPuzzleCore(
		cores.WithService(pflags.GetServiceName()),
		httppuzzle.WithCoreHttpCORS(),
		httppuzzle.WithCoreHttpPuzzle("/api", engine),
	)

	plog.PanicError(cores.Start(srv, port()))
}
