package main

import (
	"context"
	"timestampsourcename"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	camera "go.viam.com/rdk/components/camera"
)

func main() {
	err := realMain()
	if err != nil {
		panic(err)
	}
}

func realMain() error {
	ctx := context.Background()
	logger := logging.NewLogger("cli")

	deps := resource.Dependencies{}
	// can load these from a remote machine if you need

	cfg := timestampsourcename.Config{}

	thing, err := timestampsourcename.NewTimestampSourceNames(ctx, deps, camera.Named("foo"), &cfg, logger)
	if err != nil {
		return err
	}
	defer thing.Close(ctx)

	return nil
}
