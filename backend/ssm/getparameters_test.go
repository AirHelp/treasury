package ssm

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/AirHelp/treasury/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// batchingSSMClient serves any name it is asked for and remembers how the names
// were split between calls
type batchingSSMClient struct {
	ClientInterface
	batches [][]string
}

func (b *batchingSSMClient) GetParameters(_ context.Context, input *ssm.GetParametersInput, _ ...func(*ssm.Options)) (*ssm.GetParametersOutput, error) {
	if len(input.Names) > maxKeysPerCall {
		return nil, fmt.Errorf("asked for %d names on one call, SSM allows %d", len(input.Names), maxKeysPerCall)
	}
	b.batches = append(b.batches, input.Names)

	parameters := make([]ssmtypes.Parameter, 0, len(input.Names))
	for _, name := range input.Names {
		value := "value of " + name
		parameters = append(parameters, ssmtypes.Parameter{Name: &name, Value: &value})
	}
	return &ssm.GetParametersOutput{Parameters: parameters}, nil
}

func TestClient_GetObjectsByKeysBatches(t *testing.T) {
	// 23 keys have to be split into 10 + 10 + 3
	keys := make([]string, 0, 23)
	for i := range 23 {
		keys = append(keys, fmt.Sprintf("test/webapp/KEY_%02d", i))
	}

	svc := &batchingSSMClient{}
	got, err := (&Client{svc: svc}).GetObjects(&types.GetObjectsInput{Keys: keys})
	if err != nil {
		t.Fatal(err)
	}

	wantBatches := []int{10, 10, 3}
	if len(svc.batches) != len(wantBatches) {
		t.Fatalf("GetObjects() made %d calls, want %d", len(svc.batches), len(wantBatches))
	}
	for i, want := range wantBatches {
		if len(svc.batches[i]) != want {
			t.Errorf("call %d carried %d names, want %d", i+1, len(svc.batches[i]), want)
		}
	}

	if len(got.Secrets) != len(keys) {
		t.Fatalf("GetObjects() returned %d secrets, want %d", len(got.Secrets), len(keys))
	}
	// the leading slash SSM needs is not a part of a treasury key
	for _, key := range keys {
		value, found := got.Secrets[key]
		if !found {
			t.Fatalf("GetObjects() did not return %q", key)
		}
		if want := "value of /" + key; value != want {
			t.Errorf("GetObjects()[%q] = %q, want %q", key, value, want)
		}
	}
}

// missingSSMClient reports every name as invalid, the way SSM does for a
// parameter which does not exist
type missingSSMClient struct {
	ClientInterface
}

func (m *missingSSMClient) GetParameters(_ context.Context, input *ssm.GetParametersInput, _ ...func(*ssm.Options)) (*ssm.GetParametersOutput, error) {
	invalid := append([]string{}, input.Names...)
	sort.Strings(invalid)
	return &ssm.GetParametersOutput{InvalidParameters: invalid}, nil
}

func TestClient_GetObjectsByKeysReportsMissing(t *testing.T) {
	keys := []string{"test/webapp/GONE", "test/webapp/ALSO_GONE"}

	_, err := (&Client{svc: &missingSSMClient{}}).GetObjects(&types.GetObjectsInput{Keys: keys})
	if err == nil {
		t.Fatal("GetObjects() with unknown keys returned no error")
	}
	for _, key := range keys {
		if !strings.Contains(err.Error(), key) {
			t.Errorf("GetObjects() error = %q, want it to name %q", err, key)
		}
	}
}
