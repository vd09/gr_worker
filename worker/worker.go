package worker

import "github.com/vd09/gr_worker/domain"

type IsEligibleToStopFunc func(domain.WorkerStatus) bool

type Worker interface {
	Start()
}
