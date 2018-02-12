package ssm

import (
	"reflect"
	"testing"

	test "github.com/AirHelp/treasury/test/ssm"
	"github.com/AirHelp/treasury/types"
)

func TestClient_PutObject(t *testing.T) {
	c := &Client{
		svc: &test.MockSSMClient{},
	}
	tests := []struct {
		name    string
		input   *types.PutObjectInput
		wantErr bool
	}{
		{
			name: "correct input",
			input: &types.PutObjectInput{
				Key:         test.Key1,
				Value:       test.KeyValueMap[test.Key1],
				Application: "application",
				Environment: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.PutObject(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("Client.PutObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetObject(t *testing.T) {
	c := &Client{
		svc: &test.MockSSMClient{},
	}
	tests := []struct {
		name    string
		input   *types.GetObjectInput
		want    *types.GetObjectOutput
		wantErr bool
	}{
		{
			name: "correct values",
			input: &types.GetObjectInput{
				Key: test.Key1,
			},
			want: &types.GetObjectOutput{
				Value: test.KeyValueMap[test.Key1],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetObject(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetObjects(t *testing.T) {
	c := &Client{
		svc: &test.MockSSMClient{},
	}
	secrets := make(map[string]string)
	secrets[test.Key1] = test.KeyValueMap[test.Key1]
	secrets[test.Key2] = test.KeyValueMap[test.Key2]

	oneSecret := make(map[string]string)
	oneSecret[test.Key1] = test.KeyValueMap[test.Key1]
	tests := []struct {
		name    string
		input   *types.GetObjectsInput
		want    *types.GetObjectsOuput
		wantErr bool
	}{
		{
			name: "correct values",
			input: &types.GetObjectsInput{
				Prefix: "test/webapp/",
			},
			want:    &types.GetObjectsOuput{secrets},
			wantErr: false,
		},
		{
			name: "correct full path to key",
			input: &types.GetObjectsInput{
				Prefix: test.Key1,
			},
			want:    &types.GetObjectsOuput{oneSecret},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetObjects(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetObjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetObjects() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unSlash(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "correct input",
			args: args{
				input: "/test/webapp/cocpit_api_pass",
			},
			want: "test/webapp/cocpit_api_pass",
		},
		{
			name: "input without slash",
			args: args{
				input: "test/webapp/cocpit_api_pass",
			},
			want: "test/webapp/cocpit_api_pass",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unSlash(tt.args.input); got != tt.want {
				t.Errorf("unSlash() = %v, want %v", got, tt.want)
			}
		})
	}
}
