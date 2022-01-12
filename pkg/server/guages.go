// Copyright 2018-2022 (c) The Go Service Components authors. All rights reserved. Issued under the Apache 2.0 License.

package server // import "github.com/karlmutch/go-service/pkg/server"

import (
	"fmt"
	"os"

	"github.com/karlmutch/go-service/pkg/network"

	"github.com/dustin/go-humanize"
	"github.com/go-stack/stack"
	"github.com/karlmutch/kv" // MIT License

	"github.com/prometheus/client_golang/prometheus"
)

// This file contains a set of guages and data structures for
// exporting the current set of resource assignments to prometheus

var (
	cpuFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "resource_cpu_free_slots",
			Help: "The number of CPU slots available on a host.",
		},
		[]string{"host"},
	)
	ramFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "resource_cpu_ram_free_bytes",
			Help: "The amount of CPU accessible RAM available on a host.",
		},
		[]string{"host"},
	)
	diskFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "resource_disk_free_bytes",
			Help: "The amount of free space on the working disk of a host.",
		},
		[]string{"host"},
	)
	gpuFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "resource_gpu_free_slots",
			Help: "The the number of GPU slots available on a host.",
		},
		[]string{"host"},
	)

	host = network.GetHostName()
)

func init() {
	if errGo := prometheus.Register(cpuFree); errGo != nil {
		fmt.Fprintln(os.Stderr, kv.Wrap(errGo).With("stack", stack.Trace().TrimRuntime()))
	}
	if errGo := prometheus.Register(ramFree); errGo != nil {
		fmt.Fprintln(os.Stderr, kv.Wrap(errGo).With("stack", stack.Trace().TrimRuntime()))
	}
	if errGo := prometheus.Register(diskFree); errGo != nil {
		fmt.Fprintln(os.Stderr, kv.Wrap(errGo).With("stack", stack.Trace().TrimRuntime()))
	}
	if errGo := prometheus.Register(gpuFree); errGo != nil {
		fmt.Fprintln(os.Stderr, kv.Wrap(errGo).With("stack", stack.Trace().TrimRuntime()))
	}
}

func updateGauges(rsc *Resource) (err kv.Error) {
	ram, errGo := humanize.ParseBytes(rsc.Ram)
	if errGo != nil {
		return kv.Wrap(errGo).With("stack", stack.Trace().TrimRuntime())
	}

	hdd, errGo := humanize.ParseBytes(rsc.Hdd)
	if errGo != nil {
		return kv.Wrap(errGo).With("stack", stack.Trace().TrimRuntime())
	}

	cpuFree.With(prometheus.Labels{"host": host}).Set(float64(rsc.Cpus))
	ramFree.With(prometheus.Labels{"host": host}).Set(float64(ram))

	diskFree.With(prometheus.Labels{"host": host}).Set(float64(hdd))

	gpuFree.With(prometheus.Labels{"host": host}).Set(float64(rsc.Gpus))

	return err
}
