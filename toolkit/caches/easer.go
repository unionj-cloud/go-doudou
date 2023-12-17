package caches

import "sync"

func ease(t task, queue *sync.Map) task {
	eq := &eased{
		task: t,
		wg:   &sync.WaitGroup{},
	}
	eq.wg.Add(1)

	runner, ok := queue.LoadOrStore(t.GetId(), eq)
	et := runner.(*eased)

	// If this request is the first of its kind, we execute the Run
	if !ok {
		et.task.Run()

		queue.Delete(et.task.GetId())
		et.wg.Done()
	}

	et.wg.Wait()
	return et.task
}

type eased struct {
	task task
	wg   *sync.WaitGroup
}
