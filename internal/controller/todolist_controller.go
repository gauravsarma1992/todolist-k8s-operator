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
	"time"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	log "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	todov1 "sarmag.co/todo/api/v1"
)

type TodoListReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=todo.sarmag.co,resources=todolists,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=todo.sarmag.co,resources=todolists/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=todo.sarmag.co,resources=todolists/finalizers,verbs=update

func (r *TodoListReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	var (
		todoList todov1.TodoList
		podList  corev1.PodList
		logger   logr.Logger

		isCompleted bool
	)

	logger = log.FromContext(ctx)
	logger.Info("Reconciling TodoList")

	if err = r.Get(ctx, req.NamespacedName, &todoList); err != nil {
		logger.Error(err, "Error in fetching Todolist")
		err = client.IgnoreNotFound((err))
		return
	}

	if err = r.List(ctx, &podList); err != nil {
		logger.Error(err, "Error in fetching pods list")
		return
	}

	for _, item := range podList.Items {
		if item.GetName() != todoList.Spec.Task {
			continue
		}
		logger.Info("Pod just became available with", "name", item.GetName())
		isCompleted = true
	}

	todoList.Status.IsCompleted = isCompleted
	if err = r.Status().Update(ctx, &todoList); err != nil {
		logger.Error(err, "Error in updating TodoList", "status", isCompleted)
		return
	}

	if todoList.Status.IsCompleted == true {
		result.RequeueAfter = time.Minute * 2
	}
	return
}

func (r *TodoListReconciler) startTickerLoop(periodicReconcileCh chan event.GenericEvent) {
	var (
		ticker *time.Ticker
		count  int
	)
	ticker = time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			periodicReconcileCh <- event.GenericEvent{Object: &todov1.TodoList{ObjectMeta: metav1.ObjectMeta{Name: "jack", Namespace: "operator-namespace"}}}

			count += 1
			if count > 100 {
				return
			}
		}
	}
}

func (r *TodoListReconciler) SetupWithManager(mgr ctrl.Manager) (err error) {
	var (
		periodicReconcileCh chan event.GenericEvent
	)
	periodicReconcileCh = make(chan event.GenericEvent)
	go r.startTickerLoop(periodicReconcileCh)

	err = ctrl.NewControllerManagedBy(mgr).
		For(&todov1.TodoList{}).
		Watches(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Channel{Source: periodicReconcileCh}, &handler.EnqueueRequestForObject{}).
		Complete(r)
	return
}
