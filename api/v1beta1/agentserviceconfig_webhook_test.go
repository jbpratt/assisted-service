/*
Copyright 2020.

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

package v1beta1

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Validate AgentServiceConfig", func() {
	var (
		agentServiceConfig = &AgentServiceConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-agentserviceconfig",
				Namespace: "test-agentserviceconfig-namespace",
			},
			Spec: AgentServiceConfigSpec{
				OSImages: []OSImage{
					{
						Version: "413.92.202303190222-0",
					},
					{
						Version: "413.92.202303190222-0",
					},
				},
			},
		}
	)

	It("create fails if containing duplicate OSImages", func() {
		warn, err := agentServiceConfig.ValidateCreate()
		Expect(warn).To(BeNil())
		Expect(err).NotTo(BeNil())
	})

	It("update fails if containing duplicate OSImages", func() {
		warn, err := agentServiceConfig.ValidateCreate()
		Expect(warn).To(BeNil())
		Expect(err).NotTo(BeNil())
	})
})
