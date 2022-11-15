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

package types

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	elementalv1 "github.com/rancher/elemental-operator/api/v1beta1"
	"github.com/rancher/elemental-operator/pkg/services/syncer/config"
	"github.com/sirupsen/logrus"
)

type JSONSyncer struct {
	URI     string `json:"uri"`
	Timeout string `json:"timeout"`
}

func (j *JSONSyncer) Sync(c config.Config, s elementalv1.ManagedOSVersionChannel) ([]elementalv1.ManagedOSVersion, error) {
	logrus.Infof("Syncing '%s/%s' (JSON)", s.Namespace, s.Name)

	timeout := time.Second * 30
	if j.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(j.Timeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timeout: %w", err)
		}
	}
	client := &http.Client{
		Timeout: timeout,
	}

	logrus.Debug("Fetching JSON from ", j.URI)
	resp, err := client.Get(j.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s: %w", j.URI, err)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	res := []elementalv1.ManagedOSVersion{}

	err = json.Unmarshal(buf, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return res, nil
}