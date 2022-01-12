// Copyright 2018-2022 (c) The Go Service Components authors. All rights reserved. Issued under the Apache 2.0 License.

package server // import "github.com/karlmutch/go-service/pkg/server"

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/dustin/go-humanize"
	"github.com/go-stack/stack"
	"github.com/karlmutch/kv"
)

// Resource describes the needed resources for a runner task in a data structure that can be
// marshalled as json
//
type Resource struct {
	Cpus     uint   `json:"cpus"`
	Gpus     uint   `json:"gpus"`
	Hdd      string `json:"hdd"`
	Ram      string `json:"ram"`
	GpuMem   string `json:"gpuMem"`
	GpuCount uint   `json:"gpuCount"`
}

func (rsc Resource) String() (serialized string) {
	serialize, _ := json.Marshal(rsc)

	return string(serialize)
}

// Fit determines is a supplied resource description acting as a request can
// be satisfied by the receiver resource
//
// Asks the question does rsc fit into r
//
func (rsc *Resource) Fit(r *Resource) (didFit bool, err kv.Error) {

	lRam, errGo := humanize.ParseBytes(rsc.Ram)
	if errGo != nil {
		return false, kv.NewError("left side RAM could not be parsed").With("stack", stack.Trace().TrimRuntime())
	}

	rRam, errGo := humanize.ParseBytes(r.Ram)
	if errGo != nil {
		return false, kv.NewError("right side RAM could not be parsed").With("stack", stack.Trace().TrimRuntime())
	}

	lHdd, errGo := humanize.ParseBytes(rsc.Hdd)
	if errGo != nil {
		return false, kv.NewError("left side Hdd could not be parsed").With("stack", stack.Trace().TrimRuntime())
	}

	rHdd, errGo := humanize.ParseBytes(r.Hdd)
	if errGo != nil {
		return false, kv.NewError("right side Hdd could not be parsed").With("stack", stack.Trace().TrimRuntime())
	}

	lGpuMem, errGo := humanize.ParseBytes(rsc.GpuMem)
	// GpuMem is optional so handle the case when it does not parse and is empty
	if len(rsc.GpuMem) != 0 {
		if errGo != nil {
			return false, kv.NewError("left side gpuMem could not be parsed").With("left_mem", rsc.GpuMem).With("stack", stack.Trace().TrimRuntime())
		}
	}

	rGpuMem, errGo := humanize.ParseBytes(r.GpuMem)
	// GpuMem is optional so handle the case when it does not parse and is empty
	if len(r.GpuMem) != 0 {
		if errGo != nil {
			return false, kv.NewError("right side gpuMem could not be parsed").With("right", r.GpuMem).With("stack", stack.Trace().TrimRuntime())
		}
	}

	return rsc.Cpus <= r.Cpus && rsc.Gpus <= r.Gpus && rsc.GpuCount <= r.GpuCount && lHdd <= rHdd && lRam <= rRam && lGpuMem <= rGpuMem, nil
}

// Clone will deep copy a resource and return the copy
//
func (rsc *Resource) Clone() (r *Resource) {

	mod := bytes.Buffer{}
	enc := gob.NewEncoder(&mod)
	dec := gob.NewDecoder(&mod)

	if err := enc.Encode(rsc); err != nil {
		return nil
	}

	r = &Resource{}
	if err := dec.Decode(r); err != nil {
		return nil
	}
	return r
}
