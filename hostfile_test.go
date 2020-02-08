package hostfile

import (
	"fmt"
	"strings"
	"testing"
)

func TestUnmarshalHostfile(t *testing.T) {
	testfile, expected := fakeHostfileString(t)
	hf := unmarshalHostfile(testfile)
	if len(hf) != len(expected) {
		t.Errorf("Error parsing file. Desired %d, got %d. Printing entries below..\n%s",
			len(expected), len(hf), hfDisplay(t, hf))
		t.FailNow()
	}
	for i, v := range expected {
		if v.Hostname != hf[i].Hostname {
			t.Errorf("in item number %d, expected: '%s', got '%s'",
				i, v.Hostname, hf[i].Hostname)
			t.Fail()
		}
		if v.IPAddress != hf[i].IPAddress {
			t.Errorf("in item number %d, expected: '%s', got '%s'",
				i, v.IPAddress, hf[i].IPAddress)
			t.Fail()
		}
	}

}
func hfDisplay(t *testing.T, hf []HostfileEntry) string {
	t.Helper()
	var entries []string
	for i, v := range hf {
		entries = append(entries,
			fmt.Sprintf("%d. ip='%s' name='%s'\n", i, v.IPAddress, v.Hostname))
	}
	return fmt.Sprintf("%s\n%s\n%s",
		strings.Repeat("-", 20),
		strings.Join(entries, "\n"),
		strings.Repeat("-", 20),
	)
}

func fakeHostfileString(t *testing.T) (example string, expected []HostfileEntry) {
	t.Helper()
	return `# Copyright (c) 1993-2009 Microsoft Corp.
	#
	# This is a sample HOSTS file used by Microsoft TCP/IP for Windows.
	#
	# This file contains the mappings of IP addresses to host names. Each
	# entry should be kept on an individual line. The IP address should
	# be placed in the first column followed by the corresponding host name.
	# The IP address and the host name should be separated by at least one
	# space.
	#
	# Additionally, comments (such as these) may be inserted on individual
	# lines or following the machine name denoted by a '#' symbol.
	#
	# For example:
	#
	#      102.54.94.97     rhino.acme.com          # source server
	#       38.25.63.10     x.acme.com              # x client host
127.0.0.4 localhost4
	127.0.0.3 localhost3 #this is localhost 3
	# localhost name resolution is handled within DNS itself.
		127.0.0.1       localhost
	#	::1             localhost
		127.0.0.2       localhost2
	`, []HostfileEntry{
			HostfileEntry{IPAddress: "127.0.0.4", Hostname: "localhost4"},
			HostfileEntry{IPAddress: "127.0.0.3", Hostname: "localhost3"},
			HostfileEntry{IPAddress: "127.0.0.1", Hostname: "localhost"},
			HostfileEntry{IPAddress: "127.0.0.2", Hostname: "localhost2"},
		}
}
