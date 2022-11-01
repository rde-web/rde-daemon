package daemon

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rde-daemon/internal/config"
	"strings"
	"time"
)

type FS struct {
	sock     string
	listener net.Listener
	server   http.Server
	mux      *http.ServeMux
}

func (f *FS) Run(errChan *chan error) {
	f.sock = "fs.sock"
	f.mux = http.NewServeMux()
	f.mux.HandleFunc("/ls", makePOSTHandlerFN(f.LS))
	f.mux.HandleFunc("/cat", makePOSTHandlerFN(f.Cat))
	f.mux.HandleFunc("/touch", makePOSTHandlerFN(f.Touch))
	f.mux.HandleFunc("/override", makePOSTHandlerFN(f.Override))

	var sockLocation string = path.Join(config.Instance.SocketsPath, f.sock)
	lis, err := net.Listen("unix", sockLocation)
	if err != nil {
		*errChan <- err
		return
	}
	f.listener = lis
	log.Printf("%s listening", sockLocation)
	f.server = http.Server{
		Handler: f.mux,
	}
	errServe := f.server.Serve(f.listener)
	if errServe != nil {
		*errChan <- errServe
	}
}

func (f *FS) Shutdown() error {
	if f.listener == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	f.server.Shutdown(ctx)
	return f.listener.Close()
}

type (
	lsRequest struct {
		Path string `msgpack:"path"`
	}
	lsResponse struct {
		Files map[string]lsInfo `msgpack:"files"`
	}
	lsInfo struct {
		IsDir bool `msgpack:"is_dir"`
	}

	catRequest  lsRequest
	catResponse struct {
		Content string `msgpack:"content"`
	}

	touchRequest  lsRequest
	touchResponse struct{}

	overrideRequest struct {
		lsRequest
		catResponse
	}
	overrideResponse struct{}
)

func (FS) LS(req lsRequest) (*lsResponse, error) {
	var dir string = makePath(req.Path)
	var result lsResponse = lsResponse{
		Files: make(map[string]lsInfo),
	}
	filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		result.Files[strings.Replace(p, config.Instance.ProjectPath, "/", 1)] = lsInfo{
			IsDir: d.IsDir(),
		}
		return nil
	})
	return &result, nil
}
func (FS) Cat(req catRequest) (*catResponse, error) {
	var file string = makePath(req.Path)
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s: %v", file, err)
	}
	//@todo split large file
	return &catResponse{Content: string(content)}, nil
}

func (FS) Touch(req touchRequest) (*touchResponse, error) {
	var file string = makePath(req.Path)
	var dir string = path.Dir(file)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("cannot make dir %s: %v", dir, err)
	}
	if err := os.WriteFile(file, []byte(""), os.ModePerm); err != nil {
		return nil, fmt.Errorf("cannot make dir %s: %v", dir, err)
	}
	return &touchResponse{}, nil
}

func (FS) Override(overrideRequest) (*overrideResponse, error) {
	return nil, fmt.Errorf("method Override not implemented")
}
