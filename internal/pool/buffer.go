package pool

import (
	"bytes"
	"sync"
)

type BufferPool struct {
	pool *sync.Pool
	size int
}

func NewBufferPool(size int) *BufferPool {
	res := &BufferPool{
		pool: &sync.Pool{},
		size: size,
	}
	res.pool.New = res.newBuff
	return res
}

func (p *BufferPool) newBuff() any {
	return bytes.NewBuffer(make([]byte, p.size))
}

func (p *BufferPool) New() *bytes.Buffer {
	b := p.pool.Get().(*bytes.Buffer)
	return b
}

func (p *BufferPool) Return(b *bytes.Buffer) {
	p.pool.Put(b)
}
