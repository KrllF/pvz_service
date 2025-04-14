//go:build unit

package httph

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/handler/mocks"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestListOrdersUser_pagRet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    url.Values
		want    bool
		wantErr error
	}{
		{
			name: "good pag true",
			args: url.Values{
				"pag": []string{"true"},
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "good pag false",
			args: url.Values{
				"pag": []string{"false"},
			},
			want:    false,
			wantErr: nil,
		},
		{
			name:    "without pag",
			args:    url.Values{},
			want:    false,
			wantErr: nil,
		},
		{
			name: "bad pag",
			args: url.Values{
				"pag": []string{"not_correct"},
			},
			want:    false,
			wantErr: errors.New("некорректный параметр pag"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			have, err := pagRet(tt.args)
			assert.Equal(t, tt.want, have)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListOrdersUser_inPVZRet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    url.Values
		want    bool
		wantErr error
	}{
		{
			name: "good in_pvz true",
			args: url.Values{
				"in_pvz": []string{"true"},
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "good in_pvz false",
			args: url.Values{
				"in_pvz": []string{"false"},
			},
			want:    false,
			wantErr: nil,
		},
		{
			name:    "without in_pvz",
			args:    url.Values{},
			want:    false,
			wantErr: nil,
		},
		{
			name: "bad in_pvz",
			args: url.Values{
				"in_pvz": []string{"not_correct"},
			},
			want:    false,
			wantErr: errors.New("некорректный параметр in_pvz"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			have, err := inPVZRet(tt.args)
			assert.Equal(t, tt.want, have)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListOrdersUser_limitRet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    url.Values
		want    int64
		wantErr error
	}{
		{
			name: "good limit",
			args: url.Values{
				"limit": []string{"1"},
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "negative limit",
			args: url.Values{
				"limit": []string{"-2"},
			},
			want:    -2,
			wantErr: errors.New("отрицательный параметр limit"),
		},
		{
			name:    "without limit",
			args:    url.Values{},
			want:    -1,
			wantErr: nil,
		},
		{
			name: "bad limit",
			args: url.Values{
				"limit": []string{"not_correct"},
			},
			want:    -1,
			wantErr: errors.New("некорректный параметр limit"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			have, err := limitRet(tt.args)
			assert.Equal(t, tt.want, have)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListOrdersUser_pageRet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    url.Values
		want    int64
		wantErr error
	}{
		{
			name: "good page",
			args: url.Values{
				"page": []string{"2"},
			},
			want:    2,
			wantErr: nil,
		},
		{
			name: "negative page",
			args: url.Values{
				"page": []string{"-2"},
			},
			want:    -2,
			wantErr: errors.New("отрицательный параметр page"),
		},
		{
			name:    "without page",
			args:    url.Values{},
			want:    -1,
			wantErr: nil,
		},
		{
			name: "bad page",
			args: url.Values{
				"page": []string{"not_correct"},
			},
			want:    -1,
			wantErr: errors.New("некорректный параметр page"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			have, err := pageRet(tt.args)
			assert.Equal(t, tt.want, have)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListOrdersUser_nRet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    url.Values
		want    int64
		wantErr error
	}{
		{
			name: "good n",
			args: url.Values{
				"n": []string{"2"},
			},
			want:    2,
			wantErr: nil,
		},
		{
			name: "negative n",
			args: url.Values{
				"n": []string{"-2"},
			},
			want:    -2,
			wantErr: errors.New("отрицательный параметр n"),
		},
		{
			name:    "without n",
			args:    url.Values{},
			want:    -1,
			wantErr: nil,
		},
		{
			name: "bad n",
			args: url.Values{
				"n": []string{"not_correct"},
			},
			want:    -1,
			wantErr: errors.New("некорректный параметр n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			have, err := nRet(tt.args)
			assert.Equal(t, tt.want, have)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListOrdersUser_validateQueryParamsU(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    url.Values
		want    ValidateStructToListOrders
		wantErr error
	}{
		{
			name: "all good",
			args: url.Values{
				"pag":    []string{"true"},
				"limit":  []string{"1"},
				"page":   []string{"1"},
				"n":      []string{"1"},
				"in_pvz": []string{"true"},
			},
			want: ValidateStructToListOrders{
				n:     1,
				inPVZ: true,
				limit: 1,
				page:  1,
				pag:   true,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ret, err := validateQueryParamsU(tt.args)
			assert.Equal(t, tt.want, ret)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListOrdersUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		url      string
		mock     func(ctrl *gomock.Controller) *mocks.MockShowService
		wantCode int
	}{
		{
			name: "all good",
			url:  "/history/orders/1?in_pvz=true&n=1&limit=1&page=1&pag=true",
			mock: func(ctrl *gomock.Controller) *mocks.MockShowService {
				mockShow := mocks.NewMockShowService(ctrl)
				// тут не убрал, потому что gorilla/mux обрабатывает query
				mockShow.EXPECT().ListOrdersUserPVZ(gomock.Any(),
					int64(1), int64(1), int64(1), true, false, "").Return([]models.Order{
					{
						OrderID:       1,
						UserID:        1,
						Weight:        1,
						Price:         1,
						Packaging:     entity.Pack{PackType: "box", ExtraPack: ""},
						InStorageFrom: time.Time{},
						ShelfLife:     time.Time{},
						Status:        "accepted",
						TwoDaysOfLife: time.Time{},
						LastUpdate:    time.Time{},
					},
				}, nil)
				return mockShow
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockShow := tt.mock(ctrl)

			router := mux.NewRouter()
			serv := NewHandler(context.Background(),
				nil, nil, mockShow, nil)
			router.HandleFunc("/history/orders/{id:[0-9]+}", serv.ListOrdersUser).Methods(http.MethodGet)

			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet,
				tt.url, nil)
			router.ServeHTTP(res, req)
			assert.Equal(t, tt.wantCode, res.Code)
		})
	}
}
