package actor

// Config is a struct to hold the configuration
type Config struct {
	MinActors       int `env:"MIN_ACTORS" default:"10"`
	MaxActors       int `env:"MAX_ACTORS" default:"100"`
	ScalingInterval int `env:"SCALING_INTERVAL" default:"100"`
	UpScaleFactor   int `env:"UP_SCALE_FACTOR" default:"1"`
	DownScaleFactor int `env:"DOWN_SCALE_FACTOR" default:"1"`
	AutoScale
}

// AutoScale is a struct to hold the auto scale configuration
type AutoScale struct {
	UpScaleQueueSize   int `env:"UP_SCALE_QUEUE_SIZE" default:"100"`
	DownScaleQueueSize int `env:"DOWN_SCALE_QUEUE_SIZE" default:"100"`
}
