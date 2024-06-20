/*
Copyright 2024.

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
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	databasesv1 "github.com/leeminki/my-operator/api/v1"
)

// +kubebuilder:rbac:groups=databases.test.io,resources=databasebackups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=databases.test.io,resources=databasebackups/status,verbs=get;update;patch

// DatabaseBackupReconciler는 DatabaseBackup 객체를 관리하는 컨트롤러
type DatabaseBackupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile 함수는 DatabaseBackup 객체의 상태를 읽고 변경 사항을 반영
func (r *DatabaseBackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// DatabaseBackup 인스턴스를 가져옴
	backup := &databasesv1.DatabaseBackup{}
	err := r.Get(ctx, req.NamespacedName, backup)
	if err != nil {
		if errors.IsNotFound(err) {
			// DatabaseBackup 리소스를 찾을 수 없는 경우, 삭제된 것으로 간주하고 처리하지 않음
			return ctrl.Result{}, nil
		}
		// 리소스를 읽는 도중 에러가 발생한 경우, 재시도를 위해 에러를 반환
		return ctrl.Result{}, err
	}

	// 백업 작업이 이미 진행 중인지 확인
	if backup.Status.InProgress {
		// 이미 백업 작업이 진행 중인 경우, 재시도하지 않음
		return ctrl.Result{}, nil
	}

	// 백업 작업을 시작할 때 InProgress 상태를 true로 설정
	backup.Status.InProgress = true
	// Update 시 ResourceVersion을 확인하여 상태 업데이트 일관성을 보장
	latestBackup := backup.DeepCopy()
	err = r.Status().Update(ctx, latestBackup)
	if err != nil {
		log.Error(err, "Failed to update DatabaseBackup status to InProgress")
		return ctrl.Result{RequeueAfter: time.Minute * 1}, err
	}

	// 여기서 백업 로직을 구현
	fmt.Printf("Backing up database: %s\n", backup.Spec.DatabaseName)

	// 백업 작업이 완료되면 InProgress 상태를 false로 설정하고 LastBackupTime을 업데이트
	latestBackup.Status.LastBackupTime = metav1.Now()
	latestBackup.Status.InProgress = false
	// Update 시 ResourceVersion을 확인하여 상태 업데이트 일관성을 보장
	err = r.Status().Update(ctx, latestBackup)
	if err != nil {
		log.Error(err, "Failed to update DatabaseBackup status")
		return ctrl.Result{RequeueAfter: time.Minute * 1}, err
	}

	// 일정 시간 후에 다시 요청을 큐에 넣음 (예: 1분 후).
	return ctrl.Result{RequeueAfter: time.Minute * 1}, nil
}

// SetupWithManager 함수는 컨트롤러를 매니저에 등록
func (r *DatabaseBackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// 컨트롤러를 생성하고 DatabaseBackup 리소스를 관리하도록 설정
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasesv1.DatabaseBackup{}).
		Complete(r)
}
