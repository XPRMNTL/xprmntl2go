package xprmntl

import (
	"fmt"
	"net/http"
//	"time"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"errors"
	"os"
)
/**
 REQUEST OBJECTS
 */
/**
 STRUCT: Experiment
 */
type Experiment struct {
	Name string `json:"name"`;
	Description string `json:"description"`;
	ExpDefault bool `json:"default"`;
}

/**
 STRUCT: Config
 */
type Config struct {
	DevKey      string `json:"devKey"`;
	FeatureURL  string;
	Timeout     int;
	Experiments []Experiment `json:"experiments"`;
	Shared      *Config `json:"shared"`;
}

/**
 Config: GET functions
 */
func (c Config) getDevKey() *string {
	return &c.DevKey;
};

func (c Config) getFeatureURL() *string {
	return &c.FeatureURL;
};

func (c Config) getTimeout() int {
	return c.Timeout;
};

func (c Config) getExperiments() *[]Experiment {
	return &c.Experiments;
};

func (c Config) getSharedConfig() *Config {
	return c.Shared;
};

/**
 STRUCT: FeatureClient
 */
type FeatureClient struct {
	DevKey      *string       `json:"devKey"`;
	FeatureURL  *string       `json:"featureUrl"`;
	Timeout     int           `json:"timeout"`;
	Experiments *[]Experiment `json:"experiments"`;
	Shared      *Config `json:"shared"`;
}

/**
 FeatureClient: Constructors
 */
func New(config Config) (*FeatureClient, error) {
	devKey      := config.getDevKey();
	featureURL  := config.getFeatureURL();
	timeout     := config.getTimeout();
	experiments := config.getExperiments();
	shared      := config.getSharedConfig();

	if devKey == nil || len(*devKey) == 0 {
		envKey := os.Getenv("FEATURE_DEVKEY");
		if len(envKey) == 0 {
			return nil, errors.New("XPRMNTL: New(): XPRMNTL requires a devKey to be set");
		}
		devKey = &envKey;
	}

	if featureURL == nil || len(*featureURL) == 0 {
		envUrl := os.Getenv("FEATURE_URL");
		if len(envUrl) == 0 {
			return nil, errors.New("XPRMNTL: New(): XPRMNTL requires a featureUrl to be set");
		}
		featureURL = &envUrl;
	}

	return &FeatureClient{devKey, featureURL, timeout, experiments, shared}, nil;
}

/**
 FeatureClient: Utility
 */
func (c FeatureClient) Announce() (*AppConfig, error) {
	url := *c.FeatureURL + "/api/coupling/";
//	var timeout time.Duration = time.Duration(c.Timeout) * time.Millisecond;
	jsonBody, jsonErr := json.Marshal(c);
	if jsonErr != nil {
		return nil, errors.New("XPRMNTL: Announce(): There was an error building your JSON")
	} else {
		fmt.Println(string(jsonBody[:]))
	}
	client := &http.Client{}

	req, reqErr := http.NewRequest("POST", url, bytes.NewReader(jsonBody));
	req.Header.Add("x-feature-key", *c.DevKey);
	req.Header.Add("Content-Type", "application/json");
	if reqErr != nil {
		return nil, errors.New("XPRMNTL: Announce(): There was an error in your request")
	}
	res, resErr := client.Do(req);
	if resErr != nil {
		fmt.Println(resErr);
		return nil, errors.New("XPRMNTL: Announce(): There was an error in the server response")
	}
	
	if (res.StatusCode != 200) {
		fmt.Println(res);
		return nil, errors.New("XPRMNTL: Announce(): Server return a non-200 response");
	}
	body, bodyReadErr := ioutil.ReadAll(res.Body);
	defer res.Body.Close();
	if bodyReadErr != nil {
		return nil, errors.New("XPRMNTL: Announce(): There was an error in reading the server response body")
	}

	var response FeatureClientResponse;
	marshalErr := json.Unmarshal(body, &response);
	if marshalErr != nil {
		return nil, errors.New("XPRMNTL: (function) Announce: There was an error in parsing the JSON response")
	}
	return &response.App, nil;
}
/**
 END STRUCT
 */


/**
 RESPONSE OBJECTS
 */
/**
 STRUCT: FeatureClientResponse
 */
type FeatureClientResponse struct {
	App AppConfig
}

/**
 STRUCT: AppConfig
 */
type AppConfig struct {
	Groups interface {};
	Experiments map[string]interface {};
	Envs interface {};
}
