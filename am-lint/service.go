package main

import (
	"gopkg.in/yaml.v3"
	"sort"
)

type DataStoreUse struct {
	QueueId     string `yaml:"queue,omitempty"`
	DatabaseId  string `yaml:"database,omitempty"`
	Description string
	DataFlow    string
}

func (d *DataStoreUse) setDirection(direction string) {
	d.DataFlow = direction
}

func (d *DataStoreUse) read(node *yaml.Node, issues []Issue) []Issue {
	fields, issue := toMap(node)
	if issue != nil {
		return []Issue{*issue}
	}
	queue, found, issue := stringFieldOf(fields, "queue")
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		d.QueueId = queue
	}
	database, found, issue := stringFieldOf(fields, "database")
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		d.DatabaseId = database
	}
	description, found, issue := stringFieldOf(fields, "description")
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		d.Description = description
	}
	if d.QueueId == "" && d.DatabaseId == "" {
		issues = append(issues, *NodeError("A dataStore must have either a database or a queue", node))
	}
	issue = setDataFlow(fields, d)
	if issue != nil {
		issues = append(issues, *issue)
	}
	return issues
}

type Service struct {
	node           *yaml.Node
	Id             string
	Name           string
	DataStores     []*DataStoreUse
	Forms          []string
	TechnologyIds  []string
	TechnologiesId string
	Calls          []*Call
}

func (s *Service) setTechnologyBundleId(technologyBundle string) {
	s.TechnologiesId = technologyBundle
}

func (s *Service) setTechnologyIds(technologies []string) {
	s.TechnologyIds = technologies
}

func (s *Service) setId(id string) {
	s.Id = id
}

func (s *Service) setNode(node *yaml.Node) {
	s.node = node
}

func (s *Service) setName(name string) {
	s.Name = name
}

func (s *Service) read(id string, node *yaml.Node, issues []Issue) []Issue {
	var fields map[string]*yaml.Node
	fields, issues = namedObject(id, node, s, issues)
	datastoreNodes, _, issue := sequenceFieldOf(fields, "dataStores")
	if issue != nil {
		return append(issues, *issue)
	}
	dataStores := make([]*DataStoreUse, 0)
	for _, dataStoreNode := range datastoreNodes {
		dataStore := DataStoreUse{}
		issues = dataStore.read(dataStoreNode, issues)
		dataStores = append(dataStores, &dataStore)
	}
	s.DataStores = dataStores
	formNodes, _, issue := sequenceFieldOf(fields, "forms")
	if issue != nil {
		return append(issues, *issue)
	}
	forms := make([]string, 0)
	for _, formNode := range formNodes {
		form, issue := toString(formNode, "form")
		if issue == nil {
			forms = append(forms, form)
		} else {
			issues = append(issues, *issue)
		}
	}
	s.Forms = forms
	issues = append(issues, setTechnologies(fields, s)...)
	return issues
}

type ServiceReader struct {
}

func (s ServiceReader) read(node *yaml.Node, _ string, model *ArchitectureModel) []Issue {
	if node == nil {
		return []Issue{}
	}
	servicesById, issue := toMap(node)
	if issue != nil {
		return []Issue{*issue}
	}
	issues := make([]Issue, 0)
	services := make([]*Service, 0)
	for id, serviceNode := range servicesById {
		service := Service{}
		issues = service.read(id, serviceNode, issues)
		services = append(services, &service)
	}
	sort.Slice(services, func(i, j int) bool {
		return services[i].Name < services[j].Name
	})
	model.Services = services
	return issues
}
