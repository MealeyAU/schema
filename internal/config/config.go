package config

import (
	"flag"
	"strings"
)

type Config struct {
	// GoOutput Outputs go bindings when enabled. Usage: --go
	GoOutput bool
	// WebOutput Outputs web bindings when enabled Usage: --web
	WebOutput bool
	// AllOutput Outputs all binding types when enabled. Usage: --all
	AllOutput bool
}

func (c *Config) Init() {
	flag.BoolVar(&c.GoOutput, "go", false, "--go")
	flag.BoolVar(&c.WebOutput, "web", false, "--web")
	flag.BoolVar(&c.AllOutput, "all", false, "--all")

	flag.Parse()

	if c.AllOutput {
		// NOTE: please add any additional output targets here
		c.GoOutput = true
		c.WebOutput = true
	}
}

func (c *Config) EnabledOutputsStrings() string {
	var outputs []string

	if c.GoOutput {
		outputs = append(outputs, "go")
	}

	if c.WebOutput {
		outputs = append(outputs, "web")
	}

	return strings.Join(outputs, ",")
}
