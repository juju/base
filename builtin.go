// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package systems

import (
	"fmt"

	"github.com/juju/collections/set"

	"github.com/juju/systems/channel"
)

// Supported OS constant names for Systems.
// This list should match the ones found in juju/os except for "kubernetes".
const (
	Ubuntu       = "ubuntu"
	CentOS       = "centos"
	Windows      = "windows"
	OSX          = "osx"
	OpenSUSE     = "opensuse"
	GenericLinux = "genericlinux"
)

// validOS is a string set of valid OS names.
var validOS = set.NewStrings(Ubuntu, CentOS, Windows, OSX, OpenSUSE, GenericLinux)

// seriesToSystem is a map of series names to systems.
// This should match the ones found in juju/os except for "kubernetes".
var seriesToSystem = map[string]System{
	"precise": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("12.04/stable"),
	},
	"quantal": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("12.10/stable"),
	},
	"raring": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("13.04/stable"),
	},
	"saucy": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("13.10/stable"),
	},
	"trusty": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("14.04/stable"),
	},
	"utopic": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("14.10/stable"),
	},
	"vivid": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("15.04/stable"),
	},
	"wily": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("15.10/stable"),
	},
	"xenial": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("16.04/stable"),
	},
	"yakkety": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("16.10/stable"),
	},
	"zesty": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("17.04/stable"),
	},
	"artful": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("17.10/stable"),
	},
	"bionic": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("18.04/stable"),
	},
	"cosmic": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("18.10/stable"),
	},
	"disco": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("19.04/stable"),
	},
	"eoan": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("19.10/stable"),
	},
	"focal": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("20.04/stable"),
	},
	"groovy": System{
		OS:      Ubuntu,
		Channel: channel.MustParse("20.10/stable"),
	},
	"win2008r2": System{
		OS:      Windows,
		Channel: channel.MustParse("win2008r2/stable"),
	},
	"win2012hvr2": System{
		OS:      Windows,
		Channel: channel.MustParse("win2012hvr2/stable"),
	},
	"win2012hv": System{
		OS:      Windows,
		Channel: channel.MustParse("win2012hv/stable"),
	},
	"win2012r2": System{
		OS:      Windows,
		Channel: channel.MustParse("win2012r2/stable"),
	},
	"win2012": System{
		OS:      Windows,
		Channel: channel.MustParse("win2012/stable"),
	},
	"win2016": System{
		OS:      Windows,
		Channel: channel.MustParse("win2016/stable"),
	},
	"win2016hv": System{
		OS:      Windows,
		Channel: channel.MustParse("win2016hv/stable"),
	},
	"win2016nano": System{
		OS:      Windows,
		Channel: channel.MustParse("win2016nano/stable"),
	},
	"win2019": System{
		OS:      Windows,
		Channel: channel.MustParse("win2019/stable"),
	},
	"win7": System{
		OS:      Windows,
		Channel: channel.MustParse("win7/stable"),
	},
	"win8": System{
		OS:      Windows,
		Channel: channel.MustParse("win8/stable"),
	},
	"win81": System{
		OS:      Windows,
		Channel: channel.MustParse("win81/stable"),
	},
	"win10": System{
		OS:      Windows,
		Channel: channel.MustParse("win10/stable"),
	},
	"centos7": System{
		OS:      CentOS,
		Channel: channel.MustParse("centos7/stable"),
	},
	"centos8": System{
		OS:      CentOS,
		Channel: channel.MustParse("centos8/stable"),
	},
	"opensuseleap": System{
		OS:      OpenSUSE,
		Channel: channel.MustParse("opensuse42/stable"),
	},
	"genericlinux": System{
		OS:      GenericLinux,
		Channel: channel.MustParse("latest/stable"),
	},
}

// systemToSeries is a reverse of seriesToSystem
var systemToSeries = reverseSeriesMap()

func reverseSeriesMap() map[System]string {
	r := make(map[System]string)
	for series, system := range seriesToSystem {
		if _, ok := r[system]; ok {
			panic(fmt.Sprintf("duplicate system %q = %v", series, system))
		}
		r[system] = series
	}
	return r
}
