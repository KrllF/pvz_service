//go:build unit

package httph

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KrllF/pvz_service/internal/handler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProcessClient_ValidateRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args ProcessOrderRequest
		want error
	}{
		{
			name: "all good",
			args: ProcessOrderRequest{
				UserID:    1,
				Operation: "issue",
				OrderIDs:  []int64{1, 2},
			},
			want: nil,
		},
		{
			name: "bad user_id",
			args: ProcessOrderRequest{
				UserID:    -1,
				Operation: "issue",
				OrderIDs:  []int64{1, 2},
			},
			want: errors.New("некорректный user ID"),
		},
		{
			name: "bad operation",
			args: ProcessOrderRequest{
				UserID:    1,
				Operation: "",
				OrderIDs:  []int64{1, 2},
			},
			want: errors.New("некорректная операция"),
		},
		{
			name: "bad order_ids",
			args: ProcessOrderRequest{
				UserID:    1,
				Operation: "return",
				OrderIDs:  []int64{},
			},
			want: errors.New("пустой список order IDs"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateRequest(tt.args)
			if tt.want != nil {
				assert.EqualError(t, err, tt.want.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProcessClient(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    ProcessOrderRequest
		mockAcc func(ctrl *gomock.Controller) *mocks.MockAcceptService
		mockRet func(ctrl *gomock.Controller) *mocks.MockReturnsService
		want    int
	}{
		{
			name: "issue all good",
			args: ProcessOrderRequest{
				UserID:    1,
				Operation: "issue",
				OrderIDs:  []int64{1, 2, 3},
			},
			mockAcc: func(ctrl *gomock.Controller) *mocks.MockAcceptService {
				mockAcc := mocks.NewMockAcceptService(ctrl)
				mockAcc.EXPECT().ProcessClientIssue(context.Background(), int64(1), []int64{1, 2, 3}).Return([]int64{3}, nil)
				return mockAcc
			},
			mockRet: func(ctrl *gomock.Controller) *mocks.MockReturnsService { return nil },
			want:    http.StatusOK,
		},

		{
			name: "return all good",
			args: ProcessOrderRequest{
				UserID:    1,
				Operation: "return",
				OrderIDs:  []int64{1, 2, 3},
			},
			mockAcc: func(ctrl *gomock.Controller) *mocks.MockAcceptService { return nil },
			mockRet: func(ctrl *gomock.Controller) *mocks.MockReturnsService {
				mockRet := mocks.NewMockReturnsService(ctrl)
				mockRet.EXPECT().ProcessClientReturn(context.Background(), int64(1), []int64{1, 2, 3}).Return([]int64{3}, nil)
				return mockRet
			},
			want: http.StatusOK,
		},

		{
			name: "nil operation",
			args: ProcessOrderRequest{
				UserID:    1,
				Operation: "",
				OrderIDs:  []int64{1, 2, 3},
			},
			mockAcc: func(ctrl *gomock.Controller) *mocks.MockAcceptService { return nil },
			mockRet: func(ctrl *gomock.Controller) *mocks.MockReturnsService { return nil },
			want:    http.StatusBadRequest,
		},

		{
			name: "bad operation",
			args: ProcessOrderRequest{
				UserID:    1,
				Operation: "notcorrect",
				OrderIDs:  []int64{1, 2, 3},
			},
			mockAcc: func(ctrl *gomock.Controller) *mocks.MockAcceptService { return nil },
			mockRet: func(ctrl *gomock.Controller) *mocks.MockReturnsService { return nil },
			want:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockAcc := tt.mockAcc(ctrl)
			mockRet := tt.mockRet(ctrl)
			serv := NewHandler(context.Background(),
				mockAcc, mockRet, nil, nil)

			reqBody, err := json.Marshal(tt.args)
			assert.NoError(t, err)
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/process", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			serv.ProcessClient(res, req)
			assert.Equal(t, tt.want, res.Code)
		})
	}
}
