package server

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hudson-newey/2web-cli/src/builders/build"
	"github.com/hudson-newey/2web-cli/src/cli"
)

// TODO: Refactor this entire file because it was LLM generated
var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

//go:embed static/liveReload.js
var liveReloadScript string

var liveReloadTemplate = fmt.Sprintf(`
<script type="module">%s</script>
`, liveReloadScript)

func runDevServer(inPath string, outPath string) {
	absInPath, err := filepath.Abs(inPath)
	if err != nil {
		log.Fatalf("Failed to resolve path: %v", err)
	}

	// Check if path exists
	if _, err := os.Stat(absInPath); os.IsNotExist(err) {
		log.Fatalf("Input path does not exist: '%s'", absInPath)
	}

	absOutPath, err := filepath.Abs(outPath)
	if err != nil {
		log.Fatalf("Failed to resolve path: %v", err)
	}

	// Start file watcher in a separate goroutine
	go watchFiles(absInPath, absOutPath, outPath)

	// Setup HTTP server
	mux := http.NewServeMux()

	// WebSocket endpoint for live reload
	mux.HandleFunc("/__2web_updates", handleWebSocket)

	// Serve static files with HTML injection
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveFileWithReload(w, r, absOutPath)
	})

	port := "2000"
	addr := fmt.Sprintf(":%s", port)

	// Perform an initial build before starting the server
	buildAssets(inPath, outPath)

	fmt.Printf("🚀 2web dev server running at http://localhost:%s\n", port)
	fmt.Printf("📁 Serving files from: %s\n", absInPath)
	fmt.Println("👀 Watching for file changes...")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
	}()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func serveFileWithReload(
	w http.ResponseWriter,
	r *http.Request,
	rootPath string,
) {
	// Clean the URL path
	urlPath := r.URL.Path

	// If the url path does not include a .html extension and does not end with a
	// slash, try to serve the .html file
	if !strings.HasSuffix(urlPath, "/") && !strings.Contains(filepath.Base(urlPath), ".") {
		tryPath := urlPath + ".html"
		fullTryPath := filepath.Join(rootPath, filepath.Clean(tryPath))
		if _, err := os.Stat(fullTryPath); err == nil {
			urlPath = tryPath
		}
	}

	log.Printf("📄 Serving file: %s", urlPath)

	// Construct file path
	filePath := filepath.Join(rootPath, filepath.Clean(urlPath))

	// Security check: ensure the file is within the root path
	if !strings.HasPrefix(filePath, rootPath) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Check if it's a directory
	fileInfo, err := os.Stat(filePath)
	if err == nil && fileInfo.IsDir() {
		// Try index.html
		indexPath := filepath.Join(filePath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			filePath = indexPath
		} else {
			// List directory contents
			serveDirectory(w, filePath, urlPath)
			return
		}
	}

	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	ext := filepath.Ext(filePath)

	contentTypeMap := map[string]string{
		".html": "text/html; charset=utf-8",
		".css":  "text/css; charset=utf-8",
		".js":   "application/javascript",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".json": "application/json",
		".txt":  "text/plain; charset=utf-8",
	}

	contentType, contentTypeExists := contentTypeMap[ext]
	if contentTypeExists {
		w.Header().Set("Content-Type", contentType)
	} else {
		w.Header().Set("Content-Type", http.DetectContentType(content))
	}

	// Inject live reload script for HTML files
	if strings.Contains(contentType, "text/html") {
		content = injectLiveReload(content)
	}

	w.Write(content)
}

func serveDirectory(w http.ResponseWriter, dirPath string, urlPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		http.Error(w, "Cannot read directory", http.StatusInternalServerError)
		return
	}

	var html strings.Builder
	html.WriteString("<!DOCTYPE html><html><head><meta charset='utf-8'><title>Directory listing</title>")
	html.WriteString("<style>body{font-family:sans-serif;padding:20px;}a{display:block;padding:5px;text-decoration:none;}a:hover{background:#f0f0f0;}</style>")
	html.WriteString(liveReloadTemplate)
	html.WriteString("</head><body>")
	html.WriteString(fmt.Sprintf("<h1>Index of %s</h1><hr>", urlPath))

	if urlPath != "/" {
		html.WriteString("<a href='..'>📁 ..</a>")
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			html.WriteString(fmt.Sprintf("<a href='%s/'>📁 %s/</a>", name, name))
		} else {
			html.WriteString(fmt.Sprintf("<a href='%s'>📄 %s</a>", name, name))
		}
	}

	html.WriteString("</body></html>")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html.String()))
}

func injectLiveReload(content []byte) []byte {
	// Find the closing </head> tag and inject the script before it
	headCloseTag := []byte("</head>")

	index := bytes.Index(content, headCloseTag)
	if index == -1 {
		// If no </head> tag found, append to the bottom of the file so that we
		// don't break a <!DOCTYPE html> tag.
		return append(content, []byte(liveReloadTemplate)...)
	}

	// Insert script before </head>
	result := make([]byte, 0, len(content)+len(liveReloadTemplate))
	result = append(result, content[:index]...)
	result = append(result, []byte(liveReloadTemplate)...)
	result = append(result, content[index:]...)

	return result
}

// TODO: Remove relativeOutPath from here
// This is needed because there is a bug in the compiler where it does not work
// with absolute paths properly.
func watchFiles(inPath string, outPath string, relativeOutPath string) {
	var lastModTime time.Time

	const fileWatcherInterval = 5 * time.Millisecond
	for {
		time.Sleep(fileWatcherInterval)

		modTime := getLatestModTime(inPath)
		if !lastModTime.IsZero() && modTime.After(lastModTime) {
			handleFileChange(inPath, relativeOutPath)
		}
		lastModTime = modTime
	}
}

func handleFileChange(inPath string, outPath string) {
	cli.ClearConsole()

	buildAssets(inPath, outPath)
	fmt.Println("🔄 Asset build complete, reloading clients...")
	notifyClients()
}

func buildAssets(inPath string, outPath string) {
	fmt.Println("📦 Building assets...")
	build.BuildPath(inPath, outPath)
}

func getLatestModTime(root string) time.Time {
	var latest time.Time

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden directories and files
		if d.Name()[0] == '.' {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip node_modules and other common directories
		if d.IsDir() && (d.Name() == "node_modules" || d.Name() == ".git") {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err == nil && info.ModTime().After(latest) {
				latest = info.ModTime()
			}
		}

		return nil
	})

	return latest
}

func notifyClients() {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte("reload"))
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}
