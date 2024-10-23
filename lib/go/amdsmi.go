package _go

/*
#cgo CXXFLAGS: -Iinclude -std=c++11
#cgo CFLAGS: -Iinclude
#cgo LDFLAGS: -ldl

#include <dlfcn.h>
#include <stdio.h>
#include <stdlib.h>
#include "amdsmi.h"

// ========================
// AMD SMI Function Pointer Definitions
// ========================

// Macro to define a function pointer type
#define DEFINE_AMDSMI_FUNC_TYPE(ret, name, args) typedef ret (*name##_fp) args

// Define function pointer types for AMD SMI functions
DEFINE_AMDSMI_FUNC_TYPE(amdsmi_status_t, amdsmi_init, (uint64_t flags));
DEFINE_AMDSMI_FUNC_TYPE(amdsmi_status_t, amdsmi_shut_down, (void));
DEFINE_AMDSMI_FUNC_TYPE(amdsmi_status_t, amdsmi_get_socket_handles, (uint32_t *socket_count, amdsmi_socket_handle* socket_handles));
DEFINE_AMDSMI_FUNC_TYPE(amdsmi_status_t, amdsmi_get_socket_info, (amdsmi_socket_handle socket_handle, size_t len, char *name));
DEFINE_AMDSMI_FUNC_TYPE(amdsmi_status_t, amdsmi_get_processor_handles, (amdsmi_socket_handle socket_handle, uint32_t* processor_count, amdsmi_processor_handle* processor_handles));
DEFINE_AMDSMI_FUNC_TYPE(amdsmi_status_t, amdsmi_get_processor_type, (amdsmi_processor_handle processor_handle, processor_type_t* processor_type));
DEFINE_AMDSMI_FUNC_TYPE(amdsmi_status_t, amdsmi_get_gpu_board_info, (amdsmi_processor_handle processor_handle, amdsmi_board_info_t *board_info));
DEFINE_AMDSMI_FUNC_TYPE(amdsmi_status_t, amdsmi_get_gpu_id, (amdsmi_processor_handle processor_handle, uint16_t *id));
DEFINE_AMDSMI_FUNC_TYPE(amdsmi_status_t, amdsmi_get_gpu_device_uuid, (amdsmi_processor_handle processor_handle, unsigned int *uuid_length, char *uuid));
DEFINE_AMDSMI_FUNC_TYPE(amdsmi_status_t, amdsmi_get_gpu_vram_usage, (amdsmi_processor_handle processor_handle, amdsmi_vram_usage_t *vram_info));
DEFINE_AMDSMI_FUNC_TYPE(amdsmi_status_t, amdsmi_get_gpu_bdf_id, (amdsmi_processor_handle processor_handle, uint64_t *bdf_id));

// ========================
// AMD SMI Function Pointers Struct
// ========================

typedef struct {
    amdsmi_init_fp amdsmi_init;
    amdsmi_shut_down_fp amdsmi_shut_down;
    amdsmi_get_socket_handles_fp amdsmi_get_socket_handles;
    amdsmi_get_socket_info_fp amdsmi_get_socket_info;
    amdsmi_get_processor_handles_fp amdsmi_get_processor_handles;
    amdsmi_get_processor_type_fp amdsmi_get_processor_type;
    amdsmi_get_gpu_board_info_fp amdsmi_get_gpu_board_info;
    amdsmi_get_gpu_id_fp amdsmi_get_gpu_id;
    amdsmi_get_gpu_device_uuid_fp amdsmi_get_gpu_device_uuid;
    amdsmi_get_gpu_vram_usage_fp amdsmi_get_gpu_vram_usage;
    amdsmi_get_gpu_bdf_id_fp amdsmi_get_gpu_bdf_id;
} amdsmi_functions_t;

// ========================
// Global Variables
// ========================

static void *lib_handle = NULL;
static amdsmi_functions_t amdsmi_funcs;

// ========================
// Helper Macros
// ========================

// Macro to load a symbol and assign it to a struct member
#define LOAD_AMDSMI_SYMBOL(func_name) \
    amdsmi_funcs.func_name = (func_name##_fp)dlsym(lib_handle, #func_name); \
    if (!amdsmi_funcs.func_name) { \
        fprintf(stderr, "Error loading symbol %s: %s\n", #func_name, dlerror()); \
        dlclose(lib_handle); \
        lib_handle = NULL; \
        return 0; \
    }

// Macro to define a wrapper function
#define DEFINE_AMDSMI_WRAPPER(ret_type, wrapper_name, amdsmi_func, args, ...) \
    ret_type wrapper_name args { \
        if (amdsmi_funcs.amdsmi_func) { \
            return amdsmi_funcs.amdsmi_func(__VA_ARGS__); \
        } \
        return AMDSMI_STATUS_INVAL; \
    }

// ========================
// Library Management Functions
// ========================

// Load the AMD SMI library and resolve all required symbols
int load_amdsmi_library() {
    if (lib_handle) {
        fprintf(stderr, "libamd_smi.so is already loaded.\n");
        return 1;
    }

    lib_handle = dlopen("libamd_smi.so", RTLD_LAZY);
    if (!lib_handle) {
        fprintf(stderr, "Error loading libamd_smi.so: %s\n", dlerror());
        return 0;
    }

    // Clear any existing errors
    dlerror();

    // Load all required symbols
    LOAD_AMDSMI_SYMBOL(amdsmi_init)
    LOAD_AMDSMI_SYMBOL(amdsmi_shut_down)
    LOAD_AMDSMI_SYMBOL(amdsmi_get_socket_handles)
    LOAD_AMDSMI_SYMBOL(amdsmi_get_socket_info)
    LOAD_AMDSMI_SYMBOL(amdsmi_get_processor_handles)
    LOAD_AMDSMI_SYMBOL(amdsmi_get_processor_type)
    LOAD_AMDSMI_SYMBOL(amdsmi_get_gpu_board_info)
    LOAD_AMDSMI_SYMBOL(amdsmi_get_gpu_id)
    LOAD_AMDSMI_SYMBOL(amdsmi_get_gpu_device_uuid)
    LOAD_AMDSMI_SYMBOL(amdsmi_get_gpu_vram_usage)
    LOAD_AMDSMI_SYMBOL(amdsmi_get_gpu_bdf_id)

    return 1;
}

// Unload the AMD SMI library and reset function pointers
void unload_amdsmi_library() {
    if (lib_handle) {
        dlclose(lib_handle);
        lib_handle = NULL;

        // Reset all function pointers to NULL
        amdsmi_funcs.amdsmi_init = NULL;
        amdsmi_funcs.amdsmi_shut_down = NULL;
        amdsmi_funcs.amdsmi_get_socket_handles = NULL;
        amdsmi_funcs.amdsmi_get_socket_info = NULL;
        amdsmi_funcs.amdsmi_get_processor_handles = NULL;
        amdsmi_funcs.amdsmi_get_processor_type = NULL;
        amdsmi_funcs.amdsmi_get_gpu_board_info = NULL;
        amdsmi_funcs.amdsmi_get_gpu_id = NULL;
        amdsmi_funcs.amdsmi_get_gpu_device_uuid = NULL;
        amdsmi_funcs.amdsmi_get_gpu_vram_usage = NULL;
        amdsmi_funcs.amdsmi_get_gpu_bdf_id = NULL;
    }
}

// ========================
// Wrapper Functions
// ========================

// Define wrappers for each AMD SMI function
DEFINE_AMDSMI_WRAPPER(amdsmi_status_t, call_amdsmi_init, amdsmi_init, (uint64_t flags), flags)
DEFINE_AMDSMI_WRAPPER(amdsmi_status_t, call_amdsmi_shut_down, amdsmi_shut_down, (void))
DEFINE_AMDSMI_WRAPPER(amdsmi_status_t, call_amdsmi_get_socket_handles, amdsmi_get_socket_handles, (uint32_t *socket_count, amdsmi_socket_handle* socket_handles), socket_count, socket_handles)
DEFINE_AMDSMI_WRAPPER(amdsmi_status_t, call_amdsmi_get_socket_info, amdsmi_get_socket_info, (amdsmi_socket_handle socket_handle, size_t len, char *name), socket_handle, len, name)
DEFINE_AMDSMI_WRAPPER(amdsmi_status_t, call_amdsmi_get_processor_handles, amdsmi_get_processor_handles, (amdsmi_socket_handle socket_handle, uint32_t* processor_count, amdsmi_processor_handle* processor_handles), socket_handle, processor_count, processor_handles)
DEFINE_AMDSMI_WRAPPER(amdsmi_status_t, call_amdsmi_get_processor_type, amdsmi_get_processor_type, (amdsmi_processor_handle processor_handle, processor_type_t* processor_type), processor_handle, processor_type)
DEFINE_AMDSMI_WRAPPER(amdsmi_status_t, call_amdsmi_get_gpu_board_info, amdsmi_get_gpu_board_info, (amdsmi_processor_handle processor_handle, amdsmi_board_info_t *board_info), processor_handle, board_info)
DEFINE_AMDSMI_WRAPPER(amdsmi_status_t, call_amdsmi_get_gpu_id, amdsmi_get_gpu_id, (amdsmi_processor_handle processor_handle, uint16_t *id), processor_handle, id)
DEFINE_AMDSMI_WRAPPER(amdsmi_status_t, call_amdsmi_get_gpu_device_uuid, amdsmi_get_gpu_device_uuid, (amdsmi_processor_handle processor_handle, unsigned int *uuid_length, char *uuid), processor_handle, uuid_length, uuid)
DEFINE_AMDSMI_WRAPPER(amdsmi_status_t, call_amdsmi_get_gpu_vram_usage, amdsmi_get_gpu_vram_usage, (amdsmi_processor_handle processor_handle, amdsmi_vram_usage_t *vram_info), processor_handle, vram_info)
DEFINE_AMDSMI_WRAPPER(amdsmi_status_t, call_amdsmi_get_gpu_bdf_id, amdsmi_get_gpu_bdf_id, (amdsmi_processor_handle processor_handle, uint64_t *bdf_id), processor_handle, bdf_id)
*/
import "C"

