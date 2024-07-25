package main

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

type CachedSecret struct {
	encryptedSecret string
	addedAt         time.Time
	expiresAt       time.Time
	uses            int
}

type InMemorySecrets struct {
	mutex      sync.Mutex
	cache      map[string]CachedSecret
	maxEntries uint
}

func NewInMemorySecrets(maxEntries uint) *InMemorySecrets {
	return &InMemorySecrets{
		cache:      make(map[string]CachedSecret),
		maxEntries: maxEntries,
	}
}

func (s *InMemorySecrets) clean() {
	now := time.Now()
	var oldest *CachedSecret = nil
	for k, v := range s.cache {
		if v.expiresAt.Before(now) {
			delete(s.cache, k)
		} else if oldest == nil || v.addedAt.Before(oldest.addedAt) {
			oldest = &v
		}
	}
	if len(s.cache) > int(s.maxEntries) && oldest != nil {
		delete(s.cache, oldest.encryptedSecret)
	}
}

func (s *InMemorySecrets) Add(encryptedSecret string, uses int, exp time.Duration) string {
	id := generateID()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.clean()
	s.cache[id] = CachedSecret{encryptedSecret: encryptedSecret, addedAt: time.Now(), expiresAt: time.Now().Add(exp), uses: uses}
	return id
}

func (s *InMemorySecrets) Get(id string) (string, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	v, ok := s.cache[id]
	if !ok {
		return "", false
	}
	if v.expiresAt.Before(time.Now()) {
		delete(s.cache, id)
		return "", false
	}
	if v.uses == 1 {
		delete(s.cache, id)
	} else {
		v.uses--
	}
	return v.encryptedSecret, true
}
