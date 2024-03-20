package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	ch := make(chan Task, len(tasks))
	eCount := 0

	for _, t := range tasks {
		ch <- t
	}
	close(ch)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				mu.Lock()
				needToBreak := eCount >= m
				mu.Unlock()
				if needToBreak {
					break
				}
				task, ok := <-ch
				if !ok {
					break
				}
				err := task()
				if err != nil {
					mu.Lock()
					eCount++
					mu.Unlock()
				}
			}
		}()
	}
	wg.Wait()

	mu.Lock()
	needToReturnError := eCount >= m
	mu.Unlock()
	if needToReturnError {
		return ErrErrorsLimitExceeded
	}
	return nil
}
