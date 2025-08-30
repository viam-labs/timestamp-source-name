package timestampsourcename

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"time"

	camera "go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/components/camera/rtppassthrough"
	"go.viam.com/rdk/gostream"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/rimage"
	"go.viam.com/rdk/spatialmath"
)

var (
	TimestampSourceNames = resource.NewModel("viam-labs", "time-stamp-source-name", "timestamp-source-names")
	errUnimplemented     = errors.New("unimplemented")
)

const (
	timestampFormat = "2006-01-02T15:04:05.000Z07:00"
)

func init() {
	resource.RegisterComponent(camera.API, TimestampSourceNames,
		resource.Registration[camera.Camera, *Config]{
			Constructor: newTimeStampSourceNameTimestampSourceNames,
		},
	)
}

type Config struct {
	NImages int `json:"n_images"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit required (first return) and optional (second return) dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *Config) Validate(path string) ([]string, []string, error) {
	// Add config validation code here
	if cfg.NImages < 1 {
		return nil, nil, errors.New("n_images attribute must be greater than 0")
	}
	return nil, nil, nil
}

type timeStampSourceNameTimestampSourceNames struct {
	resource.AlwaysRebuild

	name    resource.Name
	bluePic image.Image
	nImages int

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()
}

func newTimeStampSourceNameTimestampSourceNames(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (camera.Camera, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	return NewTimestampSourceNames(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewTimestampSourceNames(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (camera.Camera, error) {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	img := image.NewRGBA(image.Rect(0, 0, 600, 400))
	blue := color.RGBA{0, 0, 255, 255}
	for y := 0; y < 400; y++ {
		for x := 0; x < 600; x++ {
			img.Set(x, y, blue)
		}
	}

	s := &timeStampSourceNameTimestampSourceNames{
		name:       name,
		bluePic:    img,
		nImages:    conf.NImages,
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}
	return s, nil
}

func (s *timeStampSourceNameTimestampSourceNames) Name() resource.Name {
	return s.name
}

func (s *timeStampSourceNameTimestampSourceNames) Stream(ctx context.Context, errHandlers ...gostream.ErrorHandler) (gostream.VideoStream, error) {
	panic("not implemented")
}

// Image returns a byte slice representing an image that tries to adhere to the MIME type hint.
// Image also may return metadata about the frame.
func (s *timeStampSourceNameTimestampSourceNames) Image(ctx context.Context, mimeType string, extra map[string]interface{}) ([]byte, camera.ImageMetadata, error) {
	theBytes, err := rimage.EncodeImage(ctx, s.bluePic, mimeType)
	if err != nil {
		return nil, camera.ImageMetadata{}, err
	}
	meta := camera.ImageMetadata{
		MimeType: mimeType,
	}
	return theBytes, meta, nil

}

// Images is used for getting simultaneous images from different imagers,
// along with associated metadata (just timestamp for now). It's not for getting a time series of images from the same imager.
// The extra parameter can be used to pass additional options to the camera resource.
func (s *timeStampSourceNameTimestampSourceNames) Images(ctx context.Context, extra map[string]interface{}) ([]camera.NamedImage, resource.ResponseMetadata, error) {
	now := time.Now()
	timestampStr := now.Format(timestampFormat)
	// 5 images, all blue, different timestamp source names
	result := make([]camera.NamedImage, s.nImages)
	for i := 0; i < s.nImages; i++ {
		result[i].SourceName = fmt.Sprintf("%s_%d", timestampStr, i)
		result[i].Image = s.bluePic
	}

	meta := resource.ResponseMetadata{
		CapturedAt: now,
	}
	return result, meta, nil
}

// NextPointCloud returns the next immediately available point cloud, not necessarily one
// a part of a sequence. In the future, there could be streaming of point clouds.
func (s *timeStampSourceNameTimestampSourceNames) NextPointCloud(ctx context.Context) (pointcloud.PointCloud, error) {
	panic("not implemented")
}

// Properties returns properties that are intrinsic to the particular
// implementation of a camera.
func (s *timeStampSourceNameTimestampSourceNames) Properties(ctx context.Context) (camera.Properties, error) {
	return camera.Properties{
		SupportsPCD: false,
	}, nil
}
func (s *timeStampSourceNameTimestampSourceNames) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	panic("not implemented")
}
func (s *timeStampSourceNameTimestampSourceNames) Geometries(ctx context.Context, extra map[string]interface{}) ([]spatialmath.Geometry, error) {
	panic("not implemented")
}
func (s *timeStampSourceNameTimestampSourceNames) SubscribeRTP(ctx context.Context, bufferSize int, packetsCB rtppassthrough.PacketCallback) (rtppassthrough.Subscription, error) {
	panic("not implemented")
}
func (s *timeStampSourceNameTimestampSourceNames) Unsubscribe(ctx context.Context, id rtppassthrough.SubscriptionID) error {
	panic("not implemented")
}

func (s *timeStampSourceNameTimestampSourceNames) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}
