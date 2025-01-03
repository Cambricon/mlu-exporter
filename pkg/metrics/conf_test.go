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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	f, err := os.CreateTemp("", "config.*.yaml")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	f.Write([]byte(`
cndev:
  temperature:
    name: "temperature test"
    help: "help message"
    labels:
      slot: "slot"
      model: "model"
`,
	))
	c := getConfOrDie(f.Name())
	expected := Conf{
		"cndev": {
			"temperature": info{
				Name: "temperature test",
				Help: "help message",
				Labels: Labels{
					"slot":  "slot",
					"model": "model",
				},
			},
		},
	}
	assert.Equal(t, c, expected)
}
