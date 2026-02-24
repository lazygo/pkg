package monitor

import (
	"context"
	"errors"
	"os"
	"runtime"

	"github.com/shirou/gopsutil/v4/process"
)

type SysMontor struct {
	CPUPercent   float64                 `json:"cpu_percent"`
	MemPercent   float32                 `json:"mem_percent"`
	MemInfo      *process.MemoryInfoStat `json:"mem_info"`
	NumFDs       int32                   `json:"num_fds"`
	NumGC        uint32                  `json:"num_gc"`
	NumGoroutine int                     `json:"num_goroutine"`
	NumThreads   int32                   `json:"num_threads"`
}

func ReportSysMonitor(ctx context.Context) (*SysMontor, error) {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return nil, err
	}

	var errs error

	cpuPercent, err := proc.CPUPercentWithContext(ctx)
	if err != nil {
		errs = errors.Join(errs, err)
	}
	memPercent, err := proc.MemoryPercentWithContext(ctx)
	if err != nil {
		errs = errors.Join(errs, err)
	}
	memInfo, err := proc.MemoryInfoWithContext(ctx)
	if err != nil {
		errs = errors.Join(errs, err)
	}
	numFDs, err := proc.NumFDsWithContext(ctx)
	if err != nil && runtime.GOOS != "windows" {
		errs = errors.Join(errs, err)
	}
	numThreads, err := proc.NumThreadsWithContext(ctx)
	if err != nil {
		errs = errors.Join(errs, err)
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	sys := &SysMontor{
		CPUPercent:   cpuPercent,
		MemPercent:   memPercent,
		MemInfo:      memInfo,
		NumFDs:       numFDs,
		NumGC:        m.NumGC,
		NumGoroutine: runtime.NumGoroutine(),
		NumThreads:   numThreads,
	}
	return sys, errs
}
