# Gormite

# Usage

```shell
go run ./pkg/gormite/runners/exec -t goose --dsn "$DATABASE_DSN"
```

Gormite todo:
* gormite better diff generation using custom sources(example: postgres database + postgres database).
* relationship:
  OneToOne - uniq_c4f5ad0c3305d17d,
  OneToMany - virtual, not owner,
  ManyToOne - idx_79b67c3d4a95524f,
  ManyToMany - table1_table2_pkey, idx_f9cfc2b5b7c79121, idx_f981029cb7c72e1 - create separate table

## Tags description:
* db - column name in db
  ```go
		package main
		
		type _ struct {
			ID int `db:"id"`
		}
	```

* pk - if set then column is primary key
    ```go
	  package main
	  
	  type _ struct {
		  ID int `db:"id" id:"true"`
	  }
  ```

* nullable - if set then column is nullable(without NOT NULL)
    ```go
	  package main
	  
	  type _ struct {
		  Code *int `db:"code" nullable:"true"`
	  }
  ```

* length - column length(useful for varchar)
    ```go
	  package main
	  
	  type _ struct {
		  Name string `db:"name" length:"64"`
	  }
  ```

* uniq - if set then column is unique(can be placed to many properties - grouped, on single column can be many uniq's, separated by `,`)
    ```go
	  package main
	  
      // CREATE UNIQUE INDEX uniq_name ON example_table (name);
	  type _ struct {
		  Name string `db:"name" uniq:"uniq_name"`
	  }
  
	  // CREATE UNIQUE INDEX uniq_name_parted ON example_table (name_part1, name2);
      type _ struct {
		  NamePart1 string `db:"name_part1" uniq:"uniq_name_parted"`
		  NamePart2 string `db:"name2" uniq:"uniq_name_parted"`
	  }
  ```

* uniq_cond - work in pair with uniq, uniq_cond is a condition for uniq, can be placed on any column(of grouped uniq's, do not set to every column of grouped uniq's, can be multiple, separated by `,`)
    ```go
	  package main
  
	  // CREATE UNIQUE INDEX uniq_name_parted ON example_table (name1, name2) WHERE is_selected = true;
      type _ struct {
  		  Checked bool `db:"is_checked"`
		  NamePart1 string `db:"name1" uniq:"uniq_name_parted" uniq_cond:"uniq_name_parted:(is_checked = true)"`
		  NamePart2 string `db:"name2" uniq:"uniq_name_parted"`
	  }
  ```

* index - if set then column is index(logic same as uniq)
    ```go
	  package main
	  
      // CREATE INDEX idx_name ON example_table (name);
	  type _ struct {
		  Name string `db:"name" index:"idx_name"`
	  }
  
	  // CREATE INDEX idx_name_parted ON example_table (name_part1, name2);
      type _ struct {
		  NamePart1 string `db:"name_part1" index:"idx_name_parted"`
		  NamePart2 string `db:"name2" index:"idx_name_parted"`
	  }
  ```

* index_cond - work in pair with index(logic same as uniq_cond)
    ```go
	  package main
  
	  // CREATE INDEX idx_name_parted ON example_table (name1, name2) WHERE is_selected = true;
      type _ struct {
  		  Checked bool `db:"is_checked"`
		  NamePart1 string `db:"name1" index:"idx_name_parted" index_cond:"idx_name_parted:(is_checked = true)"`
		  NamePart2 string `db:"name2" index:"idx_name_parted"`
	  }
  ```

* default - if set then column has default value
    ```go
	  package main
  
      type _ struct {
		  Name string `db:"name" default:"123"`
	  }
  ```

* type - if set then column has type(manually set, without struct property type checking)
    ```go
	  package main
  
  	  type ExampleType int
  
      type _ struct {
		  Type ExampleType `db:"type" type:"int"`
	  }
  ```
  Possible types: text, varchar, json, jsonb, integer
