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
	var (
		defaultRds *redis.Client
	)
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
			logger.Fatal().Err(err).Str("Redis Server", rdsCfg.Name).Any("client pint error", rdsCfg.Host).Send()
			return err //if redis server is not valid, exit
		}

		if rdsCfg.Name == "" {
			rdsCfg.Name = "default"
		}
		//set default redis server, with empty name or name is default
		if defaultRds == nil || rdsCfg.Name == "default" {
			defaultRds = rdsClient
			Servers.Set("", defaultRds)
		}

		//save to the list
		Servers.Set(rdsCfg.Name, rdsClient)
		timeCmd := rdsClient.Time(context.Background())
		logger.Info().Str("Redis Server", rdsCfg.Name).Str("Time: ", timeCmd.Val().String()).Send()
		//ping the address of redisAddress, if failed, print to log
		utils.PingServer(rdsCfg.Host, true)
	}
	return nil
}

func init() {
	config.LoadItemFromToml("Redis", &redisSources)
	AfterLoad()
}
