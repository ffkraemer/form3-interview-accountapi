package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
}

type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}

func (a *AccountData) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet:
		a.Fetch(w, r)
		return
	case r.Method == http.MethodPut:
		a.Create(w, r)
		return
	case r.Method == http.MethodDelete:
		a.Delete(w, r)
		return
	default:
		http.NotFound(w, r)
		return
	}
}

func (a *AccountData) Fetch(w http.ResponseWriter, r *http.Request) {
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	param := r.URL.Query().Get("account_id")

	if len(param) <= 0 {
		http.Error(w, "Invalid account_id", 400)
		return
	}

	url := "http://localhost:8080/v1/organisation/accounts/" + param
	resp, err := c.Get(url)

	if err != nil || resp.StatusCode != 200 {
		http.Error(w, err.Error(), resp.StatusCode)

		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Error %s", err.Error())
		return
	}

	var response AccountData
	err = json.Unmarshal(body, &response)

	if err != nil {
		fmt.Println("Failed on retrieve data", err.Error())
		return
	}

	w.Write(body)
}

func (a *AccountData) Create(w http.ResponseWriter, r *http.Request) {

}

func (a *AccountData) Delete(w http.ResponseWriter, r *http.Request) {

}

func main() {

	mux := http.NewServeMux()
	mux.Handle("/accounts/", &AccountData{})
	http.ListenAndServe("localhost:8081", mux)
}
