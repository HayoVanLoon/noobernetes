package main

import (
	"context"
	"fmt"
	pb "github.com/HayoVanLoon/go-generated/noobernetes/v1"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	counterHost = "counter-service"
	counterPort = 8080
	port        = 8080
)

func getConn() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", counterHost, counterPort), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}
	return conn, nil
}

func putTick(m string) error {
	r := &pb.PutTickRequest{Message: m}

	conn, err := getConn()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Panicf("error closing connection: %v", err)
		}
	}()

	c := pb.NewCounterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resp, err := c.PutTick(ctx, r)

	log.Printf("%v\n", resp)

	return err
}

func getTicks() (m string, err error) {
	r := &pb.GetTicksRequest{}

	conn, err := getConn()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Panicf("error closing connection: %v", err)
		}
	}()

	c := pb.NewCounterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resp, err := c.GetTicks(ctx, r)

	return fmt.Sprintf("%v / %v", resp.Ticks, resp.Tocks), nil
}

func (h *handler) loop() error {
	for {
		time.Sleep(1 * time.Second)
		h.mux.Lock()
		if h.bother() {
			var m string
			if rand.Intn(10) >= 1 {
				m = "tick"
			} else {
				m = "tock"
			}

			if err := putTick(m); err != nil {
				fmt.Println(err)
				return err
			}
		}
		h.mux.Unlock()
	}
}

type handler struct {
	bothering bool
	mux       sync.Mutex
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var m string

	if r.RequestURI == "/start" {
		m = "started"
		h.setBothering(true)
	} else if r.RequestURI == "/stop" {
		m = "paused"
		h.setBothering(false)
	} else if r.RequestURI == "/kill" {
		os.Exit(0)
	} else if r.RequestURI == "/count" {
		m2, err := getTicks()
		m = m2
		if err != nil {
			log.Print(err)
		}
	} else {
		m = "hic sunt dracones: " + r.RequestURI
	}

	log.Printf(m)
	if _, err := w.Write([]byte(m)); err != nil {
		log.Print(err)
	}
}

func (h *handler) setBothering(b bool) {
	log.Printf("acquiring lock ...")
	h.mux.Lock()
	defer h.mux.Unlock()
	h.bothering = b
	if b {
		log.Printf("started bothering")
	} else {
		log.Printf("stopped bothering")
	}
}

func (h *handler) bother() bool {
	log.Printf("botherer acquiring lock ...")
	h.mux.Lock()
	defer func() {
		h.mux.Unlock()
		log.Printf("botherer released lock ...")
	}()
	return h.bothering
}

func main() {
	h := &handler{}
	go func() {
		if err := h.loop(); err != nil {
			panic(err)
		}
	}()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		Handler:        h,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 0,
	}
	log.Fatal(s.ListenAndServe())
}
