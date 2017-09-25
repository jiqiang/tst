package apm

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/jiqiang/tst/server/apm/model"
	"github.com/parnurzeal/gorequest"
)

const (
	uaaURL              = "https://cc076c41-1318-4aab-a64e-18aa6dd254b7.predix-uaa.run.aws-usw02-pr.ice.predix.io/oauth/token"
	username            = "ems-apm-admin2"
	password            = "se3ret"
	timeout             = 5
	enterpriseSourceKey = "ENTERPRISE_da4ab60d-2f69-4bdb-af18-6cafe981af82"
	sitesAPITmpl        = "http://localhost:8008/v1/enterprises/%s/sites"
	assetsAPITmpl       = "http://localhost:8008/v1/sites/%s/assets"
)

// Token holds a token.
type token struct {
	Value string `json:"access_token"`
}

// GetToken gets a new UI service access token from remote.
func GetToken() (string, []error) {
	var t token
	request := gorequest.New()
	request.Timeout(time.Duration(timeout) * time.Second)
	request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	request.SetBasicAuth(username, password)
	request.Post(uaaURL)
	request.Send("grant_type=client_credentials")
	_, _, errs := request.EndStruct(&t)
	if errs != nil {
		return "", errs
	}
	return t.Value, nil
}

// GetSites gets all sites data.
func GetSites(token string) ([]model.Site, []error) {
	siteCollection := struct {
		Sites []model.Site `json:"content"`
	}{}
	sitesAPIEndpoint := fmt.Sprintf(sitesAPITmpl, enterpriseSourceKey)
	authorizationStr := fmt.Sprintf("Bearer %s", token)
	request := gorequest.New()
	request.Timeout(time.Duration(timeout) * time.Second)
	request.Get(sitesAPIEndpoint)
	request.Set("Accept", "application/json")
	request.Set("Authorization", authorizationStr)
	_, _, errs := request.EndStruct(&siteCollection)
	if errs != nil {
		return nil, errs
	}
	sites := []model.Site{}
	for _, site := range siteCollection.Sites {
		sites = append(sites, site)
	}
	return sites, nil
}

// GetAssetsBySite gets all assets data for a specific site.
func GetAssetsBySite(token string, siteSourceKey string) ([]model.Asset, []error) {
	assetCollection := struct {
		Assets []model.Asset `json:"content"`
	}{}
	authorizationStr := fmt.Sprintf("Bearer %s", token)
	request := gorequest.New()
	request.Timeout(time.Duration(timeout) * time.Second)
	assetsAPIEndpoint := fmt.Sprintf(assetsAPITmpl, siteSourceKey)
	request.Get(assetsAPIEndpoint)
	request.Set("Accept", "application/json")
	request.Set("Authorization", authorizationStr)
	_, _, errs := request.EndStruct(&assetCollection)
	if errs != nil {
		return nil, errs
	}
	assets := []model.Asset{}
	for _, asset := range assetCollection.Assets {
		assets = append(assets, asset)
	}
	return assets, nil
}
