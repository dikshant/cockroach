// Copyright 2024 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package operations

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/cluster"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/operation"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/option"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/registry"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/roachtestflags"
	"github.com/cockroachdb/cockroach/pkg/util/randutil"
)

type cleanupAddedIndex struct {
	db, table, index string
}

func (cl *cleanupAddedIndex) Cleanup(ctx context.Context, o operation.Operation, c cluster.Cluster) {
	conn := c.Conn(ctx, o.L(), 1, option.TenantName(roachtestflags.VirtualCluster))
	defer conn.Close()

	o.Status(fmt.Sprintf("dropping index %s", cl.index))
	_, err := conn.ExecContext(ctx, fmt.Sprintf("DROP INDEX %s.%s@%s", cl.db, cl.table, cl.index))
	if err != nil {
		o.Fatal(err)
	}
}

func runAddIndex(ctx context.Context, o operation.Operation, c cluster.Cluster) registry.OperationCleanup {
	conn := c.Conn(ctx, o.L(), 1, option.TenantName(roachtestflags.VirtualCluster))
	defer conn.Close()

	rng, _ := randutil.NewPseudoRand()
	dbName := pickRandomDB(ctx, o, conn)
	tableName := pickRandomTable(ctx, o, conn, dbName)
	rows, err := conn.QueryContext(ctx, fmt.Sprintf("SELECT column_name FROM [SHOW COLUMNS FROM %s.%s]", dbName, tableName))
	if err != nil {
		o.Fatal(err)
	}
	rows.Next()
	if !rows.Next() {
		o.Fatalf("not enough columns in table %s.%s", dbName, tableName)
	}
	var colName string
	if err := rows.Scan(&colName); err != nil {
		o.Fatal(err)
	}

	indexName := fmt.Sprintf("add_index_op_%d", rng.Uint32())
	o.Status(fmt.Sprintf("adding index to column %s", colName))
	createIndexStmt := fmt.Sprintf("CREATE INDEX %s ON %s.%s (%s)", indexName, dbName, tableName, colName)
	_, err = conn.ExecContext(ctx, createIndexStmt)
	if err != nil {
		o.Fatal(err)
		return nil
	}

	o.Status(fmt.Sprintf("index %s created", indexName))
	return &cleanupAddedIndex{
		db:    dbName,
		table: tableName,
		index: indexName,
	}
}

func registerAddIndex(r registry.Registry) {
	r.AddOperation(registry.OperationSpec{
		Name:             "add-index",
		Owner:            registry.OwnerSQLFoundations,
		Timeout:          24 * time.Hour,
		CompatibleClouds: registry.AllClouds,
		Dependencies:     []registry.OperationDependency{registry.OperationRequiresPopulatedDatabase},
		Run:              runAddIndex,
	})
}
