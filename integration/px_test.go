package integration_test

import (
	"context"
	"fmt"
	// "testing"
	"time"

	. "github.com/pingcap/check"
	"github.com/pingcap/tidb/domain"
	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tidb/session"
	"github.com/pingcap/tidb/sessionctx"
	"github.com/pingcap/tidb/store/mockstore"
	"github.com/pingcap/tidb/store/mockstore/mocktikv"
	"github.com/pingcap/tidb/util/testkit"
)

var _ = Suite(&testIntegrationSuite{})

// func Test(t *testing.T) {
// 	TestingT(t)
// }

type testIntegrationSuite struct {
	lease     time.Duration
	cluster   *mocktikv.Cluster
	mvccStore mocktikv.MVCCStore
	store     kv.Storage
	dom       *domain.Domain
	ctx       sessionctx.Context
	tk        *testkit.TestKit
}

func (s *testIntegrationSuite) TearDownTest(c *C) {
	tk := testkit.NewTestKit(c, s.store)
	tk.MustExec("use test")
	r := tk.MustQuery("show tables")
	for _, tb := range r.Rows() {
		tableName := tb[0]
		tk.MustExec(fmt.Sprintf("drop table %v", tableName))
	}
}

func (s *testIntegrationSuite) SetUpSuite(c *C) {
	var err error
	s.lease = 50 * time.Millisecond

	s.cluster = mocktikv.NewCluster()
	mocktikv.BootstrapWithMultiStores(s.cluster, 2)
	s.mvccStore = mocktikv.MustNewMVCCStore()
	s.store, err = mockstore.NewMockTikvStore(
		mockstore.WithCluster(s.cluster),
		mockstore.WithMVCCStore(s.mvccStore),
	)
	c.Assert(err, IsNil)
	session.SetSchemaLease(s.lease)
	session.SetStatsLease(0)
	s.dom, err = session.BootstrapSession(s.store)
	c.Assert(err, IsNil)

	se, err := session.CreateSession4Test(s.store)
	c.Assert(err, IsNil)
	s.ctx = se.(sessionctx.Context)
	_, err = se.Execute(context.Background(), "create database test_db")
	c.Assert(err, IsNil)
	s.tk = testkit.NewTestKit(c, s.store)
}

func (s *testIntegrationSuite) TestGiveupLeaderAndRecover(c *C) {
	tk := testkit.NewTestKit(c, s.store)
	tk.MustExec("USE test")
	tk.MustExec("create table x (id int primary key, c int);")
	tk.MustExec("insert into x values(1, 1);")
	region := s.cluster.GetAllRegions()[0]
	regionID := region.Meta.GetId()
	_, leaderID := s.cluster.GetRegion(regionID)
	s.cluster.GiveUpLeader(regionID)
	// wait 20s
	err := tk.QueryToErr("select count(*) from x")
	c.Assert(err, NotNil)
	s.cluster.ChangeLeader(regionID, leaderID)
	tk.MustQuery("select count(*) from x").Check(testkit.Rows("1"))
}
