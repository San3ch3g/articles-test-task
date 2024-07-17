package di

import (
	"articleModule/internal/pkg/storage/pg"
	"articleModule/internal/utils/config"
	"gorm.io/gorm"
	"net"
)

type Container struct {
	cfg         *config.Config
	netListener *net.Listener
	storage     *pg.Storage
	db          *gorm.DB
}

func New(cfg *config.Config) *Container {
	return &Container{cfg: cfg}
}

func (c *Container) GetDB() *gorm.DB {
	return get(&c.db, func() *gorm.DB {
		return pg.MustNewPostgresDB(c.cfg)
	})
}

func (c *Container) GetSQLStorage() *pg.Storage {
	return get(&c.storage, func() *pg.Storage {
		return pg.New(c.GetDB())
	})
}

func get[T comparable](obj *T, builder func() T) T {
	if *obj != *new(T) {
		return *obj
	}
	*obj = builder()
	return *obj
}
