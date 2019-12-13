/*
 * Cloud Hypervisor API
 *
 * Local HTTP based API for managing and inspecting a cloud-hypervisor virtual machine.
 *
 * API version: 0.3.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi
// NetConfig struct for NetConfig
type NetConfig struct {
	Tap string `json:"tap,omitempty"`
	Ip string `json:"ip"`
	Mask string `json:"mask"`
	Mac string `json:"mac"`
	Iommu bool `json:"iommu,omitempty"`
}