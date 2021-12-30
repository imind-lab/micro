/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package api

import (
	"os"
	"text/template"

	tp "github.com/imind-lab/micro/microctl/template"
)

// 生成constant
func CreatePkg(data *tp.Data) error {

	// 生成pkg.cache.go
	tpl := `/**
 *  IMindLab
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package constant

import "time"

const (
	CachePrefix   = "im_"
	CacheDay30    = time.Hour * 24 * 30
	CacheDay15    = time.Hour * 24 * 15
	CacheDay7     = time.Hour * 24 * 7
	CacheDay3     = time.Hour * 24 * 3
	CacheDay2     = time.Hour * 24 * 2
	CacheDay1     = time.Hour * 24
	CacheHour12   = time.Hour * 12
	CacheHour6    = time.Hour * 6
	CacheHour2    = time.Hour * 2
	CacheHour1    = time.Hour
	CacheMinute30 = time.Minute * 30
	CacheMinute10 = time.Minute * 10
	CacheMinute5  = time.Minute * 5
	CacheMinute1  = time.Minute
	CacheSecond20 = time.Second * 20
)
`

	t, err := template.New("constant.cache").Parse(tpl)
	if err != nil {
		return err
	}

	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/pkg/constant/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "cache.go"
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成pkg.option.go
	tpl = `/**
 *  IMindLab
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package constant

import (
	"time"
)

//CRequestTimeout 并发请求超时时间
const CRequestTimeout = time.Second * 10

const DBName = "imind"
const Realtime = false

const MQName = "business"
const GreetQueueLen = 32
`

	t, err = template.New("constant.option").Parse(tpl)
	if err != nil {
		return err
	}

	fileName = dir + "option.go"
	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成pkg.cache.go
	tpl = `package util

import (
        "{{.Domain}}/{{.Project}}/{{.Service}}/pkg/constant"
        "github.com/imind-lab/micro/util"
)

func CacheKey(keys ...string) string {
        return constant.CachePrefix + util.AppendString(keys...)
}
`

	t, err = template.New("util.cache").Parse(tpl)
	if err != nil {
		return err
	}

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/pkg/util/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + "cache.go"
	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成pkg.hash.go
	tpl = `/**
 *  IMindLab
 *
 *  Create by songli on {{.Date}}
 *  Copyright © {{.Year}} imind.tech All rights reserved.
 */

package util

import "math"

var HashBase = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

func Base62encode(id int32) string {
	baseStr := ""
	for {
		if id <= 0 {
			break
		}
		i := id % 62
		baseStr += HashBase[i]
		id = (id - i) / 62
	}
	return baseStr
}
func Base62decode(base62 string) int {
	rs := 0
	length := len(base62)
	f := flip(HashBase)
	for i := 0; i < length; i++ {
		rs += f[string(base62[i])] * int(math.Pow(62, float64(i)))
	}
	return rs
}
func flip(s []string) map[string]int {
	f := make(map[string]int)
	for index, value := range s {
		f[value] = index
	}
	return f
}
`

	t, err = template.New("util.hash").Parse(tpl)
	if err != nil {
		return err
	}

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/pkg/util/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + "hash.go"
	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成pkg.tls.go
	tpl = `package util

import (
	"crypto/tls"
	"io/ioutil"
	"log"

	"golang.org/x/net/http2"
)

func GetTLSConfig(certPemPath, certKeyPath string) *tls.Config {
	var certKeyPair *tls.Certificate
	cert, _ := ioutil.ReadFile(certPemPath)
	key, _ := ioutil.ReadFile(certKeyPath)

	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Println("TLS KeyPair err: %v\n", err)
	}

	certKeyPair = &pair

	return &tls.Config{
		Certificates: []tls.Certificate{*certKeyPair},
		NextProtos:   []string{http2.NextProtoTLS},
	}
}
`

	t, err = template.New("util.tls").Parse(tpl)
	if err != nil {
		return err
	}

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/pkg/util/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + "tls.go"
	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成pkg.probe.go
	tpl = `// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var (
	flAddr          string
	flService       string
	flUserAgent     string
	flConnTimeout   time.Duration
	flRPCTimeout    time.Duration
	flTLS           bool
	flTLSNoVerify   bool
	flTLSCACert     string
	flTLSClientCert string
	flTLSClientKey  string
	flTLSServerName string
	flVerbose       bool
	flGZIP          bool
	flSPIFFE        bool
)

const (
	// StatusInvalidArguments indicates specified invalid arguments.
	StatusInvalidArguments = 1
	// StatusConnectionFailure indicates connection failed.
	StatusConnectionFailure = 2
	// StatusRPCFailure indicates rpc failed.
	StatusRPCFailure = 3
	// StatusUnhealthy indicates rpc succeeded but indicates unhealthy service.
	StatusUnhealthy = 4
	// StatusSpiffeFailed indicates failure to retrieve credentials using spiffe workload API
	StatusSpiffeFailed = 20
)

func init() {
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	log.SetFlags(0)
	flagSet.StringVar(&flAddr, "addr", "", "(required) tcp host:port to connect")
	flagSet.StringVar(&flService, "service", "", "service name to check (default: \"\")")
	flagSet.StringVar(&flUserAgent, "user-agent", "grpc_health_probe", "user-agent header value of health check requests")
	// timeouts
	flagSet.DurationVar(&flConnTimeout, "connect-timeout", time.Second, "timeout for establishing connection")
	flagSet.DurationVar(&flRPCTimeout, "rpc-timeout", time.Second, "timeout for health check rpc")
	// tls settings
	flagSet.BoolVar(&flTLS, "tls", false, "use TLS (default: false, INSECURE plaintext transport)")
	flagSet.BoolVar(&flTLSNoVerify, "tls-no-verify", false, "(with -tls) don't verify the certificate (INSECURE) presented by the server (default: false)")
	flagSet.StringVar(&flTLSCACert, "tls-ca-cert", "", "(with -tls, optional) file containing trusted certificates for verifying server")
	flagSet.StringVar(&flTLSClientCert, "tls-client-cert", "", "(with -tls, optional) client certificate for authenticating to the server (requires -tls-client-key)")
	flagSet.StringVar(&flTLSClientKey, "tls-client-key", "", "(with -tls) client private key for authenticating to the server (requires -tls-client-cert)")
	flagSet.StringVar(&flTLSServerName, "tls-server-name", "", "(with -tls) override the hostname used to verify the server certificate")
	flagSet.BoolVar(&flVerbose, "v", false, "verbose logs")
	flagSet.BoolVar(&flGZIP, "gzip", false, "use GZIPCompressor for requests and GZIPDecompressor for response (default: false)")
	flagSet.BoolVar(&flSPIFFE, "spiffe", false, "use SPIFFE to obtain mTLS credentials")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		os.Exit(StatusInvalidArguments)
	}

	argError := func(s string, v ...interface{}) {
		log.Printf("error: "+s, v...)
		os.Exit(StatusInvalidArguments)
	}

	if flAddr == "" {
		argError("-addr not specified")
	}
	if flConnTimeout <= 0 {
		argError("-connect-timeout must be greater than zero (specified: %v)", flConnTimeout)
	}
	if flRPCTimeout <= 0 {
		argError("-rpc-timeout must be greater than zero (specified: %v)", flRPCTimeout)
	}
	if !flTLS && flTLSNoVerify {
		argError("specified -tls-no-verify without specifying -tls")
	}
	if !flTLS && flTLSCACert != "" {
		argError("specified -tls-ca-cert without specifying -tls")
	}
	if !flTLS && flTLSClientCert != "" {
		argError("specified -tls-client-cert without specifying -tls")
	}
	if !flTLS && flTLSServerName != "" {
		argError("specified -tls-server-name without specifying -tls")
	}
	if flTLSClientCert != "" && flTLSClientKey == "" {
		argError("specified -tls-client-cert without specifying -tls-client-key")
	}
	if flTLSClientCert == "" && flTLSClientKey != "" {
		argError("specified -tls-client-key without specifying -tls-client-cert")
	}
	if flTLSNoVerify && flTLSCACert != "" {
		argError("cannot specify -tls-ca-cert with -tls-no-verify (CA cert would not be used)")
	}
	if flTLSNoVerify && flTLSServerName != "" {
		argError("cannot specify -tls-server-name with -tls-no-verify (server name would not be used)")
	}

	if flVerbose {
		log.Printf("parsed options:")
		log.Printf("> addr=%s conn_timeout=%v rpc_timeout=%v", flAddr, flConnTimeout, flRPCTimeout)
		log.Printf("> tls=%v", flTLS)
		if flTLS {
			log.Printf("  > no-verify=%v ", flTLSNoVerify)
			log.Printf("  > ca-cert=%s", flTLSCACert)
			log.Printf("  > client-cert=%s", flTLSClientCert)
			log.Printf("  > client-key=%s", flTLSClientKey)
			log.Printf("  > server-name=%s", flTLSServerName)
		}
		log.Printf("> spiffe=%v", flSPIFFE)
	}
}

func buildCredentials(skipVerify bool, caCerts, clientCert, clientKey, serverName string) (credentials.TransportCredentials, error) {
	var cfg tls.Config

	if clientCert != "" && clientKey != "" {
		keyPair, err := tls.LoadX509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load tls client cert/key pair. error=%v", err)
		}
		cfg.Certificates = []tls.Certificate{keyPair}
	}

	if skipVerify {
		cfg.InsecureSkipVerify = true
	} else if caCerts != "" {
		// override system roots
		rootCAs := x509.NewCertPool()
		pem, err := ioutil.ReadFile(caCerts)
		if err != nil {
			return nil, fmt.Errorf("failed to load root CA certificates from file (%s) error=%v", caCerts, err)
		}
		if !rootCAs.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("no root CA certs parsed from file %s", caCerts)
		}
		cfg.RootCAs = rootCAs
	}
	if serverName != "" {
		cfg.ServerName = serverName
	}
	return credentials.NewTLS(&cfg), nil
}

func main() {
	retcode := 0
	defer func() { os.Exit(retcode) }()

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		sig := <-c
		if sig == os.Interrupt {
			log.Printf("cancellation received")
			cancel()
			return
		}
	}()

	opts := []grpc.DialOption{
		grpc.WithUserAgent(flUserAgent),
		grpc.WithBlock(),
	}
	if flTLS && flSPIFFE {
		log.Printf("-tls and -spiffe are mutually incompatible")
		retcode = StatusInvalidArguments
		return
	}
	if flTLS {
		creds, err := buildCredentials(flTLSNoVerify, flTLSCACert, flTLSClientCert, flTLSClientKey, flTLSServerName)
		if err != nil {
			log.Printf("failed to initialize tls credentials. error=%v", err)
			retcode = StatusInvalidArguments
			return
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else if flSPIFFE {
		spiffeCtx, _ := context.WithTimeout(ctx, flRPCTimeout)
		source, err := workloadapi.NewX509Source(spiffeCtx)
		if err != nil {
			log.Printf("failed to initialize tls credentials with spiffe. error=%v", err)
			retcode = StatusSpiffeFailed
			return
		}
		if flVerbose {
			svid, err := source.GetX509SVID()
			if err != nil {
				log.Fatalf("error getting x509 svid: %+v", err)
			}
			log.Printf("SPIFFE Verifiable Identity Document (SVID): %q", svid.ID)
		}
		creds := credentials.NewTLS(tlsconfig.MTLSClientConfig(source, source, tlsconfig.AuthorizeAny()))
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	if flGZIP {
		opts = append(opts,
			grpc.WithCompressor(grpc.NewGZIPCompressor()),
			grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
		)
	}

	if flVerbose {
		log.Print("establishing connection")
	}
	connStart := time.Now()
	dialCtx, dialCancel := context.WithTimeout(ctx, flConnTimeout)
	defer dialCancel()
	conn, err := grpc.DialContext(dialCtx, flAddr, opts...)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Printf("timeout: failed to connect service %q within %v", flAddr, flConnTimeout)
		} else {
			log.Printf("error: failed to connect service at %q: %+v", flAddr, err)
		}
		retcode = StatusConnectionFailure
		return
	}
	connDuration := time.Since(connStart)
	defer conn.Close()
	if flVerbose {
		log.Printf("connection established (took %v)", connDuration)
	}

	rpcStart := time.Now()
	rpcCtx, rpcCancel := context.WithTimeout(ctx, flRPCTimeout)
	defer rpcCancel()
	resp, err := healthpb.NewHealthClient(conn).Check(rpcCtx,
		&healthpb.HealthCheckRequest{
			Service: flService})
	if err != nil {
		if stat, ok := status.FromError(err); ok && stat.Code() == codes.Unimplemented {
			log.Printf("error: this server does not implement the grpc health protocol (grpc.health.v1.Health): %s", stat.Message())
		} else if stat, ok := status.FromError(err); ok && stat.Code() == codes.DeadlineExceeded {
			log.Printf("timeout: health rpc did not complete within %v", flRPCTimeout)
		} else {
			log.Printf("error: health rpc failed: %+v", err)
		}
		retcode = StatusRPCFailure
		return
	}
	rpcDuration := time.Since(rpcStart)

	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		log.Printf("service unhealthy (responded with %q)", resp.GetStatus().String())
		retcode = StatusUnhealthy
		return
	}
	if flVerbose {
		log.Printf("time elapsed: connect=%v rpc=%v", connDuration, rpcDuration)
	}
	log.Printf("status: %v", resp.GetStatus().String())
}
`

	t, err = template.New("util.probe").Parse(tpl)
	if err != nil {
		return err
	}

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/pkg/grpc-health-probe/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + "main.go"
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
