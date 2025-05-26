package utils

import "hash/fnv"

func HashUrl(url string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(url))
	return hash.Sum64()
}
