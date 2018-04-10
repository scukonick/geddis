# \StringsApi

All URIs are relative to *http://127.0.0.1:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetString**](StringsApi.md#GetString) | **Get** /strings/{key} | Get string value by key
[**SetString**](StringsApi.md#SetString) | **Post** /strings/{key} | Set string value defined by key


# **GetString**
> StringValue GetString($key)

Get string value by key

Returns a string value defined by specified key


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **key** | **string**| key of value to return | 

### Return type

[**StringValue**](StringValue.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/text

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SetString**
> SetString($key, $body)

Set string value defined by key

Sets string value defined by key. If TTL is omitted, ttl of this value would not be set


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **key** | **string**| key of value to set | 
 **body** | [**SetStringValueReq**](SetStringValueReq.md)| Value and ttl of the string element | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

