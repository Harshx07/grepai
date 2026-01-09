package store

import (
	"context"
	"time"
)

// Chunk represents a piece of code with its vector embedding
type Chunk struct {
	ID        string    `json:"id"`
	FilePath  string    `json:"file_path"`
	StartLine int       `json:"start_line"`
	EndLine   int       `json:"end_line"`
	Content   string    `json:"content"`
	Vector    []float32 `json:"vector"`
	Hash      string    `json:"hash"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Document represents a file with its chunks
type Document struct {
	Path      string    `json:"path"`
	Hash      string    `json:"hash"`
	ModTime   time.Time `json:"mod_time"`
	ChunkIDs  []string  `json:"chunk_ids"`
}

// SearchResult represents a search match with its relevance score
type SearchResult struct {
	Chunk Chunk   `json:"chunk"`
	Score float32 `json:"score"`
}

// VectorStore defines the interface for vector storage backends
type VectorStore interface {
	// SaveChunks stores multiple chunks atomically
	SaveChunks(ctx context.Context, chunks []Chunk) error

	// DeleteByFile removes all chunks for a given file path
	DeleteByFile(ctx context.Context, filePath string) error

	// Search finds the most similar chunks to a query vector
	Search(ctx context.Context, queryVector []float32, limit int) ([]SearchResult, error)

	// GetDocument retrieves document metadata by path
	GetDocument(ctx context.Context, filePath string) (*Document, error)

	// SaveDocument stores document metadata
	SaveDocument(ctx context.Context, doc Document) error

	// DeleteDocument removes document metadata
	DeleteDocument(ctx context.Context, filePath string) error

	// ListDocuments returns all indexed document paths
	ListDocuments(ctx context.Context) ([]string, error)

	// Load reads the store from persistent storage
	Load(ctx context.Context) error

	// Persist writes the store to persistent storage
	Persist(ctx context.Context) error

	// Close cleanly shuts down the store
	Close() error
}
