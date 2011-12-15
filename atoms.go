package main

import (
	"code.google.com/p/x-go-binding/xgb"
)

var (
	AtomNetSupported,
	AtomNetSupportingWmCheck,
	AtomNetWmName,
	AtomNetWmStateFullscreen,
	AtomNetWmState,
	AtomNetWmWindowType,
	AtomNetWmWindowTypeDock,
	AtomNetWmWindowTypeDialog,
	AtomNetWmWindowTypeUtility,
	AtomNetWmWindowTypeToolbar,
	AtomNetWmWindowTypeSplash,
	AtomNetWmDesktop,
	AtomNetWmStrutPartial,
	AtomNetClientListStacking,
	AtomNetCurrentDesktop,
	AtomNetActiveWindow,
	AtomNetWorkarea,
	AtomNetStartupId,
	AtomWmProtocols,
	AtomWmDeleteWindow,
	AtomUtf8String,
	AtomWmState,
	AtomWmClientLeader,
	AtomWmTakeFocus,
	AtomWmWindowRole xgb.Id
)

func getAtomId(name string) xgb.Id {
	a, err := conn.InternAtom(false, name)
	if err != nil {
		l.Fatal("Can't get atom id: ", err)
	}
	return a.Atom
}

func setupAtoms() {
	AtomNetSupported = getAtomId("_NET_SUPPORTED")
	AtomNetSupportingWmCheck = getAtomId("_NET_SUPPORTING_WM_CHECK")
	AtomNetWmName = getAtomId("_NET_WM_NAME")
	AtomNetWmStateFullscreen = getAtomId("_NET_WM_STATE_FULLSCREEN")
	AtomNetWmState = getAtomId("_NET_WM_STATE")
	AtomNetWmWindowType = getAtomId("_NET_WM_WINDOW_TYPE")
	AtomNetWmWindowTypeDock = getAtomId("_NET_WM_WINDOW_TYPE_DOCK")
	AtomNetWmWindowTypeDialog = getAtomId("_NET_WM_WINDOW_TYPE_DIALOG")
	AtomNetWmWindowTypeUtility = getAtomId("_NET_WM_WINDOW_TYPE_UTILITY")
	AtomNetWmWindowTypeToolbar = getAtomId("_NET_WM_WINDOW_TYPE_TOOLBAR")
	AtomNetWmWindowTypeSplash = getAtomId("_NET_WM_WINDOW_TYPE_SPLASH")
	AtomNetWmDesktop = getAtomId("_NET_WM_DESKTOP")
	AtomNetWmStrutPartial = getAtomId("_NET_WM_STRUT_PARTIAL")
	AtomNetClientListStacking = getAtomId("_NET_CLIENT_LIST_STACKING")
	AtomNetCurrentDesktop = getAtomId("_NET_CURRENT_DESKTOP")
	AtomNetActiveWindow = getAtomId("_NET_ACTIVE_WINDOW")
	AtomNetWorkarea = getAtomId("_NET_WORKAREA")
	AtomNetStartupId = getAtomId("_NET_STARTUP_ID")
	AtomWmProtocols = getAtomId("WM_PROTOCOLS")
	AtomWmDeleteWindow = getAtomId("WM_DELETE_WINDOW")
	AtomUtf8String = getAtomId("UTF8_STRING")
	AtomWmState = getAtomId("WM_STATE")
	AtomWmClientLeader = getAtomId("WM_CLIENT_LEADER")
	AtomWmTakeFocus = getAtomId("WM_TAKE_FOCUS")
	AtomWmWindowRole = getAtomId("WM_WINDOW_ROLE")
}
