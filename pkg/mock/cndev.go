// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Cambricon/mlu-exporter/pkg/cndev (interfaces: Cndev)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// Cndev is a mock of Cndev interface.
type Cndev struct {
	ctrl     *gomock.Controller
	recorder *CndevMockRecorder
}

// CndevMockRecorder is the mock recorder for Cndev.
type CndevMockRecorder struct {
	mock *Cndev
}

// NewCndev creates a new mock instance.
func NewCndev(ctrl *gomock.Controller) *Cndev {
	mock := &Cndev{ctrl: ctrl}
	mock.recorder = &CndevMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Cndev) EXPECT() *CndevMockRecorder {
	return m.recorder
}

// GetDeviceArmOsMemory mocks base method.
func (m *Cndev) GetDeviceArmOsMemory(arg0 uint) (int64, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceArmOsMemory", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceArmOsMemory indicates an expected call of GetDeviceArmOsMemory.
func (mr *CndevMockRecorder) GetDeviceArmOsMemory(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceArmOsMemory", reflect.TypeOf((*Cndev)(nil).GetDeviceArmOsMemory), arg0)
}

// GetDeviceCPUUtil mocks base method.
func (m *Cndev) GetDeviceCPUUtil(arg0 uint) (uint16, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceCPUUtil", arg0)
	ret0, _ := ret[0].(uint16)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceCPUUtil indicates an expected call of GetDeviceCPUUtil.
func (mr *CndevMockRecorder) GetDeviceCPUUtil(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceCPUUtil", reflect.TypeOf((*Cndev)(nil).GetDeviceCPUUtil), arg0)
}

// GetDeviceCRCInfo mocks base method.
func (m *Cndev) GetDeviceCRCInfo(arg0 uint) (uint64, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceCRCInfo", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceCRCInfo indicates an expected call of GetDeviceCRCInfo.
func (mr *CndevMockRecorder) GetDeviceCRCInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceCRCInfo", reflect.TypeOf((*Cndev)(nil).GetDeviceCRCInfo), arg0)
}

// GetDeviceClusterCount mocks base method.
func (m *Cndev) GetDeviceClusterCount(arg0 uint) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceClusterCount", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceClusterCount indicates an expected call of GetDeviceClusterCount.
func (mr *CndevMockRecorder) GetDeviceClusterCount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceClusterCount", reflect.TypeOf((*Cndev)(nil).GetDeviceClusterCount), arg0)
}

// GetDeviceCoreNum mocks base method.
func (m *Cndev) GetDeviceCoreNum(arg0 uint) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceCoreNum", arg0)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceCoreNum indicates an expected call of GetDeviceCoreNum.
func (mr *CndevMockRecorder) GetDeviceCoreNum(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceCoreNum", reflect.TypeOf((*Cndev)(nil).GetDeviceCoreNum), arg0)
}

// GetDeviceCount mocks base method.
func (m *Cndev) GetDeviceCount() (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceCount")
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceCount indicates an expected call of GetDeviceCount.
func (mr *CndevMockRecorder) GetDeviceCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceCount", reflect.TypeOf((*Cndev)(nil).GetDeviceCount))
}

// GetDeviceCurrentPCIeInfo mocks base method.
func (m *Cndev) GetDeviceCurrentPCIeInfo(arg0 uint) (int, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceCurrentPCIeInfo", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceCurrentPCIeInfo indicates an expected call of GetDeviceCurrentPCIeInfo.
func (mr *CndevMockRecorder) GetDeviceCurrentPCIeInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceCurrentPCIeInfo", reflect.TypeOf((*Cndev)(nil).GetDeviceCurrentPCIeInfo), arg0)
}

// GetDeviceDDRInfo mocks base method.
func (m *Cndev) GetDeviceDDRInfo(arg0 uint) (int, float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceDDRInfo", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(float64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceDDRInfo indicates an expected call of GetDeviceDDRInfo.
func (mr *CndevMockRecorder) GetDeviceDDRInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceDDRInfo", reflect.TypeOf((*Cndev)(nil).GetDeviceDDRInfo), arg0)
}

// GetDeviceECCInfo mocks base method.
func (m *Cndev) GetDeviceECCInfo(arg0 uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceECCInfo", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(uint64)
	ret3, _ := ret[3].(uint64)
	ret4, _ := ret[4].(uint64)
	ret5, _ := ret[5].(uint64)
	ret6, _ := ret[6].(uint64)
	ret7, _ := ret[7].(uint64)
	ret8, _ := ret[8].(error)
	return ret0, ret1, ret2, ret3, ret4, ret5, ret6, ret7, ret8
}

// GetDeviceECCInfo indicates an expected call of GetDeviceECCInfo.
func (mr *CndevMockRecorder) GetDeviceECCInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceECCInfo", reflect.TypeOf((*Cndev)(nil).GetDeviceECCInfo), arg0)
}

