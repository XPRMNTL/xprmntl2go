[![XPRMNTL](https://raw.githubusercontent.com/XPRMNTL/XPRMNTL.github.io/master/images/ghLogo.png)](https://github.com/XPRMNTL/XPRMNTL.github.io)
# Feature-Client.js
[![Build Status](https://travis-ci.org/jshcrowthe/xprmntl.svg?branch=master)](https://travis-ci.org/XPRMNTL/feature-client.js) [![NPM version](https://img.shields.io/npm/v/feature-client.svg)](https://www.npmjs.org/package/feature-client)

This is a GoLang library for the consumption of [XPRMNTL](https://github.com/XPRMNTL/feature) product.

```go
package main

import . "github.com/xprmntl/xprmntl"

func main() {
	config := xprmntl.Config {
  		Experiments: []xprmntl.Experiment{
  			xprmntl.Experiment{
          Name: "experimentName",
          Description: "Experiment Description",
          ExpDefault: true
  			},
  		},
  		Shared: &xprmntl.Config {
  			DevKey: "testDevKey",
  			Experiments: []xprmntl.Experiment{
  				xprmntl.Experiment{Name: "Big Scary Experiment"},
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
  
  if experiments.isSet("experimentName") {
    // Feature Specific Code
  }
```

### Installation and Importing
```sh
$ go get github.com/jshcrowthe/xprmntl
```

```go
import (
  . "github.com/jshcrowthe/xprmntl"
)
```

### Usage


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
- `featureUrl`
  - this is the URL to the XPRMNTL dashboard.
  - Defaults to `os.Getenv("FEATURE_URL")`.
- `devKey`
  - this is the devKey generated for you by the XPRMNTL dashboard.
  - Defaults to `os.Getenv("FEATURE_DEVKEY")`
- `experiments`
  - This is an array of all of your app-level experiments. This must be an Experiment object:
```go
type Experiment struct {
	Name string        `json:"name"`;
	Description string `json:"description"`;
	ExpDefault bool    `json:"default"`;
}
```
- `timeout`
  - This is the number of milliseconds after which the request should time out.
  - Must be set to a non-zero value
  - Defaults to 5000 (5s)
- `shared`
  - This object allows you to configure and accept configuration for a shared set of experiments. If, for example, you have a separate set of experiments for your site-wide theme, you would configure those here, shared among your applications.
  - This also must be a Config object (NOTE: only the DevKey and the Experiments array will be used here)

### Announcement
This step preforms the fetching of the configuration against the XPRMNTL dashboard. Any new experiments get registered and default either to `false` or to whatever you've set as your `default` for that experiment.
```go

```