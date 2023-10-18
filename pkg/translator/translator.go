// Package translator contains translator functions
package translator

import "context"

// Translator is an interface that translates text
type Translator interface {
	Translate(ctx context.Context, s string) (string, error)
}
