package main

import (
	"code.google.com/p/x-go-binding/xgb"
)

var (
	AtomUtf8String,

	AtomWmProtocols,
	AtomWmDeleteWindow,
	AtomWmState,
	AtomWmClientLeader,
	AtomWmTakeFocus,
	AtomWmWindowRole xgb.Id

	AtomNetWmName,
	AtomNetWmWindowType,
	AtomNetWmWindowTypeDock,
	AtomNetWmWindowTypeDialog,
	AtomNetWmWindowTypeUtility,
	AtomNetWmWindowTypeToolbar,
	AtomNetWmWindowTypeSplash,
	AtomNetWmDesktop,
	AtomNetWmStrut,
	AtomNetWmStrutPartial,
	AtomNetWmState,
	AtomNetWmStateModal,
	AtomNetWmStateHidden,
	AtomNetWmStateFullscreen,

	AtomNetSupported,
	AtomNetSupportingWmCheck,

	AtomNetStartupId,
	AtomNetClientListStacking,
	AtomNetActiveWindow,

	AtomNetWorkarea,
	AtomNetNumberOfDesktops,
	AtomNetVirtualRoots,
	AtomNetCurrentDesktop xgb.Id
)

func getAtomId(name string) xgb.Id {
	a, err := conn.InternAtom(false, name)
	if err != nil {
		l.Fatal("Can't get an atom id: ", err)
	}
	return a.Atom
}

func setupAtoms() {
	AtomUtf8String = getAtomId("UTF8_STRING")

	AtomWmProtocols = getAtomId("WM_PROTOCOLS")
	AtomWmDeleteWindow = getAtomId("WM_DELETE_WINDOW")
	AtomWmState = getAtomId("WM_STATE")
	AtomWmClientLeader = getAtomId("WM_CLIENT_LEADER")
	AtomWmTakeFocus = getAtomId("WM_TAKE_FOCUS")
	AtomWmWindowRole = getAtomId("WM_WINDOW_ROLE")

	AtomNetWmName = getAtomId("_NET_WM_NAME")
	AtomNetWmWindowType = getAtomId("_NET_WM_WINDOW_TYPE")
	AtomNetWmWindowTypeDock = getAtomId("_NET_WM_WINDOW_TYPE_DOCK")
	AtomNetWmWindowTypeDialog = getAtomId("_NET_WM_WINDOW_TYPE_DIALOG")
	AtomNetWmWindowTypeUtility = getAtomId("_NET_WM_WINDOW_TYPE_UTILITY")
	AtomNetWmWindowTypeToolbar = getAtomId("_NET_WM_WINDOW_TYPE_TOOLBAR")
	AtomNetWmWindowTypeSplash = getAtomId("_NET_WM_WINDOW_TYPE_SPLASH")
	AtomNetWmDesktop = getAtomId("_NET_WM_DESKTOP")
	AtomNetWmStrut = getAtomId("_NET_WM_STRUT")
	AtomNetWmStrutPartial = getAtomId("_NET_WM_STRUT_PARTIAL")
	AtomNetWmState = getAtomId("_NET_WM_STATE")
	AtomNetWmStateModal = getAtomId("_NET_WM_STATE_MODAL")
	AtomNetWmStateHidden = getAtomId("_NET_WM_STATE_HIDDEN")
	AtomNetWmStateFullscreen = getAtomId("_NET_WM_STATE_FULLSCREEN")

	AtomNetSupported = getAtomId("_NET_SUPPORTED")
	AtomNetSupportingWmCheck = getAtomId("_NET_SUPPORTING_WM_CHECK")

	AtomNetStartupId = getAtomId("_NET_STARTUP_ID")
	AtomNetClientListStacking = getAtomId("_NET_CLIENT_LIST_STACKING")
	AtomNetActiveWindow = getAtomId("_NET_ACTIVE_WINDOW")

	AtomNetWorkarea = getAtomId("_NET_WORKAREA")
	AtomNetNumberOfDesktops = getAtomId("_NET_NUMBER_OF_DESKTOPS")
	AtomNetVirtualRoots = getAtomId("_NET_VIRTUAL_ROOTS")
	AtomNetCurrentDesktop = getAtomId("_NET_CURRENT_DESKTOP")
}
