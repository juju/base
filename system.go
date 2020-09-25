// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package systems

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/juju/errors"

	"github.com/juju/systems/channel"
)

// System represents an OS/Channel or Resource.
// Systems can also be converted to and from a series string.
type System struct {
	OS       string          `json:"os,omitempty"`
	Channel  channel.Channel `json:"channel,omitempty"`
	Resource string          `json:"resource,omitempty"`
}

// Validate returns with no error when the System is valid.
func (s System) Validate() error {
	if s.OS == "" && s.Resource == "" {
		return errors.NotValidf("one of os or resource must be specified")
	}

	if s.OS != "" {
		if s.Resource != "" {
			return errors.NotValidf("resource cannot be specified with os")
		}
		if !validOS.Contains(s.OS) {
			return errors.NotValidf("os %q", s.OS)
		}
		if s.Channel == channel.Empty {
			return errors.NotValidf("missing channel")
		}
	}

	if s.Resource != "" {
		if s.Channel != channel.Empty {
			return errors.NotValidf("channel cannot be specified with resource")
		}
	}
	return nil
}

// String respresentation of the System, used for series backwards compatability.
func (s System) String() string {
	// Handle legacy series.
	if series, ok := systemToSeries[s]; ok {
		return series
	}
	// Handle new system as a series.
	str := "system"
	if s.OS != "" {
		str += fmt.Sprintf("#os=%s", s.OS)
	}
	if s.Channel != channel.Empty {
		str += fmt.Sprintf("#channel=%s", s.Channel.String())
	}
	if s.Resource != "" {
		str += fmt.Sprintf("#resource=%s", s.Resource)
	}
	return str
}

// regex to match k=v from system series strings in the form
// "system#os=ubuntu#version=18.04#resource=imagename"
var systemSeriesRegex = regexp.MustCompile(`#(os|channel|resource)=([^#]+)`)

// ParseSystemFromSeries matches legacy series like "focal" or parses a system as series string
// in the form "system#os=ubuntu#version=18.04#resource=imagename"
func ParseSystemFromSeries(s string) (System, error) {
	var err error
	if !strings.HasPrefix(s, "system") {
		system, ok := seriesToSystem[s]
		if !ok {
			return System{}, errors.NotValidf("series %s is unsupported", s)
		}
		return system, nil
	}
	propString := strings.TrimPrefix(s, "system")
	matches := systemSeriesRegex.FindAllStringSubmatch(propString, -1)
	if len(matches) == 0 {
		return System{}, errors.NotValidf("invalid system series string %q", s)
	}
	matchedCharacters := 0
	system := System{}
	for _, v := range matches {
		matchedCharacters += len(v[0])
		key := v[1]
		value := v[2]
		switch key {
		case "os":
			system.OS = value
		case "channel":
			system.Channel, err = channel.Parse(value)
			if err != nil {
				return System{}, errors.Annotatef(err, "invalid channel %q", value)
			}
		case "resource":
			system.Resource = value
		default:
			return System{}, errors.NotValidf("key %q in system series string %q", key, s)
		}
	}
	if matchedCharacters != len(propString) {
		return System{}, errors.NotValidf("system series string %q", s)
	}
	err = system.Validate()
	if err != nil {
		return System{}, errors.Annotatef(err, "invalid system series string %q", s)
	}
	return system, nil
}
