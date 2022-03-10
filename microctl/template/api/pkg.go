/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package api

import (
	"os"
	"strings"
	"text/template"

	tpl "github.com/imind-lab/micro/microctl/template"
)

// 生成google.proto
func CreatePkg(data *tpl.Data) error {

	tpl := `// Copyright (c) 2015, Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

syntax = "proto3";

package google.api;

import "google/api/http.proto";
import "google/protobuf/descriptor.proto";

option go_package = "google.golang.org/genproto/googleapis/api/annotations;annotations";
option java_multiple_files = true;
option java_outer_classname = "AnnotationsProto";
option java_package = "com.google.api";
option objc_class_prefix = "GAPI";

extend google.protobuf.MethodOptions {
  // See ${backtick}HttpRule${backtick}.
  HttpRule http = 72295728;
}
`

	tpl = strings.Replace(tpl, "${backtick}", "`", -1)

	t, err := template.New("google.api.annotations.proto").Parse(tpl)
	if err != nil {
		return err
	}

	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/pkg/proto/google/api/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "annotations.proto"

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	tpl = `// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

syntax = "proto3";

package google.api;

option cc_enable_arenas = true;
option go_package = "google.golang.org/genproto/googleapis/api/annotations;annotations";
option java_multiple_files = true;
option java_outer_classname = "HttpProto";
option java_package = "com.google.api";
option objc_class_prefix = "GAPI";


// Defines the HTTP configuration for an API service. It contains a list of
// [HttpRule][google.api.HttpRule], each specifying the mapping of an RPC method
// to one or more HTTP REST API methods.
message Http {
  // A list of HTTP configuration rules that apply to individual API methods.
  //
  // **NOTE:** All service configuration rules follow "last one wins" order.
  repeated HttpRule rules = 1;

  // When set to true, URL path parmeters will be fully URI-decoded except in
  // cases of single segment matches in reserved expansion, where "%2F" will be
  // left encoded.
  //
  // The default behavior is to not decode RFC 6570 reserved characters in multi
  // segment matches.
  bool fully_decode_reserved_expansion = 2;
}

// ${backtick}HttpRule${backtick} defines the mapping of an RPC method to one or more HTTP
// REST API methods. The mapping specifies how different portions of the RPC
// request message are mapped to URL path, URL query parameters, and
// HTTP request body. The mapping is typically specified as an
// ${backtick}google.api.http${backtick} annotation on the RPC method,
// see "google/api/annotations.proto" for details.
//
// The mapping consists of a field specifying the path template and
// method kind.  The path template can refer to fields in the request
// message, as in the example below which describes a REST GET
// operation on a resource collection of messages:
//
//
//     service Messaging {
//       rpc GetMessage(GetMessageRequest) returns (Message) {
//         option (google.api.http).get = "/v1/messages/{message_id}/{sub.subfield}";
//       }
//     }
//     message GetMessageRequest {
//       message SubMessage {
//         string subfield = 1;
//       }
//       string message_id = 1; // mapped to the URL
//       SubMessage sub = 2;    // ${backtick}sub.subfield${backtick} is url-mapped
//     }
//     message Message {
//       string text = 1; // content of the resource
//     }
//
// The same http annotation can alternatively be expressed inside the
// ${backtick}GRPC API Configuration${backtick} YAML file.
//
//     http:
//       rules:
//         - selector: <proto_package_name>.Messaging.GetMessage
//           get: /v1/messages/{message_id}/{sub.subfield}
//
// This definition enables an automatic, bidrectional mapping of HTTP
// JSON to RPC. Example:
//
// HTTP | RPC
// -----|-----
// ${backtick}GET /v1/messages/123456/foo${backtick}  | ${backtick}GetMessage(message_id: "123456" sub: SubMessage(subfield: "foo"))${backtick}
//
// In general, not only fields but also field paths can be referenced
// from a path pattern. Fields mapped to the path pattern cannot be
// repeated and must have a primitive (non-message) type.
//
// Any fields in the request message which are not bound by the path
// pattern automatically become (optional) HTTP query
// parameters. Assume the following definition of the request message:
//
//
//     service Messaging {
//       rpc GetMessage(GetMessageRequest) returns (Message) {
//         option (google.api.http).get = "/v1/messages/{message_id}";
//       }
//     }
//     message GetMessageRequest {
//       message SubMessage {
//         string subfield = 1;
//       }
//       string message_id = 1; // mapped to the URL
//       int64 revision = 2;    // becomes a parameter
//       SubMessage sub = 3;    // ${backtick}sub.subfield${backtick} becomes a parameter
//     }
//
//
// This enables a HTTP JSON to RPC mapping as below:
//
// HTTP | RPC
// -----|-----
// ${backtick}GET /v1/messages/123456?revision=2&sub.subfield=foo${backtick} | ${backtick}GetMessage(message_id: "123456" revision: 2 sub: SubMessage(subfield: "foo"))${backtick}
//
// Note that fields which are mapped to HTTP parameters must have a
// primitive type or a repeated primitive type. Message types are not
// allowed. In the case of a repeated type, the parameter can be
// repeated in the URL, as in ${backtick}...?param=A&param=B${backtick}.
//
// For HTTP method kinds which allow a request body, the ${backtick}body${backtick} field
// specifies the mapping. Consider a REST update method on the
// message resource collection:
//
//
//     service Messaging {
//       rpc UpdateMessage(UpdateMessageRequest) returns (Message) {
//         option (google.api.http) = {
//           put: "/v1/messages/{message_id}"
//           body: "message"
//         };
//       }
//     }
//     message UpdateMessageRequest {
//       string message_id = 1; // mapped to the URL
//       Message message = 2;   // mapped to the body
//     }
//
//
// The following HTTP JSON to RPC mapping is enabled, where the
// representation of the JSON in the request body is determined by
// protos JSON encoding:
//
// HTTP | RPC
// -----|-----
// ${backtick}PUT /v1/messages/123456 { "text": "Hi!" }${backtick} | ${backtick}UpdateMessage(message_id: "123456" message { text: "Hi!" })${backtick}
//
// The special name ${backtick}*${backtick} can be used in the body mapping to define that
// every field not bound by the path template should be mapped to the
// request body.  This enables the following alternative definition of
// the update method:
//
//     service Messaging {
//       rpc UpdateMessage(Message) returns (Message) {
//         option (google.api.http) = {
//           put: "/v1/messages/{message_id}"
//           body: "*"
//         };
//       }
//     }
//     message Message {
//       string message_id = 1;
//       string text = 2;
//     }
//
//
// The following HTTP JSON to RPC mapping is enabled:
//
// HTTP | RPC
// -----|-----
// ${backtick}PUT /v1/messages/123456 { "text": "Hi!" }${backtick} | ${backtick}UpdateMessage(message_id: "123456" text: "Hi!")${backtick}
//
// Note that when using ${backtick}*${backtick} in the body mapping, it is not possible to
// have HTTP parameters, as all fields not bound by the path end in
// the body. This makes this option more rarely used in practice of
// defining REST APIs. The common usage of ${backtick}*${backtick} is in custom methods
// which don't use the URL at all for transferring data.
//
// It is possible to define multiple HTTP methods for one RPC by using
// the ${backtick}additional_bindings${backtick} option. Example:
//
//     service Messaging {
//       rpc GetMessage(GetMessageRequest) returns (Message) {
//         option (google.api.http) = {
//           get: "/v1/messages/{message_id}"
//           additional_bindings {
//             get: "/v1/users/{user_id}/messages/{message_id}"
//           }
//         };
//       }
//     }
//     message GetMessageRequest {
//       string message_id = 1;
//       string user_id = 2;
//     }
//
//
// This enables the following two alternative HTTP JSON to RPC
// mappings:
//
// HTTP | RPC
// -----|-----
// ${backtick}GET /v1/messages/123456${backtick} | ${backtick}GetMessage(message_id: "123456")${backtick}
// ${backtick}GET /v1/users/me/messages/123456${backtick} | ${backtick}GetMessage(user_id: "me" message_id: "123456")${backtick}
//
// # Rules for HTTP mapping
//
// The rules for mapping HTTP path, query parameters, and body fields
// to the request message are as follows:
//
// 1. The ${backtick}body${backtick} field specifies either ${backtick}*${backtick} or a field path, or is
//    omitted. If omitted, it indicates there is no HTTP request body.
// 2. Leaf fields (recursive expansion of nested messages in the
//    request) can be classified into three types:
//     (a) Matched in the URL template.
//     (b) Covered by body (if body is ${backtick}*${backtick}, everything except (a) fields;
//         else everything under the body field)
//     (c) All other fields.
// 3. URL query parameters found in the HTTP request are mapped to (c) fields.
// 4. Any body sent with an HTTP request can contain only (b) fields.
//
// The syntax of the path template is as follows:
//
//     Template = "/" Segments [ Verb ] ;
//     Segments = Segment { "/" Segment } ;
//     Segment  = "*" | "**" | LITERAL | Variable ;
//     Variable = "{" FieldPath [ "=" Segments ] "}" ;
//     FieldPath = IDENT { "." IDENT } ;
//     Verb     = ":" LITERAL ;
//
// The syntax ${backtick}*${backtick} matches a single path segment. The syntax ${backtick}**${backtick} matches zero
// or more path segments, which must be the last part of the path except the
// ${backtick}Verb${backtick}. The syntax ${backtick}LITERAL${backtick} matches literal text in the path.
//
// The syntax ${backtick}Variable${backtick} matches part of the URL path as specified by its
// template. A variable template must not contain other variables. If a variable
// matches a single path segment, its template may be omitted, e.g. ${backtick}{var}${backtick}
// is equivalent to ${backtick}{var=*}${backtick}.
//
// If a variable contains exactly one path segment, such as ${backtick}"{var}"${backtick} or
// ${backtick}"{var=*}"${backtick}, when such a variable is expanded into a URL path, all characters
// except ${backtick}[-_.~0-9a-zA-Z]${backtick} are percent-encoded. Such variables show up in the
// Discovery Document as ${backtick}{var}${backtick}.
//
// If a variable contains one or more path segments, such as ${backtick}"{var=foo/*}"${backtick}
// or ${backtick}"{var=**}"${backtick}, when such a variable is expanded into a URL path, all
// characters except ${backtick}[-_.~/0-9a-zA-Z]${backtick} are percent-encoded. Such variables
// show up in the Discovery Document as ${backtick}{+var}${backtick}.
//
// NOTE: While the single segment variable matches the semantics of
// [RFC 6570](https://tools.ietf.org/html/rfc6570) Section 3.2.2
// Simple String Expansion, the multi segment variable **does not** match
// RFC 6570 Reserved Expansion. The reason is that the Reserved Expansion
// does not expand special characters like ${backtick}?${backtick} and ${backtick}#${backtick}, which would lead
// to invalid URLs.
//
// NOTE: the field paths in variables and in the ${backtick}body${backtick} must not refer to
// repeated fields or map fields.
message HttpRule {
  // Selects methods to which this rule applies.
  //
  // Refer to [selector][google.api.DocumentationRule.selector] for syntax details.
  string selector = 1;

  // Determines the URL pattern is matched by this rules. This pattern can be
  // used with any of the {get|put|post|delete|patch} methods. A custom method
  // can be defined using the 'custom' field.
  oneof pattern {
    // Used for listing and getting information about resources.
    string get = 2;

    // Used for updating a resource.
    string put = 3;

    // Used for creating a resource.
    string post = 4;

    // Used for deleting a resource.
    string delete = 5;

    // Used for updating a resource.
    string patch = 6;

    // The custom pattern is used for specifying an HTTP method that is not
    // included in the ${backtick}pattern${backtick} field, such as HEAD, or "*" to leave the
    // HTTP method unspecified for this rule. The wild-card rule is useful
    // for services that provide content to Web (HTML) clients.
    CustomHttpPattern custom = 8;
  }

  // The name of the request field whose value is mapped to the HTTP body, or
  // ${backtick}*${backtick} for mapping all fields not captured by the path pattern to the HTTP
  // body. NOTE: the referred field must not be a repeated field and must be
  // present at the top-level of request message type.
  string body = 7;

  // Optional. The name of the response field whose value is mapped to the HTTP
  // body of response. Other response fields are ignored. When
  // not set, the response message will be used as HTTP body of response.
  string response_body = 12;

  // Additional HTTP bindings for the selector. Nested bindings must
  // not contain an ${backtick}additional_bindings${backtick} field themselves (that is,
  // the nesting may only be one level deep).
  repeated HttpRule additional_bindings = 11;
}

// A custom pattern is used for defining custom HTTP verb.
message CustomHttpPattern {
  // The name of this custom HTTP verb.
  string kind = 1;

  // The path matched by this custom verb.
  string path = 2;
}
`

	tpl = strings.Replace(tpl, "${backtick}", "`", -1)

	t, err = template.New("google.api.http.proto").Parse(tpl)
	if err != nil {
		return err
	}

	fileName = dir + "http.proto"

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	tpl = `// Copyright 2018 Google LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

syntax = "proto3";

package google.api;

import "google/protobuf/any.proto";

option cc_enable_arenas = true;
option go_package = "google.golang.org/genproto/googleapis/api/httpbody;httpbody";
option java_multiple_files = true;
option java_outer_classname = "HttpBodyProto";
option java_package = "com.google.api";
option objc_class_prefix = "GAPI";

// Message that represents an arbitrary HTTP body. It should only be used for
// payload formats that can't be represented as JSON, such as raw binary or
// an HTML page.
//
//
// This message can be used both in streaming and non-streaming API methods in
// the request as well as the response.
//
// It can be used as a top-level request field, which is convenient if one
// wants to extract parameters from either the URL or HTTP template into the
// request fields and also want access to the raw HTTP body.
//
// Example:
//
//     message GetResourceRequest {
//       // A unique request id.
//       string request_id = 1;
//
//       // The raw HTTP body is bound to this field.
//       google.api.HttpBody http_body = 2;
//     }
//
//     service ResourceService {
//       rpc GetResource(GetResourceRequest) returns (google.api.HttpBody);
//       rpc UpdateResource(google.api.HttpBody) returns
//       (google.protobuf.Empty);
//     }
//
// Example with streaming methods:
//
//     service CaldavService {
//       rpc GetCalendar(stream google.api.HttpBody)
//         returns (stream google.api.HttpBody);
//       rpc UpdateCalendar(stream google.api.HttpBody)
//         returns (stream google.api.HttpBody);
//     }
//
// Use of this type only changes how the request and response bodies are
// handled, all other features will continue to work unchanged.
message HttpBody {
  // The HTTP Content-Type header value specifying the content type of the body.
  string content_type = 1;

  // The HTTP request/response body as raw binary.
  bytes data = 2;

  // Application specific response metadata. Must be set in the first response
  // for streaming APIs.
  repeated google.protobuf.Any extensions = 3;
}
`

	t, err = template.New("google.api.httpbody.proto").Parse(tpl)
	if err != nil {
		return err
	}

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/pkg/proto/google/api/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + "httpbody.proto"

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	tpl = `// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

syntax = "proto3";

package google.rpc;

option go_package = "google.golang.org/genproto/googleapis/rpc/code;code";
option java_multiple_files = true;
option java_outer_classname = "CodeProto";
option java_package = "com.google.rpc";
option objc_class_prefix = "RPC";


// The canonical error codes for Google APIs.
//
//
// Sometimes multiple error codes may apply.  Services should return
// the most specific error code that applies.  For example, prefer
// ${backtick}OUT_OF_RANGE${backtick} over ${backtick}FAILED_PRECONDITION${backtick} if both codes apply.
// Similarly prefer ${backtick}NOT_FOUND${backtick} or ${backtick}ALREADY_EXISTS${backtick} over ${backtick}FAILED_PRECONDITION${backtick}.
enum Code {
  // Not an error; returned on success
  //
  // HTTP Mapping: 200 OK
  OK = 0;

  // The operation was cancelled, typically by the caller.
  //
  // HTTP Mapping: 499 Client Closed Request
  CANCELLED = 1;

  // Unknown error.  For example, this error may be returned when
  // a ${backtick}Status${backtick} value received from another address space belongs to
  // an error space that is not known in this address space.  Also
  // errors raised by APIs that do not return enough error information
  // may be converted to this error.
  //
  // HTTP Mapping: 500 Internal Server Error
  UNKNOWN = 2;

  // The client specified an invalid argument.  Note that this differs
  // from ${backtick}FAILED_PRECONDITION${backtick}.  ${backtick}INVALID_ARGUMENT${backtick} indicates arguments
  // that are problematic regardless of the state of the system
  // (e.g., a malformed file name).
  //
  // HTTP Mapping: 400 Bad Request
  INVALID_ARGUMENT = 3;

  // The deadline expired before the operation could complete. For operations
  // that change the state of the system, this error may be returned
  // even if the operation has completed successfully.  For example, a
  // successful response from a server could have been delayed long
  // enough for the deadline to expire.
  //
  // HTTP Mapping: 504 Gateway Timeout
  DEADLINE_EXCEEDED = 4;

  // Some requested entity (e.g., file or directory) was not found.
  //
  // Note to server developers: if a request is denied for an entire class
  // of users, such as gradual feature rollout or undocumented whitelist,
  // ${backtick}NOT_FOUND${backtick} may be used. If a request is denied for some users within
  // a class of users, such as user-based access control, ${backtick}PERMISSION_DENIED${backtick}
  // must be used.
  //
  // HTTP Mapping: 404 Not Found
  NOT_FOUND = 5;

  // The entity that a client attempted to create (e.g., file or directory)
  // already exists.
  //
  // HTTP Mapping: 409 Conflict
  ALREADY_EXISTS = 6;

  // The caller does not have permission to execute the specified
  // operation. ${backtick}PERMISSION_DENIED${backtick} must not be used for rejections
  // caused by exhausting some resource (use ${backtick}RESOURCE_EXHAUSTED${backtick}
  // instead for those errors). ${backtick}PERMISSION_DENIED${backtick} must not be
  // used if the caller can not be identified (use ${backtick}UNAUTHENTICATED${backtick}
  // instead for those errors). This error code does not imply the
  // request is valid or the requested entity exists or satisfies
  // other pre-conditions.
  //
  // HTTP Mapping: 403 Forbidden
  PERMISSION_DENIED = 7;

  // The request does not have valid authentication credentials for the
  // operation.
  //
  // HTTP Mapping: 401 Unauthorized
  UNAUTHENTICATED = 16;

  // Some resource has been exhausted, perhaps a per-user quota, or
  // perhaps the entire file system is out of space.
  //
  // HTTP Mapping: 429 Too Many Requests
  RESOURCE_EXHAUSTED = 8;

  // The operation was rejected because the system is not in a state
  // required for the operation's execution.  For example, the directory
  // to be deleted is non-empty, an rmdir operation is applied to
  // a non-directory, etc.
  //
  // Service implementors can use the following guidelines to decide
  // between ${backtick}FAILED_PRECONDITION${backtick}, ${backtick}ABORTED${backtick}, and ${backtick}UNAVAILABLE${backtick}:
  //  (a) Use ${backtick}UNAVAILABLE${backtick} if the client can retry just the failing call.
  //  (b) Use ${backtick}ABORTED${backtick} if the client should retry at a higher level
  //      (e.g., when a client-specified test-and-set fails, indicating the
  //      client should restart a read-modify-write sequence).
  //  (c) Use ${backtick}FAILED_PRECONDITION${backtick} if the client should not retry until
  //      the system state has been explicitly fixed.  E.g., if an "rmdir"
  //      fails because the directory is non-empty, ${backtick}FAILED_PRECONDITION${backtick}
  //      should be returned since the client should not retry unless
  //      the files are deleted from the directory.
  //
  // HTTP Mapping: 400 Bad Request
  FAILED_PRECONDITION = 9;

  // The operation was aborted, typically due to a concurrency issue such as
  // a sequencer check failure or transaction abort.
  //
  // See the guidelines above for deciding between ${backtick}FAILED_PRECONDITION${backtick},
  // ${backtick}ABORTED${backtick}, and ${backtick}UNAVAILABLE${backtick}.
  //
  // HTTP Mapping: 409 Conflict
  ABORTED = 10;

  // The operation was attempted past the valid range.  E.g., seeking or
  // reading past end-of-file.
  //
  // Unlike ${backtick}INVALID_ARGUMENT${backtick}, this error indicates a problem that may
  // be fixed if the system state changes. For example, a 32-bit file
  // system will generate ${backtick}INVALID_ARGUMENT${backtick} if asked to read at an
  // offset that is not in the range [0,2^32-1], but it will generate
  // ${backtick}OUT_OF_RANGE${backtick} if asked to read from an offset past the current
  // file size.
  //
  // There is a fair bit of overlap between ${backtick}FAILED_PRECONDITION${backtick} and
  // ${backtick}OUT_OF_RANGE${backtick}.  We recommend using ${backtick}OUT_OF_RANGE${backtick} (the more specific
  // error) when it applies so that callers who are iterating through
  // a space can easily look for an ${backtick}OUT_OF_RANGE${backtick} error to detect when
  // they are done.
  //
  // HTTP Mapping: 400 Bad Request
  OUT_OF_RANGE = 11;

  // The operation is not implemented or is not supported/enabled in this
  // service.
  //
  // HTTP Mapping: 501 Not Implemented
  UNIMPLEMENTED = 12;

  // Internal errors.  This means that some invariants expected by the
  // underlying system have been broken.  This error code is reserved
  // for serious errors.
  //
  // HTTP Mapping: 500 Internal Server Error
  INTERNAL = 13;

  // The service is currently unavailable.  This is most likely a
  // transient condition, which can be corrected by retrying with
  // a backoff.
  //
  // See the guidelines above for deciding between ${backtick}FAILED_PRECONDITION${backtick},
  // ${backtick}ABORTED${backtick}, and ${backtick}UNAVAILABLE${backtick}.
  //
  // HTTP Mapping: 503 Service Unavailable
  UNAVAILABLE = 14;

  // Unrecoverable data loss or corruption.
  //
  // HTTP Mapping: 500 Internal Server Error
  DATA_LOSS = 15;
}
`

	tpl = strings.Replace(tpl, "${backtick}", "`", -1)

	t, err = template.New("google.rpc.code.proto").Parse(tpl)
	if err != nil {
		return err
	}

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/pkg/proto/google/rpc/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + "code.proto"

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	tpl = `// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

syntax = "proto3";

package google.rpc;

import "google/protobuf/duration.proto";

option go_package = "google.golang.org/genproto/googleapis/rpc/errdetails;errdetails";
option java_multiple_files = true;
option java_outer_classname = "ErrorDetailsProto";
option java_package = "com.google.rpc";
option objc_class_prefix = "RPC";


// Describes when the clients can retry a failed request. Clients could ignore
// the recommendation here or retry when this information is missing from error
// responses.
//
// It's always recommended that clients should use exponential backoff when
// retrying.
//
// Clients should wait until ${backtick}retry_delay${backtick} amount of time has passed since
// receiving the error response before retrying.  If retrying requests also
// fail, clients should use an exponential backoff scheme to gradually increase
// the delay between retries based on ${backtick}retry_delay${backtick}, until either a maximum
// number of retires have been reached or a maximum retry delay cap has been
// reached.
message RetryInfo {
  // Clients should wait at least this long between retrying the same request.
  google.protobuf.Duration retry_delay = 1;
}

// Describes additional debugging info.
message DebugInfo {
  // The stack trace entries indicating where the error occurred.
  repeated string stack_entries = 1;

  // Additional debugging information provided by the server.
  string detail = 2;
}

// Describes how a quota check failed.
//
// For example if a daily limit was exceeded for the calling project,
// a service could respond with a QuotaFailure detail containing the project
// id and the description of the quota limit that was exceeded.  If the
// calling project hasn't enabled the service in the developer console, then
// a service could respond with the project id and set ${backtick}service_disabled${backtick}
// to true.
//
// Also see RetryDetail and Help types for other details about handling a
// quota failure.
message QuotaFailure {
  // A message type used to describe a single quota violation.  For example, a
  // daily quota or a custom quota that was exceeded.
  message Violation {
    // The subject on which the quota check failed.
    // For example, "clientip:<ip address of client>" or "project:<Google
    // developer project id>".
    string subject = 1;

    // A description of how the quota check failed. Clients can use this
    // description to find more about the quota configuration in the service's
    // public documentation, or find the relevant quota limit to adjust through
    // developer console.
    //
    // For example: "Service disabled" or "Daily Limit for read operations
    // exceeded".
    string description = 2;
  }

  // Describes all quota violations.
  repeated Violation violations = 1;
}

// Describes what preconditions have failed.
//
// For example, if an RPC failed because it required the Terms of Service to be
// acknowledged, it could list the terms of service violation in the
// PreconditionFailure message.
message PreconditionFailure {
  // A message type used to describe a single precondition failure.
  message Violation {
    // The type of PreconditionFailure. We recommend using a service-specific
    // enum type to define the supported precondition violation types. For
    // example, "TOS" for "Terms of Service violation".
    string type = 1;

    // The subject, relative to the type, that failed.
    // For example, "google.com/cloud" relative to the "TOS" type would
    // indicate which terms of service is being referenced.
    string subject = 2;

    // A description of how the precondition failed. Developers can use this
    // description to understand how to fix the failure.
    //
    // For example: "Terms of service not accepted".
    string description = 3;
  }

  // Describes all precondition violations.
  repeated Violation violations = 1;
}

// Describes violations in a client request. This error type focuses on the
// syntactic aspects of the request.
message BadRequest {
  // A message type used to describe a single bad request field.
  message FieldViolation {
    // A path leading to a field in the request body. The value will be a
    // sequence of dot-separated identifiers that identify a protocol buffer
    // field. E.g., "field_violations.field" would identify this field.
    string field = 1;

    // A description of why the request element is bad.
    string description = 2;
  }

  // Describes all violations in a client request.
  repeated FieldViolation field_violations = 1;
}

// Contains metadata about the request that clients can attach when filing a bug
// or providing other forms of feedback.
message RequestInfo {
  // An opaque string that should only be interpreted by the service generating
  // it. For example, it can be used to identify requests in the service's logs.
  string request_id = 1;

  // Any data that was used to serve this request. For example, an encrypted
  // stack trace that can be sent back to the service provider for debugging.
  string serving_data = 2;
}

// Describes the resource that is being accessed.
message ResourceInfo {
  // A name for the type of resource being accessed, e.g. "sql table",
  // "cloud storage bucket", "file", "Google calendar"; or the type URL
  // of the resource: e.g. "type.googleapis.com/google.pubsub.v1.Topic".
  string resource_type = 1;

  // The name of the resource being accessed.  For example, a shared calendar
  // name: "example.com_4fghdhgsrgh@group.calendar.google.com", if the current
  // error is [google.rpc.Code.PERMISSION_DENIED][google.rpc.Code.PERMISSION_DENIED].
  string resource_name = 2;

  // The owner of the resource (optional).
  // For example, "user:<owner email>" or "project:<Google developer project
  // id>".
  string owner = 3;

  // Describes what error is encountered when accessing this resource.
  // For example, updating a cloud project may require the ${backtick}writer${backtick} permission
  // on the developer console project.
  string description = 4;
}

// Provides links to documentation or for performing an out of band action.
//
// For example, if a quota check failed with an error indicating the calling
// project hasn't enabled the accessed service, this can contain a URL pointing
// directly to the right place in the developer console to flip the bit.
message Help {
  // Describes a URL link.
  message Link {
    // Describes what the link offers.
    string description = 1;

    // The URL of the link.
    string url = 2;
  }

  // URL(s) pointing to additional information on handling the current error.
  repeated Link links = 1;
}

// Provides a localized error message that is safe to return to the user
// which can be attached to an RPC error.
message LocalizedMessage {
  // The locale used following the specification defined at
  // http://www.rfc-editor.org/rfc/bcp/bcp47.txt.
  // Examples are: "en-US", "fr-CH", "es-MX"
  string locale = 1;

  // The localized error message in the above locale.
  string message = 2;
}
`

	tpl = strings.Replace(tpl, "${backtick}", "`", -1)

	t, err = template.New("google.rpc.error.proto").Parse(tpl)
	if err != nil {
		return err
	}

	fileName = dir + "error_details.proto"

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	tpl = `// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

syntax = "proto3";

package google.rpc;

import "google/protobuf/any.proto";

option go_package = "google.golang.org/genproto/googleapis/rpc/status;status";
option java_multiple_files = true;
option java_outer_classname = "StatusProto";
option java_package = "com.google.rpc";
option objc_class_prefix = "RPC";


// The ${backtick}Status${backtick} type defines a logical error model that is suitable for different
// programming environments, including REST APIs and RPC APIs. It is used by
// [gRPC](https://github.com/grpc). The error model is designed to be:
//
// - Simple to use and understand for most users
// - Flexible enough to meet unexpected needs
//
// # Overview
//
// The ${backtick}Status${backtick} message contains three pieces of data: error code, error message,
// and error details. The error code should be an enum value of
// [google.rpc.Code][google.rpc.Code], but it may accept additional error codes if needed.  The
// error message should be a developer-facing English message that helps
// developers *understand* and *resolve* the error. If a localized user-facing
// error message is needed, put the localized message in the error details or
// localize it in the client. The optional error details may contain arbitrary
// information about the error. There is a predefined set of error detail types
// in the package ${backtick}google.rpc${backtick} that can be used for common error conditions.
//
// # Language mapping
//
// The ${backtick}Status${backtick} message is the logical representation of the error model, but it
// is not necessarily the actual wire format. When the ${backtick}Status${backtick} message is
// exposed in different client libraries and different wire protocols, it can be
// mapped differently. For example, it will likely be mapped to some exceptions
// in Java, but more likely mapped to some error codes in C.
//
// # Other uses
//
// The error model and the ${backtick}Status${backtick} message can be used in a variety of
// environments, either with or without APIs, to provide a
// consistent developer experience across different environments.
//
// Example uses of this error model include:
//
// - Partial errors. If a service needs to return partial errors to the client,
//     it may embed the ${backtick}Status${backtick} in the normal response to indicate the partial
//     errors.
//
// - Workflow errors. A typical workflow has multiple steps. Each step may
//     have a ${backtick}Status${backtick} message for error reporting.
//
// - Batch operations. If a client uses batch request and batch response, the
//     ${backtick}Status${backtick} message should be used directly inside batch response, one for
//     each error sub-response.
//
// - Asynchronous operations. If an API call embeds asynchronous operation
//     results in its response, the status of those operations should be
//     represented directly using the ${backtick}Status${backtick} message.
//
// - Logging. If some API errors are stored in logs, the message ${backtick}Status${backtick} could
//     be used directly after any stripping needed for security/privacy reasons.
message Status {
  // The status code, which should be an enum value of [google.rpc.Code][google.rpc.Code].
  int32 code = 1;

  // A developer-facing error message, which should be in English. Any
  // user-facing error message should be localized and sent in the
  // [google.rpc.Status.details][google.rpc.Status.details] field, or localized by the client.
  string message = 2;

  // A list of messages that carry the error details.  There is a common set of
  // message types for APIs to use.
  repeated google.protobuf.Any details = 3;
}
`

	tpl = strings.Replace(tpl, "${backtick}", "`", -1)

	t, err = template.New("google.rpc.status.proto").Parse(tpl)
	if err != nil {
		return err
	}

	fileName = dir + "status.proto"

	f, err = os.Create(fileName)
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
