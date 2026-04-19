package cache_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/cache"
)

func TestGet_MissOnEmpty(t *testing.T) {
	c := cache.New(time.Minute)
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected cache miss on empty cache")
	}
}

func TestSet_ThenGet(t *testing.T) {
	c := cache.New(time.Minute)
	secrets := map[string]string{"KEY": "value"}
	c.Set("secret/app", secrets)
	got, ok := c.Get("secret/app")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got["KEY"] != "value" {
		t.Errorf("expected value 'value', got %q", got["KEY"])
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("secret/app", map[string]string{"A": "1"})
	got, _ := c.Get("secret/app")
	got["A"] = "mutated"
	again, _ := c.Get("secret/app")
	if again["A"] != "1" {
		t.Error("cache entry was mutated through returned map")
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	c := cache.New(time.Millisecond)
	c.Set("secret/app", map[string]string{"X": "y"})
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestInvalidate(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("secret/app", map[string]string{"K": "v"})
	c.Invalidate("secret/app")
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestFlush(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("secret/a", map[string]string{"A": "1"})
	c.Set("secret/b", map[string]string{"B": "2"})
	c.Flush()
	_, okA := c.Get("secret/a")
	_, okB := c.Get("secret/b")
	if okA || okB {
		t.Fatal("expected all entries flushed")
	}
}

func TestEntry_IsExpired_ZeroTTL(t *testing.T) {
	e := cache.Entry{TTL: 0, FetchedAt: time.Now()}
	if !e.IsExpired() {
		t.Error("zero TTL entry should always be expired")
	}
}
