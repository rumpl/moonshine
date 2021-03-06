package build

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
	"github.com/rumpl/moonshine/pkg/moonshine"
)

const (
	// LocalNameContext      = "context"
	LocalNameDockerfile = "dockerfile"
	// keyTarget             = "target"
	keyFilename = "filename"

	// keyCacheFrom          = "cache-from"
	defaultDockerfileName = "dockerfile.lua"

// dockerignoreFilename  = ".dockerignore"
// buildArgPrefix        = "build-arg:"
// labelPrefix           = "label:"
// keyNoCache            = "no-cache"
// keyTargetPlatform     = "platform"
// keyMultiPlatform      = "multi-platform"
// keyImageResolveMode   = "image-resolve-mode"
// keyGlobalAddHosts     = "add-hosts"
// keyForceNetwork       = "force-network-mode"
// keyOverrideCopyImage  = "override-copy-image" // remove after CopyOp implemented
)

func Build(ctx context.Context, c client.Client) (*client.Result, error) {
	cfg, err := GetDockerfile(ctx, c)
	if err != nil {
		return nil, errors.Wrap(err, "getting moonshine")
	}
	st, err := moonshine.DockerLuaToLLB(cfg)
	if err != nil {
		return nil, err
	}

	def, err := st.Marshal(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal local source")
	}
	res, err := c.Solve(ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve dockerfile")
	}
	ref, err := res.SingleRef()
	if err != nil {
		return nil, err
	}

	res.SetRef(ref)

	return res, nil
}

func GetDockerfile(ctx context.Context, c client.Client) (string, error) {
	opts := c.BuildOpts().Opts
	filename := opts[keyFilename]
	if filename == "" {
		filename = defaultDockerfileName
	}

	name := "load moonshine"
	if filename != "docker.lua" {
		name += " from " + filename
	}

	src := llb.Local(LocalNameDockerfile,
		llb.IncludePatterns([]string{filename}),
		llb.SessionID(c.BuildOpts().SessionID),
		llb.SharedKeyHint(defaultDockerfileName),
		dockerfile2llb.WithInternalName(name),
	)

	def, err := src.Marshal(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "failed to marshal local source")
	}

	var dtDockerfile []byte
	res, err := c.Solve(ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to resolve dockerfile")
	}

	ref, err := res.SingleRef()
	if err != nil {
		return "", err
	}

	dtDockerfile, err = ref.ReadFile(ctx, client.ReadRequest{
		Filename: filename,
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to read dockerfile")
	}

	return string(dtDockerfile), nil
}
