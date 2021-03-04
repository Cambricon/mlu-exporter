/*
 * Copyright 2021 Cambricon, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#ifndef CNPAPI_HOOKLIB_CALLBACKAPI_TYPES_H_
#define CNPAPI_HOOKLIB_CALLBACKAPI_TYPES_H_
#include <stdbool.h>
#include <stdint.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifndef __CNPAPI_TYPES_H
#define __CNPAPI_TYPES_H
#if defined(WIN32) || defined(WINDOWS)
  typedef unsigned __int64 u64_t;
  typedef __int64 i64_t;
  typedef unsigned __int32 u32_t;
  typedef unsigned __int16 u16_t;
  typedef unsigned __int8 u8_t;
  typedef signed __int32 i32_t;
  typedef signed __int16 i16_t;
  typedef signed __int8 i8_t;
  typedef u64_t camb_size_t;

#else /*!WIN32 || WINDOWS*/
  typedef uint64_t u64_t;
  typedef int64_t i64_t;
  typedef uint32_t u32_t;
  typedef uint16_t u16_t;
  typedef uint8_t u8_t;
  typedef int32_t i32_t;
  typedef int16_t i16_t;
  typedef int8_t i8_t;
  typedef u64_t camb_size_t;

#endif /*WIN32||WINDOWS*/
#endif /*__CNPAPI_TYPES_H*/

typedef enum {
  CNPAPI_API_ENTER = 0,
  CNPAPI_API_EXIT = 1
} cnpapi_CallbackSite;

typedef struct {
  /**
   * The point where the callback was issued.
   */
  cnpapi_CallbackSite callbackSite;
  /**
   * Name of the API function.
   */
  const char * functionName;
  /**
   * Pointer to the arguments for each API function.
   * See more structure definitions of the parameters in cndrv_params.h, cnrt_params.h, cnml_params.h and cnnl_params.h.
   */
  const void * functionParams;
  /**
   * Pointer to the return value of the API function.
   * This field is only valid within CNPAPI_API_EXIT.
   */
  const void* functionReturnValue;
  /**
   * Reserved for future use.
   */
  const char * symbolName;
  /**
   * The activity record correlation ID for this callback.
   * When this field is 0, it doesn't make sense.
   */
  uint64_t correlationId;
  /**
   * Reserved for future use.
   */
  uint64_t reserved1;
  /**
   * Reserved for future use.
   */
  uint64_t reserved2;
  /**
   * Reserved for future use.
   */
  uint64_t reserved3;
} cnpapi_CallbackData;

typedef enum {
  CNPAPI_CB_DOMAIN_CNRT_API = 0,
  CNPAPI_CB_DOMAIN_CNML_API = 1,
  CNPAPI_CB_DOMAIN_CNDRV_API = 2,
  CNPAPI_CB_DOMAIN_CNNL_API = 3,
  CNPAPI_CB_DOMAIN_CNPX_API = 4,
  CNPAPI_CB_DOMAIN_CNNL_EXTRA_API = 5,
  CNPAPI_CB_DOMAIN_SIZE = 6,
  CNPAPI_CB_DOMAIN_FORCE_INT = 0x7fffffff
} cnpapi_CallbackDomain;
typedef i32_t cnpapi_CallbackId;

typedef void (*cnpapi_CallbackFunc)(void * userdata,
                                   cnpapi_CallbackDomain domain,
                                   cnpapi_CallbackId cbid,
                                   const void* cbdata);
typedef void* cnpapi_SubscriberHandle;
#ifdef __cplusplus
}
#endif

#endif  // CNPAPI_HOOKLIB_CALLBACKAPI_TYPES_H_
