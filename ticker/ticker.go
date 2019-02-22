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
	address     = "localhost"
	counterPort = 8042
	port = 8080
)

func putTick(m string) error {
	r := &pb.PutTickRequest{Message: m}

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", address, counterPort), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Panicf("error closing connection: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c := pb.NewCounterClient(conn)

	resp, err := c.PutTick(ctx, r)

	fmt.Printf("%v\n", resp)

	return err
}

func (h *handler) bother() error {
	for {
		time.Sleep(1 * time.Second)
		if h.bothering {
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
	}
}

type handler struct {
	bothering bool
	mux       sync.Mutex
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := "hic sunt dracones"
	if r.RequestURI == "/start" {
		m = "started"
		h.setBothering(true)
	} else if r.RequestURI == "/stop" {
		m = "paused"
		h.setBothering(false)
	} else if r.RequestURI == "/kill" {
		os.Exit(0)
	}

	if _, err := w.Write([]byte(m)); err != nil {
		log.Print(err)
	}
}

func (h *handler) setBothering(b bool) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.bothering = b
}

func main() {
	h := &handler{}
	go func() {
		if err := h.bother(); err != nil {
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
