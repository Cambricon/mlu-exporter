/*
 * Copyright 2022 Cambricon, Inc.
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

#ifndef CNPAPI_ACTIVITY_API_H_
#define CNPAPI_ACTIVITY_API_H_
#include <stdint.h>
#include <stddef.h>
#include "callbackapi_types.h"
#include "cnpapi_types.h"
#ifdef __cplusplus
extern "C" {
#endif
typedef enum {
  CNPAPI_ACTIVITY_TYPE_UNKNOWN = 0,
  CNPAPI_ACTIVITY_TYPE_KERNEL,
  CNPAPI_ACTIVITY_TYPE_BANGC,
  CNPAPI_ACTIVITY_TYPE_CNDRV_API,
  CNPAPI_ACTIVITY_TYPE_CNRT_API,
  CNPAPI_ACTIVITY_TYPE_CNML_API,
  CNPAPI_ACTIVITY_TYPE_CNNL_API,
  CNPAPI_ACTIVITY_TYPE_OVERHEAD,
  CNPAPI_ACTIVITY_TYPE_RESERVED_0,
  CNPAPI_ACTIVITY_TYPE_MEMCPY,
  CNPAPI_ACTIVITY_TYPE_MEMSET,
  CNPAPI_ACTIVITY_TYPE_MEMCPY_PTOP,
  CNPAPI_ACTIVITY_TYPE_CNNL_EXTRA_API,
  CNPAPI_ACTIVITY_TYPE_NAME,
  CNPAPI_ACTIVITY_TYPE_ATOMIC_OPERATION,
  CNPAPI_ACTIVITY_TYPE_NOTIFIER,
  CNPAPI_ACTIVITY_TYPE_CNCL_API,
  CNPAPI_ACTIVITY_TYPE_SOURCELOCATOR,
  CNPAPI_ACTIVITY_TYPE_FUNCTION,
  CNPAPI_ACTIVITY_TYPE_MODULE,
  CNPAPI_ACTIVITY_TYPE_PC_SAMPLING,
  CNPAPI_ACTIVITY_TYPE_QUEUE,
  CNPAPI_ACTIVITY_TYPE_CONTEXT,
  CNPAPI_ACTIVITY_TYPE_DEVICE,
  CNPAPI_ACTIVITY_TYPE_PCIE,
  CNPAPI_ACTIVITY_TYPE_PC_SAMPLING_RECORD_INFO,
  CNPAPI_ACTIVITY_TYPE_RESERVED_1,
  CNPAPI_ACTIVITY_TYPE_EXTERNAL_CORRELATION,
  CNPAPI_ACTIVITY_TYPE_RESERVED_2,
  CNPAPI_ACTIVITY_TYPE_SIZE,
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
  const void *return_value;
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
  /**
   * A notifier object requires that we identify device, context,
   * and notifier ID.
   */
  struct {
    uint64_t device_id;
    uint64_t context_id;
    uint64_t notifier_id;
  } dcn;
} cnpapiActivityObjectTypeId;

typedef enum {
  CNPAPI_ACTIVITY_OBJECT_UNKNOWN = 0,
  CNPAPI_ACTIVITY_OBJECT_PROCESS = 1,
  CNPAPI_ACTIVITY_OBJECT_THREAD = 2,
  CNPAPI_ACTIVITY_OBJECT_DEVICE = 3,
  CNPAPI_ACTIVITY_OBJECT_CONTEXT = 4,
  CNPAPI_ACTIVITY_OBJECT_QUEUE = 5,
  CNPAPI_ACTIVITY_OBJECT_NOTIFIER = 6
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
  const char *name;
  uint64_t queue_id;
  /* cnrt correlation id */
  uint64_t runtime_correlation_id;
  /* the value of this field is equivalent to KernelClass. */
  uint64_t kernel_type;
  /* The dimension of x. */
  uint32_t dimx;
  /* The dimension of y. */
  uint32_t dimy;
  /* The dimension of z. */
  uint32_t dimz;
  /* context id */
  uint64_t context_id;
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
  /* context id */
  uint64_t context_id;
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
  /* context id */
  uint64_t context_id;
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
  /* context id */
  uint64_t context_id;
} cnpapiActivityMemset;

