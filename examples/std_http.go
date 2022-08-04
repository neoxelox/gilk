package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/neoxelox/gilk"
)

func main() {
	gilk.SkippedStackFrames = 2
	go gilk.Serve(":8000")

	http.HandleFunc("/users", getUsers)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	ctx, endContext := gilk.NewContext(r.Context(), r.URL.Path, r.Method)
	defer endContext()

	r = r.WithContext(ctx)

	time.Sleep(200 * time.Millisecond)

	users := queryWithContext(ctx, `
		SELECT * FROM "users" WHERE "name" LIKE $1 AND "username" LIKE $2;`,
		"Alex", "Neoxelox")

	time.Sleep(200 * time.Millisecond)
	w.Write([]byte(users))
}

func queryWithContext(ctx context.Context, query string, args ...interface{}) string {
	ctx, endQuery := gilk.NewQuery(ctx, query, args...)
	defer endQuery()

	time.Sleep(150 * time.Millisecond)

	return "query executed"
}
