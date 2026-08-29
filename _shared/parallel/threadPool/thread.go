package threadPool

func newThread(pool *ThreadPool) thread {
	t := thread{
		job:  func() {},
		pool: pool,
	}
	return t
}

type thread struct {
	job  job
	pool *ThreadPool
}

func (m *thread) nextJob(job job) {
	m.job = job
	m.execute()
}

func (m *thread) execute() {
	m.job()
}
