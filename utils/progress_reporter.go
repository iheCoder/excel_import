package util

import (
	"fmt"
	"sync"
	"time"
)

type ProgressMode int

const (
	ProgressModeNormal ProgressMode = iota
	ProgressModeDynamic
)

type ProgressReporter struct {
	total, progress      int
	enable               bool
	startTime            time.Time
	mu                   *sync.Mutex
	mode                 ProgressMode
	totalChangeCompleted bool
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
func (p *ProgressReporter) CommitProgress(delta int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.progress += delta
	if p.enable {
		fmt.Printf("\rProgress: %d/%d, Cost: %v \n", p.progress, p.total, time.Since(p.startTime))
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
