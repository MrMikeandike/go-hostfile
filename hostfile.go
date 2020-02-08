package hostfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// DefaultHostfilePath is the default hostfile location
	// in most versions of windows
	DefaultHostfilePath = "C:\\Windows\\System32\\drivers\\etc\\hosts"
)

// Open is a convenient way of creating a new Hostfile struct
func Open(filepath ...string) (Hostfile, error) {
	var fp string
	if len(filepath) == 0 {
		fp = DefaultHostfilePath
	} else if len(filepath) == 1 {
		fp = filepath[0]
	} else {
		return Hostfile{}, fmt.Errorf("too many strings entered as filepath parameter. Please enter 0 or 1 strings")
	}
	hf := Hostfile{Path: fp}
	return hf, hf.IsValidPath()

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
func (hf *Hostfile) List() ([]Entry, error) {
	hfString, err := getfileString(hf.Path)
	if err != nil {
		return nil, err
	}
	return list(hfString), nil
}

func list(hf string) []Entry {
	return unmarshalHostfile(hf)
}

// Get retrieves all items where the IP AND Hostname matches
func (hf *Hostfile) Get() []Entry {
	// TODO
	return nil
}

// GetByIP retrieves all items where the IP matches the given parameter
// Impliments List method of Hostfile
func (hf *Hostfile) GetByIP() []Entry {
	// TODO
	return nil
}

// GetByHostname retrieves all items where the Hostname matches the given parameter
// Impliments List method of Hostfile
func (hf *Hostfile) GetByHostname() []Entry {
	// TODO
	return nil
}

// Add adds a single entry to a hostfile
// Impliments List method of Hostfile
func (hf *Hostfile) Add(entry Entry) error {
	// TODO
	/*
		Things to think about
			- Should it error when item exist?
			- should it return the index of the item?
			- What happens if the ip OR host already exist?
			  - should it error?
			  - should it allow you to specify before it after existing items?
			  - should it just do nothing?
			  - should all of the above have parameters that let you choose?
	*/
	return nil
}

// Remove removes a single entry from a hostfile where IP AND Hostname matches given parameters
// Impliments List method of Hostfile
func (hf *Hostfile) Remove(entry Entry) error {
	return nil
}

// RemoveByIP removes a single entry from a hostfile where IP matches given parameter
// Impliments List and Remove methods of Hostfile
func (hf *Hostfile) RemoveByIP(entry Entry) error {
	return nil
}

// RemoveByHostname removes a single entry from a hostfile where Hostname matches given parameter
// Impliments List and Remove methods of Hostfile
func (hf *Hostfile) RemoveByHostname(entry Entry) error {
	return nil
}

// IsValidPath tests the hostfile for common issues, such as file not existing,
// and file being a directory
func (hf *Hostfile) IsValidPath() error {
	info, err := os.Stat(hf.Path)
	if info.IsDir() {
		return fmt.Errorf("Error while opening file: given path is directory")
	} else if err == nil {
		return nil
	} else if os.IsNotExist(err) {
		return fmt.Errorf("Error while opening file: File doesn't exist")
	} else {
		return fmt.Errorf("Error while opening file: Unknown error:\n%v", err)
	}

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
