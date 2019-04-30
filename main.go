package main

import (
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"fmt"
	"hash"
	"log"
	"runtime"
	"sync"

	_ "github.com/ianlancetaylor/cgosymbolizer"
	"github.com/pkg/profile"

	"github.com/spacemonkeygo/openssl"
)

var _ hash.Hash = (*openSSLSHA256Hash)(nil)

type openSSLSHA256Hash struct {
	*openssl.SHA256Hash
}

func (h *openSSLSHA256Hash) BlockSize() int { return sha256.BlockSize }

func (h *openSSLSHA256Hash) Reset() {
	if err := h.SHA256Hash.Reset(); err != nil {
		panic(err)
	}
}

func (h *openSSLSHA256Hash) Size() int { return sha256.Size }

func (h *openSSLSHA256Hash) Sum(p []byte) []byte {
	s, err := h.SHA256Hash.Sum()
	if err != nil {
		panic(err)
	}
	return append(p, s[:]...)
}

func newOpenSSLSHA256Hash() hash.Hash {
	h, err := openssl.NewSHA256Hash()
	if err != nil {
		panic(err)
	}
	return &openSSLSHA256Hash{
		SHA256Hash: h,
	}
}

func compute(pool *sync.Pool, data []byte) ([]byte, error) {
	h := pool.Get().(hash.Hash)
	defer pool.Put(h)

	h.Reset()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func main() {

	useProfiler := flag.Bool("profiler", false, "enable CPU profiler")
	flag.Parse()

	// setup profiler
	if *useProfiler {
		options := []func(*profile.Profile){
			profile.ProfilePath("."),
			profile.NoShutdownHook,
			profile.CPUProfile,
		}
		profiler := profile.Start(options...)
		defer profiler.Stop()
	}

	// setup input
	const dataSize = 1024 * 4
	input := make([]byte, dataSize)
	_, err := rand.Read(input)
	if err != nil {
		log.Fatal(err)
	}

	// setup hash pool to reuse hashers
	hashPool := &sync.Pool{
		New: func() interface{} { return newOpenSSLSHA256Hash() },
	}

	inputChan := make(chan []byte)
	outputChan := make(chan []byte)
	tasks := 1024 * 1024
	workers := runtime.NumCPU()

	for i := 0; i < workers; i++ {
		go func() {
			for {
				data, ok := <-inputChan
				if !ok {
					return
				}
				output, err := compute(hashPool, data)
				if err != nil {
					log.Fatal(err)
				}
				outputChan <- output
			}
		}()
	}

	go func() {
		for i := 0; i < tasks; i++ {
			inputChan <- input
		}
		close(inputChan)
	}()

	var output []byte
	for i := 0; i < tasks; i++ {
		output = <-outputChan
	}
	fmt.Println(output)
}
