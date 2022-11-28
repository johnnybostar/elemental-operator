/*
Copyright © 2022 SUSE LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"testing"

	elementalv1 "github.com/rancher/elemental-operator/api/v1beta1"
	managementv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestUnauthenticatedResponse(t *testing.T) {
	testCase := []struct {
		config *elementalv1.Config
		regUrl string
	}{
		{
			config: nil,
			regUrl: "https://rancher/url1",
		},
		{
			config: &elementalv1.Config{
				Elemental: elementalv1.Elemental{
					Registration: elementalv1.Registration{
						EmulateTPM:      true,
						EmulatedTPMSeed: 127,
						NoSMBIOS:        true,
					},
				},
			},
			regUrl: "https://rancher/url2",
		},
	}

	for _, test := range testCase {
		scheme := runtime.NewScheme()
		elementalv1.AddToScheme(scheme)
		managementv3.AddToScheme(scheme)

		i := &InventoryServer{
			Context: context.Background(),
			Client:  fake.NewClientBuilder().Build(),
		}

		registration := elementalv1.MachineRegistration{}
		registration.Spec.Config = test.config
		registration.Status.RegistrationURL = test.regUrl

		buffer := new(bytes.Buffer)

		err := i.unauthenticatedResponse(&registration, buffer)
		assert.NilError(t, err, err)

		conf := elementalv1.Config{}
		err = yaml.NewDecoder(buffer).Decode(&conf)
		assert.NilError(t, err, strings.TrimSpace(buffer.String()))
		assert.Equal(t, conf.Elemental.Registration.URL, test.regUrl)

		confReg := conf.Elemental.Registration
		testReg := elementalv1.Registration{}
		if test.config != nil {
			testReg = test.config.Elemental.Registration
		}
		assert.Equal(t, confReg.EmulateTPM, testReg.EmulateTPM)
		assert.Equal(t, confReg.EmulatedTPMSeed, testReg.EmulatedTPMSeed)
		assert.Equal(t, confReg.NoSMBIOS, testReg.NoSMBIOS)
	}
}

func TestInitNewInventory(t *testing.T) {
	const alphanum = "[0-9a-fA-F]"
	// m  '-'  8 alphanum chars  '-'  3 blocks of 4 alphanum chars  '-'  12 alphanum chars
	mUUID := regexp.MustCompile("^m-" + alphanum + "{8}-(" + alphanum + "{4}-){3}" + alphanum + "{12}")
	// e.g., m-66588488-3eb6-4a6d-b642-c994f128c6f1

	testCase := []struct {
		config       *elementalv1.Config
		initName     string
		expectedName string
	}{
		{
			config: &elementalv1.Config{
				Elemental: elementalv1.Elemental{
					Registration: elementalv1.Registration{
						NoSMBIOS: false,
					},
				},
			},
			initName:     "custom-name",
			expectedName: "custom-name",
		},

		{
			config: &elementalv1.Config{
				Elemental: elementalv1.Elemental{
					Registration: elementalv1.Registration{
						NoSMBIOS: false,
					},
				},
			},
			expectedName: "m-${System Information/UUID}",
		},
		{
			config: &elementalv1.Config{
				Elemental: elementalv1.Elemental{
					Registration: elementalv1.Registration{
						NoSMBIOS: true,
					},
				},
			},
		},
		{
			config:       nil,
			expectedName: "m-${System Information/UUID}",
		},
	}

	for _, test := range testCase {
		registration := &elementalv1.MachineRegistration{
			Spec: elementalv1.MachineRegistrationSpec{
				MachineName: test.initName,
				Config:      test.config,
			},
		}

		inventory := &elementalv1.MachineInventory{}
		initInventory(inventory, registration)

		if test.config != nil && test.config.Elemental.Registration.NoSMBIOS {
			assert.Check(t, mUUID.Match([]byte(inventory.Name)), inventory.Name+" is not UUID based")
		} else {
			assert.Equal(t, inventory.Name, test.expectedName)
		}
	}
}

func TestBuildName(t *testing.T) {
	data := map[string]interface{}{
		"level1A": map[string]interface{}{
			"level2A": "level2AValue",
			"level2B": map[string]interface{}{
				"level3A": "level3AValue",
			},
		},
		"level1B": "level1BValue",
	}

	testCase := []struct {
		Format string
		Output string
	}{
		{
			Format: "${level1B}",
			Output: "level1BValue",
		},
		{
			Format: "${level1B",
			Output: "m-level1B",
		},
		{
			Format: "a${level1B",
			Output: "a-level1B",
		},
		{
			Format: "${}",
			Output: "m",
		},
		{
			Format: "${",
			Output: "m-",
		},
		{
			Format: "a${",
			Output: "a-",
		},
		{
			Format: "${level1A}",
			Output: "m",
		},
		{
			Format: "a${level1A}c",
			Output: "ac",
		},
		{
			Format: "a${level1A}",
			Output: "a",
		},
		{
			Format: "${level1A}c",
			Output: "c",
		},
		{
			Format: "a${level1A/level2A}c",
			Output: "alevel2AValuec",
		},
		{
			Format: "a${level1A/level2B/level3A}c",
			Output: "alevel3AValuec",
		},
		{
			Format: "a${level1A/level2B/level3A}c${level1B}",
			Output: "alevel3AValueclevel1BValue",
		},
	}

	for _, testCase := range testCase {
		assert.Equal(t, testCase.Output, buildStringFromSmbiosData(data, testCase.Format))
	}
}

func TestMergeInventoryLabels(t *testing.T) {
	testCase := []struct {
		data     []byte            // labels to add to the inventory
		labels   map[string]string // labels already in the inventory
		fail     bool
		expected map[string]string
	}{
		{
			[]byte(`{"key2":"val2"}`),
			map[string]string{"key1": "val1"},
			false,
			map[string]string{"key1": "val1"},
		},
		{
			[]byte(`{"key2":2}`),
			map[string]string{"key1": "val1"},
			true,
			map[string]string{"key1": "val1"},
		},
		{
			[]byte(`{"key2":"val2", "key3":"val3"}`),
			map[string]string{"key1": "val1", "key3": "previous_val", "key4": "val4"},
			false,
			map[string]string{"key1": "val1", "key3": "previous_val"},
		},
		{
			[]byte{},
			map[string]string{"key1": "val1"},
			true,
			map[string]string{"key1": "val1"},
		},
		{
			[]byte(`{"key2":"val2"}`),
			nil,
			false,
			map[string]string{},
		},
	}

	for _, test := range testCase {
		inventory := &elementalv1.MachineInventory{}
		inventory.Labels = test.labels

		err := mergeInventoryLabels(inventory, test.data)
		if test.fail {
			assert.Assert(t, err != nil)
		} else {
			assert.Equal(t, err, nil)
		}
		for k, v := range test.expected {
			val, ok := inventory.Labels[k]
			assert.Equal(t, ok, true)
			assert.Equal(t, v, val)
		}

	}
}
