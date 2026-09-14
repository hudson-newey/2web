package constants

const EnvCacheOverride string = "__2WEB_CACHE_PATH"
const EnvDebugOverride string = "__2WEB_DEBUG_PATH"
const EnvCiOverride string = "CI"

// EnvCacheMaxSize overrides the build cache size (in bytes) above which the
// cache is vacuumed after a build.
const EnvCacheMaxSize string = "__2WEB_CACHE_MAX_SIZE"

// EnvCacheMaxAge overrides the age (in days) after which a build cache entry
// that wasn't touched by a build becomes a vacuum candidate.
const EnvCacheMaxAge string = "__2WEB_CACHE_MAX_AGE"
