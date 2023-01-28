package main

import (
	"gopkg.in/yaml.v3"
	"sort"
)

type DataStoreUse struct {
	QueueId     string `yaml:"queue,omitempty"`
	DatabaseId  string `yaml:"database,omitempty"`
	Description string
	DataFlow    DataFlow
}

func (d *DataStoreUse) read(node *yaml.Node) []Issue {
	fields, issue := toMap(node)
	if issue != nil {
		return []Issue{*issue}
	}
	issues := make([]Issue, 0)
	issues = append(issues, d.readDataStore(node, fields)...)
	issues = append(issues, setDescription(fields, d)...)
	issues = append(issues, setDataFlow(node, fields, d)...)
	return issues
}

func (d *DataStoreUse) readDataStore(node *yaml.Node, fields map[string]*yaml.Node) []Issue {
	issues := make([]Issue, 0)
	database, found, issue := stringFieldOf(fields, "database")
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		d.DatabaseId = database
	}
	queue, found, issue := stringFieldOf(fields, "queue")
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		d.QueueId = queue
	}
	if d.QueueId == "" && d.DatabaseId == "" {
		issues = append(issues, *NodeError("A dataStore must be a database or a queue", node))
	} else if d.QueueId != "" && d.DatabaseId != "" {
		issues = append(issues, *NodeError("A dataStore can be either a database or a queue, but not both", node))
	}
	return issues
}

func (d *DataStoreUse) setDescription(description string) {
	d.Description = description
}

func (d *DataStoreUse) setDataFlow(dataFlow DataFlow) {
	d.DataFlow = dataFlow
}

type Form struct {
	node  *yaml.Node
	Id    string
	Name  string
	State State
}

func (f *Form) setNode(node *yaml.Node) {
	f.node = node
}

func (f *Form) setId(id string) {
	f.Id = id
}

func (f *Form) setName(name string) {
	f.Name = name
}

func (f *Form) setState(state State) {
	f.State = state
}

type Service struct {
	node               *yaml.Node
	Id                 string
	Name               string
	DataStores         []*DataStoreUse
	Forms              []*Form
	Calls              []*Call
	TechnologyIds      []string
	TechnologyBundleId string
	State              State
	Technologies       []*Technology
}

func (s *Service) getNode() *yaml.Node {
	return s.node
}

func (s *Service) getTechnologyIds() []string {
	return s.TechnologyIds
}

func (s *Service) getTechnologyBundleId() string {
	return s.TechnologyBundleId
}

func (s *Service) setTechnologies(technologies []*Technology) {
	s.Technologies = technologies
}

func (s *Service) read(id string, node *yaml.Node) []Issue {
	var fields map[string]*yaml.Node
	fields, issues := namedObject(node, id, s)
	issues = append(issues, s.readDataStores(fields)...)
	issues = append(issues, s.readForms(fields)...)
	issues = append(issues, s.readCalls(fields)...)
	issues = append(issues, setTechnologies(fields, s)...)
	issues = append(issues, setState(node, fields, s)...)
	return issues
}

func (s *Service) setNode(node *yaml.Node) {
	s.node = node
}

func (s *Service) setId(id string) {
	s.Id = id
}

func (s *Service) setName(name string) {
	s.Name = name
}

func (s *Service) readDataStores(fields map[string]*yaml.Node) []Issue {
	nodes, _, issue := sequenceFieldOf(fields, "dataStores")
	if issue != nil {
		return []Issue{*issue}
	}
	issues := make([]Issue, 0)
	dataStores := make([]*DataStoreUse, 0)
	for _, node := range nodes {
		dataStore := DataStoreUse{}
		dataStores = append(dataStores, &dataStore)
		issues = append(issues, dataStore.read(node)...)
	}
	s.DataStores = dataStores
	return issues
}

func (s *Service) readCalls(fields map[string]*yaml.Node) []Issue {
	callNodes, _, issue := sequenceFieldOf(fields, "calls")
	if issue != nil {
		return []Issue{*issue}
	}
	issues := make([]Issue, 0)
	calls := make([]*Call, 0)
	for _, callNode := range callNodes {
		call := Call{}
		calls = append(calls, &call)
		issues = append(issues, call.read(callNode)...)
	}
	s.Calls = calls
	return issues
}

func (s *Service) readForms(fields map[string]*yaml.Node) []Issue {
	issues := make([]Issue, 0)
	formMaps, found, issue := mapFieldOf(fields, "forms")
	if found && issue == nil {
		forms := make([]*Form, 0)
		for formId, formNode := range formMaps {
			form := Form{}
			forms = append(forms, &form)
			formFields, formIssues := namedObject(formNode, formId, &form)
			issues = append(issues, formIssues...)
			issues = append(issues, setState(formNode, formFields, &form)...)
		}
		s.Forms = forms
	} else if found {
		formNodes, _, issue := sequenceFieldOf(fields, "forms")
		if issue == nil {
			forms := make([]*Form, 0)
			for _, formNode := range formNodes {
				name, issue := toString(formNode, "form")
				if issue == nil {
					form := Form{}
					forms = append(forms, &form)
					form.node = formNode
					form.Id = name
					form.Name = name
					form.State = Ok
				} else {
					issues = append(issues, *issue)
				}
			}
			s.Forms = forms
		} else {
			issues = append(issues, *issue)
		}
	}
	return issues
}

func (s *Service) setTechnologyBundleId(technologyBundle string) {
	s.TechnologyBundleId = technologyBundle
}

func (s *Service) setTechnologyIds(technologies []string) {
	s.TechnologyIds = technologies
}

func (s *Service) setState(state State) {
	s.State = state
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
		services = append(services, &service)
		issues = append(issues, service.read(id, serviceNode)...)
	}
	sort.Slice(services, func(i, j int) bool {
		return services[i].Name < services[j].Name
	})
	model.Services = services
	return issues
}

type ServiceConnector struct {
}

func (s ServiceConnector) connect(model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	for _, service := range model.Services {
		issues = append(issues, connectTechnologies(service, model)...)
		for _, call := range service.Calls {
			issues = append(issues, connectTechnologies(call, model)...)
		}
	}
	return issues
}
