# \ArraysApi

All URIs are relative to *http://127.0.0.1:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetArrByIndex**](ArraysApi.md#GetArrByIndex) | **Get** /arrays/{key}/{index} | Get string value by key
[**GetArray**](ArraysApi.md#GetArray) | **Get** /arrays/{key} | Get array value by key
[**SetArray**](ArraysApi.md#SetArray) | **Post** /arrays/{key} | Set array value defined by key


# **GetArrByIndex**
> StringValue GetArrByIndex($key, $index)

Get string value by key

Returns an array value defined by specified key


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **key** | **string**| key of value to return | 
 **index** | **int32**| index of element in array to return | 

### Return type

[**StringValue**](StringValue.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/text

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetArray**
> Array GetArray($key)

Get array value by key

Returns an array value defined by specified key


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **key** | **string**| key of value to return | 

### Return type

[**Array**](Array.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SetArray**
> SetArray($key, $body)

Set array value defined by key

Sets string value defined by key. If TTL is omitted, ttl of this value would not be set


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **key** | **string**| key of value to set | 
 **body** | [**SetArrayReq**](SetArrayReq.md)| Value and ttl of the string element | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

