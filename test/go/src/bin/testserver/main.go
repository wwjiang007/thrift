/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/apache/thrift/test/go/src/common"
)

var host = flag.String("host", "localhost", "Host to connect")
var port = flag.Int64("port", 9090, "Port number to connect")
var domain_socket = flag.String("domain-socket", "", "Domain Socket (e.g. /tmp/ThriftTest.thrift), instead of host and port")
var transport = flag.String("transport", "buffered", "Transport: buffered, framed, http, zlib")
var _ = flag.Bool("zlib", false, "For compatibility. Ignored.")
var protocol = flag.String("protocol", "binary", "Protocol: binary, compact, json, header")
var ssl = flag.Bool("ssl", false, "Encrypted Transport using SSL")
var certPath = flag.String("certPath", "keys", "Directory that contains SSL certificates")

func main() {
	flag.Parse()

	processor, serverTransport, transportFactory, protocolFactory, _, err := common.GetServerParams(*host, *port, *domain_socket, *transport, *protocol, *ssl, *certPath, common.PrintingHandler)

	if err != nil {
		log.Fatalf("Unable to process server params: %v", err)
	}

	if *transport == "http" {
		http.HandleFunc("/", thrift.NewThriftHandlerFunc(processor, protocolFactory, protocolFactory))

		if *ssl {
			err := http.ListenAndServeTLS(fmt.Sprintf(":%d", *port),
				*certPath+"/server.pem", *certPath+"/server.key", nil)

			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
				fmt.Println(err)
				return
			}
		}
	} else {
		server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
		if err = server.Listen(); err != nil {
			return
		}
		go server.AcceptLoop()
		server.Serve()
	}
}
