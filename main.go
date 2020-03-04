// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package main

import (
	"context"
	"net/http"

	"github.com/drone/drone-go/plugin/admission"
	"github.com/obcode/drone-admit-members/plugin"

	"github.com/google/go-github/v28/github"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// spec provides the plugin settings.
type spec struct {
	Bind   string `envconfig:"DRONE_BIND"`
	Debug  bool   `envconfig:"DRONE_DEBUG"`
	Secret string `envconfig:"DRONE_SECRET"`

	Token    string `envconfig:"DRONE_GITHUB_TOKEN"`
	Endpoint string `envconfig:"DRONE_GITHUB_ENDPOINT" default:"https://api.github.com/"`
	Orgs     string `envconfig:"DRONE_GITHUB_ORGS"`
	Admins   string `envconfig:"DRONE_GITHUB_ADMINS"`
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)

	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {

		logrus.SetLevel(logrus.DebugLevel)
	}

	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}

	if spec.Bind == "" {
		spec.Bind = ":3000"
	}

	// creates the github client transport used
	// to authenticate API requests.
	trans := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: spec.Token},
	))

	// create the github client
	client, err := github.NewEnterpriseClient(spec.Endpoint, spec.Endpoint, trans)
	if err != nil {
		logrus.Fatal(err)
	}

	handler := admission.Handler(
		plugin.New(
			client,
			spec.Orgs,
			spec.Admins,
		),
		spec.Secret,
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}
