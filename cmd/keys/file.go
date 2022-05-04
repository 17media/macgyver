package keys

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"os"
)

func parseYaml(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = parseYaml(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = parseYaml(v)
		}
	}
	return i
}

func convertYamlToMap(data []byte) map[string]interface{} {
	var body interface{}
	if err := yaml.Unmarshal(data, &body); err != nil {
		panic(err)
	}
	body = parseYaml(body)
	return body.(map[string]interface{})
}

func operationInMap(v interface{}, f func(path string) (interface{}, error)) (interface{}, error) {
	merge := func(strMap map[string]interface{}, k string, v interface{}) (bool, error) {
		switch v.(type) {
		case map[string]interface{}, map[interface{}]interface{}:
		default:
			return false, nil
		}

		k2, err := operationInMap(k, f)
		if err != nil {
			return false, err
		}
		switch yamlOrJson := k2.(type) {
		case string:
			if yamlOrJson == k {
				break
			}
			m := map[string]interface{}{}
			if err := yaml.Unmarshal([]byte(yamlOrJson), &m); err != nil {
				return false, err
			}
			for mk, mv := range m {
				strMap[mk] = mv
			}
			return true, nil
		}
		return false, nil
	}

	var castedValue interface{}
	switch typedValue := v.(type) {
	case string:
		return f(typedValue)
	case map[interface{}]interface{}:
		strMap := map[string]interface{}{}
		for k, v := range typedValue {
			strMap[fmt.Sprintf("%v", k)] = v
		}
		extends := map[string]interface{}{}
		var deleted []string
		for k, v := range strMap {
			ok, err := merge(extends, k, v)
			if ok {
				deleted = append(deleted, k)
				continue
			}
			if err != nil {
				return nil, err
			}

			v2, err := operationInMap(v, f)
			if err != nil {
				return nil, err
			}
			strMap[k] = v2
		}
		for _, k := range deleted {
			delete(strMap, k)
		}
		for k, v := range extends {
			strMap[k] = v
		}
		return strMap, nil

	case map[string]interface{}:
		extends := map[string]interface{}{}
		var deleted []string
		for k, v := range typedValue {
			ok, err := merge(extends, k, v)
			if ok {
				deleted = append(deleted, k)
				continue
			}

			v2, err := operationInMap(v, f)
			if err != nil {
				return nil, err
			}
			typedValue[k] = v2
		}
		for _, k := range deleted {
			delete(typedValue, k)
		}
		for k, v := range extends {
			typedValue[k] = v
		}
		return typedValue, nil
	case []interface{}:
		var a []interface{}
		for i := range typedValue {
			res, err := operationInMap(typedValue[i], f)
			if err != nil {
				return nil, err
			}
			a = append(a, res)
		}
		castedValue = a
	case []string:
		var a []interface{}
		for i := range typedValue {
			res, err := f(typedValue[i])
			if err != nil {
				return nil, err
			}
			a = append(a, res)
		}
		castedValue = a
	default:
		castedValue = typedValue
	}
	return castedValue, nil
}

func RandomCreateFile(data []byte) (string, error) {
	id, _ := uuid.NewUUID()
	fileName := id.String()
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer file.Close()
	if err != nil {
		return "", err
	}
	if _, err := file.Write(data); err != nil {
		return "", err
	}
	return fileName, nil
}

func RemoveFile(name string) error {
	err := os.Remove(name)
	return err
}
