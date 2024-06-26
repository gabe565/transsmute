package kemono

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

//nolint:gochecknoglobals
var creatorCache *ttlcache.Cache[creatorCacheKey, *Creator]

func initCreatorCache() {
	creatorCache = ttlcache.New[creatorCacheKey, *Creator](
		ttlcache.WithTTL[creatorCacheKey, *Creator](24*time.Hour),
		ttlcache.WithDisableTouchOnHit[creatorCacheKey, *Creator](),
	)
	go creatorCache.Start()
}
