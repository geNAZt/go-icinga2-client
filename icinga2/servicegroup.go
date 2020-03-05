package icinga2

import (
	"fmt"
	"net/url"
)

type ServiceGroup struct {
	DisplayName string `json:"display_name"`
	Package     string `json:"package"`
	Notes       string `json:"notes"`
	NotesURL    string `json:"notes_url"`
	Vars        Vars   `json:"vars"`
	Zone        string `json:"zone,omitempty"`
}

func (group *ServiceGroup) FullName() string {
	return group.Name
}

type ServiceGroupResults struct {
	Results []struct {
		ServiceGroup ServiceGroup `json:"attrs"`
	} `json:"results"`
}

type ServiceGroupCreate struct {
	Attrs ServiceGroup `json:"attrs"`
}

func (s *WebClient) GetServiceGroup(name string) (ServiceGroup, error) {
	var serviceGroupResults ServiceGroupResults
	resp, err := s.napping.Get(s.URL+"/v1/objects/servicegroups/"+name, nil, &serviceGroupResults, nil)
	if err != nil {
		return ServiceGroup{}, err
	}

	if resp.HttpResponse().StatusCode != 200 {
		return ServiceGroup{}, fmt.Errorf("Did not get 200 OK")
	}

	return serviceGroupResults.Results[0].ServiceGroup, nil
}

func (s *WebClient) CreateServiceGroup(serviceGroup ServiceGroup) error {
	serviceCreate := ServiceGroupCreate{Attrs: serviceGroup}
	err := s.CreateObject("/servicegroups/"+serviceGroup.FullName(), serviceCreate)
	return err
}

func (s *WebClient) ListServiceGroups() ([]ServiceGroup, error) {
	var serviceResults ServiceGroupResults
	var serviceGroups []ServiceGroup

	_, err := s.napping.Get(s.URL+"/v1/objects/services/", nil, &serviceResults, nil)
	if err != nil {
		return nil, err
	}

	for _, result := range serviceResults.Results {
		if s.Zone == "" || s.Zone == result.ServiceGroup.Zone {
			serviceGroups = append(serviceGroups, result.ServiceGroup)
		}
	}

	return serviceGroups, nil
}

func (s *WebClient) DeleteServiceGroup(name string) error {
	_, err := s.napping.Delete(s.URL+"/v1/objects/servicegroups/"+name, &url.Values{"cascade": []string{"1"}}, nil, nil)
	return err
}

func (s *WebClient) UpdateServiceGroup(serviceGroup ServiceGroup) error {
	serviceUpdate := ServiceGroupCreate{Attrs: serviceGroup}

	err := s.UpdateObject("/services/"+serviceGroup.FullName(), serviceUpdate)
	return err
}
