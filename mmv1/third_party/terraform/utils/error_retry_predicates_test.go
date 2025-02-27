package google

import (
	"strconv"
	"testing"

	"google.golang.org/api/googleapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestIsAppEngineRetryableError_operationInProgress(t *testing.T) {
	err := googleapi.Error{
		Code: 409,
		Body: "Operation is already in progress",
	}
	isRetryable, _ := isAppEngineRetryableError(&err)
	if !isRetryable {
		t.Errorf("Error not detected as retryable")
	}
}

func TestIsAppEngineRetryableError_p4saPropagation(t *testing.T) {
	err := googleapi.Error{
		Code: 404,
		Body: "Unable to retrieve P4SA: [service-111111111111@gcp-gae-service.iam.gserviceaccount.com] from GAIA. Could be GAIA propagation delay or request from deleted apps.",
	}
	isRetryable, _ := isAppEngineRetryableError(&err)
	if !isRetryable {
		t.Errorf("Error not detected as retryable")
	}
}

func TestIsAppEngineRetryableError_missingPage(t *testing.T) {
	err := googleapi.Error{
		Code: 404,
		Body: "Missing page",
	}
	isRetryable, _ := isAppEngineRetryableError(&err)
	if isRetryable {
		t.Errorf("Error incorrectly detected as retryable")
	}
}

func TestIsAppEngineRetryableError_serverError(t *testing.T) {
	err := googleapi.Error{
		Code: 500,
		Body: "Unable to retrieve P4SA because of a bad thing happening",
	}
	isRetryable, _ := isAppEngineRetryableError(&err)
	if isRetryable {
		t.Errorf("Error incorrectly detected as retryable")
	}
}

func TestIsCommonRetryableErrorCode_retryableErrorCode(t *testing.T) {
	codes := []int{429, 500, 502, 503}
	for _, code := range codes {
		code := code
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			err := googleapi.Error{
				Code: code,
				Body: "some text describing error",
			}
			isRetryable, _ := isCommonRetryableErrorCode(&err)
			if !isRetryable {
				t.Errorf("Error not detected as retryable")
			}
		})
	}
}

func TestIsCommonRetryableErrorCode_otherError(t *testing.T) {
	err := googleapi.Error{
		Code: 404,
		Body: "Some unretryable issue",
	}
	isRetryable, _ := isCommonRetryableErrorCode(&err)
	if isRetryable {
		t.Errorf("Error incorrectly detected as retryable")
	}
}

func TestIsOperationReadQuotaError_quotaExceeded(t *testing.T) {
	err := googleapi.Error{
		Code: 403,
		Body: "Quota exceeded for quota metric 'OperationReadGroup' and limit 'Operation read requests per minute' of service 'compute.googleapis.com' for consumer 'project_number:11111111'.",
	}
	isRetryable, _ := is403QuotaExceededPerMinuteError(&err)
	if !isRetryable {
		t.Errorf("Error not detected as retryable")
	}
}

func TestIs403QuotaExceededPerMinuteError_perMinuteQuotaExceeded(t *testing.T) {
	err := googleapi.Error{
		Code: 403,
		Body: "Quota exceeded for quota metric 'Queries' and limit 'Queries per minute' of service 'compute.googleapis.com' for consumer 'project_number:11111111'.",
	}
	isRetryable, _ := is403QuotaExceededPerMinuteError(&err)
	if !isRetryable {
		t.Errorf("Error not detected as retryable")
	}
}

func TestIs403QuotaExceededPerMinuteError_perDayQuotaExceededNotRetryable(t *testing.T) {
	err := googleapi.Error{
		Code: 403,
		Body: "Quota exceeded for quota metric 'Queries' and limit 'Queries per day' of service 'compute.googleapis.com' for consumer 'project_number:11111111'.",
	}
	isRetryable, _ := is403QuotaExceededPerMinuteError(&err)
	if isRetryable {
		t.Errorf("Error incorrectly detected as retryable")
	}
}

func TestGRPCRetryable(t *testing.T) {
	code := codes.FailedPrecondition
	err := status.Error(code, "is retryable")
	isRetryable, _ := isBigTableRetryableError(err)
	if !isRetryable {
		t.Errorf("Error not detected as retryable")
	}
}

func TestGRPCNotRetryable(t *testing.T) {
	code := codes.InvalidArgument
	err := status.Error(code, "is noto retryable")
	isRetryable, _ := isBigTableRetryableError(err)
	if isRetryable {
		t.Errorf("Error incorrectly detected as retryable")
	}
}
