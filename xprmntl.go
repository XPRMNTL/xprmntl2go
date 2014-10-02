// feature-client.go is a client implementation of the XPRMNTL service
package xprmntl2go

import (
	"log"
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

type Experiment struct {
	Name string        `json:"name"`;
	Description string `json:"description"`;
	ExpDefault bool    `json:"default"`;
}

type Config struct {
	DevKey      string        `json:"devKey"`;
	FeatureURL  string				`json:"featureUrl"`;
	Timeout     int						`json:"timeout"`;
	Reference   string        `json:"reference"`;
	Experiments []*Experiment `json:"experiments"`;
	Shared      *Config       `json:"shared"`;
}

// Get a pointer to the config devKey
func (c *Config) getDevKey() *string {
	return &c.DevKey;
};

// Get a pointer to the config featureURL
func (c *Config) getFeatureURL() *string {
	return &c.FeatureURL;
};

// Get the config timeout
func (c *Config) getTimeout() int {
	return c.Timeout;
};

// Get a pointer to the config reference
func (c *Config) getReference() *string {
	return &c.Reference;
};

// Get the config experiments list
func (c *Config) getExperiments() []*Experiment {
	return c.Experiments;
};

// Get the shared config
func (c *Config) getSharedConfig() *Config {
	return c.Shared;
};

type FeatureClient struct {
	DevKey       *string       `json:"devKey"`;
	FeatureURL   *string       `json:"featureUrl"`;
	Timeout      int           `json:"timeout"`;
	Reference    *string       `json:"reference"`;
	Experiments  []*Experiment `json:"experiments"`;
	Shared       *Config       `json:"shared"`;
}

// Creates a New FeatureClient object from the provided config object. Will utilize ENV defaults as necessary
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

	if experiments == nil || len(experiments) == 0 {
		return nil, errors.New("XPRMNTL: New(): Cannot register experiments without `experiments`. Please see the docs.")
	}

	return &FeatureClient{devKey, featureURL, timeout, reference, experiments, shared}, nil;
}

// Get a pointer to the name specific experiment in a FeatureClient
func (c *FeatureClient) getExp(name string) *Experiment {
	for i := 0; i < len(c.Experiments); i++ {
		if c.Experiments[i].Name == name { return c.Experiments[i]; }
	}
	return nil;
}

// Initialize the Feature Client dashboard (if available), if not return empty app that falls back to defaults
func (c *FeatureClient) Announce() (*AppConfig, error) {
	var response FeatureClientResponse;
	response.App.SetDefault(c);

	url := *c.FeatureURL + "/api/coupling/";
	var timeout time.Duration = time.Duration(c.Timeout) * time.Millisecond;
	jsonBody, jsonErr := json.Marshal(c);
	if jsonErr != nil {
		log.Print(jsonErr);
		return &response.App, errors.New("XPRMNTL: Announce(): There was an error building your JSON")
	}

	client := &http.Client{
		Timeout: timeout,
	}

	req, reqErr := http.NewRequest("POST", url, bytes.NewReader(jsonBody));
	if reqErr != nil {
		log.Print(reqErr);
		return &response.App, errors.New("XPRMNTL: Announce(): There was an error in your request");
	}

	req.Header.Add("x-feature-key", *c.DevKey);
	req.Header.Add("Content-Type", "application/json");


	res, resErr := client.Do(req);
	if resErr != nil {
		log.Print(resErr);
		return &response.App, errors.New("XPRMNTL: Announce(): There was an error in the server response")
	}
	
	if (res.StatusCode != 200) {
		log.Print(res);
		return &response.App, errors.New("XPRMNTL: Announce(): Server return a non-200 response");
	}
	body, bodyReadErr := ioutil.ReadAll(res.Body);
	defer res.Body.Close();
	if bodyReadErr != nil {
		log.Print(bodyReadErr);
		return &response.App, errors.New("XPRMNTL: Announce(): There was an error in reading the server response body")
	}
	marshalErr := json.Unmarshal(body, &response);
	if marshalErr != nil {
		log.Print(marshalErr);
		return &response.App, errors.New("XPRMNTL: (function) Announce: There was an error in parsing the JSON response")
	}
	response.App.SetReference(c.Reference);
	return &response.App, nil;
}

type FeatureClientResponse struct {
	App AppConfig
}

type AppConfig struct {
	Reference   string;
	Groups      interface {};
	Experiments map[string]interface {};
	Default     *FeatureClient;
	userID      int;
	req 				*http.Request;
	resWriter   *http.ResponseWriter;
}

// Set the Reference for an app
func (app *AppConfig) SetReference(reference *string) {
	app.Reference = *reference;
}

// Set the default config for an App takes a FeatureClient object
func (app *AppConfig) SetDefault(config *FeatureClient) {
	app.Default = config;
}

// Initialize the app so that it can set and check cookies for specific values
func (app *AppConfig) Initialize(w *http.ResponseWriter, req *http.Request) {
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
				log.Print(reflect.TypeOf(app.Experiments[experimentName]));
			}
		}
	}
	expDefault := app.Default.getExp(experimentName);
	if expDefault != nil {
		return expDefault.ExpDefault;
	}
	return false;
};

// Parse the experiment variants and determine if an experiment is on or off.
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
			log.Print("Parsing Group", val);
		}
	}
	return true;
}
