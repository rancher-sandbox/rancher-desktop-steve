package schema

import (
	"fmt"
	"net/http"

	"github.com/rancher/naok/pkg/accesscontrol"
	"github.com/rancher/naok/pkg/attributes"
	"github.com/rancher/norman/pkg/api/builtin"
	"github.com/rancher/norman/pkg/types"
	"k8s.io/apiserver/pkg/authentication/user"
)

func newSchemas() (*types.Schemas, error) {
	s, err := types.NewSchemas(builtin.Schemas)
	if err != nil {
		return nil, err
	}
	s.DefaultMapper = func() types.Mapper {
		return newDefaultMapper()
	}

	return s, nil
}

func (c *Collection) Schemas(user user.Info) (*types.Schemas, error) {
	access := c.as.AccessFor(user)
	return c.schemasForSubject("", access)
}

func (c *Collection) schemasForSubject(subjectKey string, access *accesscontrol.AccessSet) (*types.Schemas, error) {
	result, err := newSchemas()
	if err != nil {
		return nil, err
	}

	if _, err := result.AddSchemas(c.baseSchema); err != nil {
		return nil, err
	}

	for _, s := range c.schemas {
		gr := attributes.GR(s)

		if gr.Resource == "" {
			if err := result.AddSchema(*s); err != nil {
				return nil, err
			}
			continue
		}

		verbs := attributes.Verbs(s)
		verbAccess := accesscontrol.AccessListMap{}

		for _, verb := range verbs {
			a := access.AccessListFor(verb, gr)
			if len(a) > 0 {
				verbAccess[verb] = a
			}
		}

		if len(verbAccess) == 0 {
			continue
		}

		s = s.DeepCopy()
		attributes.SetAccess(s, verbAccess)
		if verbAccess.AnyVerb("list", "get") {
			s.ResourceMethods = append(s.ResourceMethods, http.MethodGet)
			s.CollectionMethods = append(s.CollectionMethods, http.MethodGet)
		}
		if verbAccess.AnyVerb("delete") {
			s.ResourceMethods = append(s.ResourceMethods, http.MethodDelete)
		}
		if verbAccess.AnyVerb("update") {
			s.ResourceMethods = append(s.ResourceMethods, http.MethodPut)
		}
		if verbAccess.AnyVerb("create") {
			s.CollectionMethods = append(s.CollectionMethods, http.MethodPost)
		}

		c.applyTemplates(s)

		if err := result.AddSchema(*s); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c *Collection) applyTemplates(schema *types.Schema) {
	templates := []*Template{
		c.templates[schema.ID],
		c.templates[fmt.Sprintf("%s/%s", attributes.Group(schema), attributes.Kind(schema))],
		c.templates[""],
	}

	for _, t := range templates {
		if t == nil {
			continue
		}
		if schema.Mapper == nil {
			schema.Mapper = t.Mapper
		}
		if schema.Formatter == nil {
			schema.Formatter = t.Formatter
		}
		if schema.Store == nil {
			schema.Store = t.Store
		}
	}
}
