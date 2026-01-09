package indexer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	DefaultChunkSize    = 512
	DefaultChunkOverlap = 50
)

type ChunkInfo struct {
	ID        string
	FilePath  string
	StartLine int
	EndLine   int
	Content   string
	Hash      string
}

type Chunker struct {
	chunkSize int
	overlap   int
}

func NewChunker(chunkSize, overlap int) *Chunker {
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}
	if overlap < 0 {
		overlap = DefaultChunkOverlap
	}
	if overlap >= chunkSize {
		overlap = chunkSize / 10
	}

	return &Chunker{
		chunkSize: chunkSize,
		overlap:   overlap,
	}
}

func (c *Chunker) Chunk(filePath string, content string) []ChunkInfo {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return nil
	}

	var chunks []ChunkInfo
	chunkIndex := 0

	// Estimate tokens per line (rough approximation: 1 line ≈ 10 tokens on average)
	// This is a simplification; a proper tokenizer would be more accurate
	tokensPerLine := 10
	linesPerChunk := c.chunkSize / tokensPerLine
	if linesPerChunk < 1 {
		linesPerChunk = 1
	}
	overlapLines := c.overlap / tokensPerLine
	if overlapLines < 1 {
		overlapLines = 1
	}

	for start := 0; start < len(lines); {
		end := start + linesPerChunk
		if end > len(lines) {
			end = len(lines)
		}

		// Create chunk content
		chunkLines := lines[start:end]
		chunkContent := strings.Join(chunkLines, "\n")

		// Skip empty chunks
		if strings.TrimSpace(chunkContent) == "" {
			start = end - overlapLines
			if start < 0 || start >= end {
				start = end
			}
			continue
		}

		// Generate chunk ID
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%d:%d:%s", filePath, start, end, chunkContent)))
		chunkID := fmt.Sprintf("%s_%d", filePath, chunkIndex)

		chunks = append(chunks, ChunkInfo{
			ID:        chunkID,
			FilePath:  filePath,
			StartLine: start + 1, // 1-indexed
			EndLine:   end,       // Inclusive, 1-indexed
			Content:   chunkContent,
			Hash:      hex.EncodeToString(hash[:8]), // Short hash for efficiency
		})

		chunkIndex++

		// Move to next chunk with overlap
		start = end - overlapLines
		if start <= chunks[len(chunks)-1].StartLine-1 {
			start = end // Prevent infinite loop
		}
	}

	return chunks
}

// ChunkWithContext adds surrounding context to improve embedding quality
func (c *Chunker) ChunkWithContext(filePath string, content string) []ChunkInfo {
	chunks := c.Chunk(filePath, content)

	// Add file path context to each chunk
	for i := range chunks {
		chunks[i].Content = fmt.Sprintf("File: %s\n\n%s", filePath, chunks[i].Content)
	}

	return chunks
}

// EstimateTokens provides a rough token count (simple word-based estimation)
func EstimateTokens(text string) int {
	words := strings.Fields(text)
	// Rough estimation: 1 word ≈ 1.3 tokens on average for code
	return int(float64(len(words)) * 1.3)
}
