package providers

import (
	"sync"
)

type ProviderStat struct {
	Total   int
	Success int
	Rate    float64
	mutex   sync.Mutex
}

var ProviderStats sync.Map

func getOrCreateStat(provider string) *ProviderStat {
	val, ok := ProviderStats.Load(provider) // Attempt to load the existing value.
	if !ok {
		// If the value does not exist, create a new one and store it.
		val = &ProviderStat{}
		ProviderStats.Store(provider, val)
	}

	// Since we store *ProviderStat in the map, we need to cast the value back to *ProviderStat.
	return val.(*ProviderStat)
}

func updateStat(provider string, success bool) {
	go func() {
		stat := getOrCreateStat(provider)

		stat.mutex.Lock()
		defer stat.mutex.Unlock()

		stat.Total++
		if success {
			stat.Success++
		}

		// Calculate success rate
		if stat.Total > 0 {
			stat.Rate = float64(stat.Success) / float64(stat.Total)
		}
	}()
}
