package models

type Response struct {
	Data struct {
		Results []struct {
			Profile Profile `json:"user"`
		} `json:"results"`
	} `json:"data"`
}