/*
===============================================================================
Adding a New AMD SMI Function to the Cgo Block
===============================================================================

To add a new AMD SMI function, follow these steps:

1. **Define the Function Pointer Type:**
   - Use `DEFINE_AMDSMI_FUNC_TYPE` to create a typedef for the new function.
   - Example:
     `DEFINE_AMDSMI_FUNC_TYPE(amdsmi_status_t, amdsmi_new_function, (int arg1, float arg2));`

2. **Add to the Struct:**
   - Add the new function pointer to `amdsmi_functions_t`.
   - Example:
     `amdsmi_new_function_fp amdsmi_new_function;`

3. **Load the Symbol:**
   - In `load_amdsmi_library`, load the new symbol using `LOAD_AMDSMI_SYMBOL`.
   - Example:
     `LOAD_AMDSMI_SYMBOL(amdsmi_new_function)`

4. **Create a Wrapper Function:**
   - Define a wrapper using `DEFINE_AMDSMI_WRAPPER`.
   - Example:
     `DEFINE_AMDSMI_WRAPPER(amdsmi_status_t, call_amdsmi_new_function, amdsmi_new_function, (int arg1, float arg2), arg1, arg2)`

5. **Rebuild and Test:**
   - Rebuild your project and verify the new function works as expected.

===============================================================================
*/

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
// TODO: Use a Go type instead of C type.
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
func Init() (bool, error) {
	if C.load_amdsmi_library() == 0 {
		return false, fmt.Errorf("failed to load AMD SMI library")
	}

	return C.call_amdsmi_init(C.AMDSMI_INIT_AMD_GPUS) == C.AMDSMI_STATUS_SUCCESS, nil
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

// GetGPUBDFID retrieves the GPU BDF ID for a given processor handle.
func GetGPUBDFID(processor ProcessorHandle) (uint64, error) {
	var bdfID C.uint64_t
	ret := C.call_amdsmi_get_gpu_bdf_id(C.amdsmi_processor_handle(processor.handle), &bdfID)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		return 0, fmt.Errorf("failed to get GPU BDF ID: %v", ret)
	}

	return uint64(bdfID), nil
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
