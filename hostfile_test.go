package hostfile

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func getTestfilePath(t *testing.T) string {
	t.Helper()
	return os.TempDir() + string(os.PathSeparator) + "hosts"
}

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
func hfDisplay(t *testing.T, hf []Entry) string {
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

func fakeHostfileString(t *testing.T) (example string, expected []Entry) {
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
	`, []Entry{
			Entry{IPAddress: "127.0.0.4", Hostname: "localhost4"},
			Entry{IPAddress: "127.0.0.3", Hostname: "localhost3"},
			Entry{IPAddress: "127.0.0.1", Hostname: "localhost"},
			Entry{IPAddress: "127.0.0.2", Hostname: "localhost2"},
		}
}

func TestList(t *testing.T) {
	if strings.ToUpper(os.Getenv("COMPUTERNAME")) != "W10-MIKEANDIKE" {
		t.Skip("Skipping due to being on unknown computer where hostfile may not exist")
	}
	hf, err := Open(DefaultHostfilePath)
	if err != nil {
		t.Log(err)
		t.SkipNow()
	}
	entries, err := hf.List()
	if err != nil {
		t.Log(err)
		t.SkipNow()
	}
	t.Logf("Logging entries on test machine:\n%s", hfDisplay(t, entries))

}

func TestGet(t *testing.T) {
	hf, expected := fakeHostfileString(t)
	var got []Entry
	for _, exp := range expected {

		got = get(hf, exp.IPAddress, exp.Hostname)
		if len(got) < 1 {
			t.Errorf("expected at least one entry for ip='%s', host='%s'. got 0", exp.Hostname, exp.IPAddress)
			t.Fail()
		} else {
			for _, v := range got {
				if strings.ToUpper(exp.Hostname) != strings.ToUpper(v.Hostname) ||
					strings.ToUpper(exp.IPAddress) != strings.ToUpper(v.IPAddress) {
					t.Errorf("expected: ip='%s',host='%s'. got: ip='%s', host='%s'",
						exp.IPAddress, exp.Hostname, v.IPAddress, v.Hostname)
					t.Fail()
				}
			}
		}

		got = getByIP(hf, exp.IPAddress)
		if len(got) < 1 {

			t.Errorf("expected at least one entry for ip='%s'. got 0", exp.IPAddress)
			t.Fail()
		} else {
			for _, v := range got {
				if strings.ToUpper(exp.IPAddress) != strings.ToUpper(v.IPAddress) {
					t.Errorf("expected: ip='%s'. got: ip='%s'",
						exp.IPAddress, v.IPAddress)
					t.Fail()
				}
			}
		}

		got = getByHostname(hf, exp.Hostname)
		if len(got) < 1 {

			t.Errorf("expected at least one entry for ip='%s'. got 0", exp.IPAddress)
			t.Fail()
		} else {
			for _, v := range got {
				if strings.ToUpper(exp.Hostname) != strings.ToUpper(v.Hostname) {
					t.Errorf("expected: host='%s'. got: host='%s'",
						exp.Hostname, v.Hostname)
					t.Fail()
				}
			}
		}
	}
}

func TestSetEntries(t *testing.T) {
	testfile := getTestfilePath(t)
	_, entries := fakeHostfileString(t)
	err := setEntries(testfile, entries)
	if err != nil {
		t.Fatal(err)
	}
	hf, err := Open(testfile)
	if err != nil {
		t.Fatal(err)
	}
	writtenEntries, err := hf.List()
	if len(entries) != len(writtenEntries) {

		t.Fatal(err)
	}

}

func TestRemove(t *testing.T) {
	testfile := getTestfilePath(t)
	_, entries := fakeHostfileString(t)
	err := setEntries(testfile, entries)
	if err != nil {
		t.Skip(err)
	}
	hf, err := Open(testfile)

	rmd, err := hf.Remove(entries[1].IPAddress, entries[1].Hostname)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("items removed during test: %d", rmd)
	if rmd != 1 {
		t.Fatalf("unexpected number of removed entries. expected: %d, got: %d",
			1, rmd)
	}
	current, err := hf.List()
	if err != nil {
		t.Fatal(err)
	}
	if entries[2] != current[1] {
		t.Fatalf("expected next item in row \n%+v\nbut got\n%+v\n", entries[2], current[1])
	}

}

func TestAdd(t *testing.T) {
	testfile := getTestfilePath(t)
	_, entries := fakeHostfileString(t)
	err := setEntries(testfile, entries)
	if err != nil {
		t.Skip(err)
	}
	hf, err := Open(testfile)
	if err != nil {
		t.Skip()
	}
	err = hf.Add(Entry{Hostname: "testentry", IPAddress: "127.0.0.100"})
	if err != nil {
		t.Fatal(err)
	}
	current, err := hf.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(current) != len(entries)+1 {
		t.Fatalf("unexpected number of entries after adding.\nExpected %d, got %d", len(entries)+1, len(entries))
	}
}
