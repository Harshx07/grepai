package embedder

import "context"

// Embedder defines the interface for text embedding providers
type Embedder interface {
	// Embed converts text into a vector embedding
	Embed(ctx context.Context, text string) ([]float32, error)

	// EmbedBatch converts multiple texts into vector embeddings
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)

	// Dimensions returns the vector dimension size for this embedder
	Dimensions() int

	// Close cleanly shuts down the embedder
	Close() error
}
