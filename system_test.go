// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package systems_test

import (
	"encoding/json"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/systems"
	"github.com/juju/systems/channel"
)

type systemSuite struct {
	testing.CleanupSuite
}

var _ = gc.Suite(&systemSuite{})

func (s *systemSuite) TestBaseParsingToFromSeries(c *gc.C) {
	tests := []struct {
		system     systems.Base
		str        string
		parsedBase systems.Base
		err        string
	}{
		{systems.Base{Name: systems.Ubuntu}, "ubuntu", systems.Base{}, `invalid base string "ubuntu": channel not valid`},
		{systems.Base{Name: systems.Windows}, "windows", systems.Base{}, `invalid base string "windows": channel not valid`},
		{systems.Base{Name: "mythicalos"}, "mythicalos", systems.Base{}, `series "mythicalos" not valid`},
		{systems.Base{Name: systems.Ubuntu, Channel: channel.MustParse("20.04/stable")}, "focal", systems.Base{Name: systems.Ubuntu, Channel: channel.MustParse("20.04/stable")}, ""},
		{systems.Base{Name: systems.Ubuntu, Channel: channel.MustParse("18.04/stable")}, "bionic", systems.Base{Name: systems.Ubuntu, Channel: channel.MustParse("18.04/stable")}, ""},
		{systems.Base{Name: systems.Windows, Channel: channel.MustParse("win10/stable")}, "win10", systems.Base{Name: systems.Windows, Channel: channel.MustParse("win10/stable")}, ""},
		{systems.Base{Name: systems.Ubuntu, Channel: channel.MustParse("20.04/edge")}, "ubuntu/20.04/edge", systems.Base{Name: systems.Ubuntu, Channel: channel.MustParse("20.04/edge")}, ""},
	}
	for i, v := range tests {
		str := v.system.String()
		comment := gc.Commentf("test %d", i)
		c.Check(str, gc.Equals, v.str, comment)
		s, err := systems.ParseBaseFromSeries(str)
		if v.err != "" {
			c.Check(err, gc.ErrorMatches, v.err, comment)
		} else {
			c.Check(err, jc.ErrorIsNil, comment)
		}
		c.Check(s, jc.DeepEquals, v.parsedBase, comment)
	}
}

func (s *systemSuite) TestJSONEncoding(c *gc.C) {
	sys := systems.Base{
		Name:    systems.Ubuntu,
		Channel: channel.MustParse("20.04/stable"),
	}
	bytes, err := json.Marshal(sys)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(bytes), gc.Equals, `{"name":"ubuntu","channel":{"name":"20.04/stable","track":"20.04","risk":"stable"}}`)
	sys2 := systems.Base{}
	err = json.Unmarshal(bytes, &sys2)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(sys2, jc.DeepEquals, sys)
}
