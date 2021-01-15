// Copyright 2021 Cambricon, Inc.
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

package collector

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcVFUtil(t *testing.T) {
	tests := []struct {
		utils []uint
		vfNum int
		vfID  int
		util  float64
	}{
		{
			vfNum: 0,
		},
		{
			utils: []uint{13, 13, 13, 13, 14, 14, 14, 14, 15, 15, 15, 15, 16, 16, 16, 16},
			vfNum: 1,
			vfID:  1,
			util:  float64(14.5),
		},
		{
			utils: []uint{13, 13, 13, 13, 14, 14, 14, 14, 15, 15, 15, 15, 16, 16, 16, 16},
			vfNum: 2,
			vfID:  2,
			util:  float64(15.5),
		},
		{
			utils: []uint{13, 13, 13, 13, 14, 14, 14, 14, 15, 15, 15, 15, 16, 16, 16, 16},
			vfNum: 4,
			vfID:  3,
			util:  float64(15),
		},
	}
	for i, tt := range tests {
		t.Run("", func(t *testing.T) {
			util, err := calcVFUtil(tt.utils, tt.vfNum, tt.vfID)
			if i == 0 {
				assert.Equal(t, err, fmt.Errorf("invalid vfNum %d", tt.vfNum))
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.util, util)
		})
	}
}
