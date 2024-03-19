package main

import (
	"context"
	"testing"

	"sample-manager/constants"
	pb "sample-manager/proto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreatingAMapping(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error creating mock db: %v", err)

	defer db.Close()

	dialect := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDb, err := gorm.Open(dialect, &gorm.Config{})
	assert.Nil(t, err, "Error creating mock gorm db: %v", err)

	type args struct {
		ctx context.Context
		req *pb.CreateRequest
	}
	tests := []struct {
		name    string
		args    args
		rows    func()
		want    *pb.CreateResponse
		wantErr bool
		errorCode codes.Code
	}{
		{
			name: "Creating a mapping - Expect Success",
			args: args{
				ctx: context.Background(),
				req: &pb.CreateRequest{
					Segments: []string{"test_segment1", "test_segment2"},
					ItemId: "test_item_id",
					SampleItemId: "test_sample_item_id",
				},
			},
			rows: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{
					"id",
					"segments",
					"item_id",
					"sample_item_id",
				}).AddRow(1, pq.StringArray{"test_segment1", "test_segment2"}, "test_item_id", "test_sample_item_id")
				mock.ExpectQuery("INSERT").WillReturnRows(rows)
				mock.ExpectCommit()
			},
			want: &pb.CreateResponse{
				Message: constants.CREATE_MAPPING_SUCCESS,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rows()
			server := &Server{DB: gormDb} 

			got, err := server.CreateMapping(tt.args.ctx, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateMapping() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				statusErr, ok := status.FromError(err)
				assert.True(t, ok, "Expected gRPC status error")
				assert.Equalf(t, tt.errorCode, statusErr.Code(), "Expected %v error", tt.errorCode)
			} else {
				assert.Equalf(t, tt.want, got, "CreateMapping(%v, %v)", tt.args.ctx, tt.args.req)
			}
		})
	}
}

func TestGettingASampleID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error creating mock db: %v", err)

	defer db.Close()

	dialect := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDb, err := gorm.Open(dialect, &gorm.Config{})
	assert.Nil(t, err, "Error creating mock gorm db: %v", err)

	type args struct {
		ctx context.Context
		req *pb.GetRequest
	}
	tests := []struct {
		name    string
		args    args
		rows    func()
		want    *pb.GetResponse
		wantErr bool
		errorCode codes.Code
	}{
		{
			name: "Getting a sample ID - Expect Success",
			args: args{
				ctx: context.Background(),
				req: &pb.GetRequest{
					Clm: []string{"test_segment1", "test_segment2"},
					ItemId: "test_item_id",
				},
			},
			rows: func() {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "sample_item_id", "item_id", "segments"}).
					AddRow(1, "1", "test_item_id", pq.StringArray{"test_segment1", "test_segment2"}))
			},
			want: &pb.GetResponse{
				SampleItemId: "1",
			},
			wantErr: false,
		},
		{
			name: "Getting a sample ID when mapping is doesn't exist - Expect Error",
			args: args{
				ctx: context.Background(),
				req: &pb.GetRequest{
					Clm: []string{"test_segment1", "test_segment2"},
					ItemId: "test_item_id",
				},
			},
			rows: func() {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{}))
			},
			want: nil,
			wantErr: true,
			errorCode: codes.Unavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rows()
			server := &Server{DB: gormDb} 

			got, err := server.GetSampleId(tt.args.ctx, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetSampleId() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				statusErr, ok := status.FromError(err)
				assert.True(t, ok, "Expected gRPC status error")
				assert.Equalf(t, tt.errorCode, statusErr.Code(), "Expected %v error", tt.errorCode)
			} else {
				assert.Equalf(t, tt.want, got, "GetSampleId(%v, %v)", tt.args.ctx, tt.args.req)
			}
		})
	}
}