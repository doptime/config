package redisdb

import (
	"context"
	"strconv"
	"time"

	"github.com/doptime/config/utils"

	"github.com/doptime/config"

	"github.com/doptime/doptime/dlog"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/redis/go-redis/v9"
)

type DataSource struct {
	Name     string
	Username string `env:"Username"`
	Password string `env:"Password" json:"pswd"`
	Host     string `env:"Host,required=true"`
	Port     int64  `env:"Port,required=true"`
	DB       int64  `env:"DB,required=true"`
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

var Rds cmap.ConcurrentMap[string, *redis.Client] = cmap.New[*redis.Client]()

func AfterLoad() (err error) {
	dlog.Info().Str("Checking Redis", "Start").Send()
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
			dlog.Fatal().Err(err).Any("Redis server ping error", rdsCfg.Host).Send()
			return err //if redis server is not valid, exit
		}
		//save to the list
		dlog.Info().Str("Redis Load ", "Success").Any("RedisUsername", rdsCfg.Username).Any("RedisHost", rdsCfg.Host).Any("RedisPort", rdsCfg.Port).Send()
		Rds.Set(rdsCfg.Name, rdsClient)
		timeCmd := rdsClient.Time(context.Background())
		dlog.Info().Any("Redis server time: ", timeCmd.Val().String()).Send()
		//ping the address of redisAddress, if failed, print to log
		utils.PingServer(rdsCfg.Host, true)
	}
	//check if default redis is set
	if _rds, ok := Rds.Get("default"); !ok {
		dlog.Warn().Msg("\"default\" redis server missing in Configuration. RPC will can not be received. Please ensure this is what your want")
		return
	} else {
		Rds.Set("", _rds)
		dlog.RdsClientToLog = _rds
	}
	return nil
}

func init() {
	config.LoadToml("Redis", &redisSources)
	AfterLoad()
}
