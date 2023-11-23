package util

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const ApiUrlSuffix = "/api2/json"

/*
GetPveUrl returns either the URL as specified by the `pveurl` arg,
or builds a URL from the `scheme`, `pvehost`, and `pveport` args.
*/
func GetPveUrl(c *cli.Context) string {
	var ret string
	pveUrl := c.String("pveurl")
	switch pveUrl {
	case "":
		ret = fmt.Sprint(
			c.String("scheme"),
			"://",
			c.String("pvehost"),
			":",
			c.String("pveport"),
			ApiUrlSuffix,
		)
	default:
		ret = pveUrl
	}
	return ret
}

/*
ViperPveUrl returns either the URL as specified by the `pve_url` arg,
or builds a URL from the `pve_uri_scheme`, `pve_host`, and `pve_port` args.
*/
func ViperPveUrl() string {
	var ret string
	pveUrl := viper.GetString("url")
	switch pveUrl {
	case "":
		ret = fmt.Sprint(
			viper.GetString("pve_uri_scheme"),
			"://",
			viper.GetString("pve_host"),
			":",
			viper.GetString("pve_port"),
			ApiUrlSuffix,
		)
	default:
		ret = pveUrl
	}
	return ret
}
