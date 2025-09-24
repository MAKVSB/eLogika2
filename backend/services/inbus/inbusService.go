package inbus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/modules/common"
)

// InbusClient is a singleton for calling the Inbus
type InbusClient struct {
	BaseURL string
	client  *http.Client
	token   *TokenData
}

var (
	instance *InbusClient
	once     sync.Once
)

// InitInbusClient initializes the singleton
func InitInbusClient() *InbusClient {
	once.Do(func() {
		instance = &InbusClient{
			BaseURL: initializers.GlobalAppConfig.INBUS_BASE_URL,
			client:  &http.Client{Timeout: 10 * time.Second},
		}

		instance.updateToken()
	})
	return instance
}

// GetInbusClient returns the singleton instance
func GetInbusClient() *InbusClient {
	if instance == nil {
		InitInbusClient()
	}
	return instance
}

type TokenData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func (api *InbusClient) getToken() string {
	if api.token == nil || api.token.ExpiresIn < int(time.Now().Unix()) {
		err := api.updateToken()
		if err != nil {
			panic(err)
		}
	}
	return api.token.AccessToken
}

func (api *InbusClient) updateToken() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", initializers.GlobalAppConfig.INBUS_CLIENT_ID)
	data.Set("client_secret", initializers.GlobalAppConfig.INBUS_CLIENT_SECRET)
	data.Set("scope", "edison edison/schedule")

	req, err := http.NewRequest("POST", api.BaseURL+"/oauth/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("inbus error: %s", resp.Status)
	}

	var token TokenData
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		panic(err)
	}

	api.token = &token

	return nil
}

func (api *InbusClient) DoRequest(endpoint string, params url.Values, result interface{}) error {
	return api.doRequestWithRetry(endpoint, params, result, true)
}

func (api *InbusClient) doRequestWithRetry(endpoint string, params url.Values, result interface{}, allowRetry bool) error {
	token := api.getToken()

	url := api.BaseURL + endpoint + "?" + params.Encode()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized && allowRetry {
		// try refreshing token
		err = api.updateToken()
		if err != nil {
			return fmt.Errorf("token refresh failed: %w", err)
		}
		// retry once
		return api.doRequestWithRetry(endpoint, params, result, false)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("inbus error: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	return nil
}

func (api *InbusClient) GetSemesterFromDate(date string) (*InbusSemester, *common.ErrorResponse) {
	semester := new(InbusSemester)
	err := api.DoRequest("service/edison/v1/edu/semesters/date/"+date, url.Values{}, semester)
	if err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to import data",
			Details: "Request to inbus has failed",
		}
	}
	return semester, nil
}

func (api *InbusClient) GetSubjectVersionFromcode(code string) (*[]*SubjectVersion, *common.ErrorResponse) {
	versions := new([]*SubjectVersion)

	err := api.DoRequest("service/edison/v1/edu/subjectVersions/byCodes", url.Values{
		"code": {code},
	}, versions)
	if err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to import data",
			Details: "Request to inbus has failed",
		}
	}

	return versions, nil
}

func (api *InbusClient) GetConcreteActivities(subjectVersionId uint, semesterId *uint) (*[]*ConcreteActivity, *common.ErrorResponse) {
	activities := new([]*ConcreteActivity) // or a struct

	values := url.Values{
		"subjectVersionId": {strconv.Itoa(int(subjectVersionId))},
		"currentWeek":      {"false"},
	}
	if semesterId != nil {
		values.Add("semesterId", strconv.Itoa(int(*semesterId)))
	}

	err := api.DoRequest("service/edison/v1/schedule/", values, activities)
	if err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to import data",
			Details: "Request to inbus has failed",
		}
	}

	return activities, nil
}

func (api *InbusClient) GetConcreteActivityStudents(concreteActivityId uint) (*[]*StudyRelation, *common.ErrorResponse) {
	students := new([]*StudyRelation)

	err := api.DoRequest("service/edison/v1/schedule/"+strconv.Itoa(int(concreteActivityId))+"/studyRelations", url.Values{}, students)
	if err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to import data",
			Details: "Request to inbus has failed",
		}
	}

	return students, nil
}

func (api *InbusClient) GetSubjectVersionStudents(subjectVersionId uint, semesterId uint) (*[]*StudyRelation, *common.ErrorResponse) {
	students := new([]*StudyRelation)

	err := api.DoRequest("service/edison/v1/admin/psp/subjectVersion/"+strconv.Itoa(int(subjectVersionId))+"/studyRelations", url.Values{
		"semesterId": {strconv.Itoa(int(semesterId))},
	}, students)
	if err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to import data",
			Details: "Request to inbus has failed",
		}
	}

	return students, nil
}
