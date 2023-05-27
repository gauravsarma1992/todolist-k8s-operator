/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TodoListSpec defines the desired state of TodoList
type TodoListSpec struct {
	Task string `json:"task,omitempty"`
}

type TodoListStatus struct {
	IsCompleted bool `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

type TodoList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TodoListSpec   `json:"spec,omitempty"`
	Status TodoListStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TodoListList contains a list of TodoList
type TodoListList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TodoList `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TodoList{}, &TodoListList{})
}