typedef struct cnpapiActivityName {
  cnpapiActivityType type;
  /* the object id */
  cnpapiActivityObjectTypeId object_id;
  /* the cnpx name */
  const char *name;
  /* the kind of activity object to be named */
  cnpapiActivityObjectType object_type;
} cnpapiActivityName;

typedef enum {
  CNPAPI_ACTIVITY_ATOMIC_OP_REQUEST = 0,
  CNPAPI_ACTIVITY_ATOMIC_OP_COMPARE = 1,
} cnpapiActivityAtomicOpType;

typedef enum {
  /* Compares input data1 and opPtr until opData1 == *opPtr */
  CNPAPI_ACTIVITY_FLAG_ATOMIC_COMPARE_EQUAL = 0,
  /* Compares input data1 and opPtr until opData1 <= *opPtr */
  CNPAPI_ACTIVITY_FLAG_ATOMIC_COMPARE_LESS_EQUAL = 1,
  /* Compares input data1 and opPtr until opData1 < *opPtr */
  CNPAPI_ACTIVITY_FLAG_ATOMIC_COMPARE_LESS = 2
} cnpapiActivityAtomicCompareFlag;

typedef enum {
  /* Unknown request type */
  CNPAPI_ACTIVITY_FLAG_ATOMIC_REQUEST_TYPE_UNKNOWN = 0,
  /* Default request operation, which is the same as CN_ATOMIC_REQUEST_ADD */
  CNPAPI_ACTIVITY_FLAG_ATOMIC_REQUEST_DEFAULT = 1,
  /* Atomic add, this is default operation */
  CNPAPI_ACTIVITY_FLAG_ATOMIC_ADD = 2,
  /* Sets operation address to input value */
  CNPAPI_ACTIVITY_FLAG_ATOMIC_SET = 3,
  /* Resets operation address to zero */
  CNPAPI_ACTIVITY_FLAG_ATOMIC_CLEAR = 4,
} cnpapiActivityAtomicRequestFlag;

typedef struct cnpapiActivityAtomicOperation {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_ATOMIC_OPERATION */
  cnpapiActivityType type;
  /* Operation Type : REQUEST or COMPARE */
  cnpapiActivityAtomicOpType operation_type;
  /* CNAtomicReqFlag or CNAtomicCompFlag */
  union {
    cnpapiActivityAtomicRequestFlag req_flag;
    cnpapiActivityAtomicCompareFlag com_flag;
  };
  /* cndrv correlation id */
  uint64_t correlation_id;
  /* tensor processor start timestamp */
  uint64_t start;
  uint64_t end;
  /* timestamp when mlu driver received the task,
     a value of 0 indicates that this field could not be collected.*/
  uint64_t received;
  /* timestamp when mlu driver queued the task into task buffer,
     a value of 0 indicates that this field could not be collected.*/
  uint64_t queued;
  /* device id */
  uint64_t device_id;
  /* queue id */
  uint64_t queue_id;
  /* operation target value */
  uint64_t value;
  /* context id */
  uint64_t context_id;
} cnpapiActivityAtomicOperation;
typedef enum {
  CNPAPI_ACTIVITY_NOTIFIER_WAIT = 0,
  CNPAPI_ACTIVITY_NOTIFIER_PLACE = 1,
} cnpapiActivityNotifierTaskType;

typedef struct cnpapiActivityNotifier {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_NOTIFIER */
  cnpapiActivityType type;
  /* notifier task type */
  cnpapiActivityNotifierTaskType task_type;
  /* cndrv correlation id */
  uint64_t correlation_id;
  /* notifier op start timestamp */
  uint64_t start;
  /* notifier op end timestamp */
  uint64_t end;
  /* timestamp when mlu driver received the task,
     a value of 0 indicates that this field could not be collected.*/
  uint64_t received;
  /* device id */
  uint64_t device_id;
  /* queue id */
  uint64_t queue_id;
  /* notifier id */
  uint64_t notifier_id;
  /* context id */
  uint64_t context_id;
} cnpapiActivityNotifier;

