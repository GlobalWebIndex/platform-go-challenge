package dto

import (
	"github.com/ipinfo/go-ipinfo/ipinfo"
)

type UserActivity struct {
	Hostname string `json:"host_name"`
	Org      string `json:"org"`
	Ip       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Phone    string `json:"phone"`
	Postal   string `json:"postal"`
}

func NewUserLogByIPInfo(details *ipinfo.Info) *UserActivity {
	return &UserActivity{
		Hostname: details.Hostname,
		Org:      details.Organization,
		Ip:       details.IP.String(),
		City:     details.City,
		Region:   details.Region,
		Country:  details.Country,
		Loc:      details.Location,
		Phone:    details.Phone,
		Postal:   details.Postal,
	}
}
