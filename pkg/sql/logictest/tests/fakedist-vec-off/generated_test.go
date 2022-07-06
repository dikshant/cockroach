// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

// Code generated by generate-logictest, DO NOT EDIT.

package testfakedist_vec_off

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/build/bazel"
	"github.com/cockroachdb/cockroach/pkg/security/securityassets"
	"github.com/cockroachdb/cockroach/pkg/security/securitytest"
	"github.com/cockroachdb/cockroach/pkg/server"
	"github.com/cockroachdb/cockroach/pkg/sql/logictest"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/skip"
	"github.com/cockroachdb/cockroach/pkg/testutils/testcluster"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/randutil"
)

const configIdx = 5

var logicTestDir string

func init() {
	if bazel.BuiltWithBazel() {
		var err error
		logicTestDir, err = bazel.Runfile("pkg/sql/logictest/testdata/logic_test")
		if err != nil {
			panic(err)
		}
	} else {
		logicTestDir = "../../../../sql/logictest/testdata/logic_test"
	}

}

func TestMain(m *testing.M) {
	securityassets.SetLoader(securitytest.EmbeddedAssets)
	randutil.SeedForTests()
	serverutils.InitTestServerFactory(server.TestServerFactory)
	serverutils.InitTestClusterFactory(testcluster.TestClusterFactory)
	os.Exit(m.Run())
}

func runLogicTest(t *testing.T, file string) {
	skip.UnderDeadlock(t, "times out and/or hangs")
	logictest.RunLogicTest(t, logictest.TestServerArgs{}, configIdx, filepath.Join(logicTestDir, file))
}

func TestLogic_aggregate(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "aggregate")
}

func TestLogic_alias_types(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alias_types")
}

func TestLogic_alter_column_type(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_column_type")
}

func TestLogic_alter_database_convert_to_schema(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_database_convert_to_schema")
}

func TestLogic_alter_database_owner(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_database_owner")
}

func TestLogic_alter_default_privileges_for_all_roles(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_default_privileges_for_all_roles")
}

func TestLogic_alter_default_privileges_for_schema(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_default_privileges_for_schema")
}

func TestLogic_alter_default_privileges_for_sequence(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_default_privileges_for_sequence")
}

func TestLogic_alter_default_privileges_for_table(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_default_privileges_for_table")
}

func TestLogic_alter_default_privileges_for_type(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_default_privileges_for_type")
}

func TestLogic_alter_default_privileges_in_schema(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_default_privileges_in_schema")
}

func TestLogic_alter_default_privileges_with_grant_option(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_default_privileges_with_grant_option")
}

func TestLogic_alter_primary_key(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_primary_key")
}

func TestLogic_alter_role(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_role")
}

func TestLogic_alter_role_set(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_role_set")
}

func TestLogic_alter_schema_owner(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_schema_owner")
}

func TestLogic_alter_sequence(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_sequence")
}

func TestLogic_alter_sequence_owner(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_sequence_owner")
}

func TestLogic_alter_table(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_table")
}

func TestLogic_alter_table_owner(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_table_owner")
}

func TestLogic_alter_type(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_type")
}

func TestLogic_alter_type_owner(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_type_owner")
}

func TestLogic_alter_view_owner(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "alter_view_owner")
}

func TestLogic_and_or(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "and_or")
}

func TestLogic_apply_join(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "apply_join")
}

func TestLogic_array(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "array")
}

func TestLogic_as_of(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "as_of")
}

func TestLogic_auto_span_config_reconciliation_job(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "auto_span_config_reconciliation_job")
}

func TestLogic_bit(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "bit")
}

func TestLogic_builtin_function(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "builtin_function")
}

func TestLogic_builtin_function_notenant(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "builtin_function_notenant")
}

func TestLogic_bytes(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "bytes")
}

func TestLogic_cascade(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "cascade")
}

func TestLogic_case_sensitive_names(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "case_sensitive_names")
}

func TestLogic_cast(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "cast")
}

func TestLogic_ccl(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "ccl")
}

func TestLogic_check_constraints(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "check_constraints")
}

func TestLogic_cluster_settings(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "cluster_settings")
}

func TestLogic_collatedstring(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "collatedstring")
}

func TestLogic_collatedstring_constraint(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "collatedstring_constraint")
}

func TestLogic_collatedstring_index1(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "collatedstring_index1")
}

func TestLogic_collatedstring_index2(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "collatedstring_index2")
}

func TestLogic_collatedstring_normalization(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "collatedstring_normalization")
}

func TestLogic_collatedstring_nullinindex(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "collatedstring_nullinindex")
}

func TestLogic_collatedstring_uniqueindex1(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "collatedstring_uniqueindex1")
}

func TestLogic_collatedstring_uniqueindex2(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "collatedstring_uniqueindex2")
}

