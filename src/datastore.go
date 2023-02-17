package main

import (
	"gopkg.in/yaml.v3"
	"sort"
)

type DataStore struct {
	node                  *yaml.Node
	Id                    string
	Name                  string
	Description           string
	State                 State
	TechnologyIds         []string
	TechnologyBundleId    string
	Technologies          []*Technology
	ApiTechnologyIds      []string
	ApiTechnologyBundleId string
	ApiTechnologies       []*Technology
}

func (s *DataStore) Print(printer *Printer) {
	printer.Print(s.Name)
	s.State.Print(printer)
	PrintTechnologies(s.Technologies, printer)
	printer.NewLine()
}

func (s *DataStore) getNode() *yaml.Node {
	return s.node
}

func (s *DataStore) getTechnologyIds() []string {
	return s.TechnologyIds
}

func (s *DataStore) getTechnologyBundleId() string {
	return s.TechnologyBundleId
}

func (s *DataStore) getTechnologies() []*Technology {
	return s.Technologies
}

func (s *DataStore) setTechnologies(technologies []*Technology) {
	s.Technologies = technologies
}

func (s *DataStore) setNode(node *yaml.Node) {
	s.node = node
}

func (s *DataStore) setId(id string) {
	s.Id = id
}

func (s *DataStore) setName(name string) {
	s.Name = name
}

func (s *DataStore) getDescription() string {
	return s.Description
}

func (s *DataStore) setDescription(description string) {
	s.Description = description
}

func (s *DataStore) setState(state State) {
	s.State = state
}

func (s *DataStore) setTechnologyIds(technologies []string) {
	s.TechnologyIds = technologies
}

func (s *DataStore) setTechnologyBundleId(technologyBundle string) {
	s.TechnologyBundleId = technologyBundle
}

type ApiTechnologies struct {
	dataStore *DataStore
}

func (a ApiTechnologies) getNode() *yaml.Node {
	return a.dataStore.getNode()
}

func (a ApiTechnologies) getTechnologyIds() []string {
	return a.dataStore.ApiTechnologyIds
}

func (a ApiTechnologies) getTechnologyBundleId() string {
	return a.dataStore.ApiTechnologyBundleId
}

func (a ApiTechnologies) getTechnologies() []*Technology {
	return a.dataStore.Technologies
}

func (a ApiTechnologies) setTechnologies(technologies []*Technology) {
	a.dataStore.ApiTechnologies = technologies
}

func (a ApiTechnologies) setTechnologyIds(technologies []string) {
	a.dataStore.ApiTechnologyIds = technologies
}

func (a ApiTechnologies) setTechnologyBundleId(technologyBundle string) {
	a.dataStore.ApiTechnologyBundleId = technologyBundle
}

func (s *DataStore) read(id string, node *yaml.Node) []Issue {
	var fields map[string]*yaml.Node
	fields, issues := namedObject(node, id, s)
	issues = append(issues, setDescription(fields, s)...)
	issues = append(issues, setState(node, fields, s)...)
	issues = append(issues, setTechnologies(fields, s)...)
	issues = append(issues, setTechnologiesFrom(fields, "apiTechnologies", ApiTechnologies{s})...)
	return issues
}

type DataStoreReader struct {
}

func (r DataStoreReader) read(node *yaml.Node) ([]*DataStore, []Issue) {
	dataStores := make([]*DataStore, 0)
	if node == nil {
		return dataStores, []Issue{}
	}
	dataStoresById, issue := toMap(node)
	if issue != nil {
		return dataStores, []Issue{*issue}
	}
	issues := make([]Issue, 0)
	for id, dataStoreNode := range dataStoresById {
		dataStore := DataStore{}
		dataStores = append(dataStores, &dataStore)
		issues = append(issues, dataStore.read(id, dataStoreNode)...)
	}
	sort.Slice(dataStores, func(i, j int) bool {
		return dataStores[i].Name < dataStores[j].Name
	})
	return dataStores, issues
}

type DataStoreValidator struct {
}

type DataStoreIdExtractor func(use *DataStoreUse) string

func (d DataStoreValidator) validate(model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	for _, database := range model.Databases {
		if d.isUnused(&database.DataStore, model, func(use *DataStoreUse) string { return use.DatabaseId }) {
			issues = append(issues, *NodeWarning("Database isn't used", database.node))
		}
	}
	for _, queue := range model.Queues {
		if d.isUnused(queue, model, func(use *DataStoreUse) string { return use.QueueId }) {
			issues = append(issues, *NodeWarning("Queue isn't used", queue.node))
		}
	}
	return issues
}

func (d DataStoreValidator) isUnused(dataStore *DataStore, model *ArchitectureModel, extractor DataStoreIdExtractor) bool {
	for _, service := range model.Services {
		for _, use := range service.DataStores {
			id := extractor(use)
			if dataStore.Id == id {
				return false
			}
		}
	}
	return true
}
