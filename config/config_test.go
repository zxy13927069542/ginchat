package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	c := Init("../etc")

	assert.Equal(t, "root:066311@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local", c.Mysql.Dns)
	assert.Equal(t, "0.0.0.0", c.Redis.Host)
	assert.Equal(t, "6379", c.Redis.Port)
	assert.Equal(t, "", c.Redis.Password)
	assert.Equal(t, 30, c.Redis.PoolSize)
	assert.Equal(t, 0, c.Redis.DB)
	assert.Equal(t, 30, c.Redis.MinIdleConn)
}
