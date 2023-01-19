package pprof

import (
	"fmt"
	"github.com/pyroscope-io/godeltaprof"
	"io"
	"net/http"
	"runtime"
	"strconv"
)

var deltaHeapProfiler = godeltaprof.NewHeapProfiler()
var deltaBlockProfiler = godeltaprof.NewBlockProfiler()
var deltaMutexProfiler = godeltaprof.NewMutexProfiler()

type deltaProfiler interface {
	Profile(w io.Writer) error
}

func init() {
	http.HandleFunc("/debug/pprof/delta_heap", Heap)
	http.HandleFunc("/debug/pprof/delta_block", Block)
	http.HandleFunc("/debug/pprof/delta_mutex", Mutex)
}

func Heap(w http.ResponseWriter, r *http.Request) {
	gc, _ := strconv.Atoi(r.FormValue("gc"))
	if gc > 0 {
		runtime.GC()
	}
	writeDeltaProfile(deltaHeapProfiler, "heap", w)
}

func Block(w http.ResponseWriter, r *http.Request) {
	writeDeltaProfile(deltaBlockProfiler, "block", w)
}

func Mutex(w http.ResponseWriter, r *http.Request) {
	writeDeltaProfile(deltaMutexProfiler, "mutex", w)
}

func writeDeltaProfile(p deltaProfiler, name string, w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pprof.gz"`, name))
	_ = p.Profile(w)
}
