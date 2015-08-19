package underscore

import (
	"reflect"
)

var EMPTY_GROUP = make(map[interface{}][]interface{})

func Group(source interface{}, keySelector func(interface{}) (interface{}, error)) (map[interface{}][]interface{}, error) {
	if source == nil {
		return EMPTY_GROUP, nil
	}

	sourceRV := reflect.ValueOf(source)
	switch sourceRV.Kind() {
		case reflect.Array:
		case reflect.Slice:
			if sourceRV.Len() == 0 {
				return EMPTY_GROUP, nil
			}

			dict := make(map[interface{}][]interface{})
			for i := 0; i < sourceRV.Len(); i++ {
				value := sourceRV.Index(i).Interface()
				key, err := keySelector(value)
				if err != nil {
					return EMPTY_GROUP, err
				}

				dict[key] = append(dict[key], value)
			}
			return dict, nil
		case reflect.Map:
			oldKeyRVs := sourceRV.MapKeys()
			if len(oldKeyRVs) == 0 {
				return EMPTY_GROUP, nil
			}

			dict := make(map[interface{}][]interface{})
			for i := 0; i < len(oldKeyRVs); i++ {
				value := sourceRV.MapIndex(oldKeyRVs[i]).Interface()
				key, err := keySelector(value)
				if err != nil {
					return EMPTY_GROUP, err
				}

				dict[key] = append(dict[key], value)
			}
			return dict, nil
	}
	return EMPTY_GROUP, nil
}

func GroupBy(source interface{}, field string) (map[interface{}][]interface{}, error) {
	return Group(source, func (item interface{}) (interface{}, error) {
		return getFieldValue(item, field)
	})
}

//Chain
func (this *Query) Group(keySelector func(item interface{}) (interface{}, error)) Queryer {
	if this.err != nil {
		return this
	}

	this.source, this.err = Group(this.source, keySelector)
	return this
}

func (this *Query) GroupBy(field string) Queryer {
	if this.err != nil {
		return this
	}

	this.source, this.err = GroupBy(this.source, field)
	return this
}