package main

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
	"os"

	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/pkg/errors"
	"github.com/rumpl/moonshine/pkg/build"
	"github.com/rumpl/moonshine/pkg/moonshine"
	"github.com/sirupsen/logrus"
)

var graph bool
var filename string

func main() {
	flag.BoolVar(&graph, "graph", false, "output a graph and exit")
	flag.StringVar(&filename, "filename", "dockerfile.lua", "the file to read from")
	flag.Parse()

	if graph {
		if err := printGraph(filename, os.Stdout); err != nil {
			logrus.Fatalf("fatal error: %s", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if err := grpcclient.RunFromEnvironment(appcontext.Context(), build.Build); err != nil {
		logrus.Fatalf("fatal error: %s", err)
		panic(err)
	}
}

func printGraph(filename string, out io.Writer) error {
	b, _ := ioutil.ReadFile(filename)
	st, err := moonshine.DockerLuaToLLB(string(b))
	if err != nil {
		return errors.Wrap(err, "to llb")
	}

	_, err = st.Marshal(context.Background())
	if err != nil {
		return errors.Wrap(err, "marshaling llb state")
	}

	// return llb.WriteTo(dt, out)
	return nil
}
