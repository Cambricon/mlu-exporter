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

#ifndef ACTIVITY_ACTIVITY_API_H_
#define ACTIVITY_ACTIVITY_API_H_
#include "callbackapi_types.h"  // NOLINT
#include "cnpapi.h"  // NOLINT
#ifdef __cplusplus
extern "C" {
#endif
typedef enum {
  CNPAPI_ACTIVITY_TYPE_UNKNOWN = 0,
  CNPAPI_ACTIVITY_TYPE_KERNEL = 1,
  CNPAPI_ACTIVITY_TYPE_BANGC = 2,
  CNPAPI_ACTIVITY_TYPE_CNDRV_API = 3,
  CNPAPI_ACTIVITY_TYPE_CNRT_API = 4,
  CNPAPI_ACTIVITY_TYPE_CNML_API = 5,
  CNPAPI_ACTIVITY_TYPE_CNNL_API = 6,
  CNPAPI_ACTIVITY_TYPE_OVERHEAD = 7,
  CNPAPI_ACTIVITY_TYPE_RESERVED_0 = 8,
  CNPAPI_ACTIVITY_TYPE_MEMCPY = 9,
  CNPAPI_ACTIVITY_TYPE_MEMSET = 10,
  CNPAPI_ACTIVITY_TYPE_MEMCPY_PTOP = 11,
  CNPAPI_ACTIVITY_TYPE_CNNL_EXTRA_API = 12
} cnpapiActivityType;
typedef enum {
  CNPAPI_ACTIVITY_FLAG_NONE = 0,
  CNPAPI_ACTIVITY_FLAG_MEMCPY_ASYNC = 1 << 0,
  CNPAPI_ACTIVITY_FLAG_MEMSET_ASYNC = 1 << 0
} cnpapiActivityFlag;

typedef struct cnpapiActivityAPI {
  /* activity type */
  cnpapiActivityType type;
  /* identity id */
  uint64_t correlation_id;
  /* callback api id */
  cnpapi_CallbackId cbid;
  /* task start ts */
  uint64_t start;
  /* task end ts */
  uint64_t end;
  /* pid */
  uint32_t process_id;
  /* tid */
  uint32_t thread_id;
  /* return value */
  const void * return_value;
} cnpapiActivityAPI;

typedef enum {
  /* unknown overhead */
  CNPAPI_ACTIVITY_OVERHEAD_UNKNOWN = 0,
  /* activity buffer flush overhead,
     includes internal buffer flush,
     user(activity) buffer flush and etc.*/
  CNPAPI_ACTIVITY_OVERHEAD_CNPAPI_BUFFER_FLUSH = 1,
  /* cnpapi resource creation, destruction or manipulation overhead,
     includes internal processing logic*/
  CNPAPI_ACTIVITY_OVERHEAD_CNPAPI_RESOURCE = 2
} cnpapiActivityOverheadType;

/**
 * \brief Identifiers for object kinds as specified by
 * cnpapiActivityOverheadType.
 * \see cnpapiActivityOverheadType
 */
typedef union {
  /**
   * A process object requires that we identify the process ID. A
   * thread object requires that we identify both the process and
   * thread ID.
   */
  struct {
    uint32_t process_id;
    uint32_t thread_id;
  } pt;
  /**
   * A device object requires that we identify the device ID. A
   * context object requires that we identify both the device and
   * context ID. A queue object requires that we identify device,
   * context, and queue ID.
   */
  struct {
    uint64_t device_id;
    uint64_t context_id;
    uint64_t queue_id;
  } dcq;
} cnpapiActivityObjectTypeId;

typedef enum {
  CNPAPI_ACTIVITY_OBJECT_UNKNOWN = 0,
  CNPAPI_ACTIVITY_OBJECT_PROCESS = 1,
  CNPAPI_ACTIVITY_OBJECT_THREAD = 2,
  CNPAPI_ACTIVITY_OBJECT_DEVICE = 3,
  CNPAPI_ACTIVITY_OBJECT_CONTEXT = 4,
  CNPAPI_ACTIVITY_OBJECT_QUEUE = 5
} cnpapiActivityObjectType;

typedef struct cnpapiActivityOverhead {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_OVERHEAD */
  cnpapiActivityType type;
  /* overhead type */
  cnpapiActivityOverheadType overhead_type;
  /* start ts, a value of 0 indicates that this field could not be collected. */
  uint64_t start;
  /* end ts, a value of 0 indicates that this field could not be collected.*/
  uint64_t end;
  /* the type of activity object that the overhead is associated with */
  cnpapiActivityObjectType object_type;
  /* the identifier for activity object */
  cnpapiActivityObjectTypeId object_id;
} cnpapiActivityOverhead;

typedef struct cnpapiActivityKernel {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_KERNEL */
  cnpapiActivityType type;
  /* cndrv correlation id */
  uint64_t correlation_id;
  /* tensor processor start timestamp */
  union {
  uint64_t start_ts;  // deprecated
  uint64_t start;
  };
  /* tensor processor end timestamp */
  union {
    uint64_t end_ts;  // deprecated
    uint64_t end;
  };
  /* timestamp when mlu driver received the task,
     a value of 0 indicates that this field could not be collected.*/
  uint64_t received;
  /* timestamp when mlu driver queued the task into task buffer,
     a value of 0 indicates that this field could not be collected.*/
  uint64_t queued;
  /* timestamp when mlu driver pushded the task into job scheduler,
     a value of 0 indicates that this field could not be collected.*/
  uint64_t submitted;
  /* device id */
  uint64_t device_id;
  /* kernel name, a value of 0 indicates that this field could not be collected. */
  const char * name;
  uint64_t queue_id;
  /* cnrt correlation id */
  uint64_t runtime_correlation_id;
  /* the value of this field is equivalent to MLUKernelClass. */
  uint64_t kernel_type;
} cnpapiActivityKernel;
typedef enum {
  CNPAPI_ACTIVITY_MEMCPY_TYPE_UNKNOWN = 0,
  CNPAPI_ACTIVITY_MEMCPY_TYPE_HTOD = 1,
  CNPAPI_ACTIVITY_MEMCPY_TYPE_DTOH = 2,
  CNPAPI_ACTIVITY_MEMCPY_TYPE_DTOD = 3,
  CNPAPI_ACTIVITY_MEMCPY_TYPE_HTOH = 4,
  CNPAPI_ACTIVITY_MEMCPY_TYPE_PTOP = 5
} cnpapiActivityMemcpyType;
typedef struct cnpapiActivityMemcpy {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_MEMCPY */
  cnpapiActivityType type;
  /* the flags associated with the memory copy */
  cnpapiActivityFlag flags;
  /* cndrv correlation id */
  uint64_t correlation_id;
  /* the number of bytes transferred by the memory copy. */
  uint64_t bytes;
  /* the kind of the memory copy */
  cnpapiActivityMemcpyType copy_type;
  /* memcpy start timestamp */
  uint64_t start;
  /* memcpy end timestamp */
  uint64_t end;
  /* timestamp when mlu driver received the task,
     a value of 0 indicates that this field could not be collected.*/
  uint64_t received;
  /* timestamp when mlu driver queued the task into task buffer,
     a value of 0 indicates that this field could not be collected.*/
  uint64_t queued;
  /* timestamp when mlu driver pushded the task into job scheduler,
     a value of 0 indicates that this field could not be collected. */
  uint64_t submitted;
  /* device id,
   * a value of (uint64_t)-1 indicates that this field could not be collected.  */
  uint64_t device_id;
  /* queue id,
   * a value of 0 indicates that this field could not be collected,
   * or means default queue. */
  uint64_t queue_id;
  /* cnrt correlation id */
  uint64_t runtime_correlation_id;
} cnpapiActivityMemcpy;

typedef struct cnpapiActivityMemcpyPtoP {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_MEMCPY_PTOP */
  cnpapiActivityType type;
  /* the flags associated with the memory copy */
  cnpapiActivityFlag flags;
  /* cndrv correlation id */
  uint64_t correlation_id;
  /* the number of bytes transferred by the memory copy. */
  uint64_t bytes;
  /* the kind of the memory copy */
  cnpapiActivityMemcpyType copy_type;
  /* memcpy start timestamp */
  uint64_t start;
  /* memcpy end timestamp */
  uint64_t end;
  /* timestamp when mlu driver received the task,
     a value of 0 indicates that this field could not be collected.*/
  uint64_t received;
  /* timestamp when mlu driver queued the task into task buffer,
     a value of 0 indicates that this field could not be collected.*/
  uint64_t queued;
  /* timestamp when mlu driver pushded the task into job scheduler,
     a value of 0 indicates that this field could not be collected. */
  uint64_t submitted;
  /* device id,
   * a value of (uint64_t)-1 indicates that this field could not be collected.  */
  uint64_t device_id;
  /* source device id,
   * a value of (uint64_t)-1 indicates that this field could not be collected.  */
  uint64_t src_device_id;
  /* destination device id,
   * a value of (uint64_t)-1 indicates that this field could not be collected.  */
  uint64_t dst_device_id;
  /* queue id,
   * a value of 0 indicates that this field could not be collected,
   * or means default queue. */
  uint64_t queue_id;
  /* cnrt correlation id */
  uint64_t runtime_correlation_id;
} cnpapiActivityMemcpyPtoP;

typedef struct cnpapiActivityMemset {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_MEMSET */
  cnpapiActivityType type;
  /* the flags associated with the memory set */
  cnpapiActivityFlag flags;
  /* cndrv correlation id */
  uint64_t correlation_id;
  /* the number of bytes being set by the memory set. */
  uint64_t bytes;
  /* memset start timestamp */
  uint64_t start;
  /* memset end timestamp */
  uint64_t end;
  /* timestamp when mlu driver received the task,
     a value of 0 indicates that this field could not be collected.*/
  uint64_t received;
  /* timestamp when mlu driver queued the task into task buffer,
     a value of 0 indicates that this field could not be collected.*/
  uint64_t queued;
  /* timestamp when mlu driver pushded the task into job scheduler,
     a value of 0 indicates that this field could not be collected. */
  uint64_t submitted;
  /* device id,
   * a value of (uint64_t)-1 indicates that this field could not be collected.  */
  uint64_t device_id;
  /* queue id,
   * a value of 0 indicates that this field could not be collected,
   * or means default queue. */
  uint64_t queue_id;
  /* cnrt correlation id */
  uint64_t runtime_correlation_id;
  /* the value being assigned to memory by the memory set.*/
  uint64_t value;
} cnpapiActivityMemset;
typedef struct cnpapiActivity {
  cnpapiActivityType type;
} cnpapiActivity;

typedef void (*cnpapi_request)(uint64_t **buffer, size_t *size, size_t *maxNumRecords);
typedef void (*cnpapi_complete)(uint64_t *buffer, size_t size, size_t validSize);
CNPAPI_EXPORT cnpapiResult cnpapiActivityEnable(cnpapiActivityType type);
CNPAPI_EXPORT cnpapiResult cnpapiActivityRegisterCallbacks(cnpapi_request bufferRequested,
                                                           cnpapi_complete bufferCompleted);
CNPAPI_EXPORT cnpapiResult cnpapiActivityGetNextRecord(void *buffer, size_t validSize,
                                                       cnpapiActivity **record);
CNPAPI_EXPORT cnpapiResult cnpapiActivityDisable(cnpapiActivityType type);
CNPAPI_EXPORT cnpapiResult cnpapiActivityFlushAll();
#ifdef __cplusplus
}
#endif
#endif  // ACTIVITY_ACTIVITY_API_H_

