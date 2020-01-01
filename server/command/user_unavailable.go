// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
)

func (c *Command) userUnavailable(parameters []string) (string, error) {
	var typ, usernames, start, end string
	var clear bool
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVarP(&usernames, flagUsers, flagPUsers, "", "users to show")
	fs.StringVarP(&start, flagStart, flagPStart, "", "start of the unavailability")
	fs.StringVarP(&end, flagEnd, flagPEnd, "", "end of unavailability (last day)")
	fs.BoolVar(&clear, flagClear, false, "clear all overlapping events")
	fs.StringVar(&typ, flagType, store.EventTypeOther, "event type")
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	endTime, err := time.Parse(api.DateFormat, end)
	if err != nil {
		return "", err
	}
	endTime = endTime.Add(time.Hour * 24) // start of next day
	end = endTime.Format(api.DateFormat)

	if clear {
		err = c.API.DeleteEvents(usernames, start, end)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("cleared %s to %s from %s", start, end, usernames), nil
	} else {
		event := store.Event{
			Start: start,
			End:   endTime.Format(api.DateFormat),
			Type:  typ,
		}
		err = c.API.AddEvent(usernames, event)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Added %s to %s", api.MarkdownEvent(event), usernames), nil
	}
}
