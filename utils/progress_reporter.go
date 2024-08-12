package util

import (
	"fmt"
	"sync"
	"time"
)

type ProgressMode int
type ProgressStatus int

const (
	ProgressModeNormal ProgressMode = iota
	ProgressModeDynamic

	ProgressStatusSuccess ProgressStatus = 1
	ProgressStatusFailed  ProgressStatus = 2
)

type ProgressReporter struct {
	total, progress      int
	enable               bool
	startTime            time.Time
	mu                   *sync.Mutex
	mode                 ProgressMode
	totalChangeCompleted bool
	detail               progressDetail
}

type progressDetail struct {
	success int
	failed  int
}

func NewProgressReporter(enable bool) *ProgressReporter {
	return &ProgressReporter{
		enable: enable,
		mu:     &sync.Mutex{},
		mode:   ProgressModeNormal,
	}
}

// StartProgress start progress. must be called
func (p *ProgressReporter) StartProgress(total int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.total = total
	p.progress = 0
	p.startTime = time.Now()
	if p.enable {
		fmt.Printf("Start Progress: %d/%d \n", p.progress, p.total)
	}
}

// SetProgressMode set progress mode
func (p *ProgressReporter) SetProgressMode(mode ProgressMode) {
	p.mode = mode
}

// CommitProgress commit progress
func (p *ProgressReporter) CommitProgress(delta int, status ProgressStatus) {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch status {
	case ProgressStatusSuccess:
		p.detail.success += delta
	case ProgressStatusFailed:
		p.detail.failed += delta
	}

	p.progress += delta
	if p.enable {
		fmt.Printf("\rProgress: %d/%d, Cost: %v ms \n", p.progress, p.total, time.Since(p.startTime).Milliseconds())
	}
}

// IncreaseTotal increase total in dynamic mode
func (p *ProgressReporter) IncreaseTotal(delta int) {
	if p.mode != ProgressModeDynamic {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.total += delta
}

// SetDynamicTotalCompleted set dynamic total completed
func (p *ProgressReporter) SetDynamicTotalCompleted() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.totalChangeCompleted = true
}

// CheckProgressCompleted check progress completed
func (p *ProgressReporter) CheckProgressCompleted() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.mode == ProgressModeDynamic && !p.totalChangeCompleted {
		return false
	}

	return p.progress >= p.total
}

// Report report progress info after completed
func (p *ProgressReporter) Report() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.enable || p.total == 0 {
		return
	}

	totalCost := time.Since(p.startTime).Milliseconds()
	averageCost := totalCost / int64(p.total)
	fmt.Printf("Complete Progress: %d/%d, Cost: %v ms, Average Cost: %v ms, Success: %d, Failed: %d \n", p.progress, p.total, totalCost, averageCost, p.detail.success, p.detail.failed)
}
