package main

import (
	"encoding/json"
	"net/http"
	"time"

	"tubes2/src/backend/algorithms"
	"tubes2/src/backend/parser"
)

type Request struct {
	URL         string `json:"url"`
	HtmlContent string `json:"html_content"`
	Selector    string `json:"selector"`
	Algo        string `json:"algo"`
	Limit       int    `json:"limit"`
}

type Response struct {
	Results         []*parser.Node `json:"results"`
	TraversalLog    []string       `json:"traversal_log"`
	NodeCount       int            `json:"node_count"`
	ExecutionTimeMs int64          `json:"execution_time_ms"`
	MaxDepth        int            `json:"max_depth"`
	RootTree        *parser.Node   `json:"root_tree"`
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var root *parser.Node
	var err error

	if req.URL != "" {
		root, err = parser.ParseHTML(req.URL)
	} else if req.HtmlContent != "" {
		root, err = parser.ParseHTMLText(req.HtmlContent)
	} else {
		http.Error(w, "URL or HTML Content is required", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Error parsing content", http.StatusInternalServerError)
		return
	}

	start := time.Now()

	var results []*parser.Node
	var log []string
	var count int

	switch req.Algo {
	case "BFS":
		results, log, count = algorithms.BFS(root, req.Selector, req.Limit)
	case "DFS":
		results, log, count = algorithms.DFS(root, req.Selector, req.Limit)
	default:
		http.Error(w, "Bad Request: Invalid algorithm", http.StatusBadRequest)
		return
	}

	execTime := time.Since(start).Milliseconds()
	maxDepth := parser.MaxDepth(root)

	resp := Response{
		Results:         results,
		TraversalLog:    log,
		NodeCount:       count,
		ExecutionTimeMs: execTime,
		MaxDepth:        maxDepth,
		RootTree:        root,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/api/search", corsMiddleware(searchHandler))
	http.ListenAndServe(":8080", nil)
}
