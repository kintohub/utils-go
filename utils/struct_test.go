package utils_test

import (
	"encoding/json"
	"github.com/kintohub/utils-go/utils"
	"testing"
)

func TestDeepEqual(t *testing.T) {

	type ConfigData struct {
		IsPublicApi bool `json:"isPublicApi"`
		EnvVars     []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"envVars"`
	}
	type KintoBlock struct {
		ID         *string    `json:"kintoblock_id"`
		ConfigJson ConfigData `json:"config_data"`
	}

	const TestJson = `{"kintoblock_id":"7927b549-798b-4c6c-b193-66f31e402044","config_data":{"isPublicApi":false,"envVars":[{"key":"GRPC_API_HOST","value":"omniscient-gateway.kintohub.svc.cluster.local","required":false},{"key":"GRPC_API_PORT","value":"8090","required":false},{"key":"GRPC_SERVER_PORT","value":"80","required":false},{"key":"GITHUB_CLIENT_ID","value":"Iv1.906bb07361b58772","required":false},{"key":"GITHUB_CLIENT_SECRET","value":"847aaba665f743fa28e3ed556b1a287bd3222712","required":false},{"key":"GITHUB_GITAPP_ID","value":"28062","required":false}],"replicas":0,"hardwareData":{"memory":128}},"kintoblock_build":{"build_id":"b1bc6361-a8e9-441f-8dc1-fc7cbcc191b3","subtype":"","protocol":"GRPC","port":80,"kintobranch":{"git_branch":{"name":"dev"}}},"KintoBlock":{"workspace_id":"d87be4a1-094a-4b38-bf83-b7a686121812","name":"github","type":"MICROSERVICE","env_vars":[{"key":"GITHUB_CLIENT_ID","value":"Iv1.906bb07361b58772","required":true},{"key":"GITHUB_CLIENT_SECRET","value":"847aaba665f743fa28e3ed556b1a287bd3222712","required":false},{"key":"GITHUB_GITAPP_ID","value":"28062","required":false},{"key":"GRPC_API_HOST","value":"omniscient-gateway.kintohub.svc.cluster.local","required":false},{"key":"GRPC_API_PORT","value":"8090","required":false},{"key":"GRPC_SERVER_PORT","value":"80","required":false}]}}`
	const TestJsonShuffled = `{"kintoblock_id":"7927b549-798b-4c6c-b193-66f31e402044","config_data":{"isPublicApi":false,"envVars":[{"key":"GITHUB_CLIENT_ID","value":"Iv1.906bb07361b58772","required":false},{"key":"GITHUB_CLIENT_SECRET","value":"847aaba665f743fa28e3ed556b1a287bd3222712","required":false},{"key":"GITHUB_GITAPP_ID","value":"28062","required":false},{"key":"GRPC_API_HOST","value":"omniscient-gateway.kintohub.svc.cluster.local","required":false},{"key":"GRPC_API_PORT","value":"8090","required":false},{"key":"GRPC_SERVER_PORT","value":"80","required":false}],"replicas":0,"hardwareData":{"memory":128}},"kintoblock_build":{"build_id":"b1bc6361-a8e9-441f-8dc1-fc7cbcc191b3","subtype":"","protocol":"GRPC","port":80,"kintobranch":{"git_branch":{"name":"dev"}}},"KintoBlock":{"workspace_id":"d87be4a1-094a-4b38-bf83-b7a686121812","name":"github","type":"MICROSERVICE","env_vars":[{"key":"GITHUB_CLIENT_ID","value":"Iv1.906bb07361b58772","required":true},{"key":"GITHUB_CLIENT_SECRET","value":"847aaba665f743fa28e3ed556b1a287bd3222712","required":false},{"key":"GITHUB_GITAPP_ID","value":"28062","required":false},{"key":"GRPC_API_HOST","value":"omniscient-gateway.kintohub.svc.cluster.local","required":false},{"key":"GRPC_API_PORT","value":"8090","required":false},{"key":"GRPC_SERVER_PORT","value":"80","required":false}]}}`
	const TestJsonSame = `{"kintoblock_id":"7927b549-798b-4c6c-b193-66f31e402044","config_data":{"isPublicApi":false,"envVars":[{"key":"GRPC_API_HOST","value":"omniscient-gateway.kintohub.svc.cluster.local","required":false},{"key":"GRPC_API_PORT","value":"8090","required":false},{"key":"GRPC_SERVER_PORT","value":"80","required":false},{"key":"GITHUB_CLIENT_ID","value":"Iv1.906bb07361b58772","required":false},{"key":"GITHUB_CLIENT_SECRET","value":"847aaba665f743fa28e3ed556b1a287bd3222712","required":false},{"key":"GITHUB_GITAPP_ID","value":"28062","required":false}],"replicas":0,"hardwareData":{"memory":128}},"kintoblock_build":{"build_id":"b1bc6361-a8e9-441f-8dc1-fc7cbcc191b3","subtype":"","protocol":"GRPC","port":80,"kintobranch":{"git_branch":{"name":"dev"}}},"KintoBlock":{"workspace_id":"d87be4a1-094a-4b38-bf83-b7a686121812","name":"github","type":"MICROSERVICE","env_vars":[{"key":"GITHUB_CLIENT_ID","value":"Iv1.906bb07361b58772","required":true},{"key":"GITHUB_CLIENT_SECRET","value":"847aaba665f743fa28e3ed556b1a287bd3222712","required":false},{"key":"GITHUB_GITAPP_ID","value":"28062","required":false},{"key":"GRPC_API_HOST","value":"omniscient-gateway.kintohub.svc.cluster.local","required":false},{"key":"GRPC_API_PORT","value":"8090","required":false},{"key":"GRPC_SERVER_PORT","value":"80","required":false}]}}`
	const TestJsonNoID = `{"kintoblock_id":null,"config_data":{"isPublicApi":false,"envVars":[{"key":"GITHUB_CLIENT_ID","value":"Iv1.906bb07361b58772","required":false},{"key":"GITHUB_CLIENT_SECRET","value":"847aaba665f743fa28e3ed556b1a287bd3222712","required":false},{"key":"GITHUB_GITAPP_ID","value":"28062","required":false},{"key":"GRPC_API_HOST","value":"omniscient-gateway.kintohub.svc.cluster.local","required":false},{"key":"GRPC_API_PORT","value":"8090","required":false},{"key":"GRPC_SERVER_PORT","value":"80","required":false}],"replicas":0,"hardwareData":{"memory":128}},"kintoblock_build":{"build_id":"b1bc6361-a8e9-441f-8dc1-fc7cbcc191b3","subtype":"","protocol":"GRPC","port":80,"kintobranch":{"git_branch":{"name":"dev"}}},"KintoBlock":{"workspace_id":"d87be4a1-094a-4b38-bf83-b7a686121812","name":"github","type":"MICROSERVICE","env_vars":[{"key":"GITHUB_CLIENT_ID","value":"Iv1.906bb07361b58772","required":true},{"key":"GITHUB_CLIENT_SECRET","value":"847aaba665f743fa28e3ed556b1a287bd3222712","required":false},{"key":"GITHUB_GITAPP_ID","value":"28062","required":false},{"key":"GRPC_API_HOST","value":"omniscient-gateway.kintohub.svc.cluster.local","required":false},{"key":"GRPC_API_PORT","value":"8090","required":false},{"key":"GRPC_SERVER_PORT","value":"80","required":false}]}}`
	const TestJsonEmpty = `{}`
	const TestJsonNull = `{"kintoblock_id": null}`
	block := KintoBlock{}
	blockShuffled := KintoBlock{}
	blockSame := KintoBlock{}
	blockNoID := KintoBlock{}
	blockEmpty := KintoBlock{}
	blockNull := KintoBlock{}
	blockNull2 := KintoBlock{}

	json.Unmarshal([]byte(TestJson), &block)
	json.Unmarshal([]byte(TestJsonShuffled), &blockShuffled)
	json.Unmarshal([]byte(TestJsonSame), &blockSame)
	json.Unmarshal([]byte(TestJsonNoID), &blockNoID)
	json.Unmarshal([]byte(TestJsonEmpty), &blockEmpty)
	json.Unmarshal([]byte(TestJsonNull), &blockNull)
	json.Unmarshal([]byte(TestJsonNull), &blockNull2)

	type args struct {
		vx interface{}
		vy interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "not in order",
			args: args{
				vx: block,
				vy: blockShuffled,
			},
			want: true,
		},
		{
			name: "same",
			args: args{
				vx: block,
				vy: blockSame,
			},
			want: true,
		},
		{
			name: "not in order",
			args: args{
				vx: blockShuffled,
				vy: blockSame,
			},
			want: true,
		},
		{
			name: "no id",
			args: args{
				vx: block,
				vy: blockNoID,
			},
			want: false,
		},
		{
			name: "empty",
			args: args{
				vx: block,
				vy: blockEmpty,
			},
			want: false,
		},
		{
			name: "empty",
			args: args{
				vx: blockNull,
				vy: blockNull2,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.DeepEqualStruct(tt.args.vx, tt.args.vy); got != tt.want {
				t.Errorf("DeepEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
