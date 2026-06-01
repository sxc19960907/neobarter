package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_NoFileUsesDefaults(t *testing.T) {
	viper.Reset()
	cfg, err := Load("/nonexistent/config.yaml")
	require.NoError(t, err)
	// 文件不存在时应回退到默认值
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "local", cfg.OSS.Provider)
}

func TestLoad_EnvOverride(t *testing.T) {
	viper.Reset()
	t.Setenv("NEOBARTER_DATABASE_HOST", "postgres")
	t.Setenv("NEOBARTER_REDIS_HOST", "redis")
	t.Setenv("NEOBARTER_SERVER_MODE", "release")

	cfg, err := Load("/nonexistent/config.yaml")
	require.NoError(t, err)
	assert.Equal(t, "postgres", cfg.Database.Host)
	assert.Equal(t, "redis", cfg.Redis.Host)
	assert.Equal(t, "release", cfg.Server.Mode)
}

func TestLoad_SliceAndURLEnvOverride(t *testing.T) {
	viper.Reset()
	t.Setenv("NEOBARTER_ELASTICSEARCH_ADDRESSES", "http://es1:9200,http://es2:9200")
	t.Setenv("NEOBARTER_RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")

	cfg, err := Load("/nonexistent/config.yaml")
	require.NoError(t, err)
	assert.Equal(t, []string{"http://es1:9200", "http://es2:9200"}, cfg.Elasticsearch.Addresses)
	assert.Equal(t, "amqp://guest:guest@rabbitmq:5672/", cfg.RabbitMQ.URL)
}
