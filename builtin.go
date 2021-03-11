// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package systems

import (
	"fmt"

	"github.com/juju/collections/set"

	"github.com/juju/systems/channel"
)

// Supported Name constant names for Bases.
// This list should match the ones found in juju/os except for "kubernetes".
const (
	Ubuntu       = "ubuntu"
	CentOS       = "centos"
	Windows      = "windows"
	OSX          = "osx"
	OpenSUSE     = "opensuse"
	GenericLinux = "genericlinux"
)

// validOS is a string set of valid Name names.
var validOS = set.NewStrings(Ubuntu, CentOS, Windows, OSX, OpenSUSE, GenericLinux)

// seriesToBases is a map of series names to systems.
// This should match the ones found in juju/os except for "kubernetes".
var seriesToBases = map[string]Base{
	"precise": {
		Name:    Ubuntu,
		Channel: channel.MustParse("12.04/stable"),
	},
	"quantal": {
		Name:    Ubuntu,
		Channel: channel.MustParse("12.10/stable"),
	},
	"raring": {
		Name:    Ubuntu,
		Channel: channel.MustParse("13.04/stable"),
	},
	"saucy": {
		Name:    Ubuntu,
		Channel: channel.MustParse("13.10/stable"),
	},
	"trusty": {
		Name:    Ubuntu,
		Channel: channel.MustParse("14.04/stable"),
	},
	"utopic": {
		Name:    Ubuntu,
		Channel: channel.MustParse("14.10/stable"),
	},
	"vivid": {
		Name:    Ubuntu,
		Channel: channel.MustParse("15.04/stable"),
	},
	"wily": {
		Name:    Ubuntu,
		Channel: channel.MustParse("15.10/stable"),
	},
	"xenial": {
		Name:    Ubuntu,
		Channel: channel.MustParse("16.04/stable"),
	},
	"yakkety": {
		Name:    Ubuntu,
		Channel: channel.MustParse("16.10/stable"),
	},
	"zesty": {
		Name:    Ubuntu,
		Channel: channel.MustParse("17.04/stable"),
	},
	"artful": {
		Name:    Ubuntu,
		Channel: channel.MustParse("17.10/stable"),
	},
	"bionic": {
		Name:    Ubuntu,
		Channel: channel.MustParse("18.04/stable"),
	},
	"cosmic": {
		Name:    Ubuntu,
		Channel: channel.MustParse("18.10/stable"),
	},
	"disco": {
		Name:    Ubuntu,
		Channel: channel.MustParse("19.04/stable"),
	},
	"eoan": {
		Name:    Ubuntu,
		Channel: channel.MustParse("19.10/stable"),
	},
	"focal": {
		Name:    Ubuntu,
		Channel: channel.MustParse("20.04/stable"),
	},
	"groovy": {
		Name:    Ubuntu,
		Channel: channel.MustParse("20.10/stable"),
	},
	"hirsute": {
		Name:    Ubuntu,
		Channel: channel.MustParse("21.04/stable"),
	},
	"win2008r2": {
		Name:    Windows,
		Channel: channel.MustParse("win2008r2/stable"),
	},
	"win2012hvr2": {
		Name:    Windows,
		Channel: channel.MustParse("win2012hvr2/stable"),
	},
	"win2012hv": {
		Name:    Windows,
		Channel: channel.MustParse("win2012hv/stable"),
	},
	"win2012r2": {
		Name:    Windows,
		Channel: channel.MustParse("win2012r2/stable"),
	},
	"win2012": {
		Name:    Windows,
		Channel: channel.MustParse("win2012/stable"),
	},
	"win2016": {
		Name:    Windows,
		Channel: channel.MustParse("win2016/stable"),
	},
	"win2016hv": {
		Name:    Windows,
		Channel: channel.MustParse("win2016hv/stable"),
	},
	"win2016nano": {
		Name:    Windows,
		Channel: channel.MustParse("win2016nano/stable"),
	},
	"win2019": {
		Name:    Windows,
		Channel: channel.MustParse("win2019/stable"),
	},
	"win7": {
		Name:    Windows,
		Channel: channel.MustParse("win7/stable"),
	},
	"win8": {
		Name:    Windows,
		Channel: channel.MustParse("win8/stable"),
	},
	"win81": {
		Name:    Windows,
		Channel: channel.MustParse("win81/stable"),
	},
	"win10": {
		Name:    Windows,
		Channel: channel.MustParse("win10/stable"),
	},
	"centos7": {
		Name:    CentOS,
		Channel: channel.MustParse("centos7/stable"),
	},
	"centos8": {
		Name:    CentOS,
		Channel: channel.MustParse("centos8/stable"),
	},
	"opensuseleap": {
		Name:    OpenSUSE,
		Channel: channel.MustParse("opensuse42/stable"),
	},
	"genericlinux": {
		Name:    GenericLinux,
		Channel: channel.MustParse("latest/stable"),
	},
}

// baseToSeries is a reverse of seriesToBase
var baseToSeries = reverseSeriesMap()

func reverseSeriesMap() map[Base]string {
	r := make(map[Base]string)
	for series, base := range seriesToBases {
		if _, ok := r[base]; ok {
			panic(fmt.Sprintf("duplicate base %q = %v", series, base))
		}
		r[base] = series
	}
	return r
}
