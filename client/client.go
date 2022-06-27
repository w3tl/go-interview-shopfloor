package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func getStatus(url, resource string) (string, error) {
	var status string
	resp, err := http.Get(fmt.Sprintf("%s/status?resource=%s", url, resource))
	if err != nil {
		return status, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	status = string(body)

	return status, nil
}

func setupResource(url, resource string) error {
	urlParams := fmt.Sprintf("%s/setup?resource=%s", url, resource)
	resp, err := http.Post(urlParams, "application/json", bytes.NewReader([]byte{}))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func startResource(url, resource string) error {
	urlParams := fmt.Sprintf("%s/start?resource=%s", url, resource)
	resp, err := http.Post(urlParams, "application/json", bytes.NewReader([]byte{}))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func stopResource(url, resource string) error {
	urlParams := fmt.Sprintf("%s/stop?resource=%s", url, resource)
	resp, err := http.Post(urlParams, "application/json", bytes.NewReader([]byte{}))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func registerQty(url, resource string, qty float64) error {
	urlParams := fmt.Sprintf("%s/registerQty?resource=%s&qty=%f", url, resource, qty)
	resp, err := http.Post(urlParams, "application/json", bytes.NewReader([]byte{}))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func Run(ctx context.Context, clientid, hubaddr string, resname string) error {
	go func() {
		startTicker := time.After(time.Second)
		<-startTicker
		err := startResource(hubaddr, resname)
		if err != nil {
			fmt.Printf("%s (%s): couldn't start resource\n", clientid, resname)
		}
	}()
	go func() {
		setupTicker := time.After(time.Millisecond * 500)
		<-setupTicker
		err := setupResource(hubaddr, resname)
		if err != nil {
			fmt.Printf("%s (%s): couldn't setup resource\n", clientid, resname)
		}
	}()

	statusTicker := time.After(time.Millisecond * 1000)
	processTimer := time.Duration(rand.Intn(2000)+500) * time.Millisecond
	processTicker := time.After(processTimer)
	stopTicker := time.After(time.Second * 3)

	var totalRegistered float64

	for {

		select {

		case <-statusTicker:
			status, err := getStatus(hubaddr, resname)
			if err != nil {
				fmt.Printf("%s (%s): couldn't get status: %v\n", clientid, resname, err)
			} else {
				fmt.Printf("%s (%s): %s\n", clientid, resname, status)
			}
			statusTicker = time.After(time.Millisecond * 500)
		case <-processTicker:
			qty := rand.Float64() * 100
			totalRegistered += qty
			fmt.Printf("%s (%s): Register %f\n", clientid, resname, qty)
			err := registerQty(hubaddr, resname, qty)
			if err != nil {
				return err
			}
			processTicker = time.After(processTimer)
		case <-stopTicker:
			fmt.Printf("%s (%s): Totally registered %f\n", clientid, resname, totalRegistered)
			err := stopResource(hubaddr, resname)
			if err != nil {
				fmt.Printf("%s: couldn't stop %s\n", clientid, resname)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
