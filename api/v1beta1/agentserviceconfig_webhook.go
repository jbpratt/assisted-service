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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var agentserviceconfiglog = logf.Log.WithName("agentserviceconfig-resource")

func (r *AgentServiceConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-agent-install-openshift-io-v1beta1-agentserviceconfig,mutating=false,failurePolicy=fail,sideEffects=None,groups=agent-install.openshift.io,resources=agentserviceconfigs,verbs=create;update,versions=v1beta1,name=vagentserviceconfig.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &AgentServiceConfig{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *AgentServiceConfig) ValidateCreate() (admission.Warnings, error) {
	agentserviceconfiglog.Info("validate create", "name", r.Name)

	if err := hasDuplicateOSImages(r.Spec.OSImages); err != nil {
		return nil, err
	}

	agentserviceconfiglog.Info("Successful validation")
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *AgentServiceConfig) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	agentserviceconfiglog.Info("validate update", "name", r.Name)

	if err := hasDuplicateOSImages(r.Spec.OSImages); err != nil {
		return nil, err
	}

	agentserviceconfiglog.Info("Successful validation")
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *AgentServiceConfig) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

func hasDuplicateOSImages(osImages []OSImage) error {
	seen := make(map[string]bool)
	for _, osImage := range osImages {
		if seen[osImage.Version] {
			err := fmt.Errorf("Failed validation: Found duplicate OSImages entry: %s", osImage.Version)
			agentserviceconfiglog.Info(err.Error())
			return err
		}
		seen[osImage.Version] = true
	}
	return nil
}
