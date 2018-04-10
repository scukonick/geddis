package geddiclient

import (
	"fmt"
	"log"

	"github.com/scukonick/geddis/serverxxx"
)

// Client represents CLI client for geddis
type Client struct {
	stringsAPI *swagger.StringsApi
	arraysAPI  *swagger.ArraysApi
	mapsAPI    *swagger.MapsApi
	commonAPI  *swagger.CommonApi
}

// NewClient returns newly initialized Client structure
func NewClient(basePath string) *Client {
	return &Client{
		stringsAPI: swagger.NewStringsApiWithBasePath(basePath),
		arraysAPI:  swagger.NewArraysApiWithBasePath(basePath),
		mapsAPI:    swagger.NewMapsApiWithBasePath(basePath),
		commonAPI:  swagger.NewCommonApiWithBasePath(basePath),
	}
}

// GetString gets string value identified by 'key'
// and prints it
func (c *Client) GetString(key string) {
	val, _, err := c.stringsAPI.GetString(key)
	if err != nil {
		log.Printf("failed to get value: %v", err)
		return
	}

	fmt.Println(val.Value)
}

// SetString sets string value identified by 'key'
func (c *Client) SetString(key, value string, ttl int) {
	_, err := c.stringsAPI.SetString(key, swagger.SetStringValueReq{
		Value: value,
		Ttl:   int64(ttl),
	})
	if err != nil {
		log.Printf("failed to set value: %v", err)
		return
	}
}

// SetArr sets array value identifed by 'key'
func (c *Client) SetArr(key string, value []string, ttl int) {
	_, err := c.arraysAPI.SetArray(key, swagger.SetArrayReq{
		Values: value,
		Ttl:    int64(ttl),
	})
	if err != nil {
		log.Printf("failed to set value: %v", err)
		return
	}
}

// GetArr gets array value identified by 'key' and prints it
func (c *Client) GetArr(key string) {
	val, _, err := c.arraysAPI.GetArray(key)
	if err != nil {
		log.Printf("failed to get value: %v", err)
		return
	}

	fmt.Println(val.Values)
}

// GetArrIndex gets array value identified by 'key' by 'index' and prints it
func (c *Client) GetArrIndex(key string, index int32) {
	val, _, err := c.arraysAPI.GetArrByIndex(key, index)
	if err != nil {
		log.Printf("failed to get value: %v", err)
		return
	}

	fmt.Println(val.Value)
}

// SetMap sets maps value identifed by 'key'
func (c *Client) SetMap(key string, value map[string]string, ttl int) {
	_, err := c.mapsAPI.SetMap(key, swagger.SetMapReq{
		Value: value,
		Ttl:   int64(ttl),
	})
	if err != nil {
		log.Printf("failed to set value: %v", err)
		return
	}
}

// GetMap gets map value identified by 'key' and prints it
func (c *Client) GetMap(key string) {
	val, _, err := c.mapsAPI.GetMap(key)
	if err != nil {
		log.Printf("failed to get value: %v", err)
		return
	}

	fmt.Println(val)
}

// GetMapSubKey gets map value identified by 'key' and 'subKey' and prints it
func (c *Client) GetMapSubKey(key, subKey string) {
	val, _, err := c.mapsAPI.GetMapBySubKey(key, subKey)
	if err != nil {
		log.Printf("failed to get value: %v", err)
		return
	}

	fmt.Println(val.Value)
}
