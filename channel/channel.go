// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package channel

import (
	"errors"
	"fmt"
	"strings"

	"github.com/juju/collections/set"
)

// Empty channel
var Empty Channel

// Risk type
type Risk string

// Well-known risk contants.
const (
	Unknown   Risk = ""
	Stable    Risk = "stable"
	Candidate Risk = "candidate"
	Beta      Risk = "beta"
	Edge      Risk = "edge"
)

var channelRisks = set.NewStrings(
	string(Stable),
	string(Candidate),
	string(Beta),
	string(Edge),
)
var channelRiskLevels = map[Risk]int{
	Stable:    0,
	Candidate: 1,
	Beta:      2,
	Edge:      3,
}

// Channel identifies and describes completely a store channel.
type Channel struct {
	Name   string `json:"name"`
	Track  string `json:"track"`
	Risk   Risk   `json:"risk"`
	Branch string `json:"branch,omitempty"`
}

func isSlash(r rune) bool { return r == '/' }

// TODO: currently there's some overlap between the toplevel Full, and
//       methods Clean, String, and Full. Needs further refactoring.

// Full normalizes the channel string to also include stable when a
// risk is not found in the original channel string.
func Full(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	components := strings.FieldsFunc(s, isSlash)
	switch len(components) {
	case 0:
		return "", nil
	case 1:
		if channelRisks.Contains(components[0]) {
			return "latest/" + components[0], nil
		}
		return components[0] + "/stable", nil
	case 2:
		if channelRisks.Contains(components[0]) {
			return "latest/" + strings.Join(components, "/"), nil
		}
		fallthrough
	case 3:
		return strings.Join(components, "/"), nil
	default:
		return "", errors.New("invalid channel")
	}
}

// ParseVerbatim parses a string representing a store channel.
// The channel representation is not normalized.
// Parse() should be used in most cases.
func ParseVerbatim(s string) (Channel, error) {
	if s == "" {
		return Empty, fmt.Errorf("channel name cannot be empty")
	}
	p := strings.Split(s, "/")
	var risk, track, branch *string
	switch len(p) {
	default:
		return Empty, fmt.Errorf("channel name has too many components: %s", s)
	case 3:
		track, risk, branch = &p[0], &p[1], &p[2]
	case 2:
		if channelRisks.Contains(p[0]) {
			risk, branch = &p[0], &p[1]
		} else {
			track, risk = &p[0], &p[1]
		}
	case 1:
		if channelRisks.Contains(p[0]) {
			risk = &p[0]
		} else {
			track = &p[0]
		}
	}

	ch := Channel{}
	if risk != nil {
		if !channelRisks.Contains(*risk) {
			return Empty, fmt.Errorf("invalid risk in channel name: %s", s)
		}
		ch.Risk = Risk(*risk)
	}
	if track != nil {
		if *track == "" {
			return Empty, fmt.Errorf("invalid track in channel name: %s", s)
		}
		ch.Track = *track
	}
	if branch != nil {
		if *branch == "" {
			return Empty, fmt.Errorf("invalid branch in channel name: %s", s)
		}
		ch.Branch = *branch
	}

	return ch, nil
}

// Parse parses a string representing a store channel.
// The returned channel's track, risk and name are normalized.
func Parse(s string) (Channel, error) {
	channel, err := ParseVerbatim(s)
	if err != nil {
		return Empty, err
	}
	return channel.Clean(), nil
}

// MustParse parses a channel string and returns a pointer to the channel
// or panics.
func MustParse(s string) Channel {
	c, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return c
}

// Clean returns a Channel with a normalized track, risk and name.
func (c Channel) Clean() Channel {
	track := c.Track
	risk := c.Risk

	if track == "latest" {
		track = ""
	}
	if risk == "" {
		risk = Stable
	}

	// normalized name
	name := string(risk)
	if track != "" {
		name = track + "/" + name
	}
	if c.Branch != "" {
		name = name + "/" + c.Branch
	}

	return Channel{
		Name:   name,
		Track:  track,
		Risk:   risk,
		Branch: c.Branch,
	}
}

func (c Channel) String() string {
	return c.Name
}

// Full returns the full name of the channel, inclusive the default track "latest".
func (c *Channel) Full() string {
	ch := c.String()
	full, err := Full(ch)
	if err != nil {
		// impossible
		panic("channel.String() returned a malformed channel: " + ch)
	}
	return full
}

// VerbatimTrackOnly returns whether the channel represents a track only.
func (c *Channel) VerbatimTrackOnly() bool {
	return c.Track != "" && c.Risk == "" && c.Branch == ""
}

// VerbatimRiskOnly returns whether the channel represents a risk only.
func (c *Channel) VerbatimRiskOnly() bool {
	return c.Track == "" && c.Risk != "" && c.Branch == ""
}

func riskLevel(risk Risk) int {
	level, ok := channelRiskLevels[risk]
	if ok {
		return level
	}
	return -1
}

// Match represents on which fields two channels are matching.
type Match struct {
	Track bool
	Risk  bool
}

// String returns the string represantion of the match, results can be:
//  "track:risk"
//  "track"
//  "risk"
//  ""
func (cm Match) String() string {
	matching := []string{}
	if cm.Track {
		matching = append(matching, "track")
	}
	if cm.Risk {
		matching = append(matching, "risk")
	}
	return strings.Join(matching, ":")
}

// Match returns a Match of which fields among architecture,track,risk match between c and c1 store channels, risk is matched taking channel inheritance into account and considering c the requested channel.
func (c *Channel) Match(c1 *Channel) Match {
	requestedRiskLevel := riskLevel(c.Risk)
	rl1 := riskLevel(c1.Risk)
	return Match{
		Track: c.Track == c1.Track,
		Risk:  requestedRiskLevel >= rl1,
	}
}

// Resolve resolves newChannel wrt channel, this means if newChannel
// is risk/branch only it will preserve the track of channel. It
// assumes that if both are not empty, channel is parseable.
func Resolve(channel, newChannel string) (string, error) {
	if newChannel == "" {
		return channel, nil
	}
	if channel == "" {
		return newChannel, nil
	}
	ch, err := ParseVerbatim(channel)
	if err != nil {
		return "", err
	}
	p := strings.Split(newChannel, "/")
	if channelRisks.Contains(p[0]) && ch.Track != "" {
		// risk/branch inherits the track if any
		return ch.Track + "/" + newChannel, nil
	}
	return newChannel, nil
}

// ErrPinnedTrackSwitch is returned from ResolvePinned when a track is pinned
// and a new channel moves to a different track.
var ErrPinnedTrackSwitch = errors.New("cannot switch pinned track")

// ResolvePinned resolves newChannel wrt a pinned track, newChannel
// can only be risk/branch-only or have the same track, otherwise
// ErrPinnedTrackSwitch is returned.
func ResolvePinned(track, newChannel string) (string, error) {
	if track == "" {
		return newChannel, nil
	}
	ch, err := ParseVerbatim(track)
	if err != nil || !ch.VerbatimTrackOnly() {
		return "", fmt.Errorf("invalid pinned track: %s", track)
	}
	if newChannel == "" {
		return track, nil
	}
	trackPrefix := ch.Track + "/"
	p := strings.Split(newChannel, "/")
	if channelRisks.Contains(p[0]) && ch.Track != "" {
		// risk/branch inherits the track if any
		return trackPrefix + newChannel, nil
	}
	if newChannel != track && !strings.HasPrefix(newChannel, trackPrefix) {
		// the track is pinned
		return "", ErrPinnedTrackSwitch
	}
	return newChannel, nil
}
