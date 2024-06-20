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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DatabaseBackupSpec는 DatabaseBackup 리소스의 원하는 상태를 정의
type DatabaseBackupSpec struct {
	// 데이터베이스 이름
	DatabaseName string `json:"databaseName,omitempty"`
	// 백업 일정 (cron 형식)
	Schedule string `json:"schedule,omitempty"`
	// 백업 파일을 저장할 경로
	BackupPath string `json:"backupPath,omitempty"`
}

// DatabaseBackupStatus는 DatabaseBackup 리소스의 현재 상태를 정의
type DatabaseBackupStatus struct {
	// 마지막 백업 시간이 기록
	LastBackupTime metav1.Time `json:"lastBackupTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DatabaseBackup은 DatabaseBackup API의 스키마
type DatabaseBackup struct {
	// API의 타입 메타데이터를 포함
	metav1.TypeMeta `json:",inline"`
	// 오브젝트 메타데이터를 포함 (이름, 네임스페이스, 레이블 등).
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// 사용자가 설정하는 원하는 상태 (spec)를 포함
	Spec DatabaseBackupSpec `json:"spec,omitempty"`
	// 오퍼레이터가 관리하는 현재 상태 (status)를 포함
	Status DatabaseBackupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DatabaseBackupList는 여러 DatabaseBackup 리소스를 포함하는 리스트
type DatabaseBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatabaseBackup `json:"items"`
}

func init() {
	// DatabaseBackup와 DatabaseBackupList 타입을 스키마에 등록
	SchemeBuilder.Register(&DatabaseBackup{}, &DatabaseBackupList{})
}
