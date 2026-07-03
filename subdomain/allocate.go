package subdomain

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
)

// Store checks global subdomain uniqueness.
type Store interface {
	SubdomainExists(ctx context.Context, subdomain string) (bool, error)
}

const defaultPrefix = "tunnel-"

// SemanticAllocate generates a subdomain with a meaningful prefix and an
// 8-character random alphanumeric suffix, checking global uniqueness.
// If prefix is empty, "tunnel-" is used.
func SemanticAllocate(ctx context.Context, st Store, prefix string) (string, error) {
	if prefix == "" {
		prefix = defaultPrefix
	}
	const suffixLen = 8
	for i := 0; i < 20; i++ {
		suffix, err := randomAlphanumeric(suffixLen)
		if err != nil {
			return "", err
		}
		candidate := prefix + suffix
		exists, err := st.SubdomainExists(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not allocate semantic subdomain with prefix %q after 20 attempts", prefix)
}

// RenameWithPrefix allocates a new subdomain with the given prefix.
func RenameWithPrefix(ctx context.Context, st Store, newPrefix string) (string, error) {
	return SemanticAllocate(ctx, st, newPrefix)
}

func randomAlphanumeric(n int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	out := make([]byte, n)
	for i := range out {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		out[i] = alphabet[idx.Int64()]
	}
	return string(out), nil
}