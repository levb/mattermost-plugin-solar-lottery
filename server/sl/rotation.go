// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Rotation struct {
	PluginVersion string
	RotationID    types.ID
	IsArchived    bool
	AutofillType  string
	TaskMaker     *TaskMaker

	MattermostUserIDs *types.IDSet `json:",omitempty"`
	PendingTaskIDs    *types.IDSet `json:",omitempty"`
	InProgressTaskIDs *types.IDSet `json:",omitempty"`

	loaded        bool
	expandedUsers bool
	users         *Users
	expandedTasks bool
	pending       *Tasks
	inProgress    *Tasks
}

func NewRotation() *Rotation {
	r := &Rotation{}
	r.init()
	return r
}

func (r *Rotation) init() {
	if r.MattermostUserIDs == nil {
		r.MattermostUserIDs = types.NewIDSet()
	}
	if r.PendingTaskIDs == nil {
		r.PendingTaskIDs = types.NewIDSet()
	}
	if r.InProgressTaskIDs == nil {
		r.InProgressTaskIDs = types.NewIDSet()
	}
	if r.TaskMaker == nil {
		r.TaskMaker = NewTaskMaker()
	}
	if r.users == nil {
		r.users = NewUsers()
	}
}

func (rotation *Rotation) WithMattermostUserIDs(pool *Users) *Rotation {
	newRotation := *rotation
	newRotation.MattermostUserIDs = types.NewIDSet()
	for _, id := range pool.IDs() {
		newRotation.MattermostUserIDs.Set(id)
	}
	if pool.IsEmpty() {
		pool = NewUsers()
	}
	newRotation.users = pool
	return &newRotation
}

func (r *Rotation) String() string {
	return r.Name()
}

func (r *Rotation) Name() string {
	return kvstore.NameFromID(r.RotationID)
}

func (r *Rotation) Markdown() string {
	return r.Name()
}

func (r *Rotation) MarkdownBullets() string {
	out := fmt.Sprintf("- **%s**\n", r.Name())
	out += fmt.Sprintf("  - ID: `%s`.\n", r.RotationID)
	out += fmt.Sprintf("  - Users (%v): %s.\n", r.MattermostUserIDs.Len(), r.users.MarkdownWithSkills())

	// if rotation.Autopilot.On {
	// 	out += fmt.Sprintf("  - Autopilot: **on**\n")
	// 	out += fmt.Sprintf("    - Auto-start: **%v**\n", rotation.Autopilot.StartFinish)
	// 	out += fmt.Sprintf("    - Auto-fill: **%v**, %v days prior to start\n", rotation.Autopilot.Fill, rotation.Autopilot.FillPrior)
	// 	out += fmt.Sprintf("    - Notify users in advance: **%v**, %v days prior to transition\n", rotation.Autopilot.Notify, rotation.Autopilot.NotifyPrior)
	// } else {
	// 	out += fmt.Sprintf("  - Autopilot: **off**\n")
	// }

	return out
}

func (r *Rotation) FindUsers(mattermostUserIDs *types.IDSet) []*User {
	uu := []*User{}
	for _, id := range r.MattermostUserIDs.IDs() {
		uu = append(uu, r.users.Get(id))
	}
	return uu
}