typedef struct {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_FUNCTION */
  cnpapiActivityType type;
  /* context_id of this function */
  uint64_t context_id;
  /* module_id of this function */
  uint64_t module_id;
  /* index of this function in this module */
  uint64_t function_index;
  /* global id the identify this record */
  uint64_t function_id;
  /* function name */
  const char *name;
} cnpapiActivityFunction;

typedef struct cnpapiActivitySourceLocator {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_SOURCELOCATOR */
  cnpapiActivityType type;
  /* file name */
  const char *file_name;
  /* global id the identify this record */
  uint64_t id;
  /* line number in file_ame */
  uint64_t line_number;
} cnpapiActivitySourceLocator;

typedef struct cnpapiActivityPcSampling {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_PC_SAMPLING */
  cnpapiActivityType type;
  /* The correlation ID of the API to which this result is associated */
  uint64_t correlation_id;
  /* global id the identify this record */
  uint64_t function_id;
  /* line number in file_ame */
  uint64_t pc_offset;
  /* source locator record id */
  uint64_t source_locator_id;
  /* samples of this pc */
  uint64_t samples;
} cnpapiActivityPcSampling;

typedef struct cnpapiActivityModule {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_MODULE */
  cnpapiActivityType type;
  /* module id of this module */
  uint64_t module_id;
  /* cnbin address of this module */
  uint64_t cnbin_addr;
  /* cnbin size of this module */
  uint64_t cnbin_size;
  /* context id of this module */
  uint64_t context_id;
} cnpapiActivityModule;

typedef struct cnpapiActivityContext {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_CONTEXT */
  cnpapiActivityType type;
  /* The ID of thie context */
  uint64_t context_id;
  /* device id */
  uint64_t device_id;
  /* The ID of NULL queue in this context, a value of 0 indicates that this field could not be collected */
  uint64_t null_queue_id;
} cnpapiActivityContext;

typedef struct cnpapiActivityQueue {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_QUEUE */
  cnpapiActivityType type;
  /* The ID of the context where the queue was created */
  uint64_t context_id;
  /* The correlation ID of the API to which this result is associated */
  uint64_t correlation_id;
  /* The priority of the queue */
  int priority;
  /* A unique queue ID to identify the queue */
  uint64_t queue_id;
} cnpapiActivityQueue;

typedef struct cnpapiActivity {
  cnpapiActivityType type;
} cnpapiActivity;

typedef struct cnpapiActivityDevice {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_DEVICE */
  cnpapiActivityType type;
  /* device id */
  uint32_t device_id;
  /* Major compute capability version number */
  uint32_t compute_capability_major;
  /* Minor compute capability version number */
  uint32_t compute_capability_minor;
  /* Maximum x-dimension of a block task. */
  uint32_t max_block_dimx;
  /* Maximum y-dimension of a block task. */
  uint32_t max_block_dimy;
  /* Maximum z-dimension of a block task.  */
  uint32_t max_block_dimz;
  /* Maximum number of clusters per union task. */
  uint32_t max_cluster_count_per_union_task;
  /* Number of clusters on Device. */
  uint32_t max_cluster_count;
  /* Maximum number of MLU cores per cluster. */
  uint32_t max_core_count_per_cluster;
  /* Size of L2 cache in bytes. */
  uint32_t max_l2_cache_size;
  /* Maximum size of Neural-RAM available per MLU core in bytes. */
  uint32_t max_nram_size_per_core;
  /* Maximum size of Weight-RAM available per MLU core in bytes. */
  uint32_t max_wram_size_per_core;
  /* Memory available on device for __mlu_const__ variables in a BANG C kernel in megabytes */
  uint32_t const_memory_size;
  /* Maximum size of local memory available per MLU core (i.e. __ldram__) in megabytes. */
  uint32_t local_memory_size_per_core;
  /* Maximum size of shared memory per cluster (i.e. __mlu_shared__) in bytes. */
  uint32_t shared_memory_size_per_cluster;
  /* Global memory bus width in bits. */
  uint32_t global_memory_bus_width;
  /* Total global memory size in megabytes. */
  uint32_t global_memory_size;
  /* Device has ECC support enabled */
  uint32_t ecc_enabled;
  /* Typical cluster clock frequency in kilohertz. */
  uint32_t cluster_clock_rate;
  /* Peak memory clock frequency in kHz */
  uint32_t peak_memory_clock_rate;
  const char *uuid;
  const char *name;
} cnpapiActvitiyDevice;

