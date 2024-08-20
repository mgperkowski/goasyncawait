package async

import (
	"sync"
)

type Promise struct {
	execute func(resolve func(interface{}), reject func(error))
	result  interface{}
	err     error
	done    bool
	mutex   sync.Mutex
	wg      sync.WaitGroup
}

func NewPromise(execute func(resolve func(interface{}), reject func(error))) *Promise {
	p := &Promise{
		execute: execute,
	}
	p.wg.Add(1)

	go p.execute(p.resolve, p.reject)
	return p
}

func (p *Promise) resolve(value interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.done {
		p.result = value
		p.done = true
		p.wg.Done()
	}
}

func (p *Promise) reject(err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.done {
		p.err = err
		p.done = true
		p.wg.Done()
	}
}

func (p *Promise) Await() (interface{}, error) {
	p.wg.Wait()
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.result, p.err
}

func AwaitAll(promises []*Promise) ([]interface{}, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	results := make([]interface{}, len(promises))
	var firstError error

	for i, p := range promises {
		wg.Add(1)
		go func(index int, promise *Promise) {
			defer wg.Done()
			result, err := promise.Await()

			mutex.Lock()
			defer mutex.Unlock()

			if err != nil && firstError == nil {
				firstError = err
			}

			if firstError == nil {
				results[index] = result
			}
		}(i, p)
	}

	wg.Wait()

	if firstError != nil {
		return nil, firstError
	}

	return results, nil
}

func AwaitRace(promises []*Promise) (interface{}, error) {
	resultChan := make(chan struct {
		result interface{}
		err    error
	})

	for _, p := range promises {
		go func(promise *Promise) {
			result, err := promise.Await()
			resultChan <- struct {
				result interface{}
				err    error
			}{result, err}
		}(p)
	}

	firstResult := <-resultChan
	return firstResult.result, firstResult.err
}
