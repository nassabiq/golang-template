syntax = "proto3";

package {{MODULE}}.v1;

import "google/api/annotations.proto";
import "proto/common/common.proto";

option go_package = "github.com/nassabiq/golang-template/proto/{{MODULE}}";

service {{MODULE}}Service {

  rpc List(List{{MODULE}}Request) returns (List{{MODULE}}Response) {
    option (google.api.http) = {
      get: "/{{MODULE|lower}}s"
    };
  }

  rpc GetByID(GetByIDRequest) returns ({{MODULE}}Response) {
    option (google.api.http) = {
      get: "/{{MODULE|lower}}s/{id}"
    };
  }

  rpc Create(Create{{MODULE}}Request) returns ({{MODULE}}Response) {
    option (google.api.http) = {
      post: "/{{MODULE|lower}}s"
      body: "*"
    };
  }

  rpc Update(Update{{MODULE}}Request) returns ({{MODULE}}Response) {
    option (google.api.http) = {
      put: "/{{MODULE|lower}}s/{id}"
      body: "*"
    };
  }

  rpc Delete(Delete{{MODULE}}Request) returns (Delete{{MODULE}}Response) {
    option (google.api.http) = {
      delete: "/{{MODULE|lower}}s/{id}"
    };
  }
}

// ENTITY
message {{MODULE}} {
  string id = 1;
  string name = 2;
}

// REQUEST
message List{{MODULE}}Request {
  int32 limit = 1;
  int32 offset = 2;
}

message GetByIDRequest {
  string id = 1;
}

message Create{{MODULE}}Request {
  string name = 1;
}

message Update{{MODULE}}Request {
  string id = 1;
  optional string name = 2;
}

message Delete{{MODULE}}Request {
  string id = 1;
}

// RESPONSE
message List{{MODULE}}Response {
  common.v1.MetaData metadata = 1;
  repeated {{MODULE}} data = 2;
  common.v1.Pagination pagination = 3;
}

message {{MODULE}}Response {
  common.v1.MetaData metadata = 1;
  {{MODULE}} data = 2;
}

message Delete{{MODULE}}Response {
  common.v1.MetaData metadata = 1;
}

message Empty {}
