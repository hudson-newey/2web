package debugjson

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"hudson-newey/2web/src/cli"
	jsonFile "hudson-newey/2web/src/content/json"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/debugger"
	"hudson-newey/2web/src/filesystem"
)

// debugFileVersion is bumped whenever the debug document shape changes. A
// previous debug file with a different version is discarded entirely (its
// reactive graphs were written by an older compiler and can be missing
// fields), so cached pages of the next build simply start without a graph
// until they are recompiled.
const debugFileVersion = 2

// assetInfo describes a single file in the compiled build output.
type assetInfo struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

// debugDocument is the shape of the __2web.debug.json file that is written to
// the compiled output for development builds. The browser devtools extension
// fetches this file to power its assets and reactivity panels.
type debugDocument struct {
	Version    int                      `json:"version"`
	Pages      []string                 `json:"pages"`
	Assets     []assetInfo              `json:"assets"`
	Messages   []string                 `json:"messages"`
	Reactivity []debugger.PageDebugInfo `json:"reactivity"`
}

// GenerateDebugJson writes the __2web.debug.json file to the compiled build
// output.
//
// The site paths are passed in by the caller (site.AfterAll) because this
// package can't import the site package (site.AfterAll calls into this
// package).
func GenerateDebugJson(sitePaths []string) {
	// Generated site assets (e.g. sitemap.xml and robots.txt) are written
	// through the asynchronous file writer. The asset walk below lists files
	// from disk, so the write queue must be drained first or the asset list
	// would nondeterministically miss files that are still queued.
	filesystem.WaitFileWriter()

	file := jsonFile.JsonFile{}
	file.AddContent(debuggerInfoString(sitePaths))

	pageModel := page.Page{
		Html: &file,
	}

	outPath := cli.GetEnvVars().DebugOverride
	pageModel.WriteHtml(outPath)
}

func debuggerInfoString(sitePaths []string) string {
	document := debugDocument{
		Version:    debugFileVersion,
		Pages:      pages(sitePaths),
		Assets:     assets(),
		Messages:   messages(),
		Reactivity: mergedReactivity(),
	}

	serialized, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return "{}"
	}

	return string(serialized)
}

// mergedReactivity merges this build's reactive graphs with the reactive
// graphs from the previous debug file.
//
// Pages that were served from the build cache don't get recompiled, so they
// don't contribute a graph to this run. Because a cached page's compile time
// inputs are unchanged, its previous graph is still accurate and is carried
// over so that the debug file always describes the whole site.
//
// Carried over graphs are pruned to the pages of the current build: graphs
// for pages that were deleted from the site (or that belong to a different
// input directory, e.g. after switching the dev workflow) must not survive
// forever.
func mergedReactivity() []debugger.PageDebugInfo {
	previous := readPreviousReactivity()
	merged := map[string]debugger.PageDebugInfo{}
	for _, graph := range previous {
		if debugger.HasCurrentBuildPage(graph.Page) {
			merged[graph.Page] = graph
		}
	}

	for _, graph := range debugger.GetPageDebugInfo() {
		merged[graph.Page] = graph
	}

	graphs := []debugger.PageDebugInfo{}
	for _, graph := range merged {
		graphs = append(graphs, graph)
	}

	sort.Slice(graphs, func(i, j int) bool {
		return graphs[i].Page < graphs[j].Page
	})

	return graphs
}

// readPreviousReactivity loads the reactive graphs from the previous build's
// debug file. Missing or corrupt debug files are treated as no previous data.
func readPreviousReactivity() []debugger.PageDebugInfo {
	raw, err := os.ReadFile(cli.GetEnvVars().DebugOverride)
	if err != nil {
		return []debugger.PageDebugInfo{}
	}

	var previous debugDocument
	if err := json.Unmarshal(raw, &previous); err != nil {
		return []debugger.PageDebugInfo{}
	}

	if previous.Version != debugFileVersion {
		return []debugger.PageDebugInfo{}
	}

	return previous.Reactivity
}

// pages returns the output paths of every page in the build, relative to the
// build output directory.
func pages(sitePaths []string) []string {
	outputRoot := outputDirectory()

	paths := []string{}
	for _, outputPath := range sitePaths {
		paths = append(paths, relativeToOutput(outputRoot, outputPath))
	}

	sort.Strings(paths)
	return paths
}

// assets walks the compiled build output and describes every file that was
// written (pages, stylesheets, scripts, and passthrough assets).
//
// The output directory is used as the source of truth instead of the compiler
// buffers so that the asset list stays complete even when pages are served
// from the build cache.
func assets() []assetInfo {
	outputRoot := outputDirectory()
	if outputRoot == "" {
		return []assetInfo{}
	}

	assets := []assetInfo{}
	debugFileName := filepath.Base(cli.GetEnvVars().DebugOverride)

	walkErr := filepath.WalkDir(outputRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if entry.IsDir() {
			return nil
		}

		// Don't list the debug file inside itself.
		if entry.Name() == debugFileName {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return nil
		}

		assets = append(assets, assetInfo{
			Path: relativeToOutput(outputRoot, path),
			Type: assetType(path),
			Size: info.Size(),
		})

		return nil
	})
	if walkErr != nil {
		return []assetInfo{}
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].Path < assets[j].Path
	})

	return assets
}

func messages() []string {
	debugInfos := debugger.GetDebugInfo()

	messages := []string{}
	for _, info := range debugInfos {
		messages = append(messages, info.Message)
	}

	return messages
}

// outputDirectory is the directory that the debug file is written to, which is
// the root of the compiled build output.
func outputDirectory() string {
	return filepath.Dir(cli.GetEnvVars().DebugOverride)
}

func relativeToOutput(outputRoot string, path string) string {
	relative, err := filepath.Rel(outputRoot, path)
	if err != nil {
		return filepath.ToSlash(path)
	}

	return filepath.ToSlash(relative)
}

func assetType(path string) string {
	extension := strings.ToLower(filepath.Ext(path))

	switch extension {
	case ".html", ".htm":
		return "page"
	case ".css":
		return "style"
	case ".js", ".mjs", ".ts":
		return "script"
	case ".svg":
		return "image"
	case ".json":
		return "data"
	}

	if extension != "" {
		return strings.TrimPrefix(extension, ".")
	}

	return "unknown"
}
