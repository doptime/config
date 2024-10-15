package cfgredis

import (
	"context"
	"strconv"
	"time"

	"github.com/doptime/config/utils"

	"github.com/doptime/config"
	"github.com/doptime/logger"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/redis/go-redis/v9"
)

type DataSource struct {
	Name     string
	Username string
	Password string `psw:"true"`
	Host     string
	Port     int64
	DB       int64
}

//config.toml example

// [[redisSources]]
//
//	Name = "default"
//	Username = "doptime"
//	Password = "yourpasswordhere"
//	Host = "drangonflydb.local"
//	Port = 6379
//	DB = 0
var redisSources []*DataSource

var Servers cmap.ConcurrentMap[string, *redis.Client] = cmap.New[*redis.Client]()

func AfterLoad() (err error) {
	logger.Info().Str("Checking Redis", "Start").Send()
	for _, rdsCfg := range redisSources {
		//apply configuration
		redisOption := &redis.Options{
			Addr:         rdsCfg.Host + ":" + strconv.Itoa(int(rdsCfg.Port)),
			Username:     rdsCfg.Username,
			Password:     rdsCfg.Password, // no password set
			DB:           int(rdsCfg.DB),  // use default DB
			PoolSize:     200,
			DialTimeout:  time.Second * 10,
			ReadTimeout:  -1,
			WriteTimeout: time.Second * 300,
		}
		rdsClient := redis.NewClient(redisOption)
		//test connection
		if _, err = rdsClient.Ping(context.Background()).Result(); err != nil {
			logger.Fatal().Err(err).Any("Redis server ping error", rdsCfg.Host).Send()
			return err //if redis server is not valid, exit
		}
		//save to the list
		logger.Info().Str("Redis Load ", "Success").Any("RedisUsername", rdsCfg.Username).Any("RedisHost", rdsCfg.Host).Any("RedisPort", rdsCfg.Port).Send()
		Servers.Set(rdsCfg.Name, rdsClient)
		timeCmd := rdsClient.Time(context.Background())
		logger.Info().Any("Redis server time: ", timeCmd.Val().String()).Send()
		//ping the address of redisAddress, if failed, print to log
		utils.PingServer(rdsCfg.Host, true)
	}
	//check if default redis is set
	if _rds, ok := Servers.Get("default"); !ok {
		logger.Warn().Msg("\"default\" redis server missing in Configuration. RPC will can not be received. Please ensure this is what your want")
		return
	} else {
		Servers.Set("", _rds)
		logger.RdsClientToLog = _rds
	}
	return nil
}

func init() {
	config.LoadToml("Redis", &redisSources)
	AfterLoad()
}
