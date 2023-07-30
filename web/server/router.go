package server

import (
	"GoLeaf/config"
	"GoLeaf/service"
	"github.com/gin-gonic/gin"
	"github.com/panjf2000/ants/v2"
)

var (
	sf   *service.Snowflake
	pool *ants.Pool
)

// todo 高并发场景
func init() {
	config.LoadConfig()
	workID := config.Configs.GetInt64("goLeaf.Snowflake.workID")
	sf = service.NewSnowflake(workID)
	pool, _ = ants.NewPool(100000000000)
}
func Start() {
	defer pool.Release()
	r := gin.Default()
	r.GET("/get", func(c *gin.Context) {
		_ = pool.Submit(func() {
			uniqueID := sf.GenerateID()
			c.JSON(200, gin.H{
				"ID": uniqueID,
			})
		})
	})
	pool.Waiting()
	r.Run(":" + config.Configs.GetString("goLeaf.port"))
}
