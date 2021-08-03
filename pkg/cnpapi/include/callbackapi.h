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

#ifndef CALLBACKAPI_H_
#define CALLBACKAPI_H_
#include "cnpapi_types.h"  // NOLINT
#include "callbackapi_types.h"  // NOLINT
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
