package hash

//go:generate mockgen -package=mock -destination=./mock/hasher.go wowpow/internal/pkg/hash Hasher

// Hasher interface to hash function
type Hasher interface {
	Hash(str string) (string, error)
}
