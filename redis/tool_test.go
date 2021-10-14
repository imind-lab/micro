package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	redisDB   *redis.Client
	redisMock redismock.ClientMock
}

func (s *Suite) SetupSuite() {
	s.redisDB, s.redisMock = redismock.NewClientMock()
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.redisMock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestSetHashTable() {
	type model struct {
		id   int
		name string
		age  int
	}
	tests := []struct {
		name     string
		key      string
		model    model
		expire   time.Duration
		expected error
	}{
		{"t1", "key01", model{100, "name100", 20}, time.Second, nil},
		{"t2", "key02", model{200, "name200", 20}, time.Second * 5, nil},
	}
	ctx := context.Background()

	for _, t := range tests {
		s.Run(t.name, func() {
			s.redisMock.ExpectHMSet(t.key, FlatStruct(t.model)).SetVal(true)
			s.redisMock.ExpectExpire(t.key, t.expire).SetVal(true)
			err := SetHashTable(ctx, s.redisDB, t.key, t.model, t.expire)
			require.NoError(s.T(), err)
		})
	}
}

func (s *Suite) TestSetSortedSet() {

}
