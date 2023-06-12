package bard

import "time"

type BardCache struct {
	Bards map[string]*Bard
}

func GarbageCollectCache(cache *BardCache) {
	for k, v := range cache.Bards {
		if time.Since(v.LastInteractionTime) > time.Minute*5 {
			delete(cache.Bards, k)
		}
	}
}

var cache *BardCache

func init() {
	cache = &BardCache{
		Bards: make(map[string]*Bard),
	}
	go func() {
		for {
			GarbageCollectCache(cache)
			time.Sleep(time.Minute)
		}
	}()
}
