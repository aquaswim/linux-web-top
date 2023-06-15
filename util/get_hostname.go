package util

import "os"

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Error get hostname"
	}
	return hostname
}
