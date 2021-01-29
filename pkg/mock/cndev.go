// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Cambricon/mlu-exporter/pkg/cndev (interfaces: Cndev)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// Cndev is a mock of Cndev interface
type Cndev struct {
	ctrl     *gomock.Controller
	recorder *CndevMockRecorder
}

// CndevMockRecorder is the mock recorder for Cndev
type CndevMockRecorder struct {
	mock *Cndev
}

// NewCndev creates a new mock instance
func NewCndev(ctrl *gomock.Controller) *Cndev {
	mock := &Cndev{ctrl: ctrl}
	mock.recorder = &CndevMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Cndev) EXPECT() *CndevMockRecorder {
	return m.recorder
}

// GetDeviceClusterCount mocks base method
func (m *Cndev) GetDeviceClusterCount(arg0 uint) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceClusterCount", arg0)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceClusterCount indicates an expected call of GetDeviceClusterCount
func (mr *CndevMockRecorder) GetDeviceClusterCount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceClusterCount", reflect.TypeOf((*Cndev)(nil).GetDeviceClusterCount), arg0)
}

// GetDeviceCoreNum mocks base method
func (m *Cndev) GetDeviceCoreNum(arg0 uint) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceCoreNum", arg0)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceCoreNum indicates an expected call of GetDeviceCoreNum
func (mr *CndevMockRecorder) GetDeviceCoreNum(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceCoreNum", reflect.TypeOf((*Cndev)(nil).GetDeviceCoreNum), arg0)
}

// GetDeviceCount mocks base method
func (m *Cndev) GetDeviceCount() (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceCount")
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceCount indicates an expected call of GetDeviceCount
func (mr *CndevMockRecorder) GetDeviceCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceCount", reflect.TypeOf((*Cndev)(nil).GetDeviceCount))
}

// GetDeviceFanSpeed mocks base method
func (m *Cndev) GetDeviceFanSpeed(arg0 uint) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceFanSpeed", arg0)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceFanSpeed indicates an expected call of GetDeviceFanSpeed
func (mr *CndevMockRecorder) GetDeviceFanSpeed(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceFanSpeed", reflect.TypeOf((*Cndev)(nil).GetDeviceFanSpeed), arg0)
}

// GetDeviceHealth mocks base method
func (m *Cndev) GetDeviceHealth(arg0 uint) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceHealth", arg0)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceHealth indicates an expected call of GetDeviceHealth
func (mr *CndevMockRecorder) GetDeviceHealth(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceHealth", reflect.TypeOf((*Cndev)(nil).GetDeviceHealth), arg0)
}

// GetDeviceMemory mocks base method
func (m *Cndev) GetDeviceMemory(arg0 uint) (uint, uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceMemory", arg0)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(uint)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceMemory indicates an expected call of GetDeviceMemory
func (mr *CndevMockRecorder) GetDeviceMemory(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceMemory", reflect.TypeOf((*Cndev)(nil).GetDeviceMemory), arg0)
}

// GetDeviceModel mocks base method
func (m *Cndev) GetDeviceModel(arg0 uint) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceModel", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDeviceModel indicates an expected call of GetDeviceModel
func (mr *CndevMockRecorder) GetDeviceModel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceModel", reflect.TypeOf((*Cndev)(nil).GetDeviceModel), arg0)
}

// GetDevicePCIeID mocks base method
func (m *Cndev) GetDevicePCIeID(arg0 uint) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDevicePCIeID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDevicePCIeID indicates an expected call of GetDevicePCIeID
func (mr *CndevMockRecorder) GetDevicePCIeID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDevicePCIeID", reflect.TypeOf((*Cndev)(nil).GetDevicePCIeID), arg0)
}

// GetDevicePower mocks base method
func (m *Cndev) GetDevicePower(arg0 uint) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDevicePower", arg0)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDevicePower indicates an expected call of GetDevicePower
func (mr *CndevMockRecorder) GetDevicePower(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDevicePower", reflect.TypeOf((*Cndev)(nil).GetDevicePower), arg0)
}

// GetDeviceSN mocks base method
func (m *Cndev) GetDeviceSN(arg0 uint) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceSN", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceSN indicates an expected call of GetDeviceSN
func (mr *CndevMockRecorder) GetDeviceSN(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceSN", reflect.TypeOf((*Cndev)(nil).GetDeviceSN), arg0)
}

// GetDeviceTemperature mocks base method
func (m *Cndev) GetDeviceTemperature(arg0 uint) (uint, []uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceTemperature", arg0)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].([]uint)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceTemperature indicates an expected call of GetDeviceTemperature
func (mr *CndevMockRecorder) GetDeviceTemperature(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceTemperature", reflect.TypeOf((*Cndev)(nil).GetDeviceTemperature), arg0)
}

// GetDeviceUtil mocks base method
func (m *Cndev) GetDeviceUtil(arg0 uint) (uint, []uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceUtil", arg0)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].([]uint)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceUtil indicates an expected call of GetDeviceUtil
func (mr *CndevMockRecorder) GetDeviceUtil(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceUtil", reflect.TypeOf((*Cndev)(nil).GetDeviceUtil), arg0)
}

// GetDeviceVersion mocks base method
func (m *Cndev) GetDeviceVersion(arg0 uint) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceVersion", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceVersion indicates an expected call of GetDeviceVersion
func (mr *CndevMockRecorder) GetDeviceVersion(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceVersion", reflect.TypeOf((*Cndev)(nil).GetDeviceVersion), arg0)
}

// Init mocks base method
func (m *Cndev) Init() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init")
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init
func (mr *CndevMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*Cndev)(nil).Init))
}

// Release mocks base method
func (m *Cndev) Release() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Release")
	ret0, _ := ret[0].(error)
	return ret0
}

// Release indicates an expected call of Release
func (mr *CndevMockRecorder) Release() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Release", reflect.TypeOf((*Cndev)(nil).Release))
}
