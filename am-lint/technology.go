package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"sort"
)

type Implementable interface {
	getNode() *yaml.Node
	getTechnologyIds() []string
	setTechnologyIds(technologies []string)
	getTechnologyBundleId() string
	setTechnologyBundleId(technologyBundle string)
	setTechnologies([]*Technology)
}

func PrintTechnologies(technologies []*Technology, printer *Printer) {
	if len(technologies) > 0 {
		printer.Print(" [")
		prefix := ""
		for _, tech := range technologies {
			printer.Print(prefix, tech.Name)
			prefix = ", "
		}
		printer.Print("]")
	}
}

func setTechnologies(fields map[string]*yaml.Node, implementable Implementable) []Issue {
	return setTechnologiesFrom(fields, "technologies", implementable)
}

func setTechnologiesFrom(fields map[string]*yaml.Node, technologiesField string, implementable Implementable) []Issue {
	issues := make([]Issue, 0)
	technologyBundle, found, issue := stringFieldOf(fields, technologiesField)
	if issue == nil {
		if found {
			implementable.setTechnologyBundleId(technologyBundle)
		}
		return issues
	}
	technologiesNodes, _, issue := sequenceFieldOf(fields, technologiesField)
	if issue == nil {
		technologies := make([]string, 0)
		for _, technologyNode := range technologiesNodes {
			technology, issue := toString(technologyNode, "technology")
			if issue == nil {
				technologies = append(technologies, technology)
			} else {
				issues = append(issues, *issue)
			}
		}
		sort.Slice(technologies, func(i, j int) bool {
			return technologies[i] < technologies[j]
		})
		implementable.setTechnologyIds(technologies)
	} else {
		issues = append(issues, *issue)
	}

	return issues
}

type Quadrant int64

const (
	LanguagesAndFrameworks = iota
	Platforms
	Tools
	Techniques
)

type Ring int64

const (
	Trial = iota
	Assess
	Adopt
	Hold
)

type Technology struct {
	node        *yaml.Node
	Id          string
	Name        string
	Description string
	Quadrant    Quadrant
	Ring        Ring
}

func (t *Technology) read(id string, node *yaml.Node) []Issue {
	var fields map[string]*yaml.Node
	fields, issues := namedObject(node, id, t)
	issues = append(issues, setQuadrant(node, fields, t)...)
	issues = append(issues, setRing(node, fields, t)...)
	return issues
}

func (t *Technology) setNode(node *yaml.Node) {
	t.node = node
}

func (t *Technology) setId(id string) {
	t.Id = id
}

func (t *Technology) setName(name string) {
	t.Name = name
}

func setQuadrant(owner *yaml.Node, fields map[string]*yaml.Node, technology *Technology) []Issue {
	var allowedQuadrants = []string{"languagesAndFrameworks", "platforms", "tools", "techniques"}

	value, issue := enumFieldOf(owner, fields, "quadrant", allowedQuadrants, "")
	if issue != nil {
		return []Issue{*issue}
	}
	for index, quadrant := range allowedQuadrants {
		if quadrant == value {
			technology.Quadrant = Quadrant(index)
		}
	}
	return []Issue{}
}

func setRing(owner *yaml.Node, fields map[string]*yaml.Node, technology *Technology) []Issue {
	defaultRing := "adopt"
	var allowedRings = []string{"trial", "assess", defaultRing, "hold"}

	value, issue := enumFieldOf(owner, fields, "ring", allowedRings, defaultRing)
	if issue != nil {
		return []Issue{*issue}
	}
	for index, ring := range allowedRings {
		if ring == value {
			technology.Ring = Ring(index)
		}
	}
	return []Issue{}
}

type TechnologyReader struct {
}

func (r TechnologyReader) read(node *yaml.Node, _ string, model *ArchitectureModel) []Issue {
	if node == nil {
		return []Issue{}
	}
	technologiesById, _ := toMap(node)
	issues := make([]Issue, 0)
	technologies := make([]*Technology, 0)
	for id, technologyNode := range technologiesById {
		technology := Technology{}
		technologies = append(technologies, &technology)
		issues = append(issues, technology.read(id, technologyNode)...)
	}
	sort.Slice(technologies, func(i, j int) bool {
		return technologies[i].Name < technologies[j].Name
	})
	model.Technologies = technologies
	return issues
}

type TechnologyBundle struct {
	node          *yaml.Node
	Id            string
	TechnologyIds []string
	Technologies  []*Technology
}

