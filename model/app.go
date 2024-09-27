package model

type AppInfo struct {
	SiteName    string `json:"site_name"`
	ExternalUrl string `json:"external_url"`
	SiteLogo    string `json:"site_logo"`
}
