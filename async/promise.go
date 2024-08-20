package async

import (
	"sync"
)

type Promise struct {
	execute func(resolve func(interface{}), reject func(error))
	result  chan interface{}
	err     chan error
}

func NewPromise(execute func(resolve func(interface{}), reject func(error))) *Promise {
	p := &Promise{
		execute: execute,
		result:  make(chan interface{}, 1),
		err:     make(chan error, 1),
	}
	p.Execute()
	return p
}

func (p *Promise) Execute() {
	go func() {
		p.execute(p.resolve, p.reject)
	}()
}

func (p *Promise) resolve(value interface{}) {
	p.result <- value
	close(p.result)
	close(p.err)
}

func (p *Promise) reject(err error) {
	p.err <- err
	close(p.err)
	close(p.result)
}

func (p *Promise) Await() (interface{}, error) {
	select {
	case res := <-p.result:
		return res, nil
	case err := <-p.err:
		return nil, err
	}
}

func AwaitAll(promises []*Promise) ([]interface{}, error) {
	var wg sync.WaitGroup
	results := make([]interface{}, len(promises))
	errorsChan := make(chan error, len(promises))

	for i, p := range promises {
		wg.Add(1)
		go func(i int, p *Promise) {
			defer wg.Done()
			result, err := p.Await()
			if err != nil {
				errorsChan <- err
				return
			}
			results[i] = result
		}(i, p)
	}

	wg.Wait()
	close(errorsChan)

	if len(errorsChan) > 0 {
		return nil, <-errorsChan
	}

	return results, nil
}

func AwaitRace(promises []*Promise) (interface{}, error) {
	resultChan := make(chan interface{})
	errorChan := make(chan error)

	for _, p := range promises {
		go func(p *Promise) {
			result, err := p.Await()
			if err != nil {
				errorChan <- err
			} else {
				resultChan <- result
			}
		}(p)
	}

	select {
	case res := <-resultChan:
		return res, nil
	case err := <-errorChan:
		return nil, err
	}
}
