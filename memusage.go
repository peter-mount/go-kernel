package kernel

import (
	"log"
	"runtime"
	"time"
)

// MemUsage is a Kernel service which will log on shutdown the duration of the
// process and how much memory it has used.
//
// To use simply include it as the first service when launching the kernel:
//
// func main() {
//   err := kernel.Launch( &kernel.MemUsage{}, &mylib.MyService{} )
//   if err != nil {
//     log.Fatal( err )
//   }
// }
//
// When the service stops then some statistics are logged showing how long the
// process has run, how much memory it's used and how often the garbage
// collector has run.
//
// Notes:
//
// The process duration time is from when the Start phase begins.
// If the kernel fails before that then no stats are generated. This is because
// services only get stopped if they have started.
//
type MemUsage struct {
	start time.Time
}

func (m *MemUsage) Name() string {
	return "MemUsage"
}

// Starts the service.
// The process duration time reported is from when this is called.
func (m *MemUsage) Start() error {
	m.start = time.Now()
	return nil
}

// Stops the service
func (m *MemUsage) Stop() {
	t := time.Now()
	elapsed := t.Sub(m.start)
	log.Println("Duration:", elapsed)

	var s runtime.MemStats
	runtime.ReadMemStats(&s)

	log.Printf(
		"MemUsage: Alloc %v MiB TotalAlloc %v MiB Sys %v MiB Mallocs %v Frees %v",
		bToMb(s.Alloc),
		bToMb(s.TotalAlloc),
		bToMb(s.Sys),
		s.Mallocs,
		s.Frees,
	)

	log.Printf(
		"Heap: Alloc %v MiB Sys %v MiB Idle %v MiB InUse %v MiB Released %v MiB Objects %v",
		bToMb(s.HeapAlloc),
		bToMb(s.HeapSys),
		bToMb(s.HeapIdle),
		bToMb(s.HeapInuse),
		bToMb(s.HeapReleased),
		bToMb(s.HeapObjects),
	)

	log.Printf(
		"GC: %v Forced %v CPU %0.6f%%",
		s.NumGC,
		s.NumForcedGC,
		s.GCCPUFraction*100,
	)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
