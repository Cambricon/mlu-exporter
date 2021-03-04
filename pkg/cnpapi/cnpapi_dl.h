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

#ifndef _CNPAPI_DL_H_
#define _CNPAPI_DL_H_

#include "./include/cnpapi.h"
#include "./include/cnpapi_types.h"
#include "./include/cnpapi_pmu_api.h"

#define CNPAPI_DL(x) x##_dl

extern cnpapiResult CNPAPI_DL(cnpapiInit)(void);
#endif // _CNPAPI_DL_H_
