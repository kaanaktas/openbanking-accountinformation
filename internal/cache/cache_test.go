package cache

import (
	gocache "github.com/patrickmn/go-cache"
	"testing"
	"time"
)

func TestLoadInMemory(t *testing.T) {
	tests := []struct {
		name           string
		usedCacheId    string
		queriedCacheId string
		want           bool
	}{
		{"retrieve_cache_success", "cacheId", "cacheId", true},
		{"retrieve_cache_fail", "cacheId", "wrongCacheId", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadInMemory(); &got != nil {
				_ = got.Set(tt.usedCacheId, "value", gocache.DefaultExpiration)
				if _, found := got.Get(tt.queriedCacheId); found != tt.want {
					t.Errorf("LoadInMemory() = %v, want %v", found, tt.want)
				}
			} else {
				t.Errorf("Couldn't LoadInMemory()")
			}
		})
	}
}

func TestManager_LoadInMemory(t *testing.T) {
	var c = LoadInMemory()
	c.Set("name", "test", 0)
	c.Set("age", 55, 0)

	if _, found := c.Get("name"); !found {
		t.Errorf("Couldn't find item")
	}

	if _, found := c.Get("age"); !found {
		t.Errorf("Couldn't find item")
	}

	var c1 = LoadInMemory()
	c1.Set("name2", "test2", 0)
	c1.Set("name3", "test3", 0)

	if _, found := c1.Get("name2"); !found {
		t.Errorf("Couldn't find item")
	}

	if _, found := c1.Get("name3"); !found {
		t.Errorf("Couldn't find item")
	}

	if _, found := c1.Get("name"); !found {
		t.Errorf("Couldn't find item")
	}

	if c != c1 {
		t.Errorf("caches inMemory addresses should be same")
	}
}

func Test_redis_ReadWrite(t *testing.T) {
	type args struct {
		k string
		v interface{}
		d time.Duration
	}

	// this is an integration test. That's why if there is no connection
	//with redis server, then just skip the test
	redisRef.initiateRedis()
	if err := redisRef.ping(); err != nil {
		t.SkipNow()
	}

	tests := []struct {
		name      string
		redisRef  redis
		args      args
		wantValue interface{}
		wantErr   bool
	}{
		{
			"redis_ReadWrite_success_string",
			redisRef,
			args{
				k: "key_1",
				v: "value_1",
				d: 60,
			},
			"value_1",
			false,
		},
		{
			"redis_ReadWrite_success_overwrite",
			redisRef,
			args{
				k: "key_1",
				v: "value_2",
				d: NoExpiration,
			},
			"value_2",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.redisRef.Set(tt.args.k, tt.args.v, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			v, found := tt.redisRef.Get(tt.args.k)
			if !found {
				t.Errorf("Set() the key %v doesn't exist", tt.args.k)
				return
			}

			if v.(string) != tt.wantValue {
				t.Errorf("Set() the value for key %v isn't correct. got %v want %v", tt.args.k, string(v.([]byte)), tt.wantValue)
				return
			}
		})
	}
}
