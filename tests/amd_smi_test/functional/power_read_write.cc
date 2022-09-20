/*
 * =============================================================================
 *   ROC Runtime Conformance Release License
 * =============================================================================
 * The University of Illinois/NCSA
 * Open Source License (NCSA)
 *
 * Copyright (c) 2022, Advanced Micro Devices, Inc.
 * All rights reserved.
 *
 * Developed by:
 *
 *                 AMD Research and AMD ROC Software Development
 *
 *                 Advanced Micro Devices, Inc.
 *
 *                 www.amd.com
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to
 * deal with the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 *  - Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimers.
 *  - Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimers in
 *    the documentation and/or other materials provided with the distribution.
 *  - Neither the names of <Name of Development Group, Name of Institution>,
 *    nor the names of its contributors may be used to endorse or promote
 *    products derived from this Software without specific prior written
 *    permission.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE CONTRIBUTORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 * OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 * ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS WITH THE SOFTWARE.
 *
 */

#include <stdint.h>
#include <stddef.h>

#include <iostream>
#include <bitset>
#include <string>
#include <algorithm>

#include "gtest/gtest.h"
#include "amd_smi/amd_smi.h"
#include "amd_smi_test/functional/power_read_write.h"
#include "amd_smi_test/test_common.h"


TestPowerReadWrite::TestPowerReadWrite() : TestBase() {
  set_title("AMDSMI Power Profiles Read/Write Test");
  set_description("The Power Profiles tests verify that the power profile "
                             "settings can be read and controlled properly.");
}

TestPowerReadWrite::~TestPowerReadWrite(void) {
}

void TestPowerReadWrite::SetUp(void) {
  TestBase::SetUp();

  return;
}

void TestPowerReadWrite::DisplayTestInfo(void) {
  TestBase::DisplayTestInfo();
}

void TestPowerReadWrite::DisplayResults(void) const {
  TestBase::DisplayResults();
  return;
}

void TestPowerReadWrite::Close() {
  // This will close handles opened within amdsmitst utility calls and call
  // amdsmi_shut_down(), so it should be done after other hsa cleanup
  TestBase::Close();
}

static const char *
power_profile_string(amdsmi_power_profile_preset_masks_t profile) {
  switch (profile) {
    case AMDSMI_PWR_PROF_PRST_CUSTOM_MASK:
      return "CUSTOM";
    case AMDSMI_PWR_PROF_PRST_VIDEO_MASK:
      return "VIDEO";
    case AMDSMI_PWR_PROF_PRST_POWER_SAVING_MASK:
      return "POWER SAVING";
    case AMDSMI_PWR_PROF_PRST_COMPUTE_MASK:
      return "COMPUTE";
    case AMDSMI_PWR_PROF_PRST_VR_MASK:
      return "VR";
    case AMDSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK:
      return "3D FULL SCREEN";
    case AMDSMI_PWR_PROF_PRST_BOOTUP_DEFAULT:
      return "BOOTUP DEFAULT";
    default:
      return "UNKNOWN";
  }
}

void TestPowerReadWrite::Run(void) {
  amdsmi_status_t ret;
  amdsmi_power_profile_status_t status;

  TestBase::Run();
  if (setup_failed_) {
    std::cout << "** SetUp Failed for this test. Skipping.**" << std::endl;
    return;
  }

  for (uint32_t dv_ind = 0; dv_ind < num_monitor_devs(); ++dv_ind) {
    PrintDeviceHeader(device_handles_[dv_ind]);

    ret = amdsmi_dev_power_profile_presets_get(device_handles_[dv_ind], 0, &status);
    CHK_ERR_ASRT(ret)

    // Verify api support checking functionality is working
    ret = amdsmi_dev_power_profile_presets_get(device_handles_[dv_ind], 0, nullptr);
    ASSERT_EQ(ret, AMDSMI_STATUS_INVAL);

    IF_VERB(STANDARD) {
      std::cout << "The available power profiles are:" << std::endl;
      uint64_t tmp = 1;
      while (tmp <= AMDSMI_PWR_PROF_PRST_LAST) {
        if ((tmp & status.available_profiles) == tmp) {
          std::cout << "\t" <<
              power_profile_string((amdsmi_power_profile_preset_masks_t)tmp) <<
                                                                      std::endl;
        }
        tmp = tmp << 1;
      }
      std::cout << "The current power profile is: " <<
                              power_profile_string(status.current) << std::endl;
    }

    amdsmi_power_profile_preset_masks_t orig_profile = status.current;

    // Try setting the profile to a different power profile
    amdsmi_bit_field_t diff_profiles;
    amdsmi_power_profile_preset_masks_t new_prof;
    diff_profiles = status.available_profiles & (~status.current);

    if (diff_profiles & AMDSMI_PWR_PROF_PRST_COMPUTE_MASK) {
      new_prof = AMDSMI_PWR_PROF_PRST_COMPUTE_MASK;
    } else if (diff_profiles & AMDSMI_PWR_PROF_PRST_VIDEO_MASK) {
      new_prof = AMDSMI_PWR_PROF_PRST_VIDEO_MASK;
    } else if (diff_profiles & AMDSMI_PWR_PROF_PRST_VR_MASK) {
      new_prof = AMDSMI_PWR_PROF_PRST_VR_MASK;
    } else if (diff_profiles & AMDSMI_PWR_PROF_PRST_POWER_SAVING_MASK) {
      new_prof = AMDSMI_PWR_PROF_PRST_POWER_SAVING_MASK;
    } else if (diff_profiles & AMDSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK) {
      new_prof = AMDSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK;
    } else {
      std::cout <<
        "No other non-custom power profiles to set to. Exiting." << std::endl;
      return;
    }

    ret = amdsmi_dev_power_profile_set(device_handles_[dv_ind], 0, new_prof);
    CHK_ERR_ASRT(ret)

    amdsmi_dev_perf_level_t pfl;
    ret = amdsmi_dev_perf_level_get(device_handles_[dv_ind], &pfl);
    CHK_ERR_ASRT(ret)
    ASSERT_EQ(pfl, AMDSMI_DEV_PERF_LEVEL_MANUAL);

    ret = amdsmi_dev_power_profile_presets_get(device_handles_[dv_ind], 0, &status);
    CHK_ERR_ASRT(ret)

    ASSERT_EQ(status.current, new_prof);

    ret = amdsmi_dev_perf_level_set(device_handles_[dv_ind], AMDSMI_DEV_PERF_LEVEL_AUTO);
    CHK_ERR_ASRT(ret)

    ret = amdsmi_dev_perf_level_get(device_handles_[dv_ind], &pfl);
    CHK_ERR_ASRT(ret)
    ASSERT_EQ(pfl, AMDSMI_DEV_PERF_LEVEL_AUTO);

    ret = amdsmi_dev_power_profile_presets_get(device_handles_[dv_ind], 0, &status);
    CHK_ERR_ASRT(ret)

    ASSERT_EQ(status.current, orig_profile);
  }
}
