[![XPRMNTL](https://raw.githubusercontent.com/XPRMNTL/XPRMNTL.github.io/master/images/ghLogo.png)](https://github.com/XPRMNTL/XPRMNTL.github.io)
# xprmntl2go
[![Build Status](https://travis-ci.org/XPRMNTL/xprmntl2go.svg?branch=master)](https://travis-ci.org/XPRMNTL/feature-client.js)

This is a GoLang library for the consumption of [XPRMNTL](https://github.com/XPRMNTL/feature) product.

```go
package main

import . "github.com/xprmntl/xprmntl2go"

func main() {
	config := xprmntl.Config {
		Experiments: []xprmntl.Experiment{
			xprmntl.Experiment{Name: "TestExp"},
		},
		Shared: &xprmntl.Config {
			DevKey: "testDevKey",
			Experiments: []xprmntl.Experiment{
				xprmntl.Experiment{Name: "SharedTestExp"},
			},
		},
	};
	cli, cliErr := xprmntl.New(&config);

	if cliErr != nil {
		fmt.Println(cliErr);
		return;
	}

	experiments, err := cli.Announce();
	if err != nil {
		fmt.Println(err);
		return;
	}
	handler := func(w http.ResponseWriter, req *http.Request) {
		// Initialize the experiments object
		experiments.Initialize(req, &w);
		
		if experiments.IsSet("TestExp") {
			// Execute body if experiment is set
		}
		
		t.Execute(w, experiments);
	};

  http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}
```
Template Implementation
```html
{{ if .IsSet "TestExp"  }}
<h2>Big Exp</h2>
<p>The experiment 'Big Exp' is on</p>
<img src="http://lorempixel.com/400/400" />
{{ end }}
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

### Usage
After importing the package the functions are accesible under the `xprmntl` namespace. 


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
  - This object allows you to configure and accept configuration for a shared set of experiments. If, for example, you have a separate set of experiments for your site-wide theme, you would configure those here, shared among your applications.
  - This also must be a Config object (NOTE: only the DevKey and the Experiments array will be used here)

### Announcement
This step preforms the fetching of the configuration against the XPRMNTL dashboard. Any new experiments get registered and default either to `false` or to whatever you've set as your `default` for that experiment.
```go

```
