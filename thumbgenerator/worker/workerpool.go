package worker

import (
	"database/sql"
	"go_study/thumbgenerator/provider"
	"sync"
)

const WorkersCount = 3

func WorkerPool(stopChan chan struct{}, db *sql.DB) *sync.WaitGroup {
	var wg sync.WaitGroup
	tasksChan := provider.RunTaskProvider(stopChan, db)
	for i := 0; i < WorkersCount; i++ {
		go func(i int) {
			wg.Add(1)
			Worker(tasksChan, db, i)
			wg.Done()
		}(i)
	}
	return &wg
}
