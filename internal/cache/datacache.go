package cache

import (
	"time"

	"github.com/ein-gast/go-squashsf-httpd/internal/logger"
	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
	lru "github.com/hashicorp/golang-lru/v2"
)

type Key string
type Data struct {
	Data  []byte
	Size  uint64
	Mime  string
	MTime time.Time
}

type Cache interface {
	Get(Key) (*Data, bool)
	IsStorable(uint64) bool
	Store(Key, *Data)
	Clear(Key)
	ClearAll()
}

func NewCache(log logger.Logger, cfg *settings.Settings) Cache {
	if cfg.DataCacheOff {
		return &CacheNull{}
	}

	lru, err := lru.New[Key, *Data](cfg.DataCacheCount)
	if err != nil {
		log.Msg("Cache init error:", err.Error())
		return &CacheNull{}
	}

	return &CacheRam{
		EntSize: uint64(cfg.DataCacheEntSize),
		Ents:    lru,
	}
}
