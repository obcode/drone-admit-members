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
func New(client *github.Client, org string) admission.Plugin {
	return &plugin{
		client: client,
		org:    org,
	}
}

type plugin struct {
	client *github.Client
	org    string // members of this org are granted access
}

func (p *plugin) Admit(ctx context.Context, req *admission.Request) (*drone.User, error) {
	u := req.User

	logrus.WithField("user", u.Login).
		Debugln("requesting system access")

	// check organization membership
	orgs := strings.Split(p.org, ",")

	for _, org := range orgs {
		_, _, err := p.client.Organizations.GetOrgMembership(ctx, u.Login, org)
		if err == nil {
			if u.Login == "obcode" {
				logrus.WithField("user", u.Login).
					WithField("org", p.org).
					WithField("role", "admin").
					Debugln("granted admin system access")
				u.Admin = true
				return &u, err
			}
			logrus.WithField("user", u.Login).
				WithField("org", p.org).
				Debugln("granted standard system access")
			return nil, err
		}
	}

	return nil, ErrAccessDenied
}
