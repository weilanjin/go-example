package cache

type cache struct {
	cacheStorage ThreadSafeStore
	keyFunc      KeyFunc
}
