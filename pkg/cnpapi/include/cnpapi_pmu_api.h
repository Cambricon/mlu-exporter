#ifndef CNPAPI_PMU_API_H_
#define CNPAPI_PMU_API_H_
#include <stdint.h>
#include "cnpapi_types.h"  // NOLINT
#ifdef __cplusplus
extern "C" {
#endif
typedef enum {
  CNPAPI_PMU_AUTO_FLUSH,
  CNPAPI_PMU_EXPLICIT_FLUSH
} cnpapiPmuFlushMode;
/**
 * @brief: used to enable counter by counter id.
 *
 * @param: dev_id[in]: device id.
 * @param: counter_id[in]: counter id.
 * @param: is_enable[in]: true to enable counter and false to disable counter.
 * @ret: *CNPAPI_SUCCESS*                           on success.
 *       *CNPAPI_ERROR_NOT_INITIALIZED*             cnpapi not initialized.
 *       *CNPAPI_ERROR_INVALID_PMU_COUNTER_ID*      invalid counter id.
 *       *CNPAPI_ERROR_INVALID_DEVICE_ID*           invalid device.
 *       *CNPAPI_ERROR_DRIVER_COMMUNICATION_FAILED* driver communication failed.
 *       *CNPAPI_ERROR_UNKNOWN*                     internal error.
 */
CNPAPI_EXPORT cnpapiResult cnpapiPmuEnableCounter(int dev_id, uint64_t counter_id, bool is_enable);

/**
 * @brief: used to flush counter data of a device.
 *
 * @param: dev_id[in]: device id.
 * @ret: *CNPAPI_SUCCESS*                           on success.
 *       *CNPAPI_ERROR_NOT_INITIALIZED*             cnpapi not initialized.
 *       *CNPAPI_ERROR_INVALID_DEVICE_ID*           invalid device.
 *       *CNPAPI_ERROR_DRIVER_COMMUNICATION_FAILED* driver communication failed.
 */
CNPAPI_EXPORT cnpapiResult cnpapiPmuFlushData(int dev_id);

/**
 * @brief: used to get counter id by counter name.
 *
 * @param: name[in]: counter name.
 * @param: counter_id[out]: counter id addr.
 * @ret: *CNPAPI_SUCCESS*                       on success.
 *       *CNPAPI_ERROR_NOT_INITIALIZED*         cnpapi not initialized.
 *       *CNPAPI_ERROR_INVALID_ARGUMENT*        counter_id is NULL or name is not a valid counter name.
 */
CNPAPI_EXPORT cnpapiResult cnpapiPmuGetCounterIdByName(const char *name, uint64_t *counter_id);

/**
 * @brief: used to get counter name by counter id.
 *
 * @param: counter_id[in]: counter id.
 * @param: name[out]: name dst.
 * @ret: *CNPAPI_SUCCESS*                      on success.
 *       *CNPAPI_ERROR_NOT_INITIALIZED*        cnpapi not initialized.
 *       *CNPAPI_ERROR_INVALID_ARGUMENT*       name is NULL.
 *       *CNPAPI_ERROR_INVALID_PMU_COUNTER_ID* unknown counter id.
 */
CNPAPI_EXPORT cnpapiResult cnpapiPmuGetCounterName(uint64_t counter_id, const char **name);

/**
 * @brief: used to get supported counters of mlu chip.
 *
 * @param: chip_type[in]: the card chip type.
 * @param: counter_dst[out]: counter id array pointer. the supported counter id will be copyed to this dst.
 * @param: size[in/out]: the size of given counter_dst, and the supported counter size will be set to this addr.
 * @ret: *CNPAPI_SUCCESS*                   on success.
 *       *CNPAPI_ERROR_NOT_INITIALIZED*     cnpapi not initialized.
 *       *CNPAPI_ERROR_INSUFFICIENT_MEMORY* size less than needed.
 *       *CNPAPI_ERROR_INVALID_ARGUMENT*    size or counter_dst is NULL.
 *       *CNPAPI_ERROR_INVALID_CHIP_TYPE*   unknown chip type.
 */
CNPAPI_EXPORT cnpapiResult cnpapiPmuGetCounterSupported(cnpapiChipType chip_type,
                                                        uint64_t *counter_dst, uint64_t *size);

/**
 * @brief: used to get counter value by counter id.
 *
 * @param: dev_id[in]: device id.
 * @param: counter_id[in]: counter id.
 * @param: value[out]: addr to get counter value.
 * @ret: *CNPAPI_SUCCESS*                           on success.
 *       *CNPAPI_ERROR_NOT_INITIALIZED*             cnpapi not initialized.
 *       *CNPAPI_ERROR_INVALID_ARGUMENT*            value is NULL.
 *       *CNPAPI_ERROR_INVALID_PMU_COUNTER_ID*      invalid counter id.
 *       *CNPAPI_ERROR_INVALID_DEVICE_ID*           invalid device.
 *       *CNPAPI_ERROR_PMU_COUNTER_NOT_ENABLED*     counter not enabled.
 *       *CNPAPI_ERROR_DRIVER_COMMUNICATION_FAILED* driver communication failed.
 *       *CNPAPI_ERROR_UNKNOWN*                     internal error.
 */
CNPAPI_EXPORT cnpapiResult cnpapiPmuGetCounterValue(int dev_id, uint64_t counter_id, uint64_t *value);

/**
 * @brief: used to set data flush mode, the default mode is CNPAPI_PMU_AUTO_FLUSH.
 *
 * @param: mode[in]: flush mode to set
 *                   CNPAPI_PMU_AUTO_FLUSH will automatic flush data when calling cnpapiPmuGetCounterValue,
 *                   CNPAPI_PMU_EXPLICIT_FLUSH will flush data only when calling cnpapiPmuFlushData.
 * @ret: *CNPAPI_SUCCESS*                           on success.
 *       *CNPAPI_ERROR_INVALID_ARGUMENT*            invalid mode.
 */
CNPAPI_EXPORT cnpapiResult cnpapiPmuSetFlushMode(cnpapiPmuFlushMode mode);

#ifdef __cplusplus
}
#endif
#endif  // CNPAPI_PMU_API_H_
