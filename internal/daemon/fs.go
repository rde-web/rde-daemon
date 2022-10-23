package daemon

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"rde-daemon/internal/config"
	. "rde-daemon/internal/protocol/rde_fs"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FS struct {
	UnimplementedRDEFSServer
}

func (FS) LS(ctx context.Context, req *LSRequest) (*LSResponse, error) {
	var dir string = makePath(req.GetPath())
	var result LSResponse = LSResponse{
		Files: make(map[string]*Info),
	}
	filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		result.Files[strings.Replace(p, config.Config.ProjectPath, "/", 1)] = &Info{
			IsDir: d.IsDir(),
		}
		return nil
	})
	return &result, nil
}
func (FS) Cat(ctx context.Context, req *CatRequest) (*CatResponse, error) {
	var file string = makePath(req.GetPath())
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannot read file %s: %v", file, err),
		)
	}
	//@todo split large file
	return &CatResponse{Content: string(content)}, nil
}
func (FS) Touch(ctx context.Context, req *TouchRequest) (*TouchResponse, error) {
	var file string = makePath(req.GetPath())
	var dir string = path.Dir(file)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("cannot make dir %s: %v", dir, err),
		)
	}
	if err := os.WriteFile(file, []byte(""), os.ModePerm); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("cannot make dir %s: %v", dir, err),
		)
	}
	return &TouchResponse{}, nil //@todo useless fields
}
func (FS) Override(context.Context, *OverrideRequest) (*OverrideResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Override not implemented")
}

func makePath(p string) string {
	var dir string = path.Join(config.Config.ProjectPath, p)
	if !strings.HasPrefix(dir, config.Config.ProjectPath) {
		dir = config.Config.ProjectPath
	}
	return dir
}
