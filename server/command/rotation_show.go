// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"

func (c *Command) rotationShow(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	err := c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}

	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}
	r, err := c.SL.LoadRotation(rotationID)
	if err != nil {
		return "", err
	}

	if c.outputJson {
		return md.JSONBlock(r), nil
	}
	return r.MarkdownBullets(), nil
}
