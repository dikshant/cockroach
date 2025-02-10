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
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
)

//go:embed templates/index.html
var templateFS embed.FS

type ImageMetadata struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	ImageURL  string    `json:"image_url"`
	Embedding []float32 `json:"-"`
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
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    filename TEXT NOT NULL,
    embedding vector(512) NOT NULL,
    VECTOR INDEX (embedding)
);`

func setupDatabase(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, setupSQL)
	return err
}

func insertEmbeddings(ctx context.Context, conn *pgx.Conn, embeddings []ImageMetadata) error {
	batchSize := 100
	totalEmbeddings := len(embeddings)

	for i := 0; i < totalEmbeddings; i += batchSize {
		end := i + batchSize
		if end > totalEmbeddings {
			end = totalEmbeddings
		}

		var queryBuilder strings.Builder
		queryBuilder.WriteString(`
			INSERT INTO image_embeddings (filename, embedding)
			VALUES 
		`)

		var values []interface{}
		paramOffset := 1

		for j, img := range embeddings[i:end] {
			if j > 0 {
				queryBuilder.WriteString(",")
			}
			queryBuilder.WriteString(fmt.Sprintf("($%d, $%d)",
				paramOffset, paramOffset+1))

			vectorStr := fmt.Sprintf("[%s]", formatVector(img.Embedding))
			values = append(values, img.Filename, vectorStr)
			paramOffset += 2
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

		if len(record) != 3 {
			return nil, fmt.Errorf("record on line %d: expected 3 fields, got %d", lineNum, len(record))
		}

		// Parse embedding string from CSV
		embStr := strings.Trim(record[1], "[]")
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
			Filename:  record[0],
			Embedding: embedding,
		})
	}

	return embeddings, nil
}

// New function to load files concurrently
func loadImageFiles(basePath string, filenames []string) (map[string][]byte, error) {
	fileData := make(map[string][]byte)
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // Limit concurrent file reads

	for _, filename := range filenames {
		wg.Add(1)
		go func(fname string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			filepath := filepath.Join(basePath, fname)
			data, err := os.ReadFile(filepath)
			if err != nil {
				log.Printf("Error reading file %s: %v", fname, err)
				return
			}

			mu.Lock()
			fileData[fname] = data
			mu.Unlock()
		}(filename)
	}

	wg.Wait()
	return fileData, nil
}

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

// Add this new handler function
func setupImageHandler(imagesPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		if filename == "" {
			http.Error(w, "filename is required", http.StatusBadRequest)
			return
		}

		// Prevent directory traversal
		cleanPath := filepath.Clean(filename)
		if strings.Contains(cleanPath, "..") {
			http.Error(w, "invalid filename", http.StatusBadRequest)
			return
		}

		filePath := filepath.Join(imagesPath, cleanPath)

		// Get the content type based on file extension
		contentType := "image/jpeg"
		// Set proper headers
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Cache-Control", "public, max-age=31536000")

		http.ServeFile(w, r, filePath)
	}
}

type SearchRequest struct {
	Query    string `json:"query"`
	UseIndex bool   `json:"useIndex"`
}

func setupHandlers(conn *pgx.Conn, embeddingService *TextEmbeddingService, numResults int, imagesPath string) {
	tmpl := template.Must(template.New("index.html").Funcs(templateFuncs).ParseFS(templateFS, "templates/index.html"))
	// Add the image handler
	http.HandleFunc("/api/image", setupImageHandler(imagesPath))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})

	// In your search handler
	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request SearchRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		var totalCount int
		// Use the appropriate table based on useIndex flag
		tableName := "image_embeddings"
		if !request.UseIndex {
			tableName = "image_embeddings_no_idx"
		}

		err := conn.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&totalCount)
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
			SELECT id, filename 
			FROM %s 
			ORDER BY embedding <-> '%s' 
			LIMIT %d`, tableName, vectorStr, numResults)

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
			if err := rows.Scan(&img.ID, &img.Filename); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			results = append(results, img)
		}
		ts := time.Since(startTime)

		// Instead of loading files, just create URLs
		for i := range results {
			results[i].ImageURL = fmt.Sprintf("/api/image?filename=%s",
				url.QueryEscape(results[i].Filename))
		}

		response := SearchResponse{
			Results:     results,
			Query:       request.Query,
			SQL:         displayQuery,
			TotalCount:  totalCount,
			QueryTimeMs: ts.Milliseconds(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}

func checkExistingEmbeddings(ctx context.Context, conn *pgx.Conn) (bool, error) {
	var count int
	err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM image_embeddings").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func main() {
	clipServer := flag.String("server", "http://localhost:8000", "CLIP embedding server URL")
	csvPath := flag.String("csv", "clip_embeddings.csv", "Path to embeddings CSV")
	imagesPath := flag.String("images", "./images", "Path to images directory")
	numResults := flag.Int("n", 6, "Number of results to return")
	dbURL := flag.String("db", "postgres://root@127.0.0.1:29000/defaultdb?sslmode=disable", "CockroachDB connection string")
	port := flag.String("port", "8080", "Web server port")
	flag.Parse()

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, *dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	if err := setupDatabase(ctx, conn); err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

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

	setupHandlers(conn, embeddingService, *numResults, *imagesPath)

	log.Printf("Server starting on http://localhost:%s", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal(err)
	}
}
