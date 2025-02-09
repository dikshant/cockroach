package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

//go:embed templates/index.html
var templateFS embed.FS

type ImageMetadata struct {
	PhotoID     string    `json:"photo_id"`
	PhotoURL    string    `json:"photo_url"`
	Description string    `json:"description"`
	Embedding   []float32 `json:"-"`
}

type SearchResponse struct {
	Results     []ImageMetadata `json:"results"`
	Query       string          `json:"query"`
	SQL         string          `json:"sql"`
	TotalCount  int             `json:"total_count"`
	QueryTimeMs int64           `json:"query_time_ms"`
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

		// Build multi-value INSERT statement
		var queryBuilder strings.Builder
		queryBuilder.WriteString(`
            INSERT INTO image_embeddings (photo_id, photo_url, description, embedding)
            VALUES 
        `)

		// Create parameter placeholders and collect values
		var values []interface{}
		paramOffset := 1

		for j, img := range embeddings[i:end] {
			if j > 0 {
				queryBuilder.WriteString(",")
			}
			queryBuilder.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d)",
				paramOffset, paramOffset+1, paramOffset+2, paramOffset+3))

			vectorStr := fmt.Sprintf("[%s]", formatVector(img.Embedding))
			values = append(values, img.PhotoID, img.PhotoURL, img.Description, vectorStr)
			paramOffset += 4
		}

		log.Printf("Inserting batch %d-%d of %d embeddings...", i+1, end, totalEmbeddings)
		if _, err := conn.Exec(ctx, queryBuilder.String(), values...); err != nil {
			return fmt.Errorf("failed to insert batch %d-%d: %w", i+1, end, err)
		}
	}

	return nil
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

func checkExistingEmbeddings(ctx context.Context, conn *pgx.Conn) (bool, error) {
	var count int
	err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM image_embeddings").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Add this near the top of your main.go file, after the imports
var templateFuncs = template.FuncMap{
	"safeHTML": func(s string) template.HTML {
		return template.HTML(s)
	},
	"safeJS": func(s string) template.JS {
		return template.JS(s)
	},
	"safeURL": func(s string) template.URL {
		return template.URL(s)
	},
}

// Then update the setupHandlers function to use the template functions:
func setupHandlers(conn *pgx.Conn, embeddingService *TextEmbeddingService, numResults int) {
	// Parse template with functions
	tmpl := template.Must(template.New("index.html").Funcs(templateFuncs).ParseFS(templateFS, "templates/index.html"))

	// Rest of the handler setup remains the same...
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})
	// Update the search handler in setupHandlers
	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Query string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		// Get total count of embeddings
		var totalCount int
		err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM image_embeddings").Scan(&totalCount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		queryEmbedding, err := embeddingService.GetEmbedding(ctx, request.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		vectorStr := fmt.Sprintf("[%s]", formatVector(queryEmbedding))
		sqlQuery := fmt.Sprintf(`
        SELECT photo_id, photo_url, description 
        FROM image_embeddings 
        ORDER BY embedding <-> '%s' 
        LIMIT %d`, vectorStr, numResults)

		// Create a display version of the query with truncated vector
		displayQuery := strings.Replace(sqlQuery, vectorStr, "[...vector...]", 1)

		startTime := time.Now()
		rows, err := conn.Query(ctx, sqlQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []ImageMetadata
		for rows.Next() {
			var img ImageMetadata
			if err := rows.Scan(&img.PhotoID, &img.PhotoURL, &img.Description); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			results = append(results, img)
		}
		queryTime := time.Since(startTime)

		response := SearchResponse{
			Results:     results,
			Query:       request.Query,
			SQL:         displayQuery,
			TotalCount:  totalCount,
			QueryTimeMs: queryTime.Milliseconds(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}

func main() {
	clipServer := flag.String("server", "http://localhost:8000", "CLIP embedding server URL")
	csvPath := flag.String("csv", "clip_embeddings.csv", "Path to embeddings CSV")
	numResults := flag.Int("n", 6, "Number of results to return")
	dbURL := flag.String("db", "postgres://root@127.0.0.1:29000/defaultdb?sslmode=disable", "CockroachDB connection string")
	port := flag.String("port", "8080", "Web server port")
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

	// Check if we need to load embeddings
	hasEmbeddings, err := checkExistingEmbeddings(ctx, conn)
	if err != nil {
		log.Fatalf("Failed to check existing embeddings: %v", err)
	}

	embeddingService := NewTextEmbeddingService(*clipServer)

	if !hasEmbeddings {
		log.Printf("Loading embeddings from %s", *csvPath)
		embeddings, err := loadEmbeddingsFromCSV(*csvPath)
		if err != nil {
			log.Fatalf("Failed to load embeddings: %v", err)
		}

		log.Printf("Inserting embeddings into database...")
		if err := insertEmbeddings(ctx, conn, embeddings); err != nil {
			log.Fatalf("Failed to insert embeddings: %v", err)
		}
	} else {
		log.Printf("Using existing embeddings from database")
	}

	setupHandlers(conn, embeddingService, *numResults)

	log.Printf("Server starting on http://localhost:%s", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal(err)
	}
}
