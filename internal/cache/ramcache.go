package cache

import (
	lru "github.com/hashicorp/golang-lru/v2"
)

type CacheRam struct {
	EntSize uint64
	Ents    *lru.Cache[Key, *Data]
}

func (c *CacheRam) Get(k Key) (*Data, bool) {
	return c.Ents.Get(k)
}

func (c *CacheRam) IsStorable(s uint64) bool {
	return c.EntSize >= s
}

func (c *CacheRam) Store(k Key, d *Data) {
	c.Ents.Add(k, d)
}

func (c *CacheRam) Clear(k Key) {
	c.Ents.Remove(k)
}

func (c *CacheRam) ClearAll() {
	c.Ents.Purge()
}
