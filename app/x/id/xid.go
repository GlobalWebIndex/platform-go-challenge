package id

import "github.com/rs/xid"

type XID = xid.ID

func XiD() XID {
	// github.com/rs/xid guid := xid.New() xid: 12 bytes, 20 chars, configuration free, sortable
	// guid := xid.New() guid.Machine() guid.Pid() guid.Time() guid.Counter() _ = xid.NilID()
	return xid.New()
}
