package main

import (
	"log"
	"os"
	"reflect"
	"strings"
)

func doExamine(v reflect.Value, prefix string, funcs map[string]bool) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		f := t.Field(i)
		funcs[f.Name] = !fv.IsZero()
		if !fv.IsZero() && strings.HasPrefix(f.Name, "tp_as_") {
			doExamine(fv.Elem(), f.Name+".", funcs)
		}
	}
}

func examine(value any) map[string]bool {
	v := reflect.ValueOf(value)
	funcs := map[string]bool{}
	doExamine(v, "", funcs)
	return funcs
}

func generate(name string, funcs map[string]bool) {
	if err := code.Execute(os.Stdout, map[string]any{
		"type":  name,
		"ltype": strings.ToLower(name[:1]) + name[1:],
		"name":  strings.ToLower(name[:1]),
		"funcs": funcs,
	}); err != nil {
		log.Fatalf("Failed to generate template: %s", err)
	}
}

func main() {
	// funcMap := map[string]map[string]bool{}
	for name, t := range types {
		// funcMap[name] = examine(t)
		generate(name, examine(t))
	}

	// output, err := json.Marshal(funcMap)
	// if err != nil {
	// 	log.Fatalf("Failed to encode: %s", err)
	// }
	// fmt.Printf("%s\n", output)
}
