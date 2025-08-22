package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

func main() {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	modelPath := filepath.Join(workdir, "internal", "dal", "model")
	queryPath := filepath.Join(workdir, "internal", "dal", "query")

	g := gen.NewGenerator(gen.Config{
		ModelPkgPath:      modelPath,
		OutPath:           queryPath,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	dsn := "root:123456@tcp(127.0.0.1:3306)/cloud_storage?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	g.UseDB(db)
	g.ApplyBasic(
		g.GenerateModel("repository_pool"),
		g.GenerateModel("share_basic"),
		g.GenerateModel("user_basic"),
		g.GenerateModel("user_repository"),
	)
	g.Execute()
}
