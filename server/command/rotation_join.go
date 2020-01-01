// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/pkg/errors"
)

func (c *Command) joinRotation(parameters []string) (string, error) {
	var rotationID, rotationName string
	users := ""
	graceShifts := 0
	fs := newRotationFlagSet(&rotationID, &rotationName)
	fs.StringVarP(&users, flagUsers, flagPUsers, "", "add nother users to rotation.")
	fs.IntVar(&graceShifts, flagGrace, 0, "start with N grace shifts.")
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	rotationID, err = c.parseRotationFlags(rotationID, rotationName)
	if err != nil {
		return "", err
	}
	rotation, err := c.API.LoadRotation(rotationID)
	if err != nil {
		return "", err
	}

	added, err := c.API.JoinRotation(users, rotation, graceShifts, time.Now())
	if err != nil {
		return "", errors.WithMessagef(err, "failed, %s might have been updated", api.MarkdownUserMap(added))
	}

	return fmt.Sprintf("%s joined rotation %s", api.MarkdownUserMap(added), rotation.Name), nil
}
