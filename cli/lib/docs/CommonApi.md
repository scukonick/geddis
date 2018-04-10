# \CommonApi

All URIs are relative to *http://127.0.0.1:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Delete**](CommonApi.md#Delete) | **Delete** /delete/{key} | Deletes value defined by key
[**GetKeys**](CommonApi.md#GetKeys) | **Get** /keys/{key} | returns keys starting with &#39;key&#39;


# **Delete**
> Delete($key)

Deletes value defined by key

Deletes value of any type defined by key


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **key** | **string**| key of value to set | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetKeys**
> Array GetKeys($key)

returns keys starting with 'key'

Returns array of keys matching 'key*'. If key == '*' it returns all the keys


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **key** | **string**| prefix of keys to return | 

### Return type

[**Array**](Array.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

