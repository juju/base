// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package systems

import (
	"strings"

	"github.com/juju/errors"

	"github.com/juju/systems/channel"
)

// Base represents an OS/Channel.
// Bases can also be converted to and from a series string.
type Base struct {
	Name    string          `json:"name,omitempty"`
	Channel channel.Channel `json:"channel,omitempty"`
}

// Validate returns with no error when the Base is valid.
func (s Base) Validate() error {
	if s.Name == "" {
		return errors.NotValidf("name must be specified")
	}

	if !validOS.Contains(s.Name) {
		return errors.NotValidf("os %q", s.Name)
	}
	if s.Channel == channel.Empty {
		return errors.NotValidf("channel")
	}

	return nil
}

// String respresentation of the Base, used for series backwards compatability.
func (s Base) String() string {
	// Handle legacy series.
	if series, ok := baseToSeries[s]; ok {
		return series
	}
	str := s.Name
	if s.Channel != channel.Empty {
		str += "/" + s.Channel.String()
	}
	return str
}

// ParseBaseFromSeries matches legacy series like "focal" or parses a base as series string
// in the form "os/track/risk/branch"
func ParseBaseFromSeries(s string) (Base, error) {
	var err error
	if base, ok := seriesToBases[s]; ok {
		return base, nil
	}

	// Split the first forward-slash to get name and channel.
	// E.g. "os/track/risk/branch" => ["os", "track/risk/branch"]
	segments := strings.SplitN(s, "/", 2)
	osName := segments[0]
	channelName := ""
	if len(segments) == 2 {
		channelName = segments[1]
	}

	base := Base{}
	if !validOS.Contains(osName) {
		return Base{}, errors.NotValidf("series %q", s)
	}
	base.Name = osName

	if channelName != "" {
		base.Channel, err = channel.Parse(channelName)
		if err != nil {
			return Base{}, errors.Annotatef(err, "malformed channel in base string %q", s)
		}
	}

	err = base.Validate()
	if err != nil {
		return Base{}, errors.Annotatef(err, "invalid base string %q", s)
	}
	return base, nil
}
