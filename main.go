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
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/email"
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
	minioConfFlag = pflags.Struct("minioAuth", (*minio.MinioConfig)(nil), "minio auth config")
	emailConfFlag = pflags.Struct("emailAuth", (*email.EmailConf)(nil), "email auth config")
)

func main() {
	pflags.Parse()
	configs := models.ParseConfig(
		srvConfigFlag,
		redisConfFlag,
		mysqlConfFlag,
		minioConfFlag,
		emailConfFlag,
	)
	redisClient := predis.NewRedisClient(configs.RedisConf.DialRedisPool())
	minioClient := minio.NewMinioOss(configs.SrvConf.UserApi, configs.MinioConf)
	neteasyEmailSender := email.NewNetEasySender(configs.EmailConf, redisClient)
	plog.PanicError(pgorm.RegisterSqlModelWithConf(configs.MysqlConf, dal.AllTables()...))
	plog.PanicError(pgorm.AutoMigrate(configs.MysqlConf))

	db := pgorm.GetDbByConf(configs.MysqlConf)

	billiardSrv := server.NewBilliardServer(configs.SrvConf, db, redisClient, minioClient, neteasyEmailSender)
	engine := api.SetupRouter(configs.SrvConf, redisClient, billiardSrv)
	srv := cores.NewPuzzleCore(
		cores.WithService(pflags.GetServiceName()),
		consulpuzzle.WithConsulRegister(),
		httppuzzle.WithCoreHttpCORS(),
		httppuzzle.WithCoreHttpPuzzle("/api", engine),
		cores.WithDaemonWorker(neteasyEmailSender.LoopAsyncTask),
	)

	plog.PanicError(cores.Start(srv, port()))
}
