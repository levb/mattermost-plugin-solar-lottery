// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"

	sl "github.com/mattermost/mattermost-plugin-solar-lottery/server/solarlottery"
)

func (c *Command) userUnavailable(parameters []string) (string, error) {
	var usernames, start, end string
	var clear bool
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVarP(&usernames, flagUsers, flagPUsers, "", "users to set unavailability")
	fs.StringVarP(&start, flagStart, flagPStart, "", "start of the unavailability")
	fs.StringVarP(&end, flagEnd, flagPEnd, "", "end of unavailability (last day)")
	fs.BoolVar(&clear, flagClear, false, "clear all overlapping events")
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	startTime, endTime, err := sl.ParseDatePair(start, end)
	if err != nil {
		return "", err
	}
	endTime = endTime.Add(time.Hour * 24) // start of next day
	end = endTime.Format(sl.DateFormat)

	if clear {
		err = c.SL.DeleteEvents(usernames, start, end)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("cleared %s to %s from %s", start, end, usernames), nil
	}

	event := sl.NewPersonalEvent(startTime, endTime)
	err = c.SL.AddEvent(usernames, event)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Added %s to %s", event.Markdown(), usernames), nil
}
