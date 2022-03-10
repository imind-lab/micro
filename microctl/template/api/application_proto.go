/**
 *  MindLab
 *
 *  Create by songli on {{.Year}}/02/27
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package api

import (
	"os"
	"text/template"

	tpl "github.com/imind-lab/micro/microctl/template"
)

// 生成client/service.go
func CreateApplicationProto(data *tpl.Data) error {
	var tpl = `syntax = "proto3";

package {{.Service}}_api;

option go_package = "{{.Domain}}/{{.Project}}/{{.Service}}-api/application/{{.Service}}/proto;{{.Service}}_api";

import "google/api/annotations.proto";

service {{.Svc}}Service {
    rpc Create{{.Svc}} (Create{{.Svc}}Request) returns (Create{{.Svc}}Response) {
        option (google.api.http) = {
           post: "/v1/{{.Service}}/create"
           body: "*"
        };
    }
    rpc Get{{.Svc}}ById (Get{{.Svc}}ByIdRequest) returns (Get{{.Svc}}ByIdResponse) {
        option (google.api.http) = {
           get: "/v1/{{.Service}}/one/{id}"
        };
    }
    rpc Get{{.Svc}}List (Get{{.Svc}}ListRequest) returns (Get{{.Svc}}ListResponse) {
        option (google.api.http) = {
           get: "/v1/{{.Service}}/list/{status}"
        };
    }
    rpc Update{{.Svc}}Status (Update{{.Svc}}StatusRequest) returns (Update{{.Svc}}StatusResponse) {
        option (google.api.http) = {
           post: "/v1/{{.Service}}/status"
           body: "*"
        };
    }
    rpc Delete{{.Svc}}ById (Delete{{.Svc}}ByIdRequest) returns (Delete{{.Svc}}ByIdResponse) {
        option (google.api.http) = {
           post: "/v1/{{.Service}}/del"
           body: "*"
        };
    }

    rpc Get{{.Svc}}ListByIds (Get{{.Svc}}ListByIdsRequest) returns (Get{{.Svc}}ListByIdsResponse) {
        option (google.api.http) = {
           post: "/v1/{{.Service}}/ids"
           body: "*"
        };
    }
}

message Create{{.Svc}}Request {
    // @inject_tag: validate:"required,email"
    string name = 1;
    // @inject_tag: validate:"gte=0,lte=3"
    int32 status = 2;
}

// @inject_response Create{{.Svc}}Response
message Create{{.Svc}}Response {
    int32 code = 1;
    string message = 2;
}

message Get{{.Svc}}ByIdRequest {
    int32 id = 1;
}

// @inject_response Get{{.Svc}}ByIdResponse *{{.Svc}} data
message Get{{.Svc}}ByIdResponse {
    int32 code = 1;
    string message = 2;
    {{.Svc}} data = 3;
}

message Get{{.Svc}}ListRequest {
    // @inject_tag: validate:"gte=0,lte=3"
    int32 status = 1;
    int32 lastid = 2;
    // @inject_tag: validate:"gte=5,lte=20"
    int32 page_size = 3;
    int32 page_num = 4;
}

// @inject_response Get{{.Svc}}ListResponse *{{.Svc}}List data
message Get{{.Svc}}ListResponse {
    int32 code = 1;
    string message = 2;
    {{.Svc}}List data = 3;
}

message Update{{.Svc}}StatusRequest {
    int32 id = 1;
    int32 status = 2;
}

// @inject_response Update{{.Svc}}StatusResponse
message Update{{.Svc}}StatusResponse {
    int32 code = 1;
    string message = 2;
}

message Delete{{.Svc}}ByIdRequest {
    int32 id = 1;
}

// @inject_response Delete{{.Svc}}ByIdResponse
message Delete{{.Svc}}ByIdResponse {
    int32 code = 1;
    string message = 2;
}

message {{.Svc}} {
    int32 id = 1;
    // @inject_tag: validate:"required,email"
    string name = 2;
    int32 view_num = 3;
    // @inject_tag: validate:"gte=0,lte=3"
    int32 status = 4;
    uint32 create_time = 5;
    string create_datetime = 6;
    string update_datetime = 7;
}

message {{.Svc}}List {
    int32 total = 1;
    int32 total_page = 2;
    int32 cur_page = 3;
    repeated {{.Svc}} datalist = 4;
}

message Get{{.Svc}}ListByIdsRequest {
    repeated int32 ids = 1;
}

// @inject_response Get{{.Svc}}ListByIdsResponse []*{{.Svc}} data
message Get{{.Svc}}ListByIdsResponse {
    int32 code = 1;
    string message = 2;
    repeated {{.Svc}} data = 3;
}
`

	t, err := template.New("application_proto").Parse(tpl)
	if err != nil {
		return err
	}

	t.Option()
	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/application/" + data.Service + "/proto/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + data.Service + "-api.proto"

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	return nil
}
