package layouts

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Migrate(db *sql.DB) {
	if _, err := db.Exec(`
create table if not exists migrations (
	key text primary key on conflict replace,
	value int
)`); err != nil {
		panic(err)
	}

	var version int
	err := db.QueryRow("select value from migrations where key = 'version'").Scan(&version)
	switch err {
	case sql.ErrNoRows:
		version = -1
	case nil:
	default:
		panic(err)
	}
	migrations := []string{
		`create table if not exists layouts(
			id string primary key,
			size_width numeric, size_height numeric,
			bounds_left numeric, bounds_top numeric, bounds_right numeric, bounds_bottom numeric,
			heightmap_url text
		)`,
		`create table if not exists heightmaps(
			id text primary key on conflict replace,
			heightmap blob
		)`,
		`create table layers(
			id text,
			name text,
			layer blob,
			primary key (id, name) on conflict replace
		)`,
		`insert into layers
			select id, 'heightmap.jpg', heightmap
			from heightmaps`,
		`drop table heightmaps`,
		`alter table layouts add column hillshade_url text`,
	}
	for index, migration := range migrations {
		if index > version {
			if _, err = db.Exec(migration); err != nil {
				panic(err)
			}
		}
	}

	statement, err := db.Prepare("insert into migrations (key, value) values(?, ?)")
	if err != nil {
		panic(err)
	}
	if _, err := statement.Exec("version", len(migrations)-1); err != nil {
		panic(err)
	}
}
