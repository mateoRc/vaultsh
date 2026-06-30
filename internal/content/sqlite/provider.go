package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/mateom/vaultsh/internal/content"
	_ "modernc.org/sqlite"
)

type Provider struct {
	db       *sql.DB
	filename string
}

// Open materializes the embedded database in a private temporary file and opens
// it in immutable, read-only mode. SQLite never receives write access.
func Open(database []byte) (*Provider, error) {
	if len(database) == 0 {
		return nil, errors.New("embedded database is empty")
	}

	file, err := os.CreateTemp("", "vaultsh-content-*.db")
	if err != nil {
		return nil, fmt.Errorf("create temporary database: %w", err)
	}
	filename := file.Name()
	cleanup := func() {
		file.Close()
		os.Remove(filename)
	}

	if _, err := file.Write(database); err != nil {
		cleanup()
		return nil, fmt.Errorf("write temporary database: %w", err)
	}
	if err := file.Close(); err != nil {
		os.Remove(filename)
		return nil, fmt.Errorf("close temporary database: %w", err)
	}
	if err := os.Chmod(filename, 0o400); err != nil {
		os.Remove(filename)
		return nil, fmt.Errorf("make temporary database read-only: %w", err)
	}

	dsn := (&url.URL{
		Scheme:   "file",
		Path:     filepath.ToSlash(filename),
		RawQuery: "mode=ro&immutable=1",
	}).String()
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		os.Remove(filename)
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		db.Close()
		os.Remove(filename)
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Provider{db: db, filename: filename}, nil
}

func (p *Provider) Close() error {
	return errors.Join(p.db.Close(), os.Remove(p.filename))
}

func (p *Provider) Load(ctx context.Context) (content.Catalog, error) {
	var catalog content.Catalog

	if err := p.db.QueryRowContext(ctx, "SELECT text FROM about WHERE id = 1").
		Scan(&catalog.About.Text); err != nil {
		return content.Catalog{}, fmt.Errorf("load about: %w", err)
	}
	if err := p.db.QueryRowContext(ctx, "SELECT text FROM education WHERE id = 1").
		Scan(&catalog.Education.Text); err != nil {
		return content.Catalog{}, fmt.Errorf("load education: %w", err)
	}
	if err := p.db.QueryRowContext(ctx, "SELECT text FROM skills WHERE id = 1").
		Scan(&catalog.Skills.Text); err != nil {
		return content.Catalog{}, fmt.Errorf("load skills: %w", err)
	}

	experiences, err := loadExperiences(ctx, p.db)
	if err != nil {
		return content.Catalog{}, err
	}
	catalog.Experiences = experiences

	projects, err := loadProjects(ctx, p.db)
	if err != nil {
		return content.Catalog{}, err
	}
	catalog.Projects = projects

	return catalog, nil
}

func loadExperiences(ctx context.Context, db *sql.DB) ([]content.Experience, error) {
	rows, err := db.QueryContext(ctx, "SELECT slug, text FROM experiences ORDER BY slug")
	if err != nil {
		return nil, fmt.Errorf("query experiences: %w", err)
	}
	defer rows.Close()

	var experiences []content.Experience
	for rows.Next() {
		var experience content.Experience
		if err := rows.Scan(&experience.Slug, &experience.Text); err != nil {
			return nil, fmt.Errorf("scan experience: %w", err)
		}
		experiences = append(experiences, experience)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read experiences: %w", err)
	}
	return experiences, nil
}

func loadProjects(ctx context.Context, db *sql.DB) ([]content.Project, error) {
	rows, err := db.QueryContext(ctx, "SELECT slug, text FROM projects ORDER BY slug")
	if err != nil {
		return nil, fmt.Errorf("query projects: %w", err)
	}
	defer rows.Close()

	var projects []content.Project
	for rows.Next() {
		var project content.Project
		if err := rows.Scan(&project.Slug, &project.Text); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read projects: %w", err)
	}
	return projects, nil
}
