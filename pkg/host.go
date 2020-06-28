package host

import "os"

var hostname string

func init() {
	_hostname, err := os.Hostname()
	if err != nil {
		// fmt.Println(err)
		_hostname = "Unknown_Host"
	}
	hostname = _hostname
}

//GetHostname 获取本机的hostname，如果获取不到会返回 Unknown_Host
func GetHostname() string {
	return hostname
}
