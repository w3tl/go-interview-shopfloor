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
	ResourceStatusWorking ResourceStatus = "WORKING"
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
	if r.status != ResourceStatusWorking {
		r.status = ResourceStatusWorking
	}

	return nil
}

func (r *Resource) RegisterQty(t float64) error {
	if r.status == ResourceStatusStopped {
		return ErrResourceStopped
	}

	r.workingQty += t

	return nil
}