// GetDeviceFanSpeed mocks base method.
func (m *Cndev) GetDeviceFanSpeed(arg0 uint) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceFanSpeed", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceFanSpeed indicates an expected call of GetDeviceFanSpeed.
func (mr *CndevMockRecorder) GetDeviceFanSpeed(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceFanSpeed", reflect.TypeOf((*Cndev)(nil).GetDeviceFanSpeed), arg0)
}

// GetDeviceHealth mocks base method.
func (m *Cndev) GetDeviceHealth(arg0 uint) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceHealth", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceHealth indicates an expected call of GetDeviceHealth.
func (mr *CndevMockRecorder) GetDeviceHealth(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceHealth", reflect.TypeOf((*Cndev)(nil).GetDeviceHealth), arg0)
}

// GetDeviceImageCodecUtil mocks base method.
func (m *Cndev) GetDeviceImageCodecUtil(arg0 uint) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceImageCodecUtil", arg0)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceImageCodecUtil indicates an expected call of GetDeviceImageCodecUtil.
func (mr *CndevMockRecorder) GetDeviceImageCodecUtil(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceImageCodecUtil", reflect.TypeOf((*Cndev)(nil).GetDeviceImageCodecUtil), arg0)
}

// GetDeviceMLULinkCapability mocks base method.
func (m *Cndev) GetDeviceMLULinkCapability(arg0, arg1 uint) (uint, uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceMLULinkCapability", arg0, arg1)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(uint)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceMLULinkCapability indicates an expected call of GetDeviceMLULinkCapability.
func (mr *CndevMockRecorder) GetDeviceMLULinkCapability(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceMLULinkCapability", reflect.TypeOf((*Cndev)(nil).GetDeviceMLULinkCapability), arg0, arg1)
}

// GetDeviceMLULinkCounter mocks base method.
func (m *Cndev) GetDeviceMLULinkCounter(arg0, arg1 uint) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceMLULinkCounter", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(uint64)
	ret3, _ := ret[3].(uint64)
	ret4, _ := ret[4].(uint64)
	ret5, _ := ret[5].(uint64)
	ret6, _ := ret[6].(uint64)
	ret7, _ := ret[7].(uint64)
	ret8, _ := ret[8].(uint64)
	ret9, _ := ret[9].(uint64)
	ret10, _ := ret[10].(uint64)
	ret11, _ := ret[11].(error)
	return ret0, ret1, ret2, ret3, ret4, ret5, ret6, ret7, ret8, ret9, ret10, ret11
}

// GetDeviceMLULinkCounter indicates an expected call of GetDeviceMLULinkCounter.
func (mr *CndevMockRecorder) GetDeviceMLULinkCounter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceMLULinkCounter", reflect.TypeOf((*Cndev)(nil).GetDeviceMLULinkCounter), arg0, arg1)
}

// GetDeviceMLULinkPortMode mocks base method.
func (m *Cndev) GetDeviceMLULinkPortMode(arg0, arg1 uint) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceMLULinkPortMode", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceMLULinkPortMode indicates an expected call of GetDeviceMLULinkPortMode.
func (mr *CndevMockRecorder) GetDeviceMLULinkPortMode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceMLULinkPortMode", reflect.TypeOf((*Cndev)(nil).GetDeviceMLULinkPortMode), arg0, arg1)
}

// GetDeviceMLULinkPortNumber mocks base method.
func (m *Cndev) GetDeviceMLULinkPortNumber(arg0 uint) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceMLULinkPortNumber", arg0)
	ret0, _ := ret[0].(int)
	return ret0
}

// GetDeviceMLULinkPortNumber indicates an expected call of GetDeviceMLULinkPortNumber.
func (mr *CndevMockRecorder) GetDeviceMLULinkPortNumber(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceMLULinkPortNumber", reflect.TypeOf((*Cndev)(nil).GetDeviceMLULinkPortNumber), arg0)
}

