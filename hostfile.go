package hostfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// DEFAULT_HOSTFILE_LOCATION is the default hostfile location
	// in most versions of windows
	DEFAULT_HOSTFILE_LOCATION = "C:\\Windows\\System32\\drivers\\etc\\hosts"
)

// OpenHostfile is a convenient way of creating a new Hostfile struct
func OpenHostfile(filepath ...string) (Hostfile, error) {
	var fp string
	if len(filepath) == 0 {
		fp = DEFAULT_HOSTFILE_LOCATION
	} else if len(filepath) == 1 {
		fp = filepath[0]
	} else {
		return Hostfile{}, fmt.Errorf("too many strings entered as filepath parameter. Please enter 0 or 1 strings")
	}
	hf := Hostfile{Path: fp}
	return hf, hf.TestPath()

}

// Hostfile represents all information needed to make changes to a hostfile
type Hostfile struct {
	Path string
}

// Entry represents one valid line from the hostfile
type Entry struct {
	IPAddress string
	Hostname  string
}

// List returns all existing hostfile entries
func (h *Hostfile) List() {

}

// Get retrieves a single item of a hostfile
func (h *Hostfile) Get() {

}

// Add adds a single entry to a hostfile
func (h *Hostfile) Add() {

}

// Remove removes a single entry from a hostfile
func (h *Hostfile) Remove() {

}

func getfileString(fp string) (string, error) {
	// assume its been tested already
	// get bytes
	hfBytes, err := ioutil.ReadFile(fp)
	if err != nil {
		return "", err
	}
	return string(hfBytes), nil
	// convert
}
func unmarshalHostfile(hfString string) []Entry {
	var entries []Entry
	for _, line := range strings.FieldsFunc(hfString, func(r rune) bool { return r == '\n' }) {
		line = strings.Trim(strings.ReplaceAll(line, "\t", " "), " ")
		if comment := strings.Index(line, "#"); comment != -1 {
			line = strings.Trim(line[0:comment], "\t ")

		}
		split := strings.FieldsFunc(line, func(r rune) bool { return r == ' ' || r == '\t' })

		if len(split) == 2 {
			entries = append(entries, Entry{
				IPAddress: split[0],
				Hostname:  split[1],
			})
		}
	}
	return entries

}

// TestPath tests the hostfile for common issues, such as file not existing,
// and file being a directory
func (h *Hostfile) TestPath() error {
	info, err := os.Stat(h.Path)
	if !info.IsDir() {
		return fmt.Errorf("Error while opening file: given path is directory")
	} else if err != nil {
		return nil
	} else if os.IsNotExist(err) {
		return fmt.Errorf("Error while opening file: File doesn't exist")
	} else {
		return fmt.Errorf("Error while opening file: Unknown error:\n%v", err)
	}

}
