package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach/pkg/sql/vecindex"
	"github.com/cockroachdb/cockroach/pkg/sql/vecindex/quantize"
	"github.com/cockroachdb/cockroach/pkg/sql/vecindex/vecstore"
	"github.com/cockroachdb/cockroach/pkg/util/stop"
)

const (
	dimensions       = 512
	minPartitionSize = 16
	maxPartitionSize = 128
	seed             = 42
)

type ImageMetadata struct {
	PhotoID     string
	PhotoURL    string
	Description string
}

type VectorSearchIndex struct {
	index    *vecindex.VectorIndex
	metadata []ImageMetadata
}

type TextEmbeddingService struct {
	baseURL string
	client  *http.Client
}

type embedRequest struct {
	Text string `json:"text"`
}

type embedResponse struct {
	Embedding []float32 `json:"embedding"`
}

func NewTextEmbeddingService(baseURL string) *TextEmbeddingService {
	return &TextEmbeddingService{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *TextEmbeddingService) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	reqBody := embedRequest{Text: text}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/embed-text", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status %d: %s", resp.StatusCode, string(body))
	}

	var embedResp embedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return embedResp.Embedding, nil
}

func loadEmbeddingsFromCSV(filename string) ([]ImageMetadata, [][]float32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Allow variable number of fields
	reader.FieldsPerRecord = -1

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, nil, err
	}

	var metadata []ImageMetadata
	var embeddings [][]float32

	lineNum := 1
	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("error on line %d: %w", lineNum, err)
		}

		// Ensure minimum required fields
		if len(record) < 4 {
			return nil, nil, fmt.Errorf("record on line %d: not enough fields (got %d, minimum 4)",
				lineNum, len(record))
		}

		// Parse embedding string from CSV
		embStr := strings.Trim(record[3], "[]")
		embParts := strings.Split(embStr, ",")
		embedding := make([]float32, len(embParts))

		for i, val := range embParts {
			f, err := strconv.ParseFloat(strings.TrimSpace(val), 32)
			if err != nil {
				return nil, nil, fmt.Errorf("parsing embedding on line %d: %w", lineNum, err)
			}
			embedding[i] = float32(f)
		}

		// Handle empty fields with default values
		photoID := record[0]
		photoURL := record[1]
		description := ""
		if len(record[2]) > 0 {
			description = record[2]
		}

		metadata = append(metadata, ImageMetadata{
			PhotoID:     photoID,
			PhotoURL:    photoURL,
			Description: description,
		})
		embeddings = append(embeddings, embedding)
	}

	return metadata, embeddings, nil
}

func NewVectorSearchIndex(ctx context.Context, metadata []ImageMetadata, embeddings [][]float32) (*VectorSearchIndex, error) {
	store := vecstore.NewInMemoryStore(dimensions, seed)
	quantizer := quantize.NewRaBitQuantizer(dimensions, seed)
	options := vecindex.VectorIndexOptions{
		MinPartitionSize: minPartitionSize,
		MaxPartitionSize: maxPartitionSize,
		BaseBeamSize:     8,
		Seed:             seed,
	}

	stopper := stop.NewStopper()

	index, err := vecindex.NewVectorIndex(ctx, store, quantizer, &options, stopper)
	if err != nil {
		return nil, err
	}

	// Insert embeddings into index
	txn, err := store.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer store.Commit(ctx, txn)

	for i, embedding := range embeddings {
		// Create primary key from index
		key := make([]byte, 4)
		binary.BigEndian.PutUint32(key, uint32(i))

		// Insert vector
		store.InsertVector(key, embedding)
		if err := index.Insert(ctx, txn, embedding, key); err != nil {
			return nil, err
		}
	}

	return &VectorSearchIndex{
		index:    index,
		metadata: metadata,
	}, nil
}

func (vsi *VectorSearchIndex) Search(ctx context.Context, queryEmbedding []float32, k int) ([]ImageMetadata, error) {
	txn, err := vsi.index.Store().Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer vsi.index.Store().Commit(ctx, txn)

	searchSet := vecstore.SearchSet{MaxResults: k}
	searchOptions := vecindex.SearchOptions{BaseBeamSize: 8}

	err = vsi.index.Search(ctx, txn, queryEmbedding, &searchSet, searchOptions)
	if err != nil {
		return nil, err
	}

	results := searchSet.PopResults()
	similarImages := make([]ImageMetadata, 0, len(results))

	for _, result := range results {
		idx := binary.BigEndian.Uint32(result.ChildKey.PrimaryKey)
		similarImages = append(similarImages, vsi.metadata[idx])
	}

	return similarImages, nil
}
func main() {
	// CLI flags
	clipServer := flag.String("server", "http://localhost:8000", "CLIP embedding server URL")
	csvPath := flag.String("csv", "clip_embeddings.csv", "Path to embeddings CSV")
	numResults := flag.Int("n", 5, "Number of results to return")
	flag.Parse()

	ctx := context.Background()

	// Initialize services
	embeddingService := NewTextEmbeddingService(*clipServer)

	// Load embeddings and create index (do this only once)
	log.Printf("Loading embeddings from %s", *csvPath)
	metadata, embeddings, err := loadEmbeddingsFromCSV(*csvPath)
	if err != nil {
		log.Fatalf("Failed to load embeddings: %v", err)
	}

	log.Printf("Building search index...")
	searchIndex, err := NewVectorSearchIndex(ctx, metadata, embeddings)
	if err != nil {
		log.Fatalf("Failed to build index: %v", err)
	}

	fmt.Println("\nIndex built successfully! Enter your queries (type 'exit' to quit):")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n> ")
		query, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}

		// Trim spaces and newline
		query = strings.TrimSpace(query)

		// Check for exit command
		if query == "exit" || query == "quit" {
			fmt.Println("Goodbye!")
			return
		}

		// Skip empty queries
		if query == "" {
			continue
		}

		// Get embedding for query
		queryEmbedding, err := embeddingService.GetEmbedding(ctx, query)
		if err != nil {
			log.Printf("Failed to get text embedding: %v", err)
			continue
		}

		// Search
		results, err := searchIndex.Search(ctx, queryEmbedding, *numResults)
		if err != nil {
			log.Printf("Search failed: %v", err)
			continue
		}

		// Print results
		fmt.Printf("\nResults for: %q\n\n", query)
		for i, img := range results {
			desc := img.Description
			if desc == "" {
				desc = "[No description]"
			}
			fmt.Printf("%d. %s\n   URL: %s\n\n", i+1, desc, img.PhotoURL)
		}
	}
}
