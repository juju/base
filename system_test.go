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

func (s *systemSuite) TestSystemParsingToFromSeries(c *gc.C) {
	tests := []struct {
		system       systems.System
		str          string
		parsedSystem systems.System
		err          string
	}{
		{systems.System{OS: systems.Ubuntu}, "system#os=ubuntu", systems.System{}, `invalid system series string "system#os=ubuntu": missing channel not valid`},
		{systems.System{OS: systems.Windows}, "system#os=windows", systems.System{}, `invalid system series string "system#os=windows": missing channel not valid`},
		{systems.System{OS: "mythicalos"}, "system#os=mythicalos", systems.System{}, `invalid system series string "system#os=mythicalos": os "mythicalos" not valid`},
		{systems.System{OS: systems.Ubuntu, Channel: channel.MustParse("20.04/stable"), Resource: "test-resource"}, "system#os=ubuntu#channel=20.04/stable#resource=test-resource", systems.System{}, `invalid system series string "system#os=ubuntu#channel=20.04/stable#resource=test-resource": resource cannot be specified with os not valid`},
		{systems.System{OS: systems.Ubuntu, Channel: channel.MustParse("20.04/stable")}, "focal", systems.System{OS: systems.Ubuntu, Channel: channel.MustParse("20.04/stable")}, ""},
		{systems.System{OS: systems.Ubuntu, Channel: channel.MustParse("18.04/stable")}, "bionic", systems.System{OS: systems.Ubuntu, Channel: channel.MustParse("18.04/stable")}, ""},
		{systems.System{OS: systems.Windows, Channel: channel.MustParse("win10/stable")}, "win10", systems.System{OS: systems.Windows, Channel: channel.MustParse("win10/stable")}, ""},
		{systems.System{Resource: "test"}, "system#resource=test", systems.System{Resource: "test"}, ""},
		{systems.System{OS: systems.Ubuntu, Channel: channel.MustParse("20.04/edge")}, "system#os=ubuntu#channel=20.04/edge", systems.System{OS: systems.Ubuntu, Channel: channel.MustParse("20.04/edge")}, ""},
	}
	for i, v := range tests {
		str := v.system.String()
		comment := gc.Commentf("test %d", i)
		c.Check(str, gc.Equals, v.str, comment)
		s, err := systems.ParseSystemFromSeries(str)
		if v.err != "" {
			c.Check(err, gc.ErrorMatches, v.err, comment)
		} else {
			c.Check(err, jc.ErrorIsNil, comment)
		}
		c.Check(s, jc.DeepEquals, v.parsedSystem, comment)
	}
}

func (s *systemSuite) TestJSONEncoding(c *gc.C) {
	sys := systems.System{
		OS:       systems.Ubuntu,
		Channel:  channel.MustParse("20.04/stable"),
		Resource: "resource-name",
	}
	bytes, err := json.Marshal(sys)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(bytes), gc.Equals, `{"os":"ubuntu","channel":{"name":"20.04/stable","track":"20.04","risk":"stable"},"resource":"resource-name"}`)
	sys2 := systems.System{}
	err = json.Unmarshal(bytes, &sys2)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(sys2, jc.DeepEquals, sys)
}
