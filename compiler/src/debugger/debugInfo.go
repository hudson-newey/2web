package debugger

import (
	"sync"
)

// The compiler runs page compilation on a thread pool, so every buffer in this
// package is guarded by a mutex.
var debugMutex sync.Mutex

var bufferedDebugInfo []DebugInfo = []DebugInfo{}
var bufferedPageDebugInfo []PageDebugInfo = []PageDebugInfo{}

func AddDebugInfo(info DebugInfo) {
	debugMutex.Lock()
	defer debugMutex.Unlock()

	bufferedDebugInfo = append(bufferedDebugInfo, info)
}

// GetDebugInfo returns a copy of the buffered compiler messages so that
// callers can't race on the underlying slice.
func GetDebugInfo() []DebugInfo {
	debugMutex.Lock()
	defer debugMutex.Unlock()

	messages := make([]DebugInfo, len(bufferedDebugInfo))
	copy(messages, bufferedDebugInfo)
	return messages
}

// AddPageDebugInfo records the reactive graph (variables, properties, and
// events) that was compiled for a single page.
func AddPageDebugInfo(pageInfo PageDebugInfo) {
	debugMutex.Lock()
	defer debugMutex.Unlock()

	bufferedPageDebugInfo = append(bufferedPageDebugInfo, pageInfo)
}

// GetPageDebugInfo returns a copy of the buffered reactive graphs.
func GetPageDebugInfo() []PageDebugInfo {
	debugMutex.Lock()
	defer debugMutex.Unlock()

	pages := make([]PageDebugInfo, len(bufferedPageDebugInfo))
	copy(pages, bufferedPageDebugInfo)
	return pages
}

type DebugInfo struct {
	Message string
}

// VariableInfo describes a reactive variable declaration
// (e.g. `$count = 0;` inside a <script compiled> block).
//
// ReactivityClass mirrors the compiler's reactivity levels:
// "unused", "static", "static-property", "assignment", "reactive".
type VariableInfo struct {
	Name            string   `json:"name"`
	InitialValue    string   `json:"initialValue"`
	ReactivityClass string   `json:"reactivity"`
	DependsOn       []string `json:"dependsOn"`
}

// PropertyInfo describes a reactive property binding
// (e.g. *hidden="$isOpen" or {{ $count }}).
//
// Dependencies lists the reactive variables that the property's reducer
// references.
type PropertyInfo struct {
	PropName     string   `json:"prop"`
	Reducer      string   `json:"reducer"`
	Dependencies []string `json:"dependencies"`
}

// EventInfo describes a reactive event binding
// (e.g. @click="$count = $count + 1").
//
// Sink is the variable that the event assigns to (empty when the event doesn't
// assign to a declared reactive variable) and Dependencies lists the reactive
// variables that the event reducer reads.
type EventInfo struct {
	EventName    string   `json:"event"`
	Reducer      string   `json:"reducer"`
	Sink         string   `json:"sink"`
	Dependencies []string `json:"dependencies"`
}

// PageDebugInfo is the reactive graph for a single compiled page.
type PageDebugInfo struct {
	Page       string         `json:"page"`
	Variables  []VariableInfo `json:"variables"`
	Properties []PropertyInfo `json:"properties"`
	Events     []EventInfo    `json:"events"`
}
