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
)

type cleanupTTLJob struct {
	db, table string
}

func (cl *cleanupTTLJob) Cleanup(ctx context.Context, o operation.Operation, c cluster.Cluster) {
	return
}

func runTTLJob(ctx context.Context, o operation.Operation, c cluster.Cluster) registry.OperationCleanup {
	conn := c.Conn(ctx, o.L(), 1, option.VirtualClusterName(roachtestflags.VirtualCluster))
	defer conn.Close()

	dbName := pickRandomDB(ctx, o, conn)
	tableName := pickRandomTable(ctx, o, conn, dbName)

	o.Status(fmt.Sprintf("enabling ttl on table %s.%s", dbName, tableName))
	addColStmt := fmt.Sprintf("ALTER TABLE %s.%s WITH (ttl_expire_after = '2 minutes', ttl_job_cron = '*/3 * * * *');", dbName, tableName)
	_, err := conn.ExecContext(ctx, addColStmt)
	if err != nil {
		o.Fatal(err)
		return nil
	}

	o.Status(fmt.Sprintf("ttl created on %s.%s", dbName, tableName))
	return &cleanupAddedColumn{
		db:    dbName,
		table: tableName,
	}
}

func registerRunTTLJob(r registry.Registry) {
	r.AddOperation(registry.OperationSpec{
		Name:             "run-ttl-job",
		Owner:            registry.OwnerSQLFoundations,
		Timeout:          24 * time.Hour,
		CompatibleClouds: registry.AllClouds,
		Dependencies:     []registry.OperationDependency{registry.OperationRequiresPopulatedDatabase},
		Run:              runTTLJob,
	})
}
