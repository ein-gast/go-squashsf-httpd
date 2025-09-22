package cache

type CacheNull struct {
}

func (c *CacheNull) Get(Key) (*Data, bool) {
	return nil, false
}

func (c *CacheNull) IsStorable(uint64) bool {
	return false
}

func (c *CacheNull) Store(Key, *Data) {

}

func (c *CacheNull) Clear(Key) {

}

func (c *CacheNull) ClearAll() {

}
