# \MapsApi

All URIs are relative to *http://127.0.0.1:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetMap**](MapsApi.md#GetMap) | **Get** /maps/{key} | Get map value by key
[**GetMapBySubKey**](MapsApi.md#GetMapBySubKey) | **Get** /maps/{key}/{subkey} | Get value of element in map defined by sub-key
[**SetMap**](MapsApi.md#SetMap) | **Post** /maps/{key} | Set map value defined by key


# **GetMap**
> MapValue GetMap($key)

Get map value by key

Returns an map value defined by specified key


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **key** | **string**| key of value to return | 

### Return type

[**MapValue**](MapValue.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetMapBySubKey**
> StringValue GetMapBySubKey($key, $subkey)

Get value of element in map defined by sub-key

Returns a value of map element defined by sub-key defined by specified key


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **key** | **string**| key of value to return | 
 **subkey** | **string**| key of the element inside the map | 

### Return type

[**StringValue**](StringValue.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/text

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SetMap**
> SetMap($key, $body)

Set map value defined by key

Sets map value defined by key. If TTL is omitted, ttl of this value would not be set


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **key** | **string**| key of value to set | 
 **body** | [**SetMapReq**](SetMapReq.md)| Value and ttl of the map element | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

