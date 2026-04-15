// Package config provides CLI flag parsing and configuration management
// for logpipe.
//
// It defines the [Config] struct which captures all user-supplied options,
// and the [Parse] function which reads os.Args-style argument slices and
// returns a validated configuration ready for use by the pipeline.
//
// Example usage:
//
//	cfg, err := config.Parse(os.Args[1:])
//	if err != nil {
//		log.Fatal(err)
//	}
package config
