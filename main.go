/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/tomahawk28/http-md5sum/pkg/httpsum"
)

func main() {
	var parallel uint
	var sites []string
	{
		if strings.Contains(os.Args[1], "parallel") || os.Args[1] == "-p" || os.Args[1] == "--p" {
			p, err := strconv.Atoi(os.Args[2])
			if err != nil {
				panic(err)
			}
			parallel = uint(p)
			sites = os.Args[3:]
		} else {
			parallel = 10
			sites = os.Args[1:]
		}
	}

	c := httpsum.Config{
		Client:   &http.Client{},
		Parallel: parallel,
	}

	service, err := httpsum.New(c)
	if err != nil {
		panic(err)
	}

	err = service.Ping(sites)
	if err != nil {
		panic(err)
	}
}