typedef struct cnpapiActivityPcie {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_PCIE */
  cnpapiActivityType type;
  /* The ID of the MLU Device */
  uint32_t device_id;
  /* PCI bus ID of the MLU Device */
  uint32_t pci_bus_id;
  /* PCI Device ID of the MLU Device */
  uint32_t pci_device_id;
  /* PCI domain ID of the MLU Device */
  uint32_t pci_domain_id;
} cnpapiActicityPcie;

typedef struct cnpapiActivityPcSamplingRecordInfo {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_PC_SAMPLING_RECORD_INFO */
  cnpapiActivityType type;
  /* total samples of pc sampling */
  uint64_t total_samples;
  /* dropped samples of pc sampling */
  uint64_t dropped_samples;
  /* The correlation ID of the API to which this result is associated */
  uint64_t correlation_id;
} cnpapiActivityPcSamplingRecordInfo;

typedef enum {
  CNPAPI_EXTERNAL_CORRELATION_TYPE_INVALID = 0,
  CNPAPI_EXTERNAL_CORRELATION_TYPE_UNKNOWN = 1,
  CNPAPI_EXTERNAL_CORRELATION_TYPE_CUSTOM0 = 2,
  CNPAPI_EXTERNAL_CORRELATION_TYPE_CUSTOM1 = 3,
  CNPAPI_EXTERNAL_CORRELATION_TYPE_CUSTOM2 = 4,
  /* Add new types before this line */
  CNPAPI_EXTERNAL_CORRELATION_TYPE_SIZE,
  CNPAPI_EXTERNAL_CORRELATION_TYPE_FORCE_INT = 0x7fffffff
} cnpapiActivityExternalCorrelationType;

typedef struct cnpapiActivityExternalCorrelation {
  /* activity type, must be CNPAPI_ACTIVITY_TYPE_EXTERNAL_CORRELATION */
  cnpapiActivityType type;
  /* The type of external API this record correlated to */
  cnpapiActivityExternalCorrelationType external_type;
  /*
     The correlation ID of the associated non-Cambricon API record.
     The exact field in the associated external record depends on that
     record's activity type.
  */
  uint64_t external_id;
  /* The correlation ID of the associated API record */
  uint64_t correlation_id;
} cnpapiActivityExternalCorrelation;

typedef void (*cnpapi_request)(uint64_t **buffer, size_t *size, size_t *maxNumRecords);
typedef void (*cnpapi_complete)(uint64_t *buffer, size_t size, size_t validSize);
CNPAPI_EXPORT cnpapiResult cnpapiActivityEnable(cnpapiActivityType type);
CNPAPI_EXPORT cnpapiResult cnpapiActivityRegisterCallbacks(cnpapi_request bufferRequested, cnpapi_complete bufferCompleted);
CNPAPI_EXPORT cnpapiResult cnpapiActivityGetNextRecord(void *buffer, size_t validSize, cnpapiActivity **record);
CNPAPI_EXPORT cnpapiResult cnpapiActivityDisable(cnpapiActivityType type);
CNPAPI_EXPORT cnpapiResult cnpapiActivityFlushAll();
CNPAPI_EXPORT cnpapiResult cnpapiActivityFlushPeriod(uint64_t time);
CNPAPI_EXPORT cnpapiResult cnpapiActivityPushExternalCorrelationId(cnpapiActivityExternalCorrelationType type, uint64_t id);
CNPAPI_EXPORT cnpapiResult cnpapiActivityPopExternalCorrelationId(cnpapiActivityExternalCorrelationType type, uint64_t *lastId);

#ifdef __cplusplus
}
#endif
#endif  // CNPAPI_ACTIVITY_API_H_
