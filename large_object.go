// Copyright 2013-2014 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

type LargeObject interface {
	packageName() string

	Destroy() error
	Size() (int, error)
	GetConfig() (map[interface{}]interface{}, error)
	SetCapacity(capacity int) error
	GetCapacity() (int, error)
}

// Create and manage a large object within a single bin. A stack is last in/first out (LIFO).
///
type baseLargeObject struct {
	client     *Client
	policy     *WritePolicy
	key        *Key
	binName    Value
	userModule Value
}

// Initialize large stack operator.
//
// client        client
// policy        generic configuration parameters, pass in null for defaults
// key         unique record identifier
// binName       bin name
// userModule      Lua function name that initializes list configuration parameters, pass null for default set
func newLargeObject(client *Client, policy *WritePolicy, key *Key, binName string, userModule string) *baseLargeObject {
	r := &baseLargeObject{
		client:  client,
		policy:  policy,
		key:     key,
		binName: NewStringValue(binName),
	}

	if userModule == "" {
		r.userModule = NewNullValue()
	} else {
		r.userModule = NewStringValue(userModule)
	}

	return r
}

// Delete bin containing the object.
func (lo *baseLargeObject) destroy(ifc LargeObject) error {
	_, err := lo.client.Execute(lo.policy, lo.key, ifc.packageName(), "destroy", lo.binName)
	return err
}

// Return size of object.
func (lo *baseLargeObject) size(ifc LargeObject) (int, error) {
	ret, err := lo.client.Execute(lo.policy, lo.key, ifc.packageName(), "size", lo.binName)
	if err != nil {
		return -1, err
	}
	return ret.(int), nil
}

// Return map of object configuration parameters.
func (lo *baseLargeObject) getConfig(ifc LargeObject) (map[interface{}]interface{}, error) {
	ret, err := lo.client.Execute(lo.policy, lo.key, ifc.packageName(), "get_config", lo.binName)
	return ret.(map[interface{}]interface{}), err
}

// Set maximum number of entries in the object.
//
// capacity      max entries in set
func (lo *baseLargeObject) setCapacity(ifc LargeObject, capacity int) error {
	_, err := lo.client.Execute(lo.policy, lo.key, ifc.packageName(), "set_capacity", lo.binName, NewIntegerValue(capacity))
	return err
}

// Return maximum number of entries in the object.
func (lo *baseLargeObject) getCapacity(ifc LargeObject) (int, error) {
	ret, err := lo.client.Execute(lo.policy, lo.key, ifc.packageName(), "get_capacity", lo.binName)
	if err != nil {
		return -1, err
	}
	return ret.(int), nil
}

// Return list of all objects on the stack.
func (lo *baseLargeObject) scan(ifc LargeObject) ([]interface{}, error) {
	ret, err := lo.client.Execute(lo.policy, lo.key, ifc.packageName(), "scan", lo.binName)
	if err != nil {
		return nil, err
	}
	return ret.([]interface{}), nil
}
