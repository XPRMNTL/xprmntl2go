package xprmntl

import (
	"fmt"
	"os"
	"github.com/franela/goreq"
	"time"
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
//	if len(name) == 0 { return Experiment{"", "", false}, errors.New("No experiment name defined.") };
	return Experiment{name, "", false};
}

func NewExperimentsList(experiments ...interface {}) ([]Experiment) {
	list := []Experiment {};
	for i := 0; i < len(experiments); i++ {
		switch experiments[i].(type) {
			case string: {
				experiment := NewExperiment(experiments[i].(string));
//				if err != nil {
//					return []Experiment {}, errors.New("");
//				}
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
//	if len(devKey) == 0 { return nil, errors.New("No devKey defined."); }
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
	timeout int64;
	Experiments []Experiment `json:"experiments"`;
	Shared *SharedConfig `json:"shared"`;
}

/**
 FeatureConfig: Constructors
 */
func NewFeatureConfig(experiments []Experiment) (FeatureConfig){
//	if len(experiments) == 0 { return nil, errors.New("Cannot register experiments without `experiments`. Please see the docs")};
	return FeatureConfig{os.Getenv("FEATURE_DEVKEY"), os.Getenv("FEATURE_URL"), 30000, experiments, nil};
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

func (c FeatureConfig) getTimeout() int64 {
	return c.timeout;
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
	Timeout int64 `json:"timeout"`;
	Config FeatureConfig `json:"config"`;
}
/**
 FeatureClient: Constructors
 */
func New(config FeatureConfig) (FeatureClient) {
	devKey     := config.getDevKey();
	featureURL := config.getFeatureURL();
	timeout    := config.getTimeout();

	if len(devKey) == 0 {
		devKey = os.Getenv("FEATURE_DEVKEY");
	}
	if len(featureURL) == 0 {
		devKey = os.Getenv("FEATURE_URL")
	}

//	if len(devKey) == 0 {
//		return nil, errors.New("");
//	}
//	if len(featureURL) == 0 {
//		return nil, errors.New("");
//	}
	return FeatureClient{devKey, featureURL, timeout, config};
}

/**
 FeatureClient: Utility
 */
func (c FeatureClient) Announce() {
	url := c.FeatureURL + "api/coupling/";
	var timeout time.Duration = time.Duration(c.Timeout) * time.Millisecond;
	req := goreq.Request{
		Method: "POST",
		Uri: url,
		Body: c,
		Timeout: timeout}
	req.AddHeader("x-feature-key", c.DevKey);
	res, err := req.Do();
	fmt.Println(res);
	fmt.Println(err);
	return;
}
/**
 END STRUCT
 */
