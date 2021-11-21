package v1

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var migrations = []string{
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
	`update layouts set hillshade_url = '' where hillshade_url IS NULL`,
}

func Migrate(db *sql.DB) error {
	if err := CreateMigrationTable(db); err != nil {
		return err
	}

	version, err := GetMigrationVersion(db)
	if err != nil {
		return err
	}

	if err := ApplyMigrations(db, migrations[version+1:]); err != nil {
		return err
	}

	return SetMigrationVersion(db, len(migrations)-1)
}

func CreateMigrationTable(db *sql.DB) error {
	_, err := db.Exec(`
create table if not exists migrations (
	key text primary key on conflict replace,
	value int
)`)
	return err
}

func GetMigrationVersion(db *sql.DB) (int, error) {
	var version int
	err := db.QueryRow("select value from migrations where key = 'version'").Scan(&version)
	switch err {
	case sql.ErrNoRows:
		return -1, nil
	case nil:
		return version, nil
	default:
		return 0, err
	}
}

func ApplyMigrations(db *sql.DB, migrations []string) error {
	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return err
		}
	}
	return nil
}

func SetMigrationVersion(db *sql.DB, version int) error {
	statement, err := db.Prepare("insert into migrations (key, value) values(?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec("version", version)
	return err
}
