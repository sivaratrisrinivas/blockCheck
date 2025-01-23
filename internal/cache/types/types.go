package types

// Stats represents cache statistics
type Stats struct {
	Hits   uint64
	Misses uint64
	Keys   uint64
}
