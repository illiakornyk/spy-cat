package breeds

import (
	"sync"
	"time"
)

var (
    breedCache []Breed
    cacheMutex sync.RWMutex
)

func StartBreedCache( interval time.Duration) {
    updateBreeds()
    ticker := time.NewTicker(interval)
    go func() {
        for {
            <-ticker.C
            updateBreeds()
        }
    }()
}

func updateBreeds() {
    breeds, err := FetchBreeds()
    if err != nil {
        return
    }

    cacheMutex.Lock()
    breedCache = breeds
    cacheMutex.Unlock()
}

func IsValidBreed(breed string) bool {
    cacheMutex.RLock()
    defer cacheMutex.RUnlock()

    for _, b := range breedCache {
        if b.Name == breed {
            return true
        }
    }
    return false
}

func GetBreeds() []Breed {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return breedCache
}