func TestLogic_computed(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "computed")
}

func TestLogic_conditional(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "conditional")
}

func TestLogic_connect_privilege(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "connect_privilege")
}

func TestLogic_crdb_internal(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "crdb_internal")
}

func TestLogic_crdb_internal_default_privileges(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "crdb_internal_default_privileges")
}

func TestLogic_create_as(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "create_as")
}

func TestLogic_create_index(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "create_index")
}

func TestLogic_create_statements(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "create_statements")
}

func TestLogic_create_table(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "create_table")
}

func TestLogic_cross_join(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "cross_join")
}

func TestLogic_cursor(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "cursor")
}

func TestLogic_custom_escape_character(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "custom_escape_character")
}

func TestLogic_dangerous_statements(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "dangerous_statements")
}

func TestLogic_database(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "database")
}

func TestLogic_datetime(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "datetime")
}

func TestLogic_decimal(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "decimal")
}

func TestLogic_default(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "default")
}

func TestLogic_delete(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "delete")
}

func TestLogic_dependencies(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "dependencies")
}

func TestLogic_discard(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "discard")
}

func TestLogic_disjunction_in_join(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "disjunction_in_join")
}

func TestLogic_distinct(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "distinct")
}

func TestLogic_distinct_on(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "distinct_on")
}

func TestLogic_distsql_automatic_stats(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "distsql_automatic_stats")
}

func TestLogic_distsql_event_log(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "distsql_event_log")
}

func TestLogic_distsql_expr(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "distsql_expr")
}

func TestLogic_distsql_join(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "distsql_join")
}

func TestLogic_distsql_srfs(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "distsql_srfs")
}

func TestLogic_drop_database(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "drop_database")
}

func TestLogic_drop_index(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "drop_index")
}

func TestLogic_drop_owned_by(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "drop_owned_by")
}

func TestLogic_drop_role_with_default_privileges(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "drop_role_with_default_privileges")
}

func TestLogic_drop_role_with_default_privileges_in_schema(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "drop_role_with_default_privileges_in_schema")
}

func TestLogic_drop_schema(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "drop_schema")
}

func TestLogic_drop_sequence(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "drop_sequence")
}

func TestLogic_drop_table(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "drop_table")
}

func TestLogic_drop_type(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "drop_type")
}

func TestLogic_drop_user(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "drop_user")
}

func TestLogic_drop_view(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "drop_view")
}

func TestLogic_edge(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "edge")
}

func TestLogic_enums(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "enums")
}

func TestLogic_errors(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "errors")
}

func TestLogic_event_log(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "event_log")
}

func TestLogic_exclude_data_from_backup(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "exclude_data_from_backup")
}

func TestLogic_experimental_distsql_planning(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "experimental_distsql_planning")
}

func TestLogic_explain_analyze(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "explain_analyze")
}

func TestLogic_expression_index(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "expression_index")
}

func TestLogic_family(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "family")
}

func TestLogic_fk(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "fk")
}

func TestLogic_float(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "float")
}

func TestLogic_function_lookup(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "function_lookup")
}

func TestLogic_fuzzystrmatch(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "fuzzystrmatch")
}

func TestLogic_geospatial(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "geospatial")
}

func TestLogic_geospatial_index(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "geospatial_index")
}

func TestLogic_geospatial_meta(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "geospatial_meta")
}

func TestLogic_geospatial_zm(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "geospatial_zm")
}

func TestLogic_grant_database(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "grant_database")
}

func TestLogic_grant_in_txn(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "grant_in_txn")
}

func TestLogic_grant_on_all_sequences_in_schema(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "grant_on_all_sequences_in_schema")
}

func TestLogic_grant_on_all_tables_in_schema(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "grant_on_all_tables_in_schema")
}

func TestLogic_grant_revoke_with_grant_option(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "grant_revoke_with_grant_option")
}

func TestLogic_grant_role(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "grant_role")
}

func TestLogic_grant_schema(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "grant_schema")
}

func TestLogic_hash_join(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "hash_join")
}

func TestLogic_hash_sharded_index(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "hash_sharded_index")
}

func TestLogic_hidden_columns(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "hidden_columns")
}

func TestLogic_impure(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "impure")
}

func TestLogic_index_join(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "index_join")
}

func TestLogic_inet(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "inet")
}

func TestLogic_inflight_trace_spans(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "inflight_trace_spans")
}

func TestLogic_information_schema(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "information_schema")
}

func TestLogic_inner_join(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "inner-join")
}

func TestLogic_insert(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "insert")
}

func TestLogic_int_size(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "int_size")
}

func TestLogic_internal_executor(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "internal_executor")
}

func TestLogic_interval(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "interval")
}

func TestLogic_inverted_filter_geospatial(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "inverted_filter_geospatial")
}

