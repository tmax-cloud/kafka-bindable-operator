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

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kafkastrimziiov1beta2 "github.com/planner26/kafka-bindableOperator2/api/v1beta2"
)

type Listeners struct {
	Name string `json:"name"`
	Port int    `json:"port"`
	Tls  bool   `json:"tls"`
	Type string `json:"type"`
}

// KafkaReconciler reconciles a Kafka object
type KafkaReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kafka.strimzi.io.my.domain,resources=kafkas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kafka.strimzi.io.my.domain,resources=kafkas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kafka.strimzi.io.my.domain,resources=kafkas/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Kafka object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *KafkaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	kafkares := &kafkastrimziiov1beta2.Kafka{}

	var resultJson = make(map[string]interface{})

	err := r.Get(ctx, req.NamespacedName, kafkares)
	if err != nil {
		println(err)

	} else {

		data, err := json.Marshal(kafkares)
		err = json.Unmarshal(data, &resultJson)
		if err != nil {
			println(err)
		}

		listeners := ((resultJson["spec"].(map[string]interface{}))["kafka"].(map[string]interface{}))["listeners"].([]interface{})
		resourceName := (resultJson["metadata"].(map[string]interface{}))["name"].(string)
		resourceNamespace := (resultJson["metadata"].(map[string]interface{}))["namespace"].(string)

		secret := &v1.Secret{}
		reqSecret := ctrl.Request{}
		reqSecret.NamespacedName.Name = resourceName + "-cluster-ca-cert"
		reqSecret.NamespacedName.Namespace = resourceNamespace

		err = r.Get(ctx, reqSecret.NamespacedName, secret)
		if err != nil {
			println(err)
			return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
		} else {

			n := 0
			for n < len(listeners) {

				listener := listeners[n].(map[string]interface{})
				name := listener["name"].(string)
				port := int(listener["port"].(float64))
				openType := listener["type"].(string)
				tls := listener["tls"].(bool)

				fmt.Println("tls: ", tls)

				if !tls {
					for k := range secret.Data {
						delete(secret.Data, k)
					}
				}

				bootStrapServerKey := strings.ToUpper(name) + "_" + strings.ToUpper(openType) + "_" + "BOOTSTRAP_SERVERS"
				bootStrapServerValue := []byte(resourceName + "-kafka-bootstrap." + resourceNamespace + ".svc:" + strconv.Itoa(port))

				secret.Data[bootStrapServerKey] = bootStrapServerValue
				println(bootStrapServerKey, " : ", bootStrapServerValue)

				n++

			}

			goto CREATE
		}

	CREATE:
		{
			secret.Name = resourceName + "-service-binding-credentials"
			secret.ResourceVersion = ""
			err = r.Create(ctx, secret, &client.CreateOptions{})
			if err != nil {
				println(err)
			}
		}

	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KafkaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kafkastrimziiov1beta2.Kafka{}).
		Complete(r)
}
