package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=123 dbname=leasing_db port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect database:", err)
	}

	log.Println("✅ Connected to database")

	g := gen.NewGenerator(gen.Config{
		OutPath:           "internal/domain/query",
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.UseDB(db)

	g.WithTableNameStrategy(func(tableName string) string {
		if tableName == "schema_migrations" {
			return ""
		}
		return tableName
	})

	log.Println("🚀 Generating models...")

	schemas := []string{"mst", "account", "dealer", "leasing", "payment"}

	for _, schema := range schemas {
		log.Println("📦 Generating schema:", schema)
		db.Exec("SET search_path TO " + schema)
		g.ApplyBasic(g.GenerateAllTable()...)
	}

	g.Execute()

	log.Println("🎉 Generate complete!")
}
