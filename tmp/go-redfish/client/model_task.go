/*
 * Redfish OAPI specification
 *
 * Partial Redfish OAPI specification for a limited client
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package client
import (
	"time"
)
// Task This resource contains information about a specific Task scheduled by or being executed by a Redfish service's Task Service.
type Task struct {
	// The OData description of a payload.
	OdataContext string `json:"@odata.context,omitempty"`
	// The current ETag of the resource.
	OdataEtag string `json:"@odata.etag,omitempty"`
	// The name of the resource.
	OdataId string `json:"@odata.id"`
	// The type of a resource.
	OdataType string `json:"@odata.type"`
	// description
	Description *string `json:"Description,omitempty"`
	// The date-time stamp that the task was last completed.
	EndTime time.Time `json:"EndTime,omitempty"`
	// Indicates that the contents of the Payload should be hidden from view after the Task has been created.  When set to True, the Payload object will not be returned on GET.
	HidePayload bool `json:"HidePayload,omitempty"`
	// The name of the resource.
	Id string `json:"Id"`
	// This is an array of messages associated with the task.
	Messages []Message `json:"Messages,omitempty"`
	// The name of the resource.
	Name string `json:"Name"`
	// This is the manufacturer/provider specific extension moniker used to divide the Oem object into sections.
	Oem string `json:"Oem,omitempty"`
	Payload Payload `json:"Payload,omitempty"`
	// The date-time stamp that the task was last started.
	StartTime time.Time `json:"StartTime,omitempty"`
	// The URI of the Task Monitor for this task.
	TaskMonitor string `json:"TaskMonitor,omitempty"`
	TaskState TaskState `json:"TaskState,omitempty"`
	TaskStatus Health `json:"TaskStatus,omitempty"`
}