func TestLogic_inverted_filter_json_array(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "inverted_filter_json_array")
}

func TestLogic_inverted_index(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "inverted_index")
}

func TestLogic_inverted_index_multi_column(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "inverted_index_multi_column")
}

func TestLogic_inverted_join_geospatial(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "inverted_join_geospatial")
}

func TestLogic_inverted_join_json_array(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "inverted_join_json_array")
}

func TestLogic_inverted_join_multi_column(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "inverted_join_multi_column")
}

func TestLogic_jobs(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "jobs")
}

func TestLogic_join(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "join")
}

func TestLogic_json(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "json")
}

func TestLogic_json_builtins(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "json_builtins")
}

func TestLogic_kv_builtin_functions(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "kv_builtin_functions")
}

func TestLogic_limit(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "limit")
}

func TestLogic_locality(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "locality")
}

func TestLogic_lock_timeout(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "lock_timeout")
}

func TestLogic_lookup_join(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "lookup_join")
}

func TestLogic_lookup_join_spans(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "lookup_join_spans")
}

func TestLogic_manual_retry(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "manual_retry")
}

func TestLogic_materialized_view(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "materialized_view")
}

func TestLogic_merge_join(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "merge_join")
}

func TestLogic_multi_statement(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "multi_statement")
}

func TestLogic_name_escapes(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "name_escapes")
}

func TestLogic_namespace(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "namespace")
}

func TestLogic_new_schema_changer(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "new_schema_changer")
}

func TestLogic_no_primary_key(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "no_primary_key")
}

func TestLogic_notice(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "notice")
}

func TestLogic_numeric_references(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "numeric_references")
}

func TestLogic_on_update(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "on_update")
}

func TestLogic_operator(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "operator")
}

func TestLogic_optimizer_timeout(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "optimizer_timeout")
}

func TestLogic_order_by(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "order_by")
}

func TestLogic_ordinal_references(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "ordinal_references")
}

func TestLogic_ordinality(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "ordinality")
}

func TestLogic_overflow(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "overflow")
}

func TestLogic_overlaps(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "overlaps")
}

func TestLogic_owner(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "owner")
}

func TestLogic_parallel_stmts_compat(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "parallel_stmts_compat")
}

func TestLogic_partial_index(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "partial_index")
}

func TestLogic_partial_txn_commit(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "partial_txn_commit")
}

func TestLogic_partitioning(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "partitioning")
}

func TestLogic_pg_builtins(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "pg_builtins")
}

func TestLogic_pg_catalog_pg_default_acl(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "pg_catalog_pg_default_acl")
}

func TestLogic_pg_catalog_pg_default_acl_with_grant_option(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "pg_catalog_pg_default_acl_with_grant_option")
}

func TestLogic_pg_extension(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "pg_extension")
}

func TestLogic_pgcrypto_builtins(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "pgcrypto_builtins")
}

func TestLogic_pgoidtype(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "pgoidtype")
}

func TestLogic_poison_after_push(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "poison_after_push")
}

func TestLogic_postgres_jsonb(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "postgres_jsonb")
}

func TestLogic_postgresjoin(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "postgresjoin")
}

func TestLogic_privilege_builtins(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "privilege_builtins")
}

func TestLogic_privileges_comments(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "privileges_comments")
}

func TestLogic_privileges_table(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "privileges_table")
}

func TestLogic_propagate_input_ordering(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "propagate_input_ordering")
}

func TestLogic_reassign_owned_by(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "reassign_owned_by")
}

func TestLogic_record(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "record")
}

func TestLogic_rename_atomic(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "rename_atomic")
}

func TestLogic_rename_column(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "rename_column")
}

func TestLogic_rename_constraint(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "rename_constraint")
}

func TestLogic_rename_database(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "rename_database")
}

func TestLogic_rename_index(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "rename_index")
}

func TestLogic_rename_sequence(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "rename_sequence")
}

func TestLogic_rename_table(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "rename_table")
}

func TestLogic_rename_view(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "rename_view")
}

func TestLogic_reset(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "reset")
}

func TestLogic_returning(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "returning")
}

func TestLogic_row_level_ttl(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "row_level_ttl")
}

func TestLogic_rows_from(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "rows_from")
}

func TestLogic_run_control(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "run_control")
}

func TestLogic_save_table(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "save_table")
}

func TestLogic_scale(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "scale")
}

func TestLogic_scatter(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "scatter")
}

func TestLogic_schema(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "schema")
}

func TestLogic_schema_change_feature_flags(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "schema_change_feature_flags")
}

func TestLogic_schema_change_in_txn(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "schema_change_in_txn")
}

func TestLogic_schema_change_retry(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "schema_change_retry")
}

func TestLogic_schema_repair(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "schema_repair")
}

func TestLogic_scrub(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "scrub")
}

