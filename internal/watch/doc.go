// Package watch provides the polling engine for stackwatch.
//
// It coordinates periodic drift detection across multiple cloud stacks
// by running each Detector concurrently on a configurable interval.
// Results are streamed over a channel so that output formatters and
// notifiers can consume them without blocking the detection loop.
//
// Typical usage:
//
//	poller := watch.NewPoller(cfg, detectors)
//	go poller.Run(ctx)
//	for result := range poller.Results() {
//		formatter.Print(result)
//	}
package watch
