#ifndef CALLBACKAPI_H_
#define CALLBACKAPI_H_
#include "cnpapi_types.h"
#include "callbackapi_types.h"
#ifdef __cplusplus
extern "C" {
#endif

CNPAPI_EXPORT cnpapiResult cnpapiSubscribe(cnpapi_SubscriberHandle* subscriber,
                                  cnpapi_CallbackFunc callback,
                                  void* userdata);

CNPAPI_EXPORT cnpapiResult cnpapiUnsubscribe(cnpapi_SubscriberHandle subscriber);

CNPAPI_EXPORT cnpapiResult cnpapiEnableCallback(u32_t enable,
                                       cnpapi_SubscriberHandle subscriber,
                                       cnpapi_CallbackDomain domain,
                                       cnpapi_CallbackId cbid);

CNPAPI_EXPORT cnpapiResult cnpapiEnableDomain(u32_t enable,
                                       cnpapi_SubscriberHandle subscriber,
                                       cnpapi_CallbackDomain domain);

CNPAPI_EXPORT cnpapiResult cnpapiEnableAllDomains(u32_t enable,
                                           cnpapi_SubscriberHandle subscriber);

CNPAPI_EXPORT cnpapiResult cnpapiGetCallbackState(u32_t *enable,
                                    cnpapi_SubscriberHandle subscriber,
                                    cnpapi_CallbackDomain domain,
                                    cnpapi_CallbackId cbid);
#ifdef __cplusplus
}
#endif

#endif  // CALLBACKAPI_H_
