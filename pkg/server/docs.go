// Copyright 2018 the Service Broker Project Authors.
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

package server

import (
	"net/http"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/generator"
	"github.com/russross/blackfriday"
)

// NewDocsHandler returns a handler func that generates HTML documentation for
// the given registry.
func NewDocsHandler(registry broker.BrokerRegistry) http.HandlerFunc {
	docsPageMd := generator.CatalogDocumentation(registry)

	renderer := blackfriday.HtmlRenderer(
		blackfriday.HTML_COMPLETE_PAGE,
		"Service Broker Documents",
		"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css",
	)

	page := blackfriday.Markdown([]byte(docsPageMd), renderer, 0)
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		w.Write(page)
	}
}
