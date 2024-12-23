package handlers_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ole-larsen/green-api/internal/httpserver/handlers"
	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			name: "test err",
			args: args{
				err: nil,
			},
			wantErr: true,
		},
		{
			name: "test error1",
			args: args{
				err: errors.New("some error"),
			},
			wantErr: true,
		},
		{
			name: "test  error2",
			args: args{
				err: errors.New("some error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handlers.NewError(tt.args.err)
			if tt.args.err == nil {
				require.Nil(t, err)
			} else {
				if (err != nil) != tt.wantErr {
					t.Errorf("NewError() error = %v, wantErr %v", err, tt.wantErr)
				}

				require.Equal(t, "*handlers.Error", fmt.Sprintf("%T", err))

				require.Equal(t, fmt.Sprintf("[handlers]: %v", tt.args.err), err.Error())
			}
		})
	}
}
