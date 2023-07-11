package id

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/syncmap"
)

func TestRandom(t *testing.T) {
	t.Parallel()

	var (
		adapter     = New()
		totalTests  = 10000
		defaultSize = 26
		waitGroup   = sync.WaitGroup{}
		uniqueIDs   = syncmap.Map{}
	)

	waitGroup.Add(totalTests)

	for i := 0; i < totalTests; i++ {
		go func() {
			defer waitGroup.Done()

			randomID := adapter.Random()

			_, ok := uniqueIDs.LoadOrStore(randomID, nil)
			valueDuplicationMsg := fmt.Sprintf("Value duplicated %v", randomID)
			valueHasNotDefaultSizeMsg := fmt.Sprintf("Value has not default size (%v) - value: %v", defaultSize, randomID)

			assert.False(t, ok, valueDuplicationMsg)
			assert.Len(t, randomID, defaultSize, valueHasNotDefaultSizeMsg)
		}()
	}

	waitGroup.Wait()

	totalItems := 0
	rangeFunc := func(key, value interface{}) bool {
		totalItems++

		return true
	}

	uniqueIDs.Range(rangeFunc)
	assert.Equal(t, totalItems, totalTests)
}
