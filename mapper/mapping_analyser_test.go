package mapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetMapping(t *testing.T) {
	type args struct {
		es_url string
	}

	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Basic Test",
			args: args{
				"10.41.4.22:9200",
			},
			want: map[string]interface{}{
				"abc": "cde",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllMappings(tt.args.es_url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMapping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMapping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAliases(t *testing.T) {
	type args struct {
		es_url string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Base Case",
			args: args{
				es_url: "10.41.5.5",
			},
			want: map[string]string{
				"asdsa": "Asdad",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAliases(tt.args.es_url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAliases() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllMappings(t *testing.T) {
	type args struct {
		es_url string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "TestGetAllMappings Base Case",
			args: args{
				es_url: "10.41.5.5",
			},
			want:    map[string]interface{}{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllMappings(tt.args.es_url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMappings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllMappings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExplodeInnerFields(t *testing.T) {
	mappingsInput := []byte(`{
		"shop" : {
			"mappings" : {
				"properties" : {
					"name" : {
						"type" : "keyword"
					},
					"attr" : {
						"properties" : {
							"loc" : {
								"type" : "text"
							},
							"category" : {
								"type" : "keyword","doc_values": false
							},
							"loc1" : {
								"type" : "text"
							},
							"category1" : {
								"type" : "keyword"
							}
						}
					}
				}
			}
		},
		"product_shop" : {
			"mappings" :{
				"_routing" : "true",
				"properties" :{
				"a" : {
					"properties" : {
						"loc" : {
							"type" : "text"
						},
						"category" : {
							"type" : "keyword"
						},
						"loc1" : {
							"type" : "text"
						},
						"category1" : {
							"type" : "keyword"
						}
					}
				}
			}
			}
		}
	}`)
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, mappingsInput); err != nil {
		fmt.Println(err)
	}
	var response map[string]Mapping
	err := json.Unmarshal([]byte(mappingsInput), &response)

	if err != nil {
		log.Printf("Error %v", err)
		t.Fail()
	}

	k := ExplodeAllIndices(response)

	for a, b := range k {
		fmt.Printf("%v\n", a)
		b.Print()
		fmt.Printf("\n\n")
	}

}
