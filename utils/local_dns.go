package utils

import (
	"github.com/txn2/txeh"
)

func SetLocalDNS(ip, domain string) error {
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		return err
	}
	hosts.AddHost(ip, domain)
	//hfData := hosts.RenderHostsFile()
	//
	//// if you like to see what the outcome will
	//// look like
	//fmt.Printf("----> /etc/hosts: [%s]", hfData)

	err = hosts.Save()
	if err != nil {
		return err
	}
	return nil
}
