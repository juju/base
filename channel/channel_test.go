// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package channel_test

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/systems/channel"
)

type storeChannelSuite struct{}

var _ = gc.Suite(&storeChannelSuite{})

func (s storeChannelSuite) TestParse(c *gc.C) {
	ch, err := channel.Parse("stable")
	c.Assert(err, gc.IsNil)
	c.Check(ch, jc.DeepEquals, channel.Channel{
		Name:   "stable",
		Track:  "",
		Risk:   channel.Stable,
		Branch: "",
	})

	ch, err = channel.Parse("latest/stable")
	c.Assert(err, gc.IsNil)
	c.Check(ch, jc.DeepEquals, channel.Channel{
		Name:   "stable",
		Track:  "",
		Risk:   channel.Stable,
		Branch: "",
	})

	ch, err = channel.Parse("1.0/edge")
	c.Assert(err, gc.IsNil)
	c.Check(ch, jc.DeepEquals, channel.Channel{
		Name:   "1.0/edge",
		Track:  "1.0",
		Risk:   channel.Edge,
		Branch: "",
	})

	ch, err = channel.Parse("1.0")
	c.Assert(err, gc.IsNil)
	c.Check(ch, jc.DeepEquals, channel.Channel{
		Name:   "1.0/stable",
		Track:  "1.0",
		Risk:   channel.Stable,
		Branch: "",
	})

	ch, err = channel.Parse("1.0/beta/foo")
	c.Assert(err, gc.IsNil)
	c.Check(ch, jc.DeepEquals, channel.Channel{
		Name:   "1.0/beta/foo",
		Track:  "1.0",
		Risk:   channel.Beta,
		Branch: "foo",
	})

	ch, err = channel.Parse("candidate/foo")
	c.Assert(err, gc.IsNil)
	c.Check(ch, jc.DeepEquals, channel.Channel{
		Name:   "candidate/foo",
		Track:  "",
		Risk:   channel.Candidate,
		Branch: "foo",
	})

	ch, err = channel.Parse("candidate/foo")
	c.Assert(err, gc.IsNil)
	c.Check(ch, jc.DeepEquals, channel.Channel{
		Name:   "candidate/foo",
		Track:  "",
		Risk:   channel.Candidate,
		Branch: "foo",
	})
}

func mustParse(c *gc.C, channelStr string) channel.Channel {
	ch, err := channel.Parse(channelStr)
	c.Assert(err, gc.IsNil)
	return ch
}

func (s storeChannelSuite) TestParseVerbatim(c *gc.C) {
	ch, err := channel.ParseVerbatim("sometrack")
	c.Assert(err, gc.IsNil)
	c.Check(ch, jc.DeepEquals, channel.Channel{
		Track: "sometrack",
	})
	c.Check(ch.VerbatimTrackOnly(), gc.Equals, true)
	c.Check(ch.VerbatimRiskOnly(), gc.Equals, false)
	c.Check(mustParse(c, "sometrack"), jc.DeepEquals, ch.Clean())

	ch, err = channel.ParseVerbatim("latest")
	c.Assert(err, gc.IsNil)
	c.Check(ch, jc.DeepEquals, channel.Channel{
		Track: "latest",
	})
	c.Check(ch.VerbatimTrackOnly(), gc.Equals, true)
	c.Check(ch.VerbatimRiskOnly(), gc.Equals, false)
	c.Check(mustParse(c, "latest"), jc.DeepEquals, ch.Clean())

	ch, err = channel.ParseVerbatim("edge")
	c.Assert(err, gc.IsNil)
	c.Check(ch, jc.DeepEquals, channel.Channel{
		Risk: channel.Edge,
	})
	c.Check(ch.VerbatimTrackOnly(), gc.Equals, false)
	c.Check(ch.VerbatimRiskOnly(), gc.Equals, true)
	c.Check(mustParse(c, "edge"), jc.DeepEquals, ch.Clean())

	ch, err = channel.ParseVerbatim("latest/stable")
	c.Assert(err, gc.IsNil)
	c.Check(ch, jc.DeepEquals, channel.Channel{
		Track: "latest",
		Risk:  channel.Stable,
	})
	c.Check(ch.VerbatimTrackOnly(), gc.Equals, false)
	c.Check(ch.VerbatimRiskOnly(), gc.Equals, false)
	c.Check(mustParse(c, "latest/stable"), jc.DeepEquals, ch.Clean())

	ch, err = channel.ParseVerbatim("latest/stable/foo")
	c.Assert(err, gc.IsNil)
	c.Check(ch, jc.DeepEquals, channel.Channel{
		Track:  "latest",
		Risk:   channel.Stable,
		Branch: "foo",
	})
	c.Check(ch.VerbatimTrackOnly(), gc.Equals, false)
	c.Check(ch.VerbatimRiskOnly(), gc.Equals, false)
	c.Check(mustParse(c, "latest/stable/foo"), jc.DeepEquals, ch.Clean())
}

