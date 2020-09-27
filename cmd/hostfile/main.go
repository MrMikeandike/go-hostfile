package main

import (
	"fmt"

	"github.com/mrmikeandike/go-hostfile"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	path = kingpin.Flag("path", "path to hostfile. will use windows default if not used").
		Short('p').
		Default(hostfile.DefaultHostfilePath).
		ExistingFile()

	add = kingpin.Command("add", "add hostfile entry").
		Alias("a")
	addIP = add.Flag("ip", "ip address of entry").
		Short('i').Required().IP()
	addHostname = add.Flag("name", "hostname of entry").
			Short('n').Required().String()

	remove = kingpin.Command("remove", "remove entries from hostfile").
		Alias("r")
	removeIP = remove.Flag("ip", "ip address of entry. Required if --name is not given").
			Short('i').IP()
	removeHostname = remove.Flag("name", "hostname of entry. Required if --ip is not given").
			Short('n').String()

	list = kingpin.Command("list", "lists entries from hostfile").
		Alias("l").Default()

	get = kingpin.Command("get", "gets all entries associated with all given parameters").
		Alias("g")
	getIP = get.Flag("ip", "ip address of entry. Required if --name is not given").
		Short('i').IP()
	getHostname = get.Flag("name", "Hostname of entry. Required if --ip is not given").
			Short('n').String()
)

func addAction() {
	fmt.Printf("name: '%s', ip: '%s'\n", *addHostname, addIP.String())
	hf, err := hostfile.Open(*path)
	kingpin.FatalIfError(err, "Error while processing file")
	err = hf.Add(hostfile.Entry{Hostname: *addHostname, IPAddress: addIP.String()})
	kingpin.FatalIfError(err, "Error while trying to add the entry")
}

func removeAction() {
	removeIPStr := removeIP.String()
	if removeIPStr == "<nil>" && *removeHostname == "" {
		kingpin.Usage()
		return
	}
	hf, err := hostfile.Open(*path)
	kingpin.FatalIfError(err, "Error while processing file")

	if removeIPStr == "<nil>" {
		changed, err := hf.RemoveByHostname(*removeHostname)
		kingpin.FatalIfError(err, "error while removing entry")
		fmt.Printf("number of records removed: '%d'\n", changed)

	} else if *removeHostname == "" {
		changed, err := hf.RemoveByIP(removeIPStr)
		kingpin.FatalIfError(err, "error while removing entry")
		fmt.Printf("number of records removed: '%d'\n", changed)

	} else {
		changed, err := hf.Remove(removeIPStr, *removeHostname)
		kingpin.FatalIfError(err, "error while removing entry")
		fmt.Printf("number of records removed: '%d'\n", changed)
	}
	return

}

func listAction() {
	hf, err := hostfile.Open(*path)
	kingpin.FatalIfError(err, "Error while processing file")
	entries, err := hf.List()
	kingpin.FatalIfError(err, "Error while listing entries")

	for _, e := range entries {
		fmt.Println(e.IPAddress + "," + e.Hostname)
	}
	return

}

func getAction() {
	getIPStr := getIP.String()
	if getIPStr == "<nil>" && *getHostname == "" {
		kingpin.Usage()
		return
	}
	hf, err := hostfile.Open(*path)
	kingpin.FatalIfError(err, "Error while processing file")
	if getIPStr == "<nil>" {
		entries, err := hf.GetByHostname(*getHostname)
		kingpin.FatalIfError(err, "Error getting entries by hostname")
		for _, e := range entries {
			fmt.Println(e.IPAddress + "," + e.Hostname)
		}
	} else if *getHostname == "" {
		entries, err := hf.GetByIP(getIPStr)
		kingpin.FatalIfError(err, "Error getting entries by IP")
		for _, e := range entries {
			fmt.Println(e.IPAddress + "," + e.Hostname)
		}
	} else {
		entries, err := hf.Get(getIPStr, *getHostname)
		kingpin.FatalIfError(err, "Error while getting entries")
		for _, e := range entries {
			fmt.Println(e.IPAddress + "," + e.Hostname)
		}
	}
}

func main() {
	switch kingpin.Parse() {
	case add.FullCommand():
		addAction()
	case remove.FullCommand():
		removeAction()
	case list.FullCommand():
		listAction()
	case get.FullCommand():
		getAction()
	}
}
