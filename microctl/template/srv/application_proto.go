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

package {{.Service}};

option go_package = "{{.Domain}}/{{.Project}}/{{.Service}}/application/{{.Service}}/proto;{{.Service}}";

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

    rpc Get{{.Svc}}List0 (Get{{.Svc}}List0Request) returns (Get{{.Svc}}ListResponse) {
        option (google.api.http) = {
           get: "/v1/{{.Service}}/list/0/{type}"
        };
    }

    rpc Get{{.Svc}}List1 (Get{{.Svc}}List1Request) returns (Get{{.Svc}}ListResponse) {
        option (google.api.http) = {
            get: "/v1/{{.Service}}/list/1/{type}"
        };
    }

    rpc Update{{.Svc}}Type (Update{{.Svc}}TypeRequest) returns (Update{{.Svc}}TypeResponse) {
        option (google.api.http) = {
           post: "/v1/{{.Service}}/type"
           body: "*"
        };
    }

    rpc Delete{{.Svc}}ById (Delete{{.Svc}}ByIdRequest) returns (Delete{{.Svc}}ByIdResponse) {
        option (google.api.http) = {
           post: "/v1/{{.Service}}/del"
           body: "*"
        };
    }

    rpc Get{{.Svc}}ListByStream (stream Get{{.Svc}}ListByStreamRequest) returns (stream Get{{.Svc}}ListByStreamResponse);
}
message Create{{.Svc}}Request {
    // @inject_tag: validate:"required,email"
    string name = 1;
    // @inject_tag: validate:"gte=0,lte=3"
    int32 type = 2;
}

// @inject_response Create{{.Svc}}Response
message Create{{.Svc}}Response {
    int32 code = 1;
    string msg = 2;
}

message Get{{.Svc}}ByIdRequest {
    int32 id = 1;
}

// @inject_response Get{{.Svc}}ByIdResponse *{{.Svc}} data
message Get{{.Svc}}ByIdResponse {
    int32 code = 1;
    string msg = 2;
    {{.Svc}} data = 3;
}

message Get{{.Svc}}List0Request {
    // @inject_tag: validate:"gte=0,lte=3"
    int32 type = 1;
    // @inject_tag: validate:"gte=5,lte=20"
    int32 page_size = 2;
    int32 page_num = 3;
    bool is_desc = 4;
}

message Get{{.Svc}}List1Request {
    // @inject_tag: validate:"gte=0,lte=3"
    int32 type = 1;
    // @inject_tag: validate:"gte=5,lte=20"
    int32 page_size = 2;
    int32 last_id = 3;
    bool is_desc = 4;
}

// @inject_response Get{{.Svc}}ListResponse *{{.Svc}}List data
message Get{{.Svc}}ListResponse {
    int32 code = 1;
    string msg = 2;
    {{.Svc}}List data = 3;
}

message Update{{.Svc}}TypeRequest {
    int32 id = 1;
    int32 type = 2;
}

// @inject_response Update{{.Svc}}TypeResponse
message Update{{.Svc}}TypeResponse {
    int32 code = 1;
    string msg = 2;
}

message Delete{{.Svc}}ByIdRequest {
    int32 id = 1;
}

// @inject_response Delete{{.Svc}}ByIdResponse
message Delete{{.Svc}}ByIdResponse {
    int32 code = 1;
    string msg = 2;
}

message {{.Svc}} {
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

message {{.Svc}}List {
    int32 total = 1;
    int32 total_page = 2;
    int32 cur_page = 3;
    repeated {{.Svc}} datalist = 4;
}

message Get{{.Svc}}ListByStreamRequest {
    int32 index = 1;
    int32 id = 2;
}

message Get{{.Svc}}ListByStreamResponse {
    int32 index = 1;
    {{.Svc}} result = 2;
}
`

    path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/application/" + data.Service + "/proto/"
    name := data.Service + ".proto"

    return template.CreateFile(data, tpl, path, name)
}
