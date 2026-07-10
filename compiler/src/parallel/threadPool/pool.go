package threadPool

import (
	"sync"
)

type ThreadPool struct {
	taskQueue chan job
	wg        sync.WaitGroup
}

func NewThreadPool(size int) *ThreadPool {
	pool := &ThreadPool{
		taskQueue: make(chan job, 1024),
	}

	for range size {
		go pool.work()
	}

	return pool
}

func (p *ThreadPool) Enqueue(j job) {
	p.wg.Add(1)
	p.taskQueue <- j
}

func (p *ThreadPool) Wait() {
	p.wg.Wait()
}

func (p *ThreadPool) Close() {
	close(p.taskQueue)
}

func (p *ThreadPool) work() {
	for job := range p.taskQueue {
		job()
		p.wg.Done()
	}
}
