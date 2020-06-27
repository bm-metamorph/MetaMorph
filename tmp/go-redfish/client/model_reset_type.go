/*
 * Redfish OAPI specification
 *
 * Partial Redfish OAPI specification for a limited client
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package client
// ResetType the model 'ResetType'
type ResetType string

// List of ResetType
const (
	RESETTYPE_ON ResetType = "On"
	RESETTYPE_FORCE_OFF ResetType = "ForceOff"
	RESETTYPE_GRACEFUL_SHUTDOWN ResetType = "GracefulShutdown"
	RESETTYPE_GRACEFUL_RESTART ResetType = "GracefulRestart"
	RESETTYPE_FORCE_RESTART ResetType = "ForceRestart"
	RESETTYPE_NMI ResetType = "Nmi"
	RESETTYPE_FORCE_ON ResetType = "ForceOn"
	RESETTYPE_PUSH_POWER_BUTTON ResetType = "PushPowerButton"
	RESETTYPE_POWER_CYCLE ResetType = "PowerCycle"
)