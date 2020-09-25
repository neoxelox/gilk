package gilk

import (
	gocontext "context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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
	templates        *template.Template = template.Must(template.ParseGlob("./templates/*.tpl"))
	decimalTrim                         = regexp.MustCompile(`\.[0-9]*`)
	firstLineTrim                       = regexp.MustCompile(`^\s`)
	postgresReplacer                    = regexp.MustCompile(`\$[1-9]+`)
	mySQLReplacer                       = regexp.MustCompile(`\?`)
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

func (q *query) Color() string {
	elapsed := q.EndTime.Sub(q.StartTime)

	switch {
	case elapsed <= (100 * time.Millisecond):
		return "success"
	case elapsed <= (250 * time.Millisecond):
		return "warning"
	default:
		return "danger"
	}
}

func (q *query) Duration() string {
	return decimalTrim.ReplaceAllString(q.EndTime.Sub(q.StartTime).String(), "")
}

func (q *query) Format() string {
	sql := q.Query
	sql = firstLineTrim.ReplaceAllString(sql, "")

	switch {
	case postgresReplacer.MatchString(sql):
		for index, arg := range q.Args {
			sarg := fmt.Sprintf("%v", arg)
			sql = strings.Replace(sql, "$"+strconv.Itoa(index+1), sarg, 1)
		}
		return sql
	case mySQLReplacer.MatchString(sql):
		for _, arg := range q.Args {
			sarg := fmt.Sprintf("%v", arg)
			sql = strings.Replace(sql, "?", sarg, 1)
		}
		return sql
	default:
		return sql
	}
}

// Context describes an scoped context
type context struct {
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	Queries   []query   `json:"queries"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (c *context) HasFinished() bool {
	return c.EndTime.After(c.StartTime)
}

func (c *context) MethodColor() string {
	switch c.Method {
	case http.MethodGet:
		return "success"
	case http.MethodPost:
		return "info"
	case http.MethodPut, http.MethodPatch:
		return "warning"
	case http.MethodDelete:
		return "danger"
	default:
		return "light"
	}
}

func (c *context) ContextColor() string {
	elapsed := c.EndTime.Sub(c.StartTime)

	switch {
	case elapsed <= (250 * time.Millisecond):
		return "success"
	case elapsed <= (500 * time.Millisecond):
		return "warning"
	default:
		return "danger"
	}
}

func (c *context) QueriesColor() string {
	elapsed := 0 * time.Millisecond

	for _, q := range c.Queries {
		elapsed += q.EndTime.Sub(q.StartTime)
	}

	switch {
	case elapsed <= (100 * time.Millisecond):
		return "success"
	case elapsed <= (250 * time.Millisecond):
		return "warning"
	default:
		return "danger"
	}
}

func (c *context) LenQueriesColor() string {
	queries := len(c.Queries)

	switch {
	case queries <= 10:
		return "success"
	case queries <= 15:
		return "warning"
	default:
		return "danger"
	}
}

func (c *context) ContextDuration() string {
	return decimalTrim.ReplaceAllString(c.EndTime.Sub(c.StartTime).String(), "")
}

func (c *context) QueriesDuration() string {
	elapsed := 0 * time.Millisecond

	for _, q := range c.Queries {
		elapsed += q.EndTime.Sub(q.StartTime)
	}

	return decimalTrim.ReplaceAllString(elapsed.String(), "")
}

// NewContext creates and caches a new Context of the executed scope
func NewContext(parent *gocontext.Context, path string, method string) func() {
	if Mode != Enabled || parent == nil {
		return func() {}
	}

	gilkContext := &context{
		Path:      path,
		Method:    method,
		StartTime: time.Now(),
	}

	if ok := cache.Prepend(gilkContext); !ok {
		cache.Pop()
		cache.Prepend(gilkContext)
	}

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
