package pool

type (
	task struct {
		argument interface{}
		function func(interface{})
	}
)

type (
	Pool struct {
		taskQueue chan *task
		done      chan struct{}
	}
)

func NewPool(size int) *Pool {
	if size <= 0 {
		panic("NewPool error: invalid size")
	}
	return &Pool{
		taskQueue: make(chan *task, size),
		done:      make(chan struct{}),
	}
}

func (p *Pool) Close() {
	close(p.taskQueue)
}

func (p *Pool) Start() {
	for i := 0; i < cap(p.taskQueue); i++ {
		go func() {
			for task := range p.taskQueue {
				task.function(task.argument)
				p.done <- struct{}{}
			}
		}()
	}
}

func (p *Pool) RunInBatch(function func(interface{}), arguments []interface{}) {
	if len(arguments) > cap(p.taskQueue) {
		panic("RunBatch error: invalid arguments number")
	}
	for _, argument := range arguments {
		p.taskQueue <- &task{
			function: function,
			argument: argument,
		}
	}
	for i := 0; i < len(arguments); i++ {
		<-p.done
	}
}
