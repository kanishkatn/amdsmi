package main

/*
#cgo CXXFLAGS: -I./src/amd_smi -I./include/amd_smi -DENABLE_DEBUG_LEVEL=2 -std=c++11
#cgo CFLAGS: -I./src/amd_smi -I./include/amd_smi -DENABLE_DEBUG_LEVEL=2
#cgo LDFLAGS: -L/opt/rocm/lib -L/opt/rocm/lib64 -lamd_smi -Wl,--unresolved-symbols=ignore-in-object-files

#include "amdsmi.h"

// Declare the C wrapper functions
extern amdsmi_status_t c_amdsmi_init(uint64_t flags);
extern amdsmi_status_t c_amdsmi_shut_down();
amdsmi_status_t c_amdsmi_get_socket_handles(uint32_t *socket_count,
                amdsmi_socket_handle* socket_handles);
amdsmi_status_t c_amdsmi_get_socket_info(
                amdsmi_socket_handle socket_handle,
                size_t len, char *name);
amdsmi_status_t c_amdsmi_get_processor_handles(amdsmi_socket_handle socket_handle,
                                    uint32_t* processor_count,
                                    amdsmi_processor_handle* processor_handles);
amdsmi_status_t c_amdsmi_get_processor_type(amdsmi_processor_handle processor_handle ,
              processor_type_t* processor_type);
amdsmi_status_t c_amdsmi_get_gpu_board_info(amdsmi_processor_handle processor_handle, amdsmi_board_info_t *board_info);
amdsmi_status_t c_amdsmi_get_gpu_id(amdsmi_processor_handle processor_handle,
                                uint16_t *id);
amdsmi_status_t c_amdsmi_get_gpu_device_uuid(amdsmi_processor_handle processor_handle, unsigned int *uuid_length, char *uuid);
amdsmi_status_t c_amdsmi_get_gpu_vram_usage(amdsmi_processor_handle processor_handle,
            amdsmi_vram_usage_t *vram_info);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// SocketHandle is a Go type that encapsulates the C amdsmi_socket_handle (void*).
type SocketHandle struct {
	handle unsafe.Pointer
}

// ProcessorHandle is a Go type that encapsulates the C amdsmi_processor_handle (void*).
type ProcessorHandle struct {
	handle unsafe.Pointer
}

// ProcessorType is a Go type to represent processor_type_t from C.
type ProcessorType C.processor_type_t

// BoardInfo is a Go representation of the C amdsmi_board_info_t struct.
type BoardInfo struct {
	ModelNumber      string
	ProductSerial    string
	FruID            string
	ProductName      string
	ManufacturerName string
}

// VRAM is a Go representation of the C amdsmi_vram_info_t struct.
type VRAM struct {
	Total    uint32
	Used     uint32
	Reserved [5]uint64
}

// Init initializes the AMD SMI library with GPUs.
func Init() bool {
	return C.c_amdsmi_init(C.AMDSMI_INIT_AMD_GPUS) == C.AMDSMI_STATUS_SUCCESS
}

// Shutdown shuts down the AMD SMI library.
func Shutdown() bool {
	return C.c_amdsmi_shut_down() == C.AMDSMI_STATUS_SUCCESS
}

// GetSocketHandles returns the socket handles of the GPUs.
func GetSocketHandles() ([]SocketHandle, error) {
	var socketCount C.uint32_t
	ret := C.amdsmi_get_socket_handles(&socketCount, nil)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return nil, fmt.Errorf("failed to get socket count: %d", ret)
	}

	sockets := make([]C.amdsmi_socket_handle, socketCount)
	ret = C.amdsmi_get_socket_handles(&socketCount, (*C.amdsmi_socket_handle)(unsafe.Pointer(&sockets[0])))
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return nil, fmt.Errorf("failed to get socket handles: %d", ret)
	}

	goSockets := make([]SocketHandle, socketCount)
	for i, socket := range sockets {
		goSockets[i] = SocketHandle{handle: unsafe.Pointer(socket)}
	}

	return goSockets, nil
}

// GetSocketName retrieves the socket name for a given socket handle.
func GetSocketName(socketHandle SocketHandle, maxLen int) (string, error) {
	name := make([]C.char, maxLen)

	ret := C.amdsmi_get_socket_info(C.amdsmi_socket_handle(socketHandle.handle), C.size_t(maxLen), (*C.char)(unsafe.Pointer(&name[0])))
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return "", fmt.Errorf("failed to get socket info: %d", ret)
	}

	socketInfo := C.GoString(&name[0])
	return socketInfo, nil
}

// GetProcessorHandles retrieves all processor handles for a given socket.
func GetProcessorHandles(socket SocketHandle) ([]ProcessorHandle, error) {
	var processorCount C.uint32_t
	ret := C.amdsmi_get_processor_handles(C.amdsmi_socket_handle(socket.handle), &processorCount, nil)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return nil, fmt.Errorf("failed to get processor count for socket: %d", ret)
	}

	processors := make([]C.amdsmi_processor_handle, processorCount)
	ret = C.amdsmi_get_processor_handles(C.amdsmi_socket_handle(socket.handle), &processorCount, (*C.amdsmi_processor_handle)(unsafe.Pointer(&processors[0])))
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return nil, fmt.Errorf("failed to get processor handles for socket: %d", ret)
	}

	goProcessors := make([]ProcessorHandle, processorCount)
	for i, processor := range processors {
		goProcessors[i] = ProcessorHandle{handle: unsafe.Pointer(processor)}
	}

	return goProcessors, nil
}

// GetProcessorType retrieves the type of a given processor.
func GetProcessorType(processor ProcessorHandle) (ProcessorType, error) {
	var processorType C.processor_type_t
	ret := C.amdsmi_get_processor_type(C.amdsmi_processor_handle(processor.handle), &processorType)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return 0, fmt.Errorf("failed to get processor type: %d", ret)
	}

	return ProcessorType(processorType), nil
}

// GetGPUBoardInfo retrieves the board information for a given GPU processor handle.
func GetGPUBoardInfo(processor ProcessorHandle) (BoardInfo, error) {
	var boardInfo C.amdsmi_board_info_t

	ret := C.amdsmi_get_gpu_board_info(C.amdsmi_processor_handle(processor.handle), &boardInfo)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return BoardInfo{}, fmt.Errorf("failed to get GPU board info: %d", ret)
	}

	goBoardInfo := BoardInfo{
		ModelNumber:      C.GoString(&boardInfo.model_number[0]),
		ProductSerial:    C.GoString(&boardInfo.product_serial[0]),
		FruID:            C.GoString(&boardInfo.fru_id[0]),
		ProductName:      C.GoString(&boardInfo.product_name[0]),
		ManufacturerName: C.GoString(&boardInfo.manufacturer_name[0]),
	}
	return goBoardInfo, nil
}

// GetGPUID retrieves the GPU ID for a given processor handle.
func GetGPUID(processor ProcessorHandle) (uint32, error) {
	var gpuID C.uint16_t
	ret := C.c_amdsmi_get_gpu_id(C.amdsmi_processor_handle(processor.handle), &gpuID)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return 0, fmt.Errorf("failed to get GPU ID: %v", ret)
	}

	return uint32(gpuID), nil
}

// GetGPUUUID retrieves the GPU UUID for a given processor handle.
func GetGPUUUID(processor ProcessorHandle) (string, error) {
	var uuid [38]C.char
	var length C.uint = 38
	ret := C.c_amdsmi_get_gpu_device_uuid(C.amdsmi_processor_handle(processor.handle), &length, (*C.char)(unsafe.Pointer(&uuid)))
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return "", fmt.Errorf("failed to get GPU UUID: %d", ret)
	}

	return C.GoString(&uuid[0]), nil
}

// GetGPUVRAM retrieves the GPU VRAM stats for a given processor handle.
func GetGPUVRAM(processor ProcessorHandle) (VRAM, error) {
	var vramUsage C.amdsmi_vram_usage_t
	ret := C.c_amdsmi_get_gpu_vram_usage(C.amdsmi_processor_handle(processor.handle), &vramUsage)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return VRAM{}, fmt.Errorf("failed to get GPU VRAM info: %v", ret)
	}

	goVRAM := VRAM{
		Total: uint32(vramUsage.vram_total),
		Used:  uint32(vramUsage.vram_used),
	}
	return goVRAM, nil
}