func TestLogic_secondary_index_column_families(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "secondary_index_column_families")
}

func TestLogic_select(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "select")
}

func TestLogic_select_for_update(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "select_for_update")
}

func TestLogic_select_index(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "select_index")
}

func TestLogic_select_index_flags(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "select_index_flags")
}

func TestLogic_select_search_path(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "select_search_path")
}

func TestLogic_select_table_alias(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "select_table_alias")
}

func TestLogic_sequences(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "sequences")
}

func TestLogic_sequences_distsql(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "sequences_distsql")
}

func TestLogic_sequences_regclass(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "sequences_regclass")
}

func TestLogic_serial(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "serial")
}

func TestLogic_serializable_eager_restart(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "serializable_eager_restart")
}

func TestLogic_set_local(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "set_local")
}

func TestLogic_set_role(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "set_role")
}

func TestLogic_set_schema(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "set_schema")
}

func TestLogic_set_time_zone(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "set_time_zone")
}

func TestLogic_shift(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "shift")
}

func TestLogic_show_completions(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "show_completions")
}

func TestLogic_show_create(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "show_create")
}

func TestLogic_show_create_all_schemas(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "show_create_all_schemas")
}

func TestLogic_show_create_all_tables(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "show_create_all_tables")
}

func TestLogic_show_create_all_tables_builtin(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "show_create_all_tables_builtin")
}

func TestLogic_show_create_all_types(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "show_create_all_types")
}

func TestLogic_show_default_privileges(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "show_default_privileges")
}

func TestLogic_show_fingerprints(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "show_fingerprints")
}

func TestLogic_show_indexes(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "show_indexes")
}

func TestLogic_show_transfer_state(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "show_transfer_state")
}

func TestLogic_span_builtins(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "span_builtins")
}

func TestLogic_split_at(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "split_at")
}

func TestLogic_sqllite(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "sqllite")
}

func TestLogic_sqlsmith(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "sqlsmith")
}

func TestLogic_srfs(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "srfs")
}

func TestLogic_statement_source(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "statement_source")
}

func TestLogic_storing(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "storing")
}

func TestLogic_strict_ddl_atomicity(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "strict_ddl_atomicity")
}

func TestLogic_suboperators(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "suboperators")
}

func TestLogic_subquery(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "subquery")
}

func TestLogic_subquery_correlated(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "subquery_correlated")
}

func TestLogic_system(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "system")
}

func TestLogic_system_columns(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "system_columns")
}

func TestLogic_system_namespace(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "system_namespace")
}

func TestLogic_system_privileges(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "system_privileges")
}

func TestLogic_table(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "table")
}

func TestLogic_target_names(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "target_names")
}

func TestLogic_temp_table(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "temp_table")
}

func TestLogic_temp_table_txn(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "temp_table_txn")
}

func TestLogic_tenant(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "tenant")
}

func TestLogic_time(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "time")
}

func TestLogic_timestamp(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "timestamp")
}

func TestLogic_timetz(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "timetz")
}

func TestLogic_trigram_builtins(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "trigram_builtins")
}

func TestLogic_trigram_indexes(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "trigram_indexes")
}

func TestLogic_truncate(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "truncate")
}

func TestLogic_tuple(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "tuple")
}

func TestLogic_txn(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "txn")
}

func TestLogic_txn_as_of(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "txn_as_of")
}

func TestLogic_txn_retry(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "txn_retry")
}

func TestLogic_txn_stats(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "txn_stats")
}

func TestLogic_type_privileges(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "type_privileges")
}

func TestLogic_typing(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "typing")
}

func TestLogic_udf(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "udf")
}

func TestLogic_union(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "union")
}

func TestLogic_unique(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "unique")
}

func TestLogic_update(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "update")
}

func TestLogic_update_from(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "update_from")
}

func TestLogic_upsert(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "upsert")
}

func TestLogic_uuid(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "uuid")
}

func TestLogic_values(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "values")
}

func TestLogic_vectorize_agg(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "vectorize_agg")
}

func TestLogic_vectorize_shutdown(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "vectorize_shutdown")
}

func TestLogic_vectorize_types(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "vectorize_types")
}

func TestLogic_vectorize_unsupported(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "vectorize_unsupported")
}

func TestLogic_views(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "views")
}

func TestLogic_virtual_columns(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "virtual_columns")
}

func TestLogic_void(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "void")
}

func TestLogic_where(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "where")
}

func TestLogic_window(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "window")
}

func TestLogic_with(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "with")
}

func TestLogic_zero(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "zero")
}

func TestLogic_zigzag_join(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "zigzag_join")
}

func TestLogic_zone_config(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "zone_config")
}

func TestLogic_zone_config_system_tenant(
	t *testing.T,
) {
	defer leaktest.AfterTest(t)()
	runLogicTest(t, "zone_config_system_tenant")
}
