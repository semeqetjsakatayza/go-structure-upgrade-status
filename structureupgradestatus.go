package structureupgradestatus

import (
	"log"
)

// StructureUpgrader is a callable to perform structure upgrade.
type StructureUpgrader func(existedRev int32) (structureChanged bool, err error)

// StructureUpgradeStatus tracks the status of structure upgrade operation.
type StructureUpgradeStatus struct {
	Changed   bool
	LastError error

	Logger *log.Logger
}

func (st *StructureUpgradeStatus) logF(format string, v ...interface{}) {
	if st.Logger == nil {
		return
	}
	st.Logger.Printf(format, v...)
}

// RunUpgrade calls given upgrader with given revision information and update
// status with returned result.
// If there are previous error, the upgrade callable will not be invoke.
func (st *StructureUpgradeStatus) RunUpgrade(
	structureName string,
	upgrader StructureUpgrader,
	existedRev int32) (structureChanged bool, err error) {
	if nil != st.LastError {
		return false, st.LastError
	}
	if structureChanged, err = upgrader(existedRev); nil != err {
		st.LastError = err
		st.logF("ERR: cannot upgrade [%s] from [%d]: %v", structureName, existedRev, err)
	}
	if structureChanged {
		st.Changed = true
	}
	return
}

// Call invokes given upgrade function if current status is error free.
// If there are previous error, the upgrade function will not be invoke.
// Given existedRev is for logging purpose.
// The returned shouldStop flag will be true once fnRef returns error.
func (st *StructureUpgradeStatus) Call(
	structureName string,
	fnRef func() (structureChanged bool, err error),
	existedRev int32) (shouldStop, structureChanged bool, err error) {
	if nil != st.LastError {
		return true, false, st.LastError
	}
	if structureChanged, err = fnRef(); nil != err {
		st.LastError = err
		shouldStop = true
		st.logF("ERR: cannot upgrade [%s] from [%d]: %v", structureName, existedRev, err)
	}
	if structureChanged {
		st.Changed = true
	}
	return
}

// PushUpgradeResult updates state with given upgrade result.
//
// Caller should make sure current state is error free before attempt to run
// upgrade without RunUpgrade or Call.
func (st *StructureUpgradeStatus) PushUpgradeResult(structureName string, structureChanged bool, err error) {
	if nil != st.LastError {
		return
	}
	if nil != err {
		st.LastError = err
		st.logF("ERR: cannot upgrade [%s]: %v", structureName, err)
	}
	if structureChanged {
		st.Changed = true
	}
}
