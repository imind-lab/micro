/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package srv

import (
	"github.com/imind-lab/micro/v2/microctl/template"
)

// 生成client/service.go
func CreateApplicationProto(data *template.Data) error {
	var tpl = `syntax = "proto3";

package {{.Package}};

option go_package = "{{.Domain}}/{{.Repo}}/application/{{.Name}}/proto;{{.Package}}";

import "google/api/annotations.proto";

service {{.Service}}Service {
    rpc Create{{.Service}} (Create{{.Service}}Request) returns (Create{{.Service}}Response) {
        option (google.api.http) = {
           post: "/v1/{{.Name}}/create"
           body: "*"
        };
    }

    rpc Get{{.Service}}ById (Get{{.Service}}ByIdRequest) returns (Get{{.Service}}ByIdResponse) {
        option (google.api.http) = {
           get: "/v1/{{.Name}}/one/{id}"
        };
    }

    rpc Get{{.Service}}List0 (Get{{.Service}}List0Request) returns (Get{{.Service}}ListResponse) {
        option (google.api.http) = {
           get: "/v1/{{.Name}}/list/0/{type}"
        };
    }

    rpc Get{{.Service}}List1 (Get{{.Service}}List1Request) returns (Get{{.Service}}ListResponse) {
        option (google.api.http) = {
            get: "/v1/{{.Name}}/list/1/{type}"
        };
    }

    rpc Update{{.Service}}Type (Update{{.Service}}TypeRequest) returns (Update{{.Service}}TypeResponse) {
        option (google.api.http) = {
           post: "/v1/{{.Name}}/type"
           body: "*"
        };
    }

    rpc Delete{{.Service}}ById (Delete{{.Service}}ByIdRequest) returns (Delete{{.Service}}ByIdResponse) {
        option (google.api.http) = {
           post: "/v1/{{.Name}}/del"
           body: "*"
        };
    }

    rpc Get{{.Service}}ListByStream (stream Get{{.Service}}ListByStreamRequest) returns (stream Get{{.Service}}ListByStreamResponse);
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

message Get{{.Service}}ListByStreamRequest {
    int32 index = 1;
    int32 id = 2;
}

message Get{{.Service}}ListByStreamResponse {
    int32 index = 1;
    {{.Service}} result = 2;
}
`

	path := "./" + data.Name + "/application/" + data.Name + "/proto/"
	name := data.Package + ".proto"

	return template.CreateFile(data, tpl, path, name)
}
