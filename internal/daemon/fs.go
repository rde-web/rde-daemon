package daemon

import (
	"context"
	"io/fs"
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
	var dir string
	if req != nil {
		dir = path.Join(config.Config.ProjectPath, req.Path)
	}
	if !strings.HasPrefix(dir, config.Config.ProjectPath) {
		dir = config.Config.ProjectPath
	}
	var result LSResponse = LSResponse{
		Files: make(map[string]*Info),
	}
	filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		result.Files[path.Join(p, d.Name())] = &Info{
			IsDir: d.IsDir(),
		}
		return nil
	})
	return &result, nil
}
func (FS) Cat(context.Context, *CatRequest) (*CatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Cat not implemented")
}
func (FS) Touch(context.Context, *TouchRequest) (*TouchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Touch not implemented")
}
func (FS) Override(context.Context, *OverrideRequest) (*OverrideResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Override not implemented")
}
