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
	modelPath := filepath.Join(workdir, "biz", "dal", "entity")
	queryPath := filepath.Join(workdir, "biz", "dal", "query")

	g := gen.NewGenerator(gen.Config{
		ModelPkgPath:      modelPath,
		OutPath:           queryPath,
		Mode:              gen.WithoutContext | gen.WithQueryInterface | gen.WithDefaultQuery,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	dsn := "root:123456@(127.0.0.1:3306)/cloud_storage?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn))
	g.UseDB(db)

	g.ApplyBasic(
		g.GenerateModel("repository_pool"),
		g.GenerateModel("share_basic"),
		g.GenerateModel("user_basic"),
		g.GenerateModel("user_repository"),
	)

	g.Execute()
}
