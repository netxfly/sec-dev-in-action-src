package mysql

import (
	"time"

	sqle "github.com/xsec-lab/go-mysql-server"
	"github.com/xsec-lab/go-mysql-server/auth"
	"github.com/xsec-lab/go-mysql-server/memory"
	"github.com/xsec-lab/go-mysql-server/server"
	"github.com/xsec-lab/go-mysql-server/sql"

	"sec-dev-in-action-src/honeypot/server/logger"
)

func StartMysql(addr string, flag bool) error {
	engine := sqle.NewDefault()
	engine.AddDatabase(createTestDatabase())
	engine.AddDatabase(sql.NewInformationSchemaDatabase(engine.Catalog))

	config := server.Config{
		Protocol: "tcp",
		Address:  addr,
		Auth:     auth.NewNativeSingle("root", "123456", auth.DefaultPermissions),
	}

	s, err := server.NewDefaultServer(config, engine)
	logger.Log.Warningf("start mysql service on %v", addr)
	if err != nil {
		return err
	}

	err = s.Start()
	return err
}

func createTestDatabase() *memory.Database {
	const (
		dbName    = "my_db"
		tableName = "my_table"
	)

	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tableName, sql.Schema{
		{Name: "name", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "email", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "phone_numbers", Type: sql.JSON, Nullable: false, Source: tableName},
		{Name: "created_at", Type: sql.Timestamp, Nullable: false, Source: tableName},
	})

	db.AddTable(tableName, table)
	ctx := sql.NewEmptyContext()

	_ = table.Insert(ctx, sql.NewRow("netxfly", "x@xsec.io", []string{"xsec.io"}, time.Now()))
	_ = table.Insert(ctx, sql.NewRow("sec.lu", "root@xsec.io", []string{}, time.Now()))
	_ = table.Insert(ctx, sql.NewRow("xsec.io", "jane@sec.lu", []string{"sec.lu"}, time.Now()))
	_ = table.Insert(ctx, sql.NewRow("xsec", "evilbob@sec.lu", []string{"555-666-555", "666-666-666"}, time.Now()))

	return db
}
