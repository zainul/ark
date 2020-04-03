package threshold

import (
	"fmt"
	"testing"

	"github.com/zainul/ark/storage/redis"
	"github.com/zainul/ark/storage/redis/dummyrds"
)

func Test_threshold_Attempt(t1 *testing.T) {
	type fields struct {
		timeToLive int
		max        int
		rds        redis.Redis
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Attempt Success",
			fields: fields{
				timeToLive: 10,
				max:        10,
				rds: dummyrds.New(dummyrds.Config{
					MockingMap: dummyrds.Mocker{},
				}),
			},
			args: args{
				key: "test:somekey",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &threshold{
				timeToLive: tt.fields.timeToLive,
				max:        tt.fields.max,
				rds:        tt.fields.rds,
			}
			if err := t.Attempt(tt.args.key); (err != nil) != tt.wantErr {
				t1.Errorf("Attempt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_threshold_IsAllow(t1 *testing.T) {
	type fields struct {
		timeToLive int
		max        int
		holdPrefix string
		rds        redis.Redis
	}
	type args struct {
		key         string
		failAction1 func()
		failAction2 func()
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{
			name: "Is Allow Success",
			fields: fields{
				timeToLive: 10,
				max:        5,
				rds: dummyrds.New(dummyrds.Config{
					MockingMap: dummyrds.Mocker{},
				}),
			},
			args: args{
				key: "somekey:test",
				failAction1: func() {
					fmt.Println("hail fail")
				},
				failAction2: func() {
					fmt.Println("hai slack")
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &threshold{
				timeToLive: tt.fields.timeToLive,
				max:        tt.fields.max,
				rds:        tt.fields.rds,
			}
			if got := t.IsAllow(tt.args.key, tt.args.failAction1, tt.args.failAction2); got != tt.want {
				t1.Errorf("IsAllow() = %v, want %v", got, tt.want)
			}
		})
	}
}
