package model

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	StatusRunning   = "Running"
	StatusSucceeded = "Succeeded"
	StatusFailed    = "Failed"
	StatusPending   = "Pending"
	StatusUnknown   = "Unknown"
	StatusActive    = "Active"
	StatusSuspended = "Suspended"
	StatusHealthy   = "Healthy"
	StatusAbnormal  = "Abnormal"
	StatusPartial   = "PartialAvailable"
	TimeFormat      = "2006-01-02 15:04:05"
)

func FormatTime(t *metav1.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Time.Format(TimeFormat)
}

func FormatTimeValue(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(TimeFormat)
}
