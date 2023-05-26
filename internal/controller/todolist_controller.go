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

package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	todov1 "sarmag.co/todo/api/v1"
)

// TodoListReconciler reconciles a TodoList object
type TodoListReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=todo.sarmag.co,resources=todolists,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=todo.sarmag.co,resources=todolists/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=todo.sarmag.co,resources=todolists/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TodoList object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *TodoListReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling todoList custom resource")

	// Get the TodoList resource that triggered the reconciliation request
	var todoList todov1.TodoList
	if err := r.Get(ctx, req.NamespacedName, &todoList); err != nil {
		log.Error(err, "unable to fetch TodoList")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get pods with the same name as TodoList's friend
	var podList corev1.PodList
	var isCompleted bool
	if err := r.List(ctx, &podList); err != nil {
		log.Error(err, "unable to list pods")
	} else {
		for _, item := range podList.Items {
			if item.GetName() == todoList.Spec.Task {
				log.Info("pod linked to a todoList custom resource found", "name", item.GetName())
				isCompleted = true
			}
		}
	}

	// Update TodoList' happy status
	todoList.Status.IsCompleted = isCompleted
	if err := r.Status().Update(ctx, &todoList); err != nil {
		log.Error(err, "unable to update todoList's happy status", "status", isCompleted)
		return ctrl.Result{}, err
	}
	log.Info("todoList's happy status updated", "status", isCompleted)

	log.Info("todoList custom resource reconciled")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TodoListReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&todov1.TodoList{}).
		Complete(r)
}
