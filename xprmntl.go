package xprmntl

import (
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"errors"
	"os"
	"time"
	"reflect"
	"strconv"
	"regexp"
)
/**
 REQUEST OBJECTS
 */
/**
 STRUCT: Experiment
 */
type Experiment struct {
	Name string        `json:"name"`;
	Description string `json:"description"`;
	ExpDefault bool    `json:"default"`;
}

/**
 STRUCT: Config
 */
type Config struct {
	DevKey      string       `json:"devKey"`;
	FeatureURL  string;
	Timeout     int;
	Reference   string       `json:"reference"`;
	Experiments []Experiment `json:"experiments"`;
	Shared      *Config      `json:"shared"`;
}

/**
 Config: GET functions
 */
func (c *Config) getDevKey() *string {
	return &c.DevKey;
};

func (c *Config) getFeatureURL() *string {
	return &c.FeatureURL;
};

func (c *Config) getTimeout() int {
	return c.Timeout;
};

func (c *Config) getReference() *string {
	return &c.Reference;
};

func (c *Config) getExperiments() *[]Experiment {
	return &c.Experiments;
};

func (c *Config) getSharedConfig() *Config {
	return c.Shared;
};

/**
 STRUCT: FeatureClient
 */
type FeatureClient struct {
	DevKey       *string       `json:"devKey"`;
	FeatureURL   *string       `json:"featureUrl"`;
	Timeout      int           `json:"timeout"`;
	Reference    *string       `json:"reference"`;
	Experiments  *[]Experiment `json:"experiments"`;
	Shared       *Config       `json:"shared"`;
}

/**
 FeatureClient: Constructors
 */
func New(config *Config) (*FeatureClient, error) {
	if config == nil {
		return nil, errors.New("XPRMNTL: New(): Cannot register experiments without a config. Please see docs.");
	}
	devKey      := config.getDevKey();
	featureURL  := config.getFeatureURL();
	timeout     := config.getTimeout();
	experiments := config.getExperiments();
	shared      := config.getSharedConfig();
	reference   := config.getReference();

	if devKey == nil || len(*devKey) == 0 {
		envKey := os.Getenv("FEATURE_DEVKEY");
		if len(envKey) == 0 {
			return nil, errors.New("XPRMNTL: New(): No devKey defined.");
		}
		devKey = &envKey;
	}

	if featureURL == nil || len(*featureURL) == 0 {
		envUrl := os.Getenv("FEATURE_URL");
		if len(envUrl) == 0 {
			return nil, errors.New("XPRMNTL: New(): No featureUrl defined.");
		}
		featureURL = &envUrl;
	}

	if timeout == 0 {
		timeout = 5000;
	}

	if shared == nil {
		sharedKey := os.Getenv("FEATURE_DEVKEY_SHARED");
		if len(sharedKey) > 0 {
			sharedConfig := Config{ DevKey: sharedKey };
			shared = &sharedConfig;
		}
	}

	if experiments == nil || len(*experiments) == 0 {
		return nil, errors.New("XPRMNTL: New(): Cannot register experiments without `experiments`. Please see the docs.")
	}

	return &FeatureClient{devKey, featureURL, timeout, reference, experiments, shared}, nil;
}

/**
 FeatureClient: Utility
 */
func (c *FeatureClient) Announce() (*AppConfig, error) {
	url := *c.FeatureURL + "/api/coupling/";
	var timeout time.Duration = time.Duration(c.Timeout) * time.Millisecond;
	jsonBody, jsonErr := json.Marshal(c);
	if jsonErr != nil {
		fmt.Println(jsonErr);
		return nil, errors.New("XPRMNTL: Announce(): There was an error building your JSON")
	}

	client := &http.Client{
		Timeout: timeout,
	}

	req, reqErr := http.NewRequest("POST", url, bytes.NewReader(jsonBody));
	if reqErr != nil {
		fmt.Println(reqErr);
		return nil, errors.New("XPRMNTL: Announce(): There was an error in your request");
	}

	req.Header.Add("x-feature-key", *c.DevKey);
	req.Header.Add("Content-Type", "application/json");


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
		fmt.Println(bodyReadErr);
		return nil, errors.New("XPRMNTL: Announce(): There was an error in reading the server response body")
	}
	var response FeatureClientResponse;
	marshalErr := json.Unmarshal(body, &response);
	if marshalErr != nil {
		fmt.Println(marshalErr);
		return nil, errors.New("XPRMNTL: (function) Announce: There was an error in parsing the JSON response")
	}
	response.App.SetReference(c.Reference);
	response.App.SetDefault(c);
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
	Reference   string;
	Groups      interface {};
	Experiments map[string]interface {};
	Default     *FeatureClient;
	userID      int;
	req 				*http.Request;
	resWriter   *http.ResponseWriter;
}

func (app *AppConfig) SetReference(reference *string) {
	app.Reference = *reference;
}
func (app *AppConfig) SetDefault(config *FeatureClient) {
	app.Default = config;
}
func (app *AppConfig) Initialize(req *http.Request, w *http.ResponseWriter) {
	app.req = req;
	app.resWriter = w;
}

func (app *AppConfig) IsSet(experimentName string) bool {
	if app.Experiments[experimentName] != nil {
		switch app.Experiments[experimentName].(type) {
			case bool: {
				return app.Experiments[experimentName].(bool);
			}
			case []interface {}: {
				cookie, err := app.req.Cookie("XPRMNTL-config");
				if err != nil {
					app.userID += 1;
					expire := time.Now().AddDate(0, 0, 1);
					cookie = &http.Cookie {
						Name: "XPRMNTL-config",
						Value: strconv.Itoa(app.userID % 100),
						Path: "/",
						Expires: expire,
					};
					http.SetCookie(*app.resWriter, cookie);
				}
				return parseExperimentVariants(app.Experiments[experimentName].([]interface {}), cookie);
			}
			default: {
				fmt.Println(reflect.TypeOf(app.Experiments[experimentName]));
			}
		}
	}
	// TODO: Setup function to check the Default config for experiment before returning false
	return false;
};

func parseExperimentVariants(config [] interface {}, xprmntlConfig *http.Cookie) bool {
	regex, _ := regexp.Compile(`(\d*)-(\d*)%$`);
	for i := 0; i < len(config); i++ {
		val := config[i].(string);
		if regex.Match([]byte(val)) {
			vals := regex.FindAllStringSubmatch(val, -1);

			lowerLimit, _ := strconv.Atoi(vals[0][1]);
			upperLimit, _ := strconv.Atoi(vals[0][2]);
			bucket, _     := strconv.Atoi(xprmntlConfig.Value);

			return bucket >= lowerLimit && bucket < upperLimit;
		} else {
			fmt.Println("Parsing Group", val);
		}
	}
	return true;
}
