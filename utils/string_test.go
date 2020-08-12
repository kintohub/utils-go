package utils_test

import (
	"github.com/kintohub/utils-go/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShortenUUID16(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test success",
			args: args{
				uuid: "91837cf5-ac35-4b6c-a67a-0845eabb3091",
			},
			want:    "37f974b0468e7bfd",
			wantErr: false,
		},
		{
			name: "test same output",
			args: args{
				uuid: "a67a0845-eabb-3091-9183-7cf5ac354b6c",
			},
			want:    "37f974b0468e7bfd",
			wantErr: false,
		},
		{
			name: "test same output",
			args: args{
				uuid: "a67a0845-eabb-3090-9183-7cf5ac354b6d",
			},
			want:    "37f974b0468e7bfd",
			wantErr: false,
		},
		{
			name: "test failure",
			args: args{
				uuid: "91837cf5-ac35-4b6c-a67a-0845eabb30911",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ShortenUUID16(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShortenUUID16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ShortenUUID16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortenUUID8(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test success",
			args: args{
				uuid: "91837cf5-ac35-4b6c-a67a-0845eabb3091",
			},
			want:    "71770f4d",
			wantErr: false,
		},
		{
			name: "test failure",
			args: args{
				uuid: "91837cf5-ac35-4b6c-a67a-0845eabb30911",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ShortenUUID8(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShortenUUID8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ShortenUUID8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortenHexString(t *testing.T) {
	type args struct {
		input     string
		outputLen int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test success",
			args: args{
				input:     "5f28eb526ffc3d5051ad2fd5",
				outputLen: 12,
			},
			want:    "6278baff4029",
			wantErr: false,
		},
		{
			name: "test success",
			args: args{
				input:     "5f28eb526ffc3d5051ad2fd5",
				outputLen: 2,
			},
			want:    "36",
			wantErr: false,
		},
		{
			name: "test success",
			args: args{
				input:     "5f28eb526ffc3d5051ad2fd5",
				outputLen: 6,
			},
			want:    "9d3893",
			wantErr: false,
		},
		{
			name: "test error if outputLen not dividend",
			args: args{
				input:     "5f28eb526ffc3d5051ad2fd5",
				outputLen: 14,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test error if outputLen not even",
			args: args{
				input:     "5f28eb526ffc3d5051ad2fd5",
				outputLen: 3,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ShortenHexString(tt.args.input, tt.args.outputLen)

			if (err != nil) != tt.wantErr {
				assert.Equal(t, len(got), tt.args.outputLen)
				t.Errorf("ShortenHexString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ShortenHexString() = %v, want %v", got, tt.want)
			}
		})
	}
}
