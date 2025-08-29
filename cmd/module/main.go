package main

import (
	"timestampsourcename"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
	camera "go.viam.com/rdk/components/camera"
)

func main() {
	// ModularMain can take multiple APIModel arguments, if your module implements multiple models.
	module.ModularMain(resource.APIModel{ camera.API, timestampsourcename.TimestampSourceNames})
}
