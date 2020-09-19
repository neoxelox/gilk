package gilk

import (
	gocontext "context"
	"encoding/json"
	"html/template"
	"net/http"
	"runtime"
	"time"

	"github.com/neoxelox/gilk/deque"
)

// ModeType represents the different Modes of Gilk
type modeType string

const (
	// Enabled makes contexts to be cached and served
	Enabled modeType = "Enabled"
	// Disabled makes contexts not cached nor served
	Disabled modeType = "Disabled"
)

var (
	// Mode describes the current mode Gilk is on
	Mode modeType = Enabled
)

var (
	// CacheCapacity describes the capacity of the context cache
	CacheCapacity int = 50

	cache *deque.Deque = deque.NewCapped(CacheCapacity)
)

// Reset allows to override all default configuration
func Reset() {
	cache = deque.NewCapped(CacheCapacity)
}

var (
	// SkippedStackFrames describes the number of stack frames to
	// be skipped when the caller of the query context is captured
	SkippedStackFrames = 1
)

type contextKeyType string

const (
	contextKey contextKeyType = "GilkContextKey"
)

var (
	templates *template.Template = template.Must(template.ParseGlob("./templates/*.tpl"))
)

// Query describes a query context
type query struct {
	Query      string        `json:"query"`
	Args       []interface{} `json:"args"`
	CallerFile string        `json:"caller_file"`
	CallerFunc string        `json:"caller_func"`
	CallerLine int           `json:"caller_line"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
}

// Context describes an scoped context
type context struct {
	Path      string    `json:"path"`
	Queries   []query   `json:"queries"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// NewContext creates and caches a new Context of the executed scope
func NewContext(parent *gocontext.Context, path string) func() {
	if Mode != Enabled || parent == nil {
		return func() {}
	}

	gilkContext := &context{
		Path:      path,
		StartTime: time.Now(),
	}

	cache.Prepend(gilkContext)

	*parent = gocontext.WithValue(*parent, contextKey, gilkContext)

	return func() {
		gilkContext.EndTime = time.Now()
	}
}

// NewQuery creates and caches a new Query to the Context of the executed scope
func NewQuery(ctx gocontext.Context, sql string, args ...interface{}) func() {
	if Mode != Enabled {
		return func() {}
	}

	file := ""
	function := ""
	line := -1

	if pc, fl, ln, ok := runtime.Caller(SkippedStackFrames); ok {
		file = fl
		function = runtime.FuncForPC(pc).Name()
		line = ln
	}

	gilkQuery := &query{
		Query:      sql,
		Args:       args,
		CallerFile: file,
		CallerFunc: function,
		CallerLine: line,
		StartTime:  time.Now(),
	}

	return func() {
		gilkQuery.EndTime = time.Now()

		gilkContext, ok := ctx.Value(contextKey).(*context)
		if !ok || gilkContext == nil {
			return
		}

		gilkContext.Queries = append(gilkContext.Queries, *gilkQuery)
	}
}

func getRendered(w http.ResponseWriter, r *http.Request) {
	var response []*context

	for iter := cache.IterFirst(); iter != nil; iter = iter.Next() {
		if gilkContext, ok := iter.Value.(*context); ok {
			response = append(response, gilkContext)
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(w, "index.tpl", response)
}

// Serve serves an HTML page with cache's Contexts
func Serve(addr string) error {
	if Mode != Enabled {
		return nil
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", getRendered)
	return http.ListenAndServe(addr, nil)
}

func getRaw(w http.ResponseWriter, r *http.Request) {
	var response []*context

	for iter := cache.IterFirst(); iter != nil; iter = iter.Next() {
		if gilkContext, ok := iter.Value.(*context); ok {
			response = append(response, gilkContext)
		}
	}

	serializedResponse, _ := json.MarshalIndent(response, "", "  ")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(serializedResponse)
}

// ServeRaw serves a plain JSON page with cache's Contexts
func ServeRaw(addr string) error {
	if Mode != Enabled {
		return nil
	}

	http.HandleFunc("/", getRaw)
	return http.ListenAndServe(addr, nil)
}
