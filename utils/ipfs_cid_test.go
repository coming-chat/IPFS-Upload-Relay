package utils

import (
	"os"
	"testing"
)

func TestGetIPFSCid(t *testing.T) {
	file, err := os.ReadFile("")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		file []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				file: file,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIPFSCid(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIPFSCid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}
