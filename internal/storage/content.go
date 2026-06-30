package storage

import (
	"context"
	"fmt"

	"github.com/mateom/vaultsh/internal/content"
	"github.com/mateom/vaultsh/internal/filesystem"
)

func Load(ctx context.Context, provider content.Provider) (*filesystem.Directory, error) {
	catalog, err := provider.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("load content: %w", err)
	}

	root := filesystem.NewDirectory("")
	experienceDirectory := filesystem.NewDirectory("experience")
	projectDirectory := filesystem.NewDirectory("projects")

	for _, node := range []filesystem.Node{
		filesystem.NewFile("about.txt", renderAbout(catalog.About)),
		filesystem.NewFile("education.txt", renderEducation(catalog.Education)),
		filesystem.NewFile("skills.txt", renderSkills(catalog.Skills)),
		experienceDirectory,
		projectDirectory,
	} {
		if err := root.Add(node); err != nil {
			return nil, fmt.Errorf("build root: %w", err)
		}
	}

	for _, experience := range catalog.Experiences {
		if err := experienceDirectory.Add(filesystem.NewFile(
			experience.Slug+".txt",
			renderExperience(experience),
		)); err != nil {
			return nil, fmt.Errorf("build experience directory: %w", err)
		}
	}
	for _, project := range catalog.Projects {
		if err := projectDirectory.Add(filesystem.NewFile(
			project.Slug+".txt",
			renderProject(project),
		)); err != nil {
			return nil, fmt.Errorf("build projects directory: %w", err)
		}
	}

	return root, nil
}

func renderAbout(about content.About) string             { return about.Text }
func renderEducation(education content.Education) string { return education.Text }
func renderSkills(skills content.Skills) string          { return skills.Text }
func renderExperience(experience content.Experience) string {
	return experience.Text
}
func renderProject(project content.Project) string { return project.Text }
