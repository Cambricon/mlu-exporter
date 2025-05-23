// Copyright 2020 Cambricon, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Conf map[string]map[string]info

type info struct {
	Name   string `yaml:"name"`
	Help   string `yaml:"help"`
	Labels Labels `yaml:"labels"`
	Push   bool   `yaml:"push"`
}

type Labels map[string]string

func (l Labels) keys() []string {
	keys := []string{}
	for k := range l {
		keys = append(keys, k)
	}
	return keys
}

func load(c *Conf, path string) error {
	f, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		log.Errorf("Can not find config file %s", path)
	}
	if err != nil {
		return err
	}
	return yaml.Unmarshal(f, c)
}

func getConfOrDie(path string) Conf {
	var conf Conf
	if err := load(&conf, path); err != nil {
		log.Fatal(err)
	}
	actual, err := yaml.Marshal(conf)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Get config:\n%s", actual)
	return conf
}
