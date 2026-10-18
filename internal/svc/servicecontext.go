package svc

import (
	"github.com/wangminghui0622-bit/admin-common/dal/query"
	"github.com/wangminghui0622-bit/admin-message/internal/config"
	"github.com/wangminghui0622-bit/admin-message/internal/utils/emailclient"
	"github.com/wangminghui0622-bit/admin-message/internal/utils/smssdk"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Query  *query.Query
	Redis  *redis.Redis
	Sms    *smssdk.SmsGroup
	Email  *emailclient.ClientCache
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DatabaseDSN), &gorm.Config{})
	logx.Must(err)

	redisClient := redis.MustNewRedis(redis.RedisConf{
		Host: c.RedisConf.Host,
		Type: c.RedisConf.Type,
		Pass: c.RedisConf.Pass,
	})

	return &ServiceContext{
		Config: c,
		DB:     db,
		Query:  query.Use(db),
		Redis:  redisClient,
		Sms:    smssdk.NewSmsGroup(),
		Email:  emailclient.NewClientCache(),
	}
}