func (s storeChannelSuite) TestClean(c *gc.C) {
	ch := channel.Channel{
		Track: "latest",
		Name:  "latest/stable",
		Risk:  channel.Stable,
	}

	cleanedCh := ch.Clean()
	c.Check(cleanedCh, gc.Not(jc.DeepEquals), c)
	c.Check(cleanedCh, jc.DeepEquals, channel.Channel{
		Track: "",
		Name:  "stable",
		Risk:  channel.Stable,
	})
}

func (s storeChannelSuite) TestParseErrors(c *gc.C) {
	for _, tc := range []struct {
		channel string
		err     string
		full    string
	}{
		{"", "channel name cannot be empty", ""},
		{"1.0////", "channel name has too many components: 1.0////", "1.0/stable"},
		{"1.0/cand", "invalid risk in channel name: 1.0/cand", ""},
		{"fix//hotfix", "invalid risk in channel name: fix//hotfix", ""},
		{"/stable/", "invalid track in channel name: /stable/", "latest/stable"},
		{"//stable", "invalid risk in channel name: //stable", "latest/stable"},
		{"stable/", "invalid branch in channel name: stable/", "latest/stable"},
		{"/stable", "invalid track in channel name: /stable", "latest/stable"},
	} {
		_, err := channel.Parse(tc.channel)
		c.Check(err, gc.ErrorMatches, tc.err)
		_, err = channel.ParseVerbatim(tc.channel)
		c.Check(err, gc.ErrorMatches, tc.err)
		if tc.full != "" {
			// testing Full behavior on the malformed channel
			full, err := channel.Full(tc.channel)
			c.Check(err, gc.IsNil)
			c.Check(full, gc.Equals, tc.full)
		}
	}
}

func (s *storeChannelSuite) TestString(c *gc.C) {
	tests := []struct {
		channel string
		str     string
	}{
		{"stable", "stable"},
		{"latest/stable", "stable"},
		{"1.0/edge", "1.0/edge"},
		{"1.0/beta/foo", "1.0/beta/foo"},
		{"1.0", "1.0/stable"},
		{"candidate/foo", "candidate/foo"},
	}

	for _, t := range tests {
		ch, err := channel.Parse(t.channel)
		c.Assert(err, gc.IsNil)

		c.Check(ch.String(), gc.Equals, t.str)
	}
}

func (s *storeChannelSuite) TestChannelFull(c *gc.C) {
	tests := []struct {
		channel string
		str     string
	}{
		{"stable", "latest/stable"},
		{"latest/stable", "latest/stable"},
		{"1.0/edge", "1.0/edge"},
		{"1.0/beta/foo", "1.0/beta/foo"},
		{"1.0", "1.0/stable"},
		{"candidate/foo", "latest/candidate/foo"},
	}

	for _, t := range tests {
		ch, err := channel.Parse(t.channel)
		c.Assert(err, gc.IsNil)

		c.Check(ch.Full(), gc.Equals, t.str)
	}
}

func (s *storeChannelSuite) TestFuncFull(c *gc.C) {
	tests := []struct {
		channel string
		str     string
	}{
		{"stable", "latest/stable"},
		{"latest/stable", "latest/stable"},
		{"1.0/edge", "1.0/edge"},
		{"1.0/beta/foo", "1.0/beta/foo"},
		{"1.0", "1.0/stable"},
		{"candidate/foo", "latest/candidate/foo"},
		// store behaviour compat; expect these to fail when we stop accommodating the madness :)
		{"//stable//", "latest/stable"},
		// rather weird corner case
		{"///", ""},
		// empty string is OK
		{"", ""},
	}

	for _, t := range tests {
		can, err := channel.Full(t.channel)
		c.Assert(err, gc.IsNil)
		c.Check(can, gc.Equals, t.str)
	}
}

