// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) rotationSet(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		commandAutopilot: c.rotationSetAutopilot,
		commandFill:      c.rotationSetFill,
		commandLimit:     c.rotationSetLimit,
		commandRequire:   c.rotationSetRequire,
		commandTask:      c.rotationSetTask,
	}

	return c.handleCommand(subcommands, parameters)
}

func (c *Command) rotationSetAutopilot(parameters []string) (md.MD, error) {
	off := c.assureFS().Bool("off", false, "turn off")
	create := c.assureFS().Bool("create", false, "create shifts automatically")
	createPrior := c.assureFS().Duration("create-prior", 0, "create shifts this long before their scheduled start")
	schedule := c.assureFS().Bool("schedule", false, "create shifts automatically")
	schedulePrior := c.assureFS().Duration("schedule-prior", 0, "fill and schedule shifts this long before their scheduled start")
	startFinish := c.assureFS().Bool("start-finish", false, "start and finish scheduled tasks")
	remindStart := c.assureFS().Bool("remind-start", false, "remind shift users prior to start")
	remindStartPrior := c.assureFS().Duration("remind-start-prior", 0, "remind shift users this long before the shift's start")
	remindFinish := c.assureFS().Bool("remind-finish", false, "remind shift users prior to finish")
	remindFinishPrior := c.assureFS().Duration("remind-finish-prior", 0, "remind shift users this long before the shift's finish")
	err := c.parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	var out md.Markdowner
	switch {
	case *off:
		out, err = c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
			r.AutopilotSettings = sl.AutopilotSettings{}
			return nil
		})

	default:
		out, err = c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
			r.AutopilotSettings.Create = *create
			r.AutopilotSettings.CreatePrior = *createPrior
			r.AutopilotSettings.Schedule = *schedule
			r.AutopilotSettings.SchedulePrior = *schedulePrior
			r.AutopilotSettings.StartFinish = *startFinish
			r.AutopilotSettings.RemindStart = *remindStart
			r.AutopilotSettings.RemindStartPrior = *remindStartPrior
			r.AutopilotSettings.RemindFinish = *remindFinish
			r.AutopilotSettings.RemindFinishPrior = *remindFinishPrior
			return nil
		})
	}

	return c.normalOut(out, err)
}

func (c *Command) rotationSetFill(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	seed := c.assureFS().Int64("seed", intNoValue, "seed to use")
	fuzz := c.assureFS().Int64("fuzz", intNoValue, `increase fill randomness`)
	err := c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
			if *seed != intNoValue {
				r.FillSettings.Seed = *seed
			}
			if *fuzz != intNoValue {
				r.FillSettings.Fuzz = *fuzz
			}
			return nil
		}))
}

func (c *Command) rotationSetRequire(parameters []string) (md.MD, error) {
	return c.rotationSetNeed(true, parameters)
}

func (c *Command) rotationSetLimit(parameters []string) (md.MD, error) {
	return c.rotationSetNeed(false, parameters)
}

func (c *Command) rotationSetNeed(require bool, parameters []string) (md.MD, error) {
	c.withFlagRotation()
	var skillLevel sl.SkillLevel
	c.assureFS().VarP(&skillLevel, "skill", "s", "skill, with optional level (1-4) as in `--skill=web-3`.")
	count := c.assureFS().Int("count", 1, "number of users")
	clear := c.assureFS().Bool("clear", false, "remove the skill from the list")
	err := c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
			needsToUpdate := r.TaskSettings.Limit
			if require {
				needsToUpdate = r.TaskSettings.Require
			}
			if *clear {
				needsToUpdate.Delete(skillLevel.AsID())
			} else {
				needsToUpdate.SetCountForSkillLevel(skillLevel, int64(*count))
			}
			return nil
		}))
}

func (c *Command) rotationSetTask(parameters []string) (md.MD, error) {
	c.withFlagRotation()
	dur := c.assureFS().Duration("duration", 0, "duration")
	grace := c.assureFS().Duration("grace", 0, "grace period after finishing a task")
	err := c.fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}

	rotationID, err := c.resolveRotation()
	if err != nil {
		return "", err
	}

	return c.normalOut(
		c.SL.UpdateRotation(rotationID, func(r *sl.Rotation) error {
			if *dur != 0 {
				r.TaskSettings.Duration = *dur
			}
			if *grace != 0 {
				r.TaskSettings.Grace = *grace
			}
			return nil
		}))
}
