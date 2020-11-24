package main

import (
	"context"
	"flag"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/subcommands"
	"github.com/nelhage/llama/cmd/internal/cli"
	"github.com/nelhage/llama/store/s3store"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")

	subcommands.Register(&InvokeCommand{}, "")

	subcommands.Register(&StoreCommand{}, "internals")
	subcommands.Register(&GetCommand{}, "internals")

	var state cli.GlobalState
	flag.StringVar(&state.Region, "region", "", "S3 region for commands")
	flag.StringVar(&state.Bucket, "bucket", "", "S3 bucket for the llama object store")

	flag.Parse()

	if state.Bucket == "" {
		state.Bucket = os.Getenv("LLAMA_BUCKET")
	}

	var cfg aws.Config
	if state.Region != "" {
		cfg.Region = &state.Region
	}
	state.Session = session.Must(session.NewSession(&cfg))
	state.Store = s3store.FromSession(state.Session, state.Bucket)

	ctx := context.Background()
	ctx = cli.WithState(ctx, &state)

	os.Exit(int(subcommands.Execute(ctx)))
}
