package interpreter

import (
	"os"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/goland/builtin"
)

type SoqlExecutor struct {
}

func (e *SoqlExecutor) Execute(soql string) (*builtin.Object, error) {
	client := soapforce.NewClient()
	username := os.Getenv("SALESFORCE_USERNAME")
	password := os.Getenv("SALESFORCE_PASSWORD")
	endpoint := os.Getenv("SALESFORCE_ENDPOINT")
	client.SetLoginUrl(endpoint)
	client.Login(username, password)
	result, err := client.Query(soql)
	if err != nil {
		return nil, err
	}
	return e.getListFromResponse(result.Records)
}

func (e *SoqlExecutor) getListFromResponse(records []*soapforce.SObject) (*builtin.Object, error) {
	objects := make([]*builtin.Object, len(records))
	classType := builtin.AccountType
	for i, r := range records {
		object := &builtin.Object{}
		object.ClassType = classType
		object.InstanceFields = builtin.NewObjectMap()
		object.InstanceFields.Set("id", newString(r.Id))
		for k, v := range r.Fields {
			switch val := v.(type) {
			case string:
				object.InstanceFields.Set(k, newString(val))
			}
		}
		objects[i] = object
	}
	// TODO: implement
	list := &builtin.Object{
		ClassType:      builtin.ListType,
		InstanceFields: builtin.NewObjectMap(),
		GenericType:    []*builtin.ClassType{classType},
		Extra: map[string]interface{}{
			"records": objects,
		},
	}
	return list, nil
}