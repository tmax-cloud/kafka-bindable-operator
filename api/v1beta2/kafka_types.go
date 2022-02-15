/*
Copyright 2022.

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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KafkaSpec defines the desired state of Kafka

type Listeners struct {
	Name string `json:"name"`
	Port int    `json:"port"`
	Tls  bool   `json:"tls"`
	Type string `json:"type"`
}

type Kafkastruct struct {
	Listeners []Listeners `json:"listeners"`
	Version   string      `json:"version"`
	Replicas  int         `json:"replicas"`
}

type Zookeeperstruct struct {
	Replicas int           `json:"replicas"`
	Storage  storagestruct `json:"storage"`
}

type storagestruct struct {
	Type string `json:"type"`
}

type EntityOperator struct {
	TopicOperator string `json:"topicOperator"`
	UserOperator  string `json:"userOperator"`
}

type KafkaSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Kafka. Edit kafka_types.go to remove/update
	Kafka     Kafkastruct     `json:"kafka"`
	Zookeeper Zookeeperstruct `json:"zookeeper"`
	// EntityOperator EntityOperator  `json:"entityOperator"`
}

// KafkaStatus defines the observed state of Kafka
type KafkaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Kafka is the Schema for the kafkas API
type Kafka struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KafkaSpec   `json:"spec,omitempty"`
	Status KafkaStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KafkaList contains a list of Kafka
type KafkaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kafka `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kafka{}, &KafkaList{})
}