func (s *storeChannelSuite) TestFuncFullErr(c *gc.C) {
	_, err := channel.Full("foo/bar/baz/quux")
	c.Check(err, gc.ErrorMatches, "invalid channel")
}

func (s *storeChannelSuite) TestMatch(c *gc.C) {
	tests := []struct {
		req string
		c1  string
		res string
	}{
		{"stable", "stable", "track:risk"},
		{"stable", "beta", "track"},
		{"beta", "stable", "track:risk"},
		{"stable", "edge", "track"},
		{"edge", "stable", "track:risk"},
		{"1.0/stable", "1.0/edge", "track"},
		{"1.0/edge", "stable", "risk"},
		{"1.0/edge", "stable", "risk"},
		{"1.0/stable", "stable", "risk"},
		{"1.0/stable", "beta", ""},
		{"1.0/stable", "2.0/beta", ""},
		{"2.0/stable", "2.0/beta", "track"},
		{"1.0/stable", "2.0/beta", ""},
	}

	for _, t := range tests {
		req, err := channel.Parse(t.req)
		c.Assert(err, gc.IsNil)
		c1, err := channel.Parse(t.c1)
		c.Assert(err, gc.IsNil)

		c.Check(req.Match(&c1).String(), gc.Equals, t.res)
	}
}

func (s *storeChannelSuite) TestResolve(c *gc.C) {
	tests := []struct {
		channel string
		new     string
		result  string
		expErr  string
	}{
		{"", "", "", ""},
		{"", "edge", "edge", ""},
		{"track/foo", "", "track/foo", ""},
		{"stable", "", "stable", ""},
		{"stable", "edge", "edge", ""},
		{"stable/branch1", "edge/branch2", "edge/branch2", ""},
		{"track", "track", "track", ""},
		{"track", "beta", "track/beta", ""},
		{"track/stable", "beta", "track/beta", ""},
		{"track/stable", "stable/branch", "track/stable/branch", ""},
		{"track/stable", "track/edge/branch", "track/edge/branch", ""},
		{"track/stable", "track/candidate", "track/candidate", ""},
		{"track/stable", "track/stable/branch", "track/stable/branch", ""},
		{"track1/stable", "track2/stable", "track2/stable", ""},
		{"track1/stable", "track2/stable/branch", "track2/stable/branch", ""},
		{"track/foo", "track/stable/branch", "", "invalid risk in channel name: track/foo"},
	}

	for _, t := range tests {
		r, err := channel.Resolve(t.channel, t.new)
		tcomm := gc.Commentf("%#v", t)
		if t.expErr == "" {
			c.Assert(err, gc.IsNil, tcomm)
			c.Check(r, gc.Equals, t.result, tcomm)
		} else {
			c.Assert(err, gc.ErrorMatches, t.expErr, tcomm)
		}
	}
}

func (s *storeChannelSuite) TestResolvePinned(c *gc.C) {
	tests := []struct {
		track  string
		new    string
		result string
		expErr string
	}{
		{"", "", "", ""},
		{"", "anytrack/stable", "anytrack/stable", ""},
		{"track/foo", "", "", "invalid pinned track: track/foo"},
		{"track", "", "track", ""},
		{"track", "track", "track", ""},
		{"track", "beta", "track/beta", ""},
		{"track", "stable/branch", "track/stable/branch", ""},
		{"track", "track/edge/branch", "track/edge/branch", ""},
		{"track", "track/candidate", "track/candidate", ""},
		{"track", "track/stable/branch", "track/stable/branch", ""},
		{"track1", "track2/stable", "track2/stable", "cannot switch pinned track"},
		{"track1", "track2/stable/branch", "track2/stable/branch", "cannot switch pinned track"},
	}
	for _, t := range tests {
		r, err := channel.ResolvePinned(t.track, t.new)
		tcomm := gc.Commentf("%#v", t)
		if t.expErr == "" {
			c.Assert(err, gc.IsNil, tcomm)
			c.Check(r, gc.Equals, t.result, tcomm)
		} else {
			c.Assert(err, gc.ErrorMatches, t.expErr, tcomm)
		}
	}
}