// GetDeviceMLULinkSpeedInfo mocks base method.
func (m *Cndev) GetDeviceMLULinkSpeedInfo(arg0, arg1 uint) (float32, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceMLULinkSpeedInfo", arg0, arg1)
	ret0, _ := ret[0].(float32)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceMLULinkSpeedInfo indicates an expected call of GetDeviceMLULinkSpeedInfo.
func (mr *CndevMockRecorder) GetDeviceMLULinkSpeedInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceMLULinkSpeedInfo", reflect.TypeOf((*Cndev)(nil).GetDeviceMLULinkSpeedInfo), arg0, arg1)
}

// GetDeviceMLULinkStatus mocks base method.
func (m *Cndev) GetDeviceMLULinkStatus(arg0, arg1 uint) (int, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceMLULinkStatus", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceMLULinkStatus indicates an expected call of GetDeviceMLULinkStatus.
func (mr *CndevMockRecorder) GetDeviceMLULinkStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceMLULinkStatus", reflect.TypeOf((*Cndev)(nil).GetDeviceMLULinkStatus), arg0, arg1)
}

// GetDeviceMLULinkVersion mocks base method.
func (m *Cndev) GetDeviceMLULinkVersion(arg0, arg1 uint) (uint, uint, uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceMLULinkVersion", arg0, arg1)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(uint)
	ret2, _ := ret[2].(uint)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetDeviceMLULinkVersion indicates an expected call of GetDeviceMLULinkVersion.
func (mr *CndevMockRecorder) GetDeviceMLULinkVersion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceMLULinkVersion", reflect.TypeOf((*Cndev)(nil).GetDeviceMLULinkVersion), arg0, arg1)
}

// GetDeviceMemory mocks base method.
func (m *Cndev) GetDeviceMemory(arg0 uint) (int64, int64, int64, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceMemory", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(int64)
	ret3, _ := ret[3].(int64)
	ret4, _ := ret[4].(error)
	return ret0, ret1, ret2, ret3, ret4
}

// GetDeviceMemory indicates an expected call of GetDeviceMemory.
func (mr *CndevMockRecorder) GetDeviceMemory(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceMemory", reflect.TypeOf((*Cndev)(nil).GetDeviceMemory), arg0)
}

// GetDeviceModel mocks base method.
func (m *Cndev) GetDeviceModel(arg0 uint) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceModel", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDeviceModel indicates an expected call of GetDeviceModel.
func (mr *CndevMockRecorder) GetDeviceModel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceModel", reflect.TypeOf((*Cndev)(nil).GetDeviceModel), arg0)
}

// GetDeviceNUMANodeID mocks base method.
func (m *Cndev) GetDeviceNUMANodeID(arg0 uint) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceNUMANodeID", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceNUMANodeID indicates an expected call of GetDeviceNUMANodeID.
func (mr *CndevMockRecorder) GetDeviceNUMANodeID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceNUMANodeID", reflect.TypeOf((*Cndev)(nil).GetDeviceNUMANodeID), arg0)
}

// GetDevicePCIeInfo mocks base method.
func (m *Cndev) GetDevicePCIeInfo(arg0 uint) (int, uint, uint, uint16, uint16, uint, uint, uint, uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDevicePCIeInfo", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(uint)
	ret2, _ := ret[2].(uint)
	ret3, _ := ret[3].(uint16)
	ret4, _ := ret[4].(uint16)
	ret5, _ := ret[5].(uint)
	ret6, _ := ret[6].(uint)
	ret7, _ := ret[7].(uint)
	ret8, _ := ret[8].(uint)
	ret9, _ := ret[9].(error)
	return ret0, ret1, ret2, ret3, ret4, ret5, ret6, ret7, ret8, ret9
}

// GetDevicePCIeInfo indicates an expected call of GetDevicePCIeInfo.
func (mr *CndevMockRecorder) GetDevicePCIeInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDevicePCIeInfo", reflect.TypeOf((*Cndev)(nil).GetDevicePCIeInfo), arg0)
}

// GetDevicePower mocks base method.
func (m *Cndev) GetDevicePower(arg0 uint) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDevicePower", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDevicePower indicates an expected call of GetDevicePower.
func (mr *CndevMockRecorder) GetDevicePower(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDevicePower", reflect.TypeOf((*Cndev)(nil).GetDevicePower), arg0)
}

// GetDeviceProcessUtil mocks base method.
func (m *Cndev) GetDeviceProcessUtil(arg0 uint) ([]uint32, []uint32, []uint32, []uint32, []uint32, []uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceProcessUtil", arg0)
	ret0, _ := ret[0].([]uint32)
	ret1, _ := ret[1].([]uint32)
	ret2, _ := ret[2].([]uint32)
	ret3, _ := ret[3].([]uint32)
	ret4, _ := ret[4].([]uint32)
	ret5, _ := ret[5].([]uint32)
	ret6, _ := ret[6].(error)
	return ret0, ret1, ret2, ret3, ret4, ret5, ret6
}

