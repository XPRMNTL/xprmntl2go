package xprmntl

import (
	"fmt"
	"os"
	"net/http"
//	"time"
	"encoding/json"
	"bytes"
	"io/ioutil"
//	"errors"
)

/**
 STRUCT: Experiment
 */
type Experiment struct {
	Name string `json:"name"`;
	Description string `json:"description"`;
	ExpDefault bool `json:"default"`;
}

func NewExperiment(name string) (Experiment) {
	return Experiment{name, "", false};
}

func NewExperimentsList(experiments ...interface {}) ([]Experiment) {
	list := []Experiment {};
	for i := 0; i < len(experiments); i++ {
		switch experiments[i].(type) {
			case string: {
				experiment := NewExperiment(experiments[i].(string));
				list = append(list, experiment);
			}
			case Experiment: {
				list = append(list, experiments[i].(Experiment));
			}
		}
	}
	return list;
}
/**
 END STRUCT
 */

/**
 STRUCT: SharedConfig
 */
type SharedConfig struct {
	DevKey string `json:"devKey"`;
	Experiments []Experiment `json:"experiments"`;
}

func NullSharedConfig() (SharedConfig) {
	return SharedConfig{"", []Experiment{}};
}

func NewSharedConfig(devKey string) (SharedConfig) {
	return SharedConfig{devKey, []Experiment{}};
}
/**
 END STRUCT
 */

/**
 STRUCT: FeatureConfig
 */
type FeatureConfig struct {
	devKey string;
	featureURL string;
	timeout int;
	experiments []Experiment `json:"experiments"`;
	shared SharedConfig `json:"shared"`;
}

/**
 FeatureConfig: Constructors
 */
func NewFeatureConfig(experiments []Experiment) (FeatureConfig){
	return FeatureConfig{os.Getenv("FEATURE_DEVKEY"), os.Getenv("FEATURE_URL"), 30000, experiments, NullSharedConfig()};
}

/**
 FeatureConfig: GET functions
 */
func (c FeatureConfig) getDevKey() string {
	return c.devKey;
};

func (c FeatureConfig) getFeatureURL() string {
	return c.featureURL;
};

func (c FeatureConfig) getTimeout() int {
	return c.timeout;
};

func (c FeatureConfig) getExperiments() []Experiment {
	return c.experiments;
};

func (c FeatureConfig) getSharedConfig() SharedConfig {
	return c.shared;
};

/**
 END STRUCT
 */

/**
 STRUCT: FeatureClient
 */
type FeatureClient struct {
	DevKey string `json:"devKey"`;
	FeatureURL string `json:"featureUrl"`;
	Timeout int `json:"timeout"`;
	Experiments []Experiment `json:"experiments"`;
	Shared SharedConfig `json:"shared"`;
}
/**
 FeatureClient: Constructors
 */
func New(config FeatureConfig) (FeatureClient) {
	devKey       := config.getDevKey();
	featureURL   := config.getFeatureURL();
	timeout      := config.getTimeout();
	experiments  := config.getExperiments();
	sharedConfig := config.getSharedConfig();

	if len(devKey) == 0 {
		devKey = os.Getenv("FEATURE_DEVKEY");
	}
	if len(featureURL) == 0 {
		devKey = os.Getenv("FEATURE_URL")
	}
	return FeatureClient{devKey, featureURL, timeout, experiments, sharedConfig};
}

/**
 FeatureClient: Utility
 */
func (c FeatureClient) Announce() App {
	url := c.FeatureURL + "api/coupling/";
//	var timeout time.Duration = time.Duration(c.Timeout) * time.Millisecond;
	jsonBody, err := json.Marshal(c);
	if err != nil {
		fmt.Println("There was an Error");
	}
	client := &http.Client{}

	req, reqErr := http.NewRequest("POST", url, bytes.NewReader(jsonBody));
	req.Header.Add("x-feature-key", c.DevKey);
	req.Header.Add("Content-Type", "application/json");
	if reqErr != nil {
		fmt.Println("There was an Error");
	}
	res, resErr := client.Do(req);
	defer res.Body.Close();
	body, bodyReadErr := ioutil.ReadAll(res.Body);
	if resErr != nil {
		fmt.Println("There was an Error");
	}
	if bodyReadErr != nil {
		fmt.Println("There was an Error");
	}

	var response FeatureClientResponse;
//	fmt.Println(string(body[:]))
	marshalErr := json.Unmarshal(body, &response);
	if marshalErr != nil {
		fmt.Println("json.Unmarshal ERROR");
	}
	return response.App;
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
	App App
}

type App struct {
	Groups interface {};
	Experiments map[string]interface {};
	Envs interface {};
}

/**
 END STRUCT
 */
