//go:build unit

package accept

import (
	"testing"

	"github.com/KrllF/pvz_service/internal/service/order/accept/mocks"
	"github.com/KrllF/pvz_service/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	t.Run("all good", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mocks.NewMockRepository(ctrl)
		cache := mocks.NewMockCache(ctrl)
		tx := mocks.NewMockTXManager(ctrl)
		logg, _ := logger.NewLogger(zapcore.DebugLevel)
		serv := NewService(repo, cache, tx, logg)
		assert.Equal(t, repo, serv.repository)
		assert.Equal(t, cache, serv.cache)
		assert.Equal(t, tx, serv.txManager)
	})
}
