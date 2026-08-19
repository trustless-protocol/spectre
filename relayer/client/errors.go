package client

import "errors"

// ErrPacketAlreadyReceived marks the proof result that is evidence about a
// packet rather than the endpoint: the timeout receipt slot is non-empty, so
// the destination already received it and no timeout can be valid.
var ErrPacketAlreadyReceived = errors.New("packet already received on the destination")
