#include "amdsmi.h"

extern "C" {;
    amdsmi_status_t c_amdsmi_init(uint64_t flags) {
        return amdsmi_init(flags);
    }
    
    amdsmi_status_t c_amdsmi_shut_down() {
        return amdsmi_shut_down();
    }

    amdsmi_status_t c_amdsmi_get_socket_handles(uint32_t *socket_count,
                amdsmi_socket_handle* socket_handles) {
        return amdsmi_get_socket_handles(socket_count, socket_handles);
    }

    amdsmi_status_t c_amdsmi_get_socket_info(
                amdsmi_socket_handle socket_handle,
                size_t len, char *name) {
        return amdsmi_get_socket_info(socket_handle, len, name);
    }

    amdsmi_status_t c_amdsmi_get_processor_handles(amdsmi_socket_handle socket_handle,
                                    uint32_t* processor_count,
                                    amdsmi_processor_handle* processor_handles) {
        return amdsmi_get_processor_handles(socket_handle, processor_count, processor_handles);
    }

    amdsmi_status_t c_amdsmi_get_processor_type(amdsmi_processor_handle processor_handle ,
              processor_type_t* processor_type) {
        return amdsmi_get_processor_type(processor_handle, processor_type);
    }

    amdsmi_status_t c_amdsmi_get_gpu_board_info(amdsmi_processor_handle processor_handle, amdsmi_board_info_t *board_info) {
        return amdsmi_get_gpu_board_info(processor_handle, board_info);
    }

    amdsmi_status_t c_amdsmi_get_gpu_id(amdsmi_processor_handle processor_handle,
                                uint16_t *id) {
        return amdsmi_get_gpu_id(processor_handle, id);
    }

    amdsmi_status_t c_amdsmi_get_gpu_device_uuid(amdsmi_processor_handle processor_handle, unsigned int *uuid_length, char *uuid) {
        return amdsmi_get_gpu_device_uuid(processor_handle, uuid_length, uuid);
    }

    amdsmi_status_t c_amdsmi_get_gpu_vram_usage(amdsmi_processor_handle processor_handle,
            amdsmi_vram_usage_t *vram_info) {
        return amdsmi_get_gpu_vram_usage(processor_handle, vram_info);
    }
}