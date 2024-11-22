// Package status provides the status of a job
package status

const (
	None            = "NONE"
	Submitted       = "SUBMITTED"
	Held            = "HELD"
	Queued          = "QUEUED"
	Running         = "RUNNING"
	Deleted         = "DELETED"
	Post_Processing = "POST_PROCESSING"
	Cancelled       = "CANCELLED"
	Failed          = "FAILED"
	Success         = "SUCCESS"
	Partial         = "PARTIAL"
	Unknown         = "UNKNOWN"
	Prepared        = "PREPARED"
	Created         = "CREATED"
)
