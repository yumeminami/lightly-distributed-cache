package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

//管理不同 key 的请求(call)
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call
}


//针对相同的 key，无论 Do 被调用多少次，函数 fn 都只会被调用一次，等待 fn 调用结束了，返回返回值或错误。
//使用了信号量
//wg.Add(1) 锁加1。
//wg.Wait() 阻塞，直到锁被释放。
//wg.Done() 锁减1。
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}