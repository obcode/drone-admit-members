// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"errors"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/admission"

	"github.com/google/go-github/v28/github"
	"github.com/sirupsen/logrus"
)

// ErrAccessDenied is returned if access is denied.
var ErrAccessDenied = errors.New("admission: access denied")

// New returns a new admission plugin.
func New(client *github.Client, orgs string, admins string) admission.Plugin {
	return &plugin{
		client: client,
		orgs:   orgs,
		admins: admins,
	}
}

type plugin struct {
	client *github.Client
	orgs   string // members of this orgs are granted access
	admins string
}

func (p *plugin) Admit(ctx context.Context, req *admission.Request) (*drone.User, error) {
	u := req.User

	logrus.WithField("user", u.Login).
		Debugln("requesting system access")

	// check organization membership
	orgs := strings.Split(p.orgs, ",")
	admins := strings.Split(p.admins, ",")

	for _, org := range orgs {
		_, _, err := p.client.Organizations.GetOrgMembership(ctx, u.Login, org)
		if err == nil {
			for _, admin := range admins {
				if u.Login == admin {
					logrus.WithField("user", u.Login).
						WithField("org", org).
						WithField("role", "admin").
						Debugln("granted admin system access")

					u.Admin = true

					return &u, err
				}
			}

			logrus.WithField("user", u.Login).
				WithField("org", orgs).
				Debugln("granted standard system access")

			return nil, err
		}
	}

	return nil, ErrAccessDenied
}