func (b *TechnologyBundle) read(id string, node *yaml.Node) []Issue {
	b.node = node
	b.Id = id
	technologyIdNodes, issue := toSequence(node, "technologies")
	if issue != nil {
		return []Issue{*issue}
	}
	issues := make([]Issue, 0)
	technologyIds := make([]string, 0)
	for _, technologyIdNode := range technologyIdNodes {
		technologyId, issue := toString(technologyIdNode, "technology")
		if issue != nil {
			issues = append(issues, *issue)
		} else {
			technologyIds = append(technologyIds, technologyId)
		}
	}
	sort.Slice(technologyIds, func(i, j int) bool {
		return technologyIds[i] < technologyIds[j]
	})
	b.TechnologyIds = technologyIds
	return issues
}

type TechnologyBundleReader struct {
}

func (t TechnologyBundleReader) read(node *yaml.Node, _ string, model *ArchitectureModel) []Issue {
	if node == nil {
		return []Issue{}
	}
	technologyBundlesById, _ := toMap(node)
	issues := make([]Issue, 0)
	technologyBundles := make([]*TechnologyBundle, 0)
	for id, technologyBundleNode := range technologyBundlesById {
		technologyBundle := TechnologyBundle{}
		technologyBundles = append(technologyBundles, &technologyBundle)
		issues = append(issues, technologyBundle.read(id, technologyBundleNode)...)
	}
	sort.Slice(technologyBundles, func(i, j int) bool {
		return technologyBundles[i].Id < technologyBundles[j].Id
	})
	model.TechnologyBundles = technologyBundles
	return issues
}

type TechnologyBundleConnector struct {
}

func (c TechnologyBundleConnector) connect(model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	for _, technologyBundle := range model.TechnologyBundles {
		issues = append(issues, c.resolveTechnologies(technologyBundle, model)...)
	}
	return issues
}

func (c TechnologyBundleConnector) resolveTechnologies(bundle *TechnologyBundle, model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	if len(bundle.Technologies) > 0 || len(bundle.TechnologyIds) == 0 {
		return issues
	}
	bundleTechnologies := make([]*Technology, 0)
	for _, id := range bundle.TechnologyIds {
		technologies, issue := technologiesById(bundle.node, id, model)
		if issue != nil {
			issues = append(issues, *issue)
		} else {
			bundleTechnologies = appendUnique(bundleTechnologies, technologies)
		}
	}
	sort.Slice(bundleTechnologies, func(i, j int) bool {
		return bundleTechnologies[i].Name < bundleTechnologies[j].Name
	})
	bundle.Technologies = bundleTechnologies
	return issues
}

func appendUnique(existing []*Technology, toAdd []*Technology) []*Technology {
	result := make([]*Technology, 0)
	result = append(result, existing...)
	for _, technology := range toAdd {
		_, found := lookUp(existing, technology.Id)
		if !found {
			result = append(result, technology)
		}
	}
	return result
}

func lookUpTechnology(model *ArchitectureModel, id string) (*Technology, bool) {
	result, found := lookUp(model.Technologies, id)
	return result, found
}

func lookUp(candidates []*Technology, id string) (*Technology, bool) {
	for _, candidate := range candidates {
		if id == candidate.Id {
			return candidate, true
		}
	}
	return nil, false
}

func connectTechnologies(implementable Implementable, model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	technologies := make([]*Technology, 0)
	id := implementable.getTechnologyBundleId()
	if id != "" {
		issues, technologies = addTechnologies(implementable, model, id, issues, technologies)
	}
	for _, id := range implementable.getTechnologyIds() {
		issues, technologies = addTechnologies(implementable, model, id, issues, technologies)
	}
	sort.Slice(technologies, func(i, j int) bool {
		return technologies[i].Name < technologies[j].Name
	})
	implementable.setTechnologies(technologies)
	return issues
}

func addTechnologies(implementable Implementable, model *ArchitectureModel, id string, issues []Issue, technologies []*Technology) ([]Issue, []*Technology) {
	foundTechnologies, issue := technologiesById(implementable.getNode(), id, model)
	if issue != nil {
		issues = append(issues, *issue)
	} else {
		technologies = appendUnique(technologies, foundTechnologies)
	}
	return issues, technologies
}

func technologiesById(owner *yaml.Node, id string, model *ArchitectureModel) ([]*Technology, *Issue) {
	bundle, found := model.findTechnologyBundleById(id)
	if found {
		return bundle.Technologies, nil
	}
	technology, found := lookUpTechnology(model, id)
	if found {
		return []*Technology{technology}, nil
	}
	return []*Technology{}, NodeError(fmt.Sprintf("Unknown technology '%v'", id), owner)
}
