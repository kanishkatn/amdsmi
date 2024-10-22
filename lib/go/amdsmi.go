package _go

/*
#cgo CXXFLAGS: -Iinclude -std=c++11
#cgo CFLAGS: -Iinclude
#cgo LDFLAGS: -ldl

#include <dlfcn.h>
#include <stdio.h>
#include <stdlib.h>
#include "amdsmi.h"

static void *lib_handle = NULL;

// Function pointer types for the AMD SMI functions
typedef amdsmi_status_t (*amdsmi_init_fp)(uint64_t flags);
typedef amdsmi_status_t (*amdsmi_shut_down_fp)();
typedef amdsmi_status_t (*amdsmi_get_socket_handles_fp)(uint32_t *socket_count, amdsmi_socket_handle* socket_handles);
typedef amdsmi_status_t (*amdsmi_get_socket_info_fp)(amdsmi_socket_handle socket_handle, size_t len, char *name);
typedef amdsmi_status_t (*amdsmi_get_processor_handles_fp)(amdsmi_socket_handle socket_handle, uint32_t* processor_count, amdsmi_processor_handle* processor_handles);
typedef amdsmi_status_t (*amdsmi_get_processor_type_fp)(amdsmi_processor_handle processor_handle, processor_type_t* processor_type);
typedef amdsmi_status_t (*amdsmi_get_gpu_board_info_fp)(amdsmi_processor_handle processor_handle, amdsmi_board_info_t *board_info);
typedef amdsmi_status_t (*amdsmi_get_gpu_id_fp)(amdsmi_processor_handle processor_handle, uint16_t *id);
typedef amdsmi_status_t (*amdsmi_get_gpu_device_uuid_fp)(amdsmi_processor_handle processor_handle, unsigned int *uuid_length, char *uuid);
typedef amdsmi_status_t (*amdsmi_get_gpu_vram_usage_fp)(amdsmi_processor_handle processor_handle, amdsmi_vram_usage_t *vram_info);

// Dynamically loaded function pointers
static amdsmi_init_fp go_amdsmi_init = NULL;
static amdsmi_shut_down_fp go_amdsmi_shut_down = NULL;
static amdsmi_get_socket_handles_fp go_amdsmi_get_socket_handles = NULL;
static amdsmi_get_socket_info_fp go_amdsmi_get_socket_info = NULL;
static amdsmi_get_processor_handles_fp go_amdsmi_get_processor_handles = NULL;
static amdsmi_get_processor_type_fp go_amdsmi_get_processor_type = NULL;
static amdsmi_get_gpu_board_info_fp go_amdsmi_get_gpu_board_info = NULL;
static amdsmi_get_gpu_id_fp go_amdsmi_get_gpu_id = NULL;
static amdsmi_get_gpu_device_uuid_fp go_amdsmi_get_gpu_device_uuid = NULL;
static amdsmi_get_gpu_vram_usage_fp go_amdsmi_get_gpu_vram_usage = NULL;

// Load the library and resolve symbols
int load_amdsmi_library() {
    lib_handle = dlopen("libamd_smi.so", RTLD_LAZY);
    if (!lib_handle) {
        fprintf(stderr, "Error loading libamd_smi.so: %s\n", dlerror());
        return 0;  // Library not available
    }

    go_amdsmi_init = (amdsmi_init_fp)dlsym(lib_handle, "amdsmi_init");
    go_amdsmi_shut_down = (amdsmi_shut_down_fp)dlsym(lib_handle, "amdsmi_shut_down");
    go_amdsmi_get_socket_handles = (amdsmi_get_socket_handles_fp)dlsym(lib_handle, "amdsmi_get_socket_handles");
    go_amdsmi_get_socket_info = (amdsmi_get_socket_info_fp)dlsym(lib_handle, "amdsmi_get_socket_info");
    go_amdsmi_get_processor_handles = (amdsmi_get_processor_handles_fp)dlsym(lib_handle, "amdsmi_get_processor_handles");
    go_amdsmi_get_processor_type = (amdsmi_get_processor_type_fp)dlsym(lib_handle, "amdsmi_get_processor_type");
    go_amdsmi_get_gpu_board_info = (amdsmi_get_gpu_board_info_fp)dlsym(lib_handle, "amdsmi_get_gpu_board_info");
    go_amdsmi_get_gpu_id = (amdsmi_get_gpu_id_fp)dlsym(lib_handle, "amdsmi_get_gpu_id");
    go_amdsmi_get_gpu_device_uuid = (amdsmi_get_gpu_device_uuid_fp)dlsym(lib_handle, "amdsmi_get_gpu_device_uuid");
    go_amdsmi_get_gpu_vram_usage = (amdsmi_get_gpu_vram_usage_fp)dlsym(lib_handle, "amdsmi_get_gpu_vram_usage");

    if (!go_amdsmi_init || !go_amdsmi_shut_down || !go_amdsmi_get_socket_handles || !go_amdsmi_get_socket_info ||
        !go_amdsmi_get_processor_handles || !go_amdsmi_get_processor_type || !go_amdsmi_get_gpu_board_info ||
        !go_amdsmi_get_gpu_id || !go_amdsmi_get_gpu_device_uuid || !go_amdsmi_get_gpu_vram_usage) {
        fprintf(stderr, "Error resolving libamd_smi symbols: %s\n", dlerror());
        dlclose(lib_handle);
        return 0;
    }

    return 1;
}

// Unload the AMD SMI library
void unload_amdsmi_library() {
    if (lib_handle) {
        dlclose(lib_handle);
    }
}

// Wrapper for amdsmi_init
amdsmi_status_t call_amdsmi_init(uint64_t flags) {
    if (go_amdsmi_init) {
        return go_amdsmi_init(flags);
    }
    return AMDSMI_STATUS_INVAL;
}

// Wrapper for amdsmi_shut_down
amdsmi_status_t call_amdsmi_shut_down() {
    if (go_amdsmi_shut_down) {
        return go_amdsmi_shut_down();
    }
    return AMDSMI_STATUS_INVAL;
}

// Wrapper for amdsmi_get_socket_handles
amdsmi_status_t call_amdsmi_get_socket_handles(uint32_t *socket_count, amdsmi_socket_handle* socket_handles) {
    if (go_amdsmi_get_socket_handles) {
        return go_amdsmi_get_socket_handles(socket_count, socket_handles);
    }
    return AMDSMI_STATUS_INVAL;
}

// Wrapper for amdsmi_get_socket_info
amdsmi_status_t call_amdsmi_get_socket_info(amdsmi_socket_handle socket_handle, size_t len, char *name) {
    if (go_amdsmi_get_socket_info) {
        return go_amdsmi_get_socket_info(socket_handle, len, name);
    }
    return AMDSMI_STATUS_INVAL;
}

// Wrapper for amdsmi_get_processor_handles
amdsmi_status_t call_amdsmi_get_processor_handles(amdsmi_socket_handle socket_handle, uint32_t* processor_count, amdsmi_processor_handle* processor_handles) {
    if (go_amdsmi_get_processor_handles) {
        return go_amdsmi_get_processor_handles(socket_handle, processor_count, processor_handles);
    }
    return AMDSMI_STATUS_INVAL;
}

// Wrapper for amdsmi_get_processor_type
amdsmi_status_t call_amdsmi_get_processor_type(amdsmi_processor_handle processor_handle, processor_type_t* processor_type) {
    if (go_amdsmi_get_processor_type) {
        return go_amdsmi_get_processor_type(processor_handle, processor_type);
    }
    return AMDSMI_STATUS_INVAL;
}

// Wrapper for amdsmi_get_gpu_board_info
amdsmi_status_t call_amdsmi_get_gpu_board_info(amdsmi_processor_handle processor_handle, amdsmi_board_info_t *board_info) {
    if (go_amdsmi_get_gpu_board_info) {
        return go_amdsmi_get_gpu_board_info(processor_handle, board_info);
    }
    return AMDSMI_STATUS_INVAL;
}

// Wrapper for amdsmi_get_gpu_id
amdsmi_status_t call_amdsmi_get_gpu_id(amdsmi_processor_handle processor_handle, uint16_t *id) {
    if (go_amdsmi_get_gpu_id) {
        return go_amdsmi_get_gpu_id(processor_handle, id);
    }
    return AMDSMI_STATUS_INVAL;
}

// Wrapper for amdsmi_get_gpu_device_uuid
amdsmi_status_t call_amdsmi_get_gpu_device_uuid(amdsmi_processor_handle processor_handle, unsigned int *uuid_length, char *uuid) {
    if (go_amdsmi_get_gpu_device_uuid) {
        return go_amdsmi_get_gpu_device_uuid(processor_handle, uuid_length, uuid);
    }
    return AMDSMI_STATUS_INVAL;
}

// Wrapper for amdsmi_get_gpu_vram_usage
amdsmi_status_t call_amdsmi_get_gpu_vram_usage(amdsmi_processor_handle processor_handle, amdsmi_vram_usage_t *vram_info) {
    if (go_amdsmi_get_gpu_vram_usage) {
        return go_amdsmi_get_gpu_vram_usage(processor_handle, vram_info);
    }
    return AMDSMI_STATUS_INVAL;
}
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
	if C.load_amdsmi_library() == 0 {
		return false
	}

	return C.call_amdsmi_init(C.AMDSMI_INIT_AMD_GPUS) == C.AMDSMI_STATUS_SUCCESS
}

// Shutdown shuts down the AMD SMI library.
func Shutdown() bool {
	return C.call_amdsmi_shut_down() == C.AMDSMI_STATUS_SUCCESS
}

// GetSocketHandles returns the socket handles of the GPUs.
func GetSocketHandles() ([]SocketHandle, error) {
	var socketCount C.uint32_t
	ret := C.call_amdsmi_get_socket_handles(&socketCount, nil)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return nil, fmt.Errorf("failed to get socket count: %d", ret)
	}

	sockets := make([]C.amdsmi_socket_handle, socketCount)
	ret = C.call_amdsmi_get_socket_handles(&socketCount, (*C.amdsmi_socket_handle)(unsafe.Pointer(&sockets[0])))
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

	ret := C.call_amdsmi_get_socket_info(C.amdsmi_socket_handle(socketHandle.handle), C.size_t(maxLen), (*C.char)(unsafe.Pointer(&name[0])))
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return "", fmt.Errorf("failed to get socket info: %d", ret)
	}

	socketInfo := C.GoString(&name[0])
	return socketInfo, nil
}

// GetProcessorHandles retrieves all processor handles for a given socket.
func GetProcessorHandles(socket SocketHandle) ([]ProcessorHandle, error) {
	var processorCount C.uint32_t
	ret := C.call_amdsmi_get_processor_handles(C.amdsmi_socket_handle(socket.handle), &processorCount, nil)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return nil, fmt.Errorf("failed to get processor count for socket: %d", ret)
	}

	processors := make([]C.amdsmi_processor_handle, processorCount)
	ret = C.call_amdsmi_get_processor_handles(C.amdsmi_socket_handle(socket.handle), &processorCount, (*C.amdsmi_processor_handle)(unsafe.Pointer(&processors[0])))
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
	ret := C.call_amdsmi_get_processor_type(C.amdsmi_processor_handle(processor.handle), &processorType)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return 0, fmt.Errorf("failed to get processor type: %d", ret)
	}

	return ProcessorType(processorType), nil
}

// GetGPUBoardInfo retrieves the board information for a given GPU processor handle.
func GetGPUBoardInfo(processor ProcessorHandle) (BoardInfo, error) {
	var boardInfo C.amdsmi_board_info_t

	ret := C.call_amdsmi_get_gpu_board_info(C.amdsmi_processor_handle(processor.handle), &boardInfo)
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
	ret := C.call_amdsmi_get_gpu_id(C.amdsmi_processor_handle(processor.handle), &gpuID)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return 0, fmt.Errorf("failed to get GPU ID: %v", ret)
	}

	return uint32(gpuID), nil
}

// GetGPUUUID retrieves the GPU UUID for a given processor handle.
func GetGPUUUID(processor ProcessorHandle) (string, error) {
	var uuid [38]C.char
	var length C.uint = 38
	ret := C.call_amdsmi_get_gpu_device_uuid(C.amdsmi_processor_handle(processor.handle), &length, (*C.char)(unsafe.Pointer(&uuid)))
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return "", fmt.Errorf("failed to get GPU UUID: %d", ret)
	}

	return C.GoString(&uuid[0]), nil
}

// GetGPUVRAM retrieves the GPU VRAM stats for a given processor handle.
func GetGPUVRAM(processor ProcessorHandle) (VRAM, error) {
	var vramUsage C.amdsmi_vram_usage_t
	ret := C.call_amdsmi_get_gpu_vram_usage(C.amdsmi_processor_handle(processor.handle), &vramUsage)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return VRAM{}, fmt.Errorf("failed to get GPU VRAM info: %v", ret)
	}

	goVRAM := VRAM{
		Total: uint32(vramUsage.vram_total),
		Used:  uint32(vramUsage.vram_used),
	}
	return goVRAM, nil
}
