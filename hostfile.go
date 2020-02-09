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
func (hf *Hostfile) Get(ip string, hostname string) ([]Entry, error) {
	hfString, err := getfileString(hf.Path)
	if err != nil {
		return nil, err
	}
	return get(hfString, ip, hostname), nil

}

func get(hf string, ip string, hn string) []Entry {
	entries := list(hf)
	var got []Entry
	for _, e := range entries {
		if strings.ToUpper(e.Hostname) == strings.ToUpper(hn) &&
			strings.ToUpper(e.IPAddress) == strings.ToUpper(ip) {
			got = append(got, e)
		}
	}

	return got
}

// GetByIP retrieves all items where the IP matches the given parameter
// Impliments List method of Hostfile
func (hf *Hostfile) GetByIP(ip string) ([]Entry, error) {
	hfString, err := getfileString(hf.Path)
	if err != nil {
		return nil, err
	}
	return getByIP(hfString, ip), nil
}

func getByIP(hf string, ip string) []Entry {
	entries := list(hf)
	var got []Entry
	for _, e := range entries {
		if strings.ToUpper(e.IPAddress) == strings.ToUpper(ip) {
			got = append(got, e)
		}
	}

	return got
}

// GetByHostname retrieves all items where the Hostname matches the given parameter
// Impliments List method of Hostfile
func (hf *Hostfile) GetByHostname(hostname string) ([]Entry, error) {
	hfString, err := getfileString(hf.Path)
	if err != nil {
		return nil, err
	}
	return getByHostname(hfString, hostname), nil
}

func getByHostname(hf string, hn string) []Entry {
	entries := list(hf)
	var got []Entry
	for _, e := range entries {
		if strings.ToUpper(e.Hostname) == strings.ToUpper(hn) {
			got = append(got, e)
		}
	}

	return got
}

// Add adds a single entry to a hostfile
// Impliments List method of Hostfile
func (hf *Hostfile) Add(entry Entry) error {
	// TODO: Verify entry
	// TODO: think about things below
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
	// Not appending file directly in case we want to impliment checks for
	// existing entries

	entries, err := hf.List()
	if err != nil {
		return err
	}
	entries = append(entries, entry)
	return setEntries(hf.Path, entries)
}

// Remove removes a single entry from a hostfile where IP AND Hostname matches given parameters
// Impliments List method of Hostfile
func (hf *Hostfile) Remove(ip string, hostname string) (int, error) {
	hfString, err := getfileString(hf.Path)
	if err != nil {
		return -1, err
	}
	keep, removed := remove(hfString, ip, hostname)
	err = setEntries(hf.Path, keep)
	if err != nil {
		return -1, err
	}
	return removed, nil
}
func remove(hf string, ip string, hn string) ([]Entry, int) {
	entries := list(hf)
	var keep []Entry
	var removed int
	for _, e := range entries {
		if strings.ToUpper(e.IPAddress) != strings.ToUpper(ip) ||
			strings.ToUpper(e.Hostname) != strings.ToUpper(hn) {
			keep = append(keep, e)
		} else {
			removed++
		}
	}
	return keep, removed
}

// RemoveByIP removes a single entry from a hostfile where IP matches given parameter
// Impliments List and Remove methods of Hostfile
func (hf *Hostfile) RemoveByIP(ip string) (int, error) {
	hfString, err := getfileString(hf.Path)
	if err != nil {
		return -1, err
	}
	keep, removed := removeByIP(hfString, ip)
	err = setEntries(hf.Path, keep)
	if err != nil {
		return -1, err
	}
	return removed, nil
}
func removeByIP(hf string, ip string) ([]Entry, int) {
	entries := list(hf)
	var keep []Entry
	var removed int
	for _, e := range entries {
		if strings.ToUpper(e.IPAddress) != strings.ToUpper(ip) {
			keep = append(keep, e)
		} else {
			removed++
		}
	}
	return keep, removed
}

// RemoveByHostname removes a single entry from a hostfile where Hostname matches given parameter
// Impliments List and Remove methods of Hostfile
func (hf *Hostfile) RemoveByHostname(hostname string) (int, error) {
	hfString, err := getfileString(hf.Path)
	if err != nil {
		return -1, err
	}
	keep, removed := removeByHostname(hfString, hostname)
	err = setEntries(hf.Path, keep)
	if err != nil {
		return -1, err
	}
	return removed, nil
}
func removeByHostname(hf string, hn string) ([]Entry, int) {
	entries := list(hf)
	var keep []Entry
	var removed int
	for _, e := range entries {
		if strings.ToUpper(e.Hostname) != strings.ToUpper(hn) {
			keep = append(keep, e)
		} else {
			removed++
		}
	}
	return keep, removed
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
	for _, line := range strings.FieldsFunc(hfString, func(r rune) bool { return r == '\n' || r == '\r' }) {
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

func setEntries(fp string, entries []Entry) error {
	var lines = []string{
		"# Entry format is 'IPADDRESS HOSTNAME'. lines starting with '#' are comment lines and ignored",
	}
	for _, e := range entries {
		lines = append(lines, fmt.Sprintf("%s   %s", e.IPAddress, e.Hostname))
	}
	lines = append(lines, "")

	return ioutil.WriteFile(fp, []byte(strings.Join(lines, "\r\n")), 0644)
}
