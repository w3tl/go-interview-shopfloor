package resource

import (
	"errors"
	"fmt"
)

var (
	ErrResourceStopped = errors.New("Resource is stopped")
)

type ResourceStatus string

const (
	ResourceStatusStopped ResourceStatus = "STOPPED"
	ResourceStatusSetup   ResourceStatus = "SETUP"
	ResourceStatusProcess ResourceStatus = "PROCESS"
)

type Resource struct {
	name       string
	status     ResourceStatus
	workingQty float64
}

func New(name string) *Resource {
	return &Resource{
		name:   name,
		status: ResourceStatusStopped,
	}
}

func (r *Resource) Status() string {
	return fmt.Sprintf("Current status: %s, quantity registered: %f", r.status, r.workingQty)
}

func (r *Resource) Name() string {
	return r.name
}

func (r *Resource) Stop() error {
	if r.status != ResourceStatusStopped {
		r.status = ResourceStatusStopped
	}

	return nil
}

func (r *Resource) Start() error {
	if r.status != ResourceStatusProcess {
		r.status = ResourceStatusProcess
	}

	return nil
}
func (r *Resource) Setup() error {
	if r.status != ResourceStatusSetup {
		r.status = ResourceStatusSetup
	}

	return nil
}
func (r *Resource) Process() error {
	if r.status != ResourceStatusProcess {
		r.status = ResourceStatusProcess
	}

	return nil
}

func (r *Resource) SetupQty(t float64) error {
	if r.status == ResourceStatusStopped {
		return ErrResourceStopped
	}

	r.workingQty = t

	return nil
}

func (r *Resource) ProcessQty(t float64) error {
	if r.status == ResourceStatusStopped {
		return ErrResourceStopped
	}

	r.workingQty += t

	return nil
}
