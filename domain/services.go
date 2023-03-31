package domain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

type DomainServices struct{}

func (s *DomainServices) GetDomains() ([]map[string]interface{}, error) {
	url := "https://api.liqu.id/v1/domains"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("9815", "c8331f82f06c11ffe5ad342b684f04c4")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	json.Unmarshal(body, &result)

	return result, nil
}

func (s *DomainServices) GetDetailManageDomain(domain string) (map[string]interface{}, error) {
	path := "domains/details-by-name"
	url := fmt.Sprintf("%s/%s?%s=%s", viper.GetString("THIRD_PARTY.URL_DOMAIN"), path, "domain_name", domain)
	req, _ := http.NewRequest("GET", url, nil)

	req.SetBasicAuth("9815", "c8331f82f06c11ffe5ad342b684f04c4")
	client := &http.Client{}
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	return result, nil
}

func (s *DomainServices) GetBalanceAccount() map[string]interface{} {
	path := "account/balance"
	url := fmt.Sprintf("%s/%s", viper.GetString("THIRD_PARTY.URL_DOMAIN"), path)
	req, _ := http.NewRequest("GET", url, nil)

	req.SetBasicAuth("9815", "c8331f82f06c11ffe5ad342b684f04c4")
	client := &http.Client{}
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	newObj := map[string]interface{}{
		"balance": string(body),
	}

	return newObj
}

func (s *DomainServices) GetAvailabiltyDomain(keyword string) (map[string]interface{}, error) {
	path := "domains/suggestion"
	params := "tlds=ac.id%2Cbiz.id%2Cco.id%2Ccom%2Cid%2Cmy.id%2Cor.id%2Cponpes.id%2Csch.id%2Cweb.id%2Cxyz&limit=10&hyphen_allowed=true&add_related=false"
	url := fmt.Sprintf("%s/%s?%s=%s&%s", viper.GetString("THIRD_PARTY.URL_DOMAIN"), path, "keyword", keyword, params)
	req, _ := http.NewRequest("GET", url, nil)

	req.SetBasicAuth("9815", "c8331f82f06c11ffe5ad342b684f04c4")
	client := &http.Client{}
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)

	var result map[string]map[string]string
	json.Unmarshal(body, &result)

	newObj := make(map[string]interface{})
	for objName, objProps := range result {
		for propName, propValue := range objProps {
			newKey := fmt.Sprintf("%s.%s", objName, propName)
			newObj[newKey] = propValue
		}
	}

	return newObj, nil
}

func (s *DomainServices) GetPriceDomain() map[string]interface{} {
	path := "account/prices"
	url := fmt.Sprintf("%s/%s", viper.GetString("THIRD_PARTY.URL_DOMAIN"), path)
	req, _ := http.NewRequest("GET", url, nil)

	req.SetBasicAuth("9815", "c8331f82f06c11ffe5ad342b684f04c4")
	client := &http.Client{}
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	var newObj map[string]interface{}

	json.Unmarshal(body, &newObj)
	return newObj
}
