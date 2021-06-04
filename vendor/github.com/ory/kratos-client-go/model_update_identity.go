/*
 * Ory Kratos API
 *
 * Documentation for all public and administrative Ory Kratos APIs. Public and administrative APIs are exposed on different ports. Public APIs can face the public internet without any protection while administrative APIs should never be exposed without prior authorization. To protect the administative API port you should use something like Nginx, Ory Oathkeeper, or any other technology capable of authorizing incoming requests. 
 *
 * API version: v0.6.3-alpha.1
 * Contact: hi@ory.sh
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
)

// UpdateIdentity struct for UpdateIdentity
type UpdateIdentity struct {
	// SchemaID is the ID of the JSON Schema to be used for validating the identity's traits. If set will update the Identity's SchemaID.
	SchemaId *string `json:"schema_id,omitempty"`
	// Traits represent an identity's traits. The identity is able to create, modify, and delete traits in a self-service manner. The input will always be validated against the JSON Schema defined in `schema_id`.
	Traits map[string]interface{} `json:"traits"`
}

// NewUpdateIdentity instantiates a new UpdateIdentity object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpdateIdentity(traits map[string]interface{}) *UpdateIdentity {
	this := UpdateIdentity{}
	this.Traits = traits
	return &this
}

// NewUpdateIdentityWithDefaults instantiates a new UpdateIdentity object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpdateIdentityWithDefaults() *UpdateIdentity {
	this := UpdateIdentity{}
	return &this
}

// GetSchemaId returns the SchemaId field value if set, zero value otherwise.
func (o *UpdateIdentity) GetSchemaId() string {
	if o == nil || o.SchemaId == nil {
		var ret string
		return ret
	}
	return *o.SchemaId
}

// GetSchemaIdOk returns a tuple with the SchemaId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UpdateIdentity) GetSchemaIdOk() (*string, bool) {
	if o == nil || o.SchemaId == nil {
		return nil, false
	}
	return o.SchemaId, true
}

// HasSchemaId returns a boolean if a field has been set.
func (o *UpdateIdentity) HasSchemaId() bool {
	if o != nil && o.SchemaId != nil {
		return true
	}

	return false
}

// SetSchemaId gets a reference to the given string and assigns it to the SchemaId field.
func (o *UpdateIdentity) SetSchemaId(v string) {
	o.SchemaId = &v
}

// GetTraits returns the Traits field value
func (o *UpdateIdentity) GetTraits() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}

	return o.Traits
}

// GetTraitsOk returns a tuple with the Traits field value
// and a boolean to check if the value has been set.
func (o *UpdateIdentity) GetTraitsOk() (map[string]interface{}, bool) {
	if o == nil  {
		return nil, false
	}
	return o.Traits, true
}

// SetTraits sets field value
func (o *UpdateIdentity) SetTraits(v map[string]interface{}) {
	o.Traits = v
}

func (o UpdateIdentity) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.SchemaId != nil {
		toSerialize["schema_id"] = o.SchemaId
	}
	if true {
		toSerialize["traits"] = o.Traits
	}
	return json.Marshal(toSerialize)
}

type NullableUpdateIdentity struct {
	value *UpdateIdentity
	isSet bool
}

func (v NullableUpdateIdentity) Get() *UpdateIdentity {
	return v.value
}

func (v *NullableUpdateIdentity) Set(val *UpdateIdentity) {
	v.value = val
	v.isSet = true
}

func (v NullableUpdateIdentity) IsSet() bool {
	return v.isSet
}

func (v *NullableUpdateIdentity) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUpdateIdentity(val *UpdateIdentity) *NullableUpdateIdentity {
	return &NullableUpdateIdentity{value: val, isSet: true}
}

func (v NullableUpdateIdentity) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUpdateIdentity) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


