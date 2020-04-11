// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/pkg/errors"
)

func (sl *sl) LoadTask(taskID types.ID) (*Task, error) {
	task, err := sl.loadTask(taskID)
	if err != nil {
		return nil, err
	}
	err = sl.expandTaskUsers(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (sl *sl) LoadTasks(taskIDs *types.IDSet) (*Tasks, error) {
	tasks := NewTasks()
	for _, id := range taskIDs.IDs() {
		t, err := sl.LoadTask(id)
		if err != nil {
			return nil, err
		}
		tasks.Set(t)
	}
	return tasks, nil
}

func (sl *sl) loadTask(taskID types.ID) (*Task, error) {
	t := NewTask("")
	err := sl.Store.Entity(KeyTask).Load(taskID, t)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load %s", taskID)
	}
	return t, nil
}

func (sl *sl) storeTask(task *Task) error {
	task.PluginVersion = sl.conf.PluginVersion
	err := sl.Store.Entity(KeyTask).Store(task.TaskID, task)
	if err != nil {
		return errors.Wrapf(err, "failed to store %s", task.String())
	}
	return nil
}

func (sl *sl) expandTaskUsers(task *Task) error {
	users, err := sl.LoadUsers(task.MattermostUserIDs)
	if err != nil {
		return err
	}
	task.Users = users
	return nil
}

var allowedAssignTaskStates = map[bool]*types.IDSet{
	false: types.NewIDSet(TaskStatePending),
	true:  types.NewIDSet(TaskStatePending, TaskStateScheduled, TaskStateStarted),
}

func (sl *sl) assignTask(task *Task, users *Users, force bool) (*Users, error) {
	if !allowedAssignTaskStates[force].Contains(task.State) {
		out := "can not "
		if force {
			out += "force "
		}
		return nil, errors.Errorf("%s assign to task in state %s", out, task.State)
	}

	limit := NewNeeds(task.Limit.AsArray()...)
	require := NewNeeds(task.Require.AsArray()...)
	added := NewUsers()
	for _, user := range users.AsArray() {
		if task.MattermostUserIDs.Contains(user.MattermostUserID) {
			continue
		}

		var failed *Needs
		limit, _, failed = limit.CheckLimits(user)
		if !failed.IsEmpty() && !force {
			return nil, errors.Errorf("user %s failed max constraints %s", user.Markdown(), failed.MarkdownSkillLevels())
		}
		require = require.CheckRequired(user)

		task.MattermostUserIDs.Set(user.MattermostUserID)
		if task.Users != nil {
			task.Users.Set(user)
		}
		added.Set(user)
	}
	return added, nil
}

func (sl *sl) unassignTask(task *Task, users *Users, force bool) (*Users, error) {
	if !allowedAssignTaskStates[force].Contains(task.State) {
		out := "can not "
		if force {
			out += "force "
		}
		return nil, errors.Errorf("%s unassign to task in state %s", out, task.State)
	}

	removed := NewUsers()
	for _, user := range users.AsArray() {
		if !task.MattermostUserIDs.Contains(user.MattermostUserID) {
			return nil, errors.Wrapf(kvstore.ErrNotFound,
				"User %s is not assigned to task %s", user.Markdown(), task.Markdown())
		}

		task.MattermostUserIDs.Delete(user.MattermostUserID)
		if task.Users != nil {
			task.Users.Delete(user.MattermostUserID)
		}
		removed.Set(user)
	}
	return removed, nil
}

func (sl *sl) fillTask(r *Rotation, task *Task) (added *Users, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrapf(err, "failed to fill task %s", task.Markdown())
		}
	}()

	// Autofill is only allowed on pending tasks
	if task.State != TaskStatePending {
		return nil, errors.Wrap(ErrWrongState, string(task.State))
	}

	filler, err := sl.taskFiller(r)
	if err != nil {
		return nil, err
	}
	added, err = filler.FillTask(r, task, sl.Logger)
	if err != nil {
		return nil, err
	}

	for _, user := range added.AsArray() {
		task.MattermostUserIDs.Set(user.MattermostUserID)
		task.Users.Set(user)
	}

	return added, nil
}

var validPriorStates = map[types.ID]*types.IDSet{
	TaskStatePending:   types.NewIDSet("none"),
	TaskStateScheduled: types.NewIDSet(TaskStatePending),
	TaskStateStarted:   types.NewIDSet(TaskStateScheduled),
}

func (sl *sl) transitionTask(r *Rotation, task *Task, now types.Time, to types.ID) error {
	if task.State == to {
		return nil
	}

	priorStates, ok := validPriorStates[to]
	if ok && !priorStates.Contains(task.State) {
		return errors.Errorf("fail to transition from %s to %s, only allowed for %s", task.State, to, priorStates.IDs())
	}

	// Touch LastServed for all users.
	for _, user := range task.Users.AsArray() {
		user.LastServed.Set(r.RotationID, now.Unix())
		_, err := sl.storeUserWelcomeNew(user)
		if err != nil {
			return err
		}
	}

	switch to {
	case TaskStatePending:
		sl.announceRotationUsers(r, func(user *User, _ *Rotation) {
			sl.dmUserTaskPending(user, task)
		})
	case TaskStateScheduled:
		sl.announceTaskUsers(task, sl.dmUserTaskScheduled)
	case TaskStateStarted:
		sl.announceTaskUsers(task, sl.dmUserTaskStarted)
	case TaskStateFinished:
		sl.announceTaskUsers(task, sl.dmUserTaskFinished)
	}

	task.State = to
	return sl.storeTask(task)
}
