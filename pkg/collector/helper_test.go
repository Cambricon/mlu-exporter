package collector

import (
	"testing"

	"github.com/Cambricon/mlu-exporter/pkg/mock"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEnsureMLUAllOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mcndev := mock.NewCndev(ctrl)
	mcndev.EXPECT().Init(false).Return(nil).AnyTimes()
	mcndev.EXPECT().Init(true).Return(nil).AnyTimes()
	mcndev.EXPECT().GetDeviceCount().Return(uint(2), nil).AnyTimes()
	mcndev.EXPECT().GenerateDeviceHandleMap(uint(2)).Return(nil).AnyTimes()

	patches := gomonkey.ApplyFunc(fetchMLUCounts, func() (uint, error) {
		return 2, nil
	})
	defer patches.Reset()
	t.Run("All ok", func(t *testing.T) {
		mluInfo := &MLUStatMap{}
		mluInfo.StatMap.Store("uuid-0", MLUStat{model: "mlu0"})
		patches.ApplyFunc(isDriverRunning, func() bool {
			return false
		})
		patches.ApplyFunc(collectMLUInfo, func(info *MLUStatMap) {
			info.InProblem.Store(false)
			info.StatMap.Store("uuid-0", MLUStat{model: "mlu1"})
		})

		for _, stat := range mluInfo.Range {
			assert.Equal(t, "mlu0", stat.model)
		}

		EnsureMLUAllOK(mcndev, mluInfo, true)

		assert.False(t, mluInfo.InProblem.Load())
		for _, stat := range mluInfo.Range {
			assert.Equal(t, "mlu1", stat.model)
		}
	})

	t.Run("Label missing", func(t *testing.T) {
		mluInfo := &MLUStatMap{}

		cn := 0
		patches.ApplyFunc(isDriverRunning, func() bool {
			cn++
			return cn >= 3
		})

		patches.ApplyFunc(collectMLUInfo, func(info *MLUStatMap) {
			if cn < 3 {
				info.InProblem.Store(true)
				return
			}
			info.InProblem.Store(false)
		})

		EnsureMLUAllOK(mcndev, mluInfo, true)
		assert.True(t, mluInfo.InProblem.Load())

		EnsureMLUAllOK(mcndev, mluInfo, false)
		assert.False(t, mluInfo.InProblem.Load())
	})
}
