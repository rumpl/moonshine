package moonshine

import (
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

func NewImageConfig() *specs.Image {
	img := &specs.Image{
		Architecture: "amd64",
		OS:           "linux",
		RootFS: specs.RootFS{
			Type: "layers",
		},
		Config: specs.ImageConfig{
			WorkingDir: "/",
		},
	}

	return img
}
