package mentionbot

import (
	"math/rand"
	"time"
)

type idsStore struct {
	expires time.Time
	ids     []int64
}

func (store *idsStore) setIds(ids []int64, d time.Duration) {
	if d == 0 {
		d = 15 * time.Minute
	}
	store.ids = ids
	store.expires = time.Now().Add(d)
}

func (store *idsStore) pickIds() (ids []int64) {
	if time.Now().After(store.expires) {
		return
	}
	// shuffle
	n := len(store.ids)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		store.ids[i], store.ids[j] = store.ids[j], store.ids[i]
	}

	maxNum := 1000
	if len(store.ids) < maxNum {
		maxNum = len(store.ids)
	}
	return store.ids[0:maxNum]
}

func (current *rateLimitStatus) waitSeconds(last *rateLimitStatus) int64 {
	var wait int64 = 10
	if diff := int(last.Remaining) - int(current.Remaining); diff > 0 {
		num := int(current.Remaining) / diff
		if num == 0 {
			num++
		}
		w := (current.Reset - time.Now().Unix()) / int64(num)
		if w > wait {
			wait = w
		}
	}
	return wait
}
