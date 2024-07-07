package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var Migrations = []*gormigrate.Migration{
	{
		ID: "init",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
				CREATE TABLE parent_product (
					sku_id int PRIMARY KEY,
					name varchar,
					description varchar
				);

				CREATE TABLE child_product(
					sku_id int PRIMARY KEY,
					parent_sku int,
					size varchar,
					color varchar,
					price int,
					quantity int,
					CONSTRAINT 
						fk_parent
						FOREIGN KEY(parent_sku) 
						REFERENCES parent_product(sku_id)
						ON DELETE CASCADE
				)
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`
				DELETE TABLE parent_product;
				DELETE TABLE child_product;
			`).Error
		},
	},
}
