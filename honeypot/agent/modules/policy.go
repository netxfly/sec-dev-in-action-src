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

package modules

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"sec-dev-in-action-src/honeypot/agent/models"
	"sec-dev-in-action-src/honeypot/agent/vars"
)

// read policy data from local yaml file
func ReadPolicyFromYaml() (data *models.PolicyData, err error) {
	data = new(models.PolicyData)
	var content []byte

	content, err = ioutil.ReadFile(filepath.Join(vars.CurDir, "conf", "policy.yaml"))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, data)
	return data, err

}

func LoadPolicy() (*models.PolicyData, error) {
	var err error
	vars.PolicyData, err = ReadPolicyFromYaml()
	if err != nil {
		return nil, err
	}

	if len(vars.PolicyData.Service) > 0 {
		for _, service := range vars.Services {
			vars.Services = append(vars.Services, service)
		}
	}

	if len(vars.PolicyData.Policy) > 0 {
		vars.HoneypotPolicy = vars.PolicyData.Policy[0]
	}

	return vars.PolicyData, err
}
