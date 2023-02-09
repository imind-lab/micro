/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package helper

import (
    "context"
    "net/http"
    "strconv"

    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc/metadata"
    "google.golang.org/protobuf/proto"
)

func InitMeta(_ context.Context, req *http.Request) metadata.MD {
    md := make(map[string]string, 4)
    md["ua"] = req.UserAgent()
    md["ip"] = req.RemoteAddr
    md["referer"] = req.Referer()
    md["phone"] = req.Form.Get("phone")
    return metadata.New(md)
}

func ResponseHeaderMatcher(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
    headers := w.Header()
    if location, ok := headers["Grpc-Metadata-Location"]; ok {
        w.Header().Set("Location", location[0])
        w.WriteHeader(http.StatusFound)
    }

    md, ok := runtime.ServerMetadataFromContext(ctx)
    if !ok {
        return nil
    }

    // set http status code
    if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
        code, err := strconv.Atoi(vals[0])
        if err != nil {
            return err
        }
        // delete the headers to not expose any grpc-metadata in http response
        delete(md.HeaderMD, "x-http-code")
        delete(w.Header(), "Grpc-Metadata-X-Http-Code")
        w.WriteHeader(code)
    }

    return nil
}
