package opentick

import (
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AdjCache(t *testing.T) {
	fdb.MustAPIVersion(FdbVersion)
	var db = fdb.MustOpenDefault()
	DropDatabase(db, "test")
	CreateDatabase(db, "test")
	_, err := Execute(db, "test", "insert into _adj_ values(1, 1, 0.25, 4)", nil)
	assert.Equal(t, nil, err)
	_, err = Execute(db, "test", "insert into _adj_ values(1, 3, 0.5, 2)", nil)
	_, err = Execute(db, "test", "insert into _adj_ values(1, 5, 0.2, 5)", nil)
	x := adjCache.get(db, "test", 1)
	assert.Equal(t, "[{1 0.025 40 4 0.25} {3 0.1 10 8 0.125} {5 0.2 5 40 0.025}]", fmt.Sprint(x))
	_, err = Execute(db, "test", "create table bar(a int, b timestamp, c double, d double, vol double, primary key(a, b))", nil)
	assert.Equal(t, nil, err)
	_, err = Execute(db, "test", "insert into bar values(1, 100, 1, 1, 1)", nil)
	assert.Equal(t, nil, err)
	_, err = Execute(db, "test", "insert into bar values(1, 99, 1.5, 1.5, 1.5)", nil)
	assert.Equal(t, nil, err)
	_, err = Execute(db, "test", "insert into bar values(1, 5, 1, 1, 1)", nil)
	assert.Equal(t, nil, err)
	_, err = Execute(db, "test", "insert into bar values(1, 4, 1, 1, 1)", nil)
	assert.Equal(t, nil, err)
	_, err = Execute(db, "test", "insert into bar values(1, 3, 1, 1, 1)", nil)
	assert.Equal(t, nil, err)
	_, err = Execute(db, "test", "insert into bar values(1, 2, 1, 1, 1)", nil)
	assert.Equal(t, nil, err)
	_, err = Execute(db, "test", "insert into bar values(1, 0, 1, 1, 1)", nil)
	assert.Equal(t, nil, err)
	ret, _ := Execute(db, "test", "select b, adj(c), adj(d), adj(vol) from bar where a=1 and b=100", nil)
	assert.Equal(t, "[[[100 0] 1 1 1]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c), adj(d), adj(vol) from bar where a=1 and b=5", nil)
	assert.Equal(t, "[[[5 0] 1 1 1]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c), adj(d), adj(vol) from bar where a=1 and b=4", nil)
	assert.Equal(t, "[[[4 0] 0.2 0.2 5]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c), adj(d), adj(vol) from bar where a=1 and b=3", nil)
	assert.Equal(t, "[[[3 0] 0.2 0.2 5]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c), adj(d), adj(vol) from bar where a=1 and b=2", nil)
	assert.Equal(t, "[[[2 0] 0.1 0.1 10]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c), adj(d), adj(vol) from bar where a=1 and b=0", nil)
	assert.Equal(t, "[[[0 0] 0.025 0.025 40]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c), adj(d), adj(vol) from bar where a=1", nil)
	assert.Equal(t, "[[[0 0] 0.025 0.025 40] [[2 0] 0.1 0.1 10] [[3 0] 0.2 0.2 5] [[4 0] 0.2 0.2 5] [[5 0] 1 1 1] [[99 0] 1.5 1.5 1.5] [[100 0] 1 1 1]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c), adj(d), adj(vol) from bar where a=1 limit -10", nil)
	assert.Equal(t, "[[[100 0] 1 1 1] [[99 0] 1.5 1.5 1.5] [[5 0] 1 1 1] [[4 0] 0.2 0.2 5] [[3 0] 0.2 0.2 5] [[2 0] 0.1 0.1 10] [[0 0] 0.025 0.025 40]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c), adj(d), adj(vol) from bar where a=1 limit -1", nil)
	assert.Equal(t, "[[[100 0] 1 1 1]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c), adj(d), adj(vol) from bar where a=1 limit 1", nil)
	assert.Equal(t, "[[[0 0] 0.025 0.025 40]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c, true), adj(d, true), adj(vol, true) from bar where a=1 and b=0", nil)
	assert.Equal(t, "[[[0 0] 1 1 1]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c, true), adj(d, true), adj(vol, true) from bar where a=1 and b=2", nil)
	assert.Equal(t, "[[[2 0] 4 4 0.25]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c, true), adj(d, true), adj(vol, true) from bar where a=1 and b=3", nil)
	assert.Equal(t, "[[[3 0] 8 8 0.125]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c, true), adj(d, true), adj(vol, true) from bar where a=1 and b=4", nil)
	assert.Equal(t, "[[[4 0] 8 8 0.125]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c, true), adj(d, true), adj(vol, true) from bar where a=1 and b=5", nil)
	assert.Equal(t, "[[[5 0] 40 40 0.025]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c, true), adj(d, true), adj(vol, true) from bar where a=1 and b=100", nil)
	assert.Equal(t, "[[[100 0] 40 40 0.025]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c, true), adj(d, true), adj(vol, true) from bar where a=1", nil)
	assert.Equal(t, "[[[0 0] 1 1 1] [[2 0] 4 4 0.25] [[3 0] 8 8 0.125] [[4 0] 8 8 0.125] [[5 0] 40 40 0.025] [[99 0] 60 60 0.037500000000000006] [[100 0] 40 40 0.025]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c, true), adj(d, true), adj(vol, true) from bar where a=1 limit -10", nil)
	assert.Equal(t, "[[[100 0] 40 40 0.025] [[99 0] 60 60 0.037500000000000006] [[5 0] 40 40 0.025] [[4 0] 8 8 0.125] [[3 0] 8 8 0.125] [[2 0] 4 4 0.25] [[0 0] 1 1 1]]", fmt.Sprint(ret))
	ret, _ = Execute(db, "test", "select b, adj(c), adj(d), adj(vol) from bar where a=1 limit -1", nil)
	ret, err = Execute(db, "test", "select b, adj(c), adj(d, true), adj(vol, true) from bar where a=1", nil)
	assert.NotEqual(t, nil, err)
	Execute(db, "", "drop table test.test", nil)
}
