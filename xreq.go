// Copyright 2014 Garrett D'Amore
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sp

import ()

// XReq is an implementation of the XReq protocol.
type XReq struct {
	handle ProtocolHandle
}

// Init implements the Protocol Init method.
func (p *XReq) Init(handle ProtocolHandle) {
	p.handle = handle
}

// Process implements the Protocol Process method.
func (p *XReq) Process() {

	h := p.handle

	if msg := h.PullDown(); msg != nil {
		// Send sends unmolested.  If we can't due to lack of a
		// connected peer, we drop it.  (Req protocol resends, but
		// we don't in XReq.)  Note that it is expected that the
		// application will have written the request ID into the
		// header at minimum, but possibly a full backtrace.  We
		// don't bother to check.  (XXX: Perhaps we should, and
		// drop any message that lacks at least a minimal header?)
		h.Send(msg)
	}

	if msg, _, _ := h.Recv(); msg != nil {
		// When we receive a message, we expect to have the request
		// ID in the header.  We strip that out into the header.
		if msg.trimUint32() == nil {
			// If app can't receive (should never happen), just
			// drop it.  App will need to resend request.
			h.PushUp(msg)
		}
	}
}

// Name implements the Protocol Name method.  It returns "XReq".
func (*XReq) Name() string {
	return "XReq"
}

// Number implements the Protocol Number method.
func (*XReq) Number() uint16 {
	return ProtoReq
}

// IsRaw implements the Protocol Raw method.
func (*XReq) IsRaw() bool {
	return true
}

// ValidPeer implements the Protocol ValidPeer method.
func (*XReq) ValidPeer(peer uint16) bool {
	if peer == ProtoRep {
		return true
	}
	return false
}

// RecvHook implements the Protocol RecvHook method.  It is a no-op.
func (*XReq) RecvHook(*Message) bool {
	return true
}

// SendHook implements the Protocol SendHook method.  It is a no-op.
func (*XReq) SendHook(*Message) bool {
	return true
}
