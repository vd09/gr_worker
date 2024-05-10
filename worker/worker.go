package worker

type IsEligibleToStopFunc func(isCtxDone bool) bool

type Worker interface {
	Start()
}
