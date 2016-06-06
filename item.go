package ttlcache

import (
	"sync"
	"time"
)

const (
	ItemNotExpire           time.Duration = -1
	ItemExpireWithGlobalTTL time.Duration = 0
)

func newItem(key string, data interface{}, ttl time.Duration) *Item {
	item := &Item{
		Data: data,
		Ttl:  ttl,
		Key:  key,
	}
	item.touch()
	return item
}

type Item struct {
	Key        string
	Data       interface{}
	Ttl        time.Duration
	ExpireAt   time.Time
	mutex      sync.Mutex `json:"-"`
	QueueIndex int
}

// Reset the item expiration time
func (item *Item) touch() {
	item.mutex.Lock()
	if item.Ttl > 0 {
		item.ExpireAt = time.Now().Add(item.Ttl)
	}
	item.mutex.Unlock()
}

// Verify if the item is expired
func (item *Item) expired() bool {
	item.mutex.Lock()
	if item.Ttl <= 0 {
		item.mutex.Unlock()
		return false
	}
	expired := item.ExpireAt.Before(time.Now())
	item.mutex.Unlock()
	return expired
}
