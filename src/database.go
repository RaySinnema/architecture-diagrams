package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Database struct {
	DataStore
	Views []*View
}

func (d *Database) hasView(viewId string) bool {
	for _, candidate := range d.Views {
		if candidate.Id == viewId {
			return true
		}
	}
	return false
}

func (d *Database) readViews(fields map[string]*yaml.Node) []Issue {
	issues := make([]Issue, 0)
	viewNodes, _, issue := sequenceFieldOf(fields, "views")
	if issue == nil {
		views := make([]*View, 0)
		for _, viewNode := range viewNodes {
			name, issue := toString(viewNode, "view")
			if issue == nil {
				view := View{}
				views = append(views, &view)
				view.node = viewNode
				view.Id = name
				view.Name = name
				view.State = Ok
				view.ImplementedBy = d
			} else {
				issues = append(issues, *issue)
			}
		}
		d.Views = views
	} else {
		issues = append(issues, *issue)
	}
	return issues
}

type DatabaseReader struct {
}

func (d DatabaseReader) read(node *yaml.Node, _ string, model *ArchitectureModel) []Issue {
	databases := make([]*Database, 0)
	dataStores, issues := DataStoreReader{}.read(node)
	for _, dataStore := range dataStores {
		database := Database{DataStore: *dataStore}
		databases = append(databases, &database)
		issues = append(issues, d.readDatabase(database.node, &database)...)
	}
	model.Databases = databases
	return issues
}

func (d DatabaseReader) readDatabase(node *yaml.Node, database *Database) []Issue {
	fields, issue := toMap(node)
	if issue != nil {
		return []Issue{*issue}
	}
	return database.readViews(fields)
}

type DatabaseConnector struct {
}

func (d DatabaseConnector) connect(model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	for _, database := range model.Databases {
		issues = append(issues, connectTechnologies(database, model)...)
		issues = append(issues, connectTechnologies(ApiTechnologies{&database.DataStore}, model)...)
	}
	return issues
}

type DatabaseValidator struct {
}

func (d DatabaseValidator) validate(model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	issues = append(issues, d.validateViewsAreUnique(model.Databases)...)
	return issues
}

func (d DatabaseValidator) validateViewsAreUnique(databases []*Database) []Issue {
	issues := make([]Issue, 0)
	views := map[string]string{}
	for _, database := range databases {
		for _, view := range database.Views {
			owner, found := views[view.Id]
			if found {
				issues = append(issues, *NodeError(fmt.Sprintf("View '%v' is already defined in database '%v'",
					view.Id, owner), database.node))
			} else {
				views[view.Id] = database.Id
			}
		}
	}

	return issues
}
