//go:build integration

package postgres_test

import (
	"context"
	"github.com/Permify/permify/internal/repositories"
	pg "github.com/Permify/permify/internal/repositories/postgres"
	"github.com/Permify/permify/pkg/database"
	PQDatabase "github.com/Permify/permify/pkg/database/postgres"
	"github.com/Permify/permify/pkg/logger"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSchemaReaderHeadVersion_Integration(t *testing.T) {
	ctx := context.Background()

	l := logger.New("fatal")

	err := repositories.Migrate(cfg, l)
	require.NoError(t, err)

	var db database.Database
	db, err = PQDatabase.New(cfg.URI,
		PQDatabase.MaxOpenConnections(cfg.MaxOpenConnections),
		PQDatabase.MaxIdleConnections(cfg.MaxIdleConnections),
		PQDatabase.MaxConnectionIdleTime(cfg.MaxConnectionIdleTime),
		PQDatabase.MaxConnectionLifeTime(cfg.MaxConnectionLifetime),
	)
	require.NoError(t, err)

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a TenantWriter instance
	schemaWriter := pg.NewSchemaWriter(db.(*PQDatabase.Postgres), l)
	schemaReader := pg.NewSchemaReader(db.(*PQDatabase.Postgres), l)

	v := xid.New().String()
	schemas := []repositories.SchemaDefinition{
		{TenantID: "t1", EntityType: "entity2", SerializedDefinition: []byte("def2"), Version: v},
	}

	// Test the CreateTenant method
	err = schemaWriter.WriteSchema(ctx, schemas)
	require.NoError(t, err)

	version, err := schemaReader.HeadVersion(ctx, "t1")
	require.NoError(t, err)
	require.Equal(t, v, version)
}