/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package api

import (
	"github.com/imind-lab/micro/v2/microctl/template"
)

// 生成client/service.go
func CreateApplicationProto(data *template.Data) error {
	var tpl = `syntax = "proto3";

package {{.Package}}_api;

option go_package = "{{.Domain}}/{{.Repo}}{{.Suffix}}/application/sample/proto;{{.Package}}_api";

import "google/api/annotations.proto";

service {{.Service}}Service {
    rpc Create{{.Service}} (Create{{.Service}}Request) returns (Create{{.Service}}Response) {
        option (google.api.http) = {
           post: "/v1/sample/create"
           body: "*"
        };
    }

    rpc Get{{.Service}}ById (Get{{.Service}}ByIdRequest) returns (Get{{.Service}}ByIdResponse) {
        option (google.api.http) = {
           get: "/v1/sample/one/{id}"
        };
    }

    rpc Get{{.Service}}List0 (Get{{.Service}}List0Request) returns (Get{{.Service}}ListResponse) {
        option (google.api.http) = {
            get: "/v1/sample/list/0/{type}"
        };
    }

    rpc Get{{.Service}}List1 (Get{{.Service}}List1Request) returns (Get{{.Service}}ListResponse) {
        option (google.api.http) = {
            get: "/v1/sample/list/1/{type}"
        };
    }

    rpc Update{{.Service}}Type (Update{{.Service}}TypeRequest) returns (Update{{.Service}}TypeResponse) {
        option (google.api.http) = {
           post: "/v1/sample/type"
           body: "*"
        };
    }

    rpc Delete{{.Service}}ById (Delete{{.Service}}ByIdRequest) returns (Delete{{.Service}}ByIdResponse) {
        option (google.api.http) = {
           post: "/v1/sample/del"
           body: "*"
        };
    }

    rpc Get{{.Service}}ListByIds (Get{{.Service}}ListByIdsRequest) returns (Get{{.Service}}ListByIdsResponse) {
        option (google.api.http) = {
           post: "/v1/sample/ids"
           body: "*"
        };
    }
}

message Create{{.Service}}Request {
    // @inject_tag: validate:"required,email"
    string name = 1;
    // @inject_tag: validate:"gte=0,lte=3"
    int32 type = 2;
}

// @inject_response Create{{.Service}}Response
message Create{{.Service}}Response {
    int32 code = 1;
    string msg = 2;
}

message Get{{.Service}}ByIdRequest {
    int32 id = 1;
}

// @inject_response Get{{.Service}}ByIdResponse *{{.Service}} data
message Get{{.Service}}ByIdResponse {
    int32 code = 1;
    string msg = 2;
    {{.Service}} data = 3;
}

message Get{{.Service}}List0Request {
    // @inject_tag: validate:"gte=0,lte=3"
    int32 type = 1;
    // @inject_tag: validate:"gte=5,lte=20"
    int32 page_size = 2;
    int32 page_num = 3;
    bool is_desc = 4;
}

message Get{{.Service}}List1Request {
    // @inject_tag: validate:"gte=0,lte=3"
    int32 type = 1;
    // @inject_tag: validate:"gte=5,lte=20"
    int32 page_size = 2;
    int32 last_id = 3;
    bool is_desc = 4;
}

// @inject_response Get{{.Service}}ListResponse *{{.Service}}List data
message Get{{.Service}}ListResponse {
    int32 code = 1;
    string msg = 2;
    {{.Service}}List data = 3;
}

message Update{{.Service}}TypeRequest {
    int32 id = 1;
    int32 type = 2;
}

// @inject_response Update{{.Service}}TypeResponse
message Update{{.Service}}TypeResponse {
    int32 code = 1;
    string msg = 2;
}

message Delete{{.Service}}ByIdRequest {
    int32 id = 1;
}

// @inject_response Delete{{.Service}}ByIdResponse
message Delete{{.Service}}ByIdResponse {
    int32 code = 1;
    string msg = 2;
}

message {{.Service}} {
    int32 id = 1;
    // @inject_tag: validate:"required,email"
    string name = 2;
    int32 view_num = 3;
    // @inject_tag: validate:"gte=0,lte=3"
    int32 type = 4;
    uint32 create_time = 5;
    string create_datetime = 6;
    string update_datetime = 7;
}

message {{.Service}}List {
    int32 total = 1;
    int32 total_page = 2;
    int32 cur_page = 3;
    repeated {{.Service}} datalist = 4;
}

message Get{{.Service}}ListByIdsRequest {
    repeated int32 ids = 1;
}

// @inject_response Get{{.Service}}ListByIdsResponse []*{{.Service}} data
message Get{{.Service}}ListByIdsResponse {
    int32 code = 1;
    string msg = 2;
    repeated {{.Service}} data = 3;
}
`

	path := "./" + data.Name + "-api/application/" + data.Name + "/proto/"
	name := data.Package + "_api.proto"

	return template.CreateFile(data, tpl, path, name)
}
