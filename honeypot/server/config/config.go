/*

Copyright (c) 2018 sec.lu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type (
	ServiceItem struct {
		Addr  string `json:"addr" yaml:"addr"`
		Proxy string `json:"proxy" yaml:"proxy"`
		Flag  bool   `json:"flag" yaml:"flag"`
	}

	ApiCnf struct {
		Addr string `json:"addr" yaml:"addr"`
		Key  string `json:"key" yaml:"key"`
	}
	ProxyCnf struct {
		Flag bool   `json:"flag" yaml:"flag"`
		Addr string `json:"addr" yaml:"addr"`
	}

	Config struct {
		Proxy    ProxyCnf               `json:"proxy" yaml:"proxy"`
		Services map[string]ServiceItem `json:"services" yaml:"services"`
		Api      ApiCnf                 `json:"api" yaml:"api"`
	}
)

func ReadConfig() (Config, error) {
	var config Config
	curDir, err := GetCurDir()
	if err != nil {
		return config, err
	}
	configFile := filepath.Join(curDir, "conf", "config.yaml")
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	fmt.Printf("service: %v, err: %v\n", config.Services, err)
	fmt.Printf("proxy: %v, err: %v\n", config.Proxy, err)
	fmt.Printf("api: %v, err: %v\n", config.Api, err)

	return config, err
}

func GetCurDir() (string, error) {
	dir, err := filepath.Abs(filepath.Dir("./"))
	if err != nil {
		return "", err
	}
	return dir, err
}
