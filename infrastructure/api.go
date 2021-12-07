package infrastructure

import (
	"encoding/json"
	"fmt"
	"go-tinder-crawler/models"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type TinderAPI struct {
	baseURL string
	token   string
}

func NewTinderAPI(token string) *TinderAPI {
	return &TinderAPI{
		baseURL: "https://api.gotinder.com",
		token:   token,
	}
}

func (a *TinderAPI) unmarshalJSON(body io.Reader) ([]models.Profile, error) {
	var response models.Response
	var profiles []models.Profile

	responseData, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, err
	}
	for _, result := range response.Data.Results {
		profile := result.Profile
		profile.YOB = profile.BirthDate.Year()
		profile.LastModified = time.Now().UTC()
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (a *TinderAPI) GetNearbyProfiles() ([]models.Profile, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/recs/core", a.baseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", a.token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	return a.unmarshalJSON(resp.Body)
}
