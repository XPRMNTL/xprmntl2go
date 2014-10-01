[![XPRMNTL](https://raw.githubusercontent.com/XPRMNTL/XPRMNTL.github.io/master/images/ghLogo.png)](https://github.com/XPRMNTL/XPRMNTL.github.io)
# xprmntl2go
[![Build Status](https://travis-ci.org/XPRMNTL/xprmntl2go.svg?branch=master)](https://travis-ci.org/XPRMNTL/feature-client.js)

This is a GoLang library for the consumption of [XPRMNTL](https://github.com/XPRMNTL/feature) product.

```go
package main

import (
	"fmt"
	"net/http"
	"github.com/xprmntl/xprmntl2go"
)

func main() {
	config := xprmntl2go.Config {
		Experiments: []*xprmntl2go.Experiment{
			&xprmntl2go.Experiment{Name: "TestExp"},
		},
		Shared: &xprmntl2go.Config {
			DevKey: "testDevKey",
			Experiments: []*xprmntl2go.Experiment{
				&xprmntl2go.Experiment{Name: "SharedTestExp"},
			},
		},
	};
	cli, cliErr := xprmntl2go.New(&config);

	if cliErr != nil {
		fmt.Println(cliErr);
		return;
	}

	experiments, err := cli.Announce();
	if err != nil {
		fmt.Println(err);
	}

	handler := func(w http.ResponseWriter, req *http.Request) {
		// Initialize the experiments object
		experiments.Initialize(&w, req);

		if experiments.IsSet("TestExp") {
			// Execute body if experiment is set
		}
	};

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

### Installation and Importing
```sh
$ go get github.com/xprmntl/xprmntl2go
```

```go
import (
  . "github.com/xprmntl/xprmntl2go"
)
```

### Configuration
Config Object:
```go
type Config struct {
	DevKey      string `json:"devKey"`;
	FeatureURL  string;
	Timeout     int;
	Experiments []Experiment `json:"experiments"`;
	Shared      *Config `json:"shared"`;
}
```
- `FeatureURL`
  - this is the URL to the XPRMNTL dashboard.
  - Defaults to `os.Getenv("FEATURE_URL")`.
- `DevKey`
  - this is the devKey generated for you by the XPRMNTL dashboard.
  - Defaults to `os.Getenv("FEATURE_DEVKEY")`
- `Experiments`
  - This is an array of all of your app-level experiments. This must be an Experiment object:
```go
type Experiment struct {
	Name string        `json:"name"`;
	Description string `json:"description"`;
	ExpDefault bool    `json:"default"`;
}
```
- `Timeout`
  - This is the number of milliseconds after which the request should time out.
  - Must be set to a non-zero value
  - Defaults to 5000 (5s)
- `Shared`
  - This object allows you to configure and accept configuration for a shared set of experiments. If, for example, 
  you have a separate set of experiments for your site-wide theme, you would configure those here, 
  shared among your applications.
  - This also must be a Config object (NOTE: only the DevKey and the Experiments array will be used here)

### Announcement
This step preforms the fetching of the configuration against the XPRMNTL dashboard. It is a feature of the 
`xprmntl2go.FeatureClient` object Any new experiments are registered and default either to `false` or to whatever 
you've set as your `default` for that experiment. If there are any errors in the Announcement the second return value
 of this will contain the error whilst the first will contain an `xprmntl2go.AppConfig` that contains the config 
 defaults
```go
experiments, err := cli.Announce();
if err != nil {
  fmt.Println(err);
}
```

### Initialization
In order to utilize the groups and buckets features of the feature dashboard you must initialize the app with the 
response writer and req object. If you don't intend to use those features you can skip this step.
 ```go
handler := func(w http.ResponseWriter, req *http.Request) {
  // Initialize the experiments object
  experiments.Initialize(&w, req);
};
 ```
 
### Usage

#### Server Usage:
```go
handler := func(w http.ResponseWriter, req *http.Request) {
  // Initialize the experiments object
  experiments.Initialize(&w, req);
  
  if experiments.IsSet("TestExp") {
    // Execute Code is TestExp is active
  } else {
    // Execute Code if TestExp is inactive
  }
};
```

#### Template Usage: 
```html
{{ if .IsSet "TestExp"  }}
<h2>TestExp</h2>
<p>The experiment 'TextExp' is on</p>
{{ end }}
```
 
