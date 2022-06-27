package resourcehub

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"shopfloor/resource"
	"strconv"
	"time"
)

var (
	ErrResourceExists    = errors.New("Resource is already under observation")
	ErrResourceNotExists = errors.New("Resource was not found")
)

type Hub struct {
	resources []*resource.Resource
}

func NewServer() *Hub {
	return &Hub{resources: make([]*resource.Resource, 0)}
}

func (h *Hub) AddResource(r *resource.Resource) error {
	for _, res := range h.resources {
		if res.Name() == r.Name() {
			return ErrResourceExists
		}
	}

	h.resources = append(h.resources, r)

	return nil
}

func (h *Hub) RemoveResource(name string) error {
	for idx, r := range h.resources {
		if r.Name() == name {
			h.resources[idx] = h.resources[len(h.resources)-1]
			h.resources[len(h.resources)-1] = nil
			h.resources = h.resources[:len(h.resources)-1]

			return nil
		}
	}

	return ErrResourceNotExists
}

func (h *Hub) GetResource(name string) *resource.Resource {
	for _, res := range h.resources {
		if res.Name() == name {
			return res
		}
	}

	return nil
}

func (h *Hub) Listen(port int) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}

	http.Handle("/status", h.status())
	http.Handle("/start", h.start())
	http.Handle("/stop", h.stop())
	http.Handle("/setup", h.setup())
	http.Handle("/process", h.process())
	http.Handle("/registerQty", h.registerQty())

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return srv
}

func (h *Hub) status() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		rname := r.URL.Query().Get("resource")
		res := h.GetResource(rname)
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, err := io.WriteString(w, string(res.Status()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return http.HandlerFunc(fn)
}

func (h *Hub) start() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		rname := r.URL.Query().Get("resource")
		res := h.GetResource(rname)
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := res.Start(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return http.HandlerFunc(fn)
}

func (h *Hub) setup() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		rname := r.URL.Query().Get("resource")
		res := h.GetResource(rname)
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := res.Setup(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return http.HandlerFunc(fn)
}
func (h *Hub) process() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		rname := r.URL.Query().Get("resource")
		res := h.GetResource(rname)
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := res.Process(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return http.HandlerFunc(fn)
}

func (h *Hub) stop() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		rname := r.URL.Query().Get("resource")
		res := h.GetResource(rname)
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := res.Stop(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return http.HandlerFunc(fn)
}

func (h *Hub) registerQty() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		rname := r.URL.Query().Get("resource")
		res := h.GetResource(rname)
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		qtystring := r.URL.Query().Get("qty")

		qty, err := strconv.ParseFloat(qtystring, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		time.Sleep(time.Duration(rand.Intn(1000)))

		if err := res.SetupQty(qty); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return http.HandlerFunc(fn)
}