// GetDeviceProcessUtil indicates an expected call of GetDeviceProcessUtil.
func (mr *CndevMockRecorder) GetDeviceProcessUtil(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceProcessUtil", reflect.TypeOf((*Cndev)(nil).GetDeviceProcessUtil), arg0)
}

// GetDeviceSN mocks base method.
func (m *Cndev) GetDeviceSN(arg0 uint) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceSN", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceSN indicates an expected call of GetDeviceSN.
func (mr *CndevMockRecorder) GetDeviceSN(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceSN", reflect.TypeOf((*Cndev)(nil).GetDeviceSN), arg0)
}

// GetDeviceTemperature mocks base method.
func (m *Cndev) GetDeviceTemperature(arg0 uint) (int, []int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceTemperature", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceTemperature indicates an expected call of GetDeviceTemperature.
func (mr *CndevMockRecorder) GetDeviceTemperature(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceTemperature", reflect.TypeOf((*Cndev)(nil).GetDeviceTemperature), arg0)
}

// GetDeviceTinyCoreUtil mocks base method.
func (m *Cndev) GetDeviceTinyCoreUtil(arg0 uint) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceTinyCoreUtil", arg0)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceTinyCoreUtil indicates an expected call of GetDeviceTinyCoreUtil.
func (mr *CndevMockRecorder) GetDeviceTinyCoreUtil(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceTinyCoreUtil", reflect.TypeOf((*Cndev)(nil).GetDeviceTinyCoreUtil), arg0)
}

// GetDeviceUUID mocks base method.
func (m *Cndev) GetDeviceUUID(arg0 uint) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceUUID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceUUID indicates an expected call of GetDeviceUUID.
func (mr *CndevMockRecorder) GetDeviceUUID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceUUID", reflect.TypeOf((*Cndev)(nil).GetDeviceUUID), arg0)
}

// GetDeviceUtil mocks base method.
func (m *Cndev) GetDeviceUtil(arg0 uint) (int, []int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceUtil", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeviceUtil indicates an expected call of GetDeviceUtil.
func (mr *CndevMockRecorder) GetDeviceUtil(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceUtil", reflect.TypeOf((*Cndev)(nil).GetDeviceUtil), arg0)
}

// GetDeviceVersion mocks base method.
func (m *Cndev) GetDeviceVersion(arg0 uint) (uint, uint, uint, uint, uint, uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceVersion", arg0)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(uint)
	ret2, _ := ret[2].(uint)
	ret3, _ := ret[3].(uint)
	ret4, _ := ret[4].(uint)
	ret5, _ := ret[5].(uint)
	ret6, _ := ret[6].(error)
	return ret0, ret1, ret2, ret3, ret4, ret5, ret6
}

// GetDeviceVersion indicates an expected call of GetDeviceVersion.
func (mr *CndevMockRecorder) GetDeviceVersion(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceVersion", reflect.TypeOf((*Cndev)(nil).GetDeviceVersion), arg0)
}

// GetDeviceVfState mocks base method.
func (m *Cndev) GetDeviceVfState(arg0 uint) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceVfState", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceVfState indicates an expected call of GetDeviceVfState.
func (mr *CndevMockRecorder) GetDeviceVfState(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceVfState", reflect.TypeOf((*Cndev)(nil).GetDeviceVfState), arg0)
}

// GetDeviceVideoCodecUtil mocks base method.
func (m *Cndev) GetDeviceVideoCodecUtil(arg0 uint) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceVideoCodecUtil", arg0)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceVideoCodecUtil indicates an expected call of GetDeviceVideoCodecUtil.
func (mr *CndevMockRecorder) GetDeviceVideoCodecUtil(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceVideoCodecUtil", reflect.TypeOf((*Cndev)(nil).GetDeviceVideoCodecUtil), arg0)
}

// Init mocks base method.
func (m *Cndev) Init() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init")
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *CndevMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*Cndev)(nil).Init))
}

// Release mocks base method.
func (m *Cndev) Release() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Release")
	ret0, _ := ret[0].(error)
	return ret0
}

// Release indicates an expected call of Release.
func (mr *CndevMockRecorder) Release() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Release", reflect.TypeOf((*Cndev)(nil).Release))
}
