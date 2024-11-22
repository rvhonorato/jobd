package errors

import (
	"reflect"
	"testing"
)

func TestNewInternalServerError(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want *RestErr
	}{
		{
			name: "create a new internal server error",
			args: args{
				message: "test message",
			},
			want: &RestErr{
				Message: "test message",
				Status:  500,
				Error:   "internal_server_error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInternalServerError(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInternalServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBadRequestError(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want *RestErr
	}{
		{
			name: "create a new bad request error",
			args: args{
				message: "test message",
			},
			want: &RestErr{
				Message: "test message",
				Status:  400,
				Error:   "bad_request",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBadRequestError(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBadRequestError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConflictError(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want *RestErr
	}{
		{
			name: "create a new conflict error",
			args: args{
				message: "test message",
			},
			want: &RestErr{
				Message: "test message",
				Status:  409,
				Error:   "conflict",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConflictError(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConflictError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNotFoundError(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want *RestErr
	}{
		{
			name: "create a new not found error",
			args: args{
				message: "test message",
			},
			want: &RestErr{
				Message: "test message",
				Status:  404,
				Error:   "not_found",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotFoundError(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStatusUnauthorized(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want *RestErr
	}{
		{
			name: "create a new unauthorized error",
			args: args{
				message: "test message",
			},
			want: &RestErr{
				Message: "test message",
				Status:  401,
				Error:   "unauthorized",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStatusUnauthorized(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStatusUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}
