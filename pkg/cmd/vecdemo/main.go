package main

import (
	"bufio"
	"bytes"
	"context"
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

	"github.com/jackc/pgx/v4"
)

type ImageMetadata struct {
	PhotoID     string
	PhotoURL    string
	Description string
	Embedding   []float32
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

// Setup database tables and indexes
const setupSQL = `
CREATE TABLE IF NOT EXISTS image_embeddings (
    id UUID default gen_random_uuid() PRIMARY KEY,
    photo_id TEXT NOT NULL,
    photo_url TEXT NOT NULL,
    description TEXT,
    embedding vector(512) NOT NULL,
	VECTOR INDEX (embedding)
);`

func setupDatabase(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, setupSQL)
	return err
}

func insertEmbeddings(ctx context.Context, conn *pgx.Conn, embeddings []ImageMetadata) error {
	batchSize := 100 // Process in smaller batches
	totalEmbeddings := len(embeddings)

	for i := 0; i < totalEmbeddings; i += batchSize {
		end := i + batchSize
		if end > totalEmbeddings {
			end = totalEmbeddings
		}

		batch := &pgx.Batch{}
		for _, img := range embeddings[i:end] {
			vectorStr := fmt.Sprintf("[%s]", formatVector(img.Embedding))
			batch.Queue(`
                INSERT INTO image_embeddings (photo_id, photo_url, description, embedding)
                VALUES ($1, $2, $3, $4);
            `, img.PhotoID, img.PhotoURL, img.Description, vectorStr)
		}

		log.Printf("Inserting batch %d-%d of %d embeddings...", i+1, end, totalEmbeddings)
		br := conn.SendBatch(ctx, batch)
		if err := br.Close(); err != nil {
			return fmt.Errorf("failed to insert batch %d-%d: %w", i+1, end, err)
		}
	}

	return nil
}

func searchSimilarImages(ctx context.Context, conn *pgx.Conn, queryEmbedding []float32, limit int) ([]ImageMetadata, error) {
	vectorStr := fmt.Sprintf("[%s]", formatVector(queryEmbedding))

	// Log the complete SQL with actual values
	fullSQL := fmt.Sprintf(`
        SELECT photo_id, photo_url, description
        FROM image_embeddings AS OF SYSTEM TIME '-1s'
        ORDER BY embedding <-> '%s'
        LIMIT %d;
    `, vectorStr, limit)

	log.Printf("Full SQL:\n%s", fullSQL)

	rows, err := conn.Query(ctx, `SELECT photo_id, photo_url, description FROM image_embeddings ORDER BY embedding <-> $1 LIMIT $2;`, vectorStr, limit)
	if err != nil {
		log.Printf("ERROR executing query: %v", err)
		panic(err)
	}
	defer rows.Close()

	var results []ImageMetadata
	for rows.Next() {
		var img ImageMetadata
		if err := rows.Scan(&img.PhotoID, &img.PhotoURL, &img.Description); err != nil {
			log.Printf("ERROR scanning row: %v", err)
			return nil, err
		}
		results = append(results, img)
	}

	log.Printf("Search completed. Found %d results", len(results))

	if len(results) == 0 {
		var count int
		err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM image_embeddings").Scan(&count)
		if err != nil {
			log.Printf("ERROR checking table count: %v", err)
		} else {
			log.Printf("Total rows in image_embeddings table: %d", count)
		}
	}

	return results, nil
}
func formatVector(v []float32) string {
	strVals := make([]string, len(v))
	for i, val := range v {
		strVals[i] = fmt.Sprintf("%f", val)
	}
	return strings.Join(strVals, ",")
}

func loadEmbeddingsFromCSV(filename string) ([]ImageMetadata, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var embeddings []ImageMetadata
	lineNum := 1

	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error on line %d: %w", lineNum, err)
		}

		if len(record) < 4 {
			return nil, fmt.Errorf("record on line %d: not enough fields", lineNum)
		}

		// Parse embedding string from CSV
		embStr := strings.Trim(record[3], "[]")
		embParts := strings.Split(embStr, ",")
		embedding := make([]float32, len(embParts))

		for i, val := range embParts {
			f, err := strconv.ParseFloat(strings.TrimSpace(val), 32)
			if err != nil {
				return nil, fmt.Errorf("parsing embedding on line %d: %w", lineNum, err)
			}
			embedding[i] = float32(f)
		}

		embeddings = append(embeddings, ImageMetadata{
			PhotoID:     record[0],
			PhotoURL:    record[1],
			Description: record[2],
			Embedding:   embedding,
		})
	}

	return embeddings, nil
}

func main() {
	clipServer := flag.String("server", "http://localhost:8000", "CLIP embedding server URL")
	csvPath := flag.String("csv", "clip_embeddings.csv", "Path to embeddings CSV")
	numResults := flag.Int("n", 5, "Number of results to return")
	dbURL := flag.String("db", "postgres://root@127.0.0.1:29001/defaultdb?sslmode=disable", "CockroachDB connection string")
	flag.Parse()

	ctx := context.Background()

	// Connect to CockroachDB
	conn, err := pgx.Connect(ctx, *dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	// Setup database
	if err := setupDatabase(ctx, conn); err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	embeddingService := NewTextEmbeddingService(*clipServer)

	log.Printf("Loading embeddings from %s", *csvPath)
	embeddings, err := loadEmbeddingsFromCSV(*csvPath)
	if err != nil {
		log.Fatalf("Failed to load embeddings: %v", err)
	}

	log.Printf("Inserting embeddings into database...")
	if err := insertEmbeddings(ctx, conn, embeddings); err != nil {
		log.Fatalf("Failed to insert embeddings: %v", err)
	}

	fmt.Println("\nDatabase ready! Enter your queries (type 'exit' to quit):")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n> ")
		query, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}

		query = strings.TrimSpace(query)
		if query == "exit" || query == "quit" {
			fmt.Println("Goodbye!")
			return
		}
		if query == "" {
			continue
		}

		queryEmbedding, err := embeddingService.GetEmbedding(ctx, query)
		if err != nil {
			log.Printf("Failed to get text embedding: %v", err)
			continue
		}

		results, err := searchSimilarImages(ctx, conn, queryEmbedding, *numResults)
		if err != nil {
			log.Printf("Search failed: %v", err)
			continue
		}

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
