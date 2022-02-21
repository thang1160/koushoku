package cache

import (
	"time"

	. "koushoku/config"

	"github.com/bluele/gcache"
)

var Cache *LRU
var TemplatesCache *LRU

func init() {
	Cache = &LRU{gcache.New(512).LRU().Expiration(time.Duration(Config.Cache.DefaultTTL)).Build()}
	TemplatesCache = &LRU{gcache.New(512).LRU().Expiration(5 * time.Minute).Build()}
}
