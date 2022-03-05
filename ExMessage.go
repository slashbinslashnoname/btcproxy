package main

import (
	"bytes"
	"encoding/binary"
)

// ex-message的magic number
const ExMessageMagicNumber uint8 = 0x7F

// types
const (
	CMD_REGISTER_WORKER            uint8 = 0x01 // Agent -> Pool
	CMD_SUBMIT_SHARE               uint8 = 0x02 // Agent -> Pool,  mining.submit(...)
	CMD_SUBMIT_SHARE_WITH_TIME     uint8 = 0x03 // Agent -> Pool,  mining.submit(..., nTime)
	CMD_UNREGISTER_WORKER          uint8 = 0x04 // Agent -> Pool
	CMD_MINING_SET_DIFF            uint8 = 0x05 // Pool  -> Agent, mining.set_difficulty(diff)
	CMD_SUBMIT_RESPONSE            uint8 = 0x10 // Pool  -> Agent, response of the submit (optional)
	CMD_SUBMIT_SHARE_WITH_VER      uint8 = 0x12 // Agent -> Pool,  mining.submit(..., nVersionMask)
	CMD_SUBMIT_SHARE_WITH_TIME_VER uint8 = 0x13 // Agent -> Pool,  mining.submit(..., nTime, nVersionMask)
	CMD_SUBMIT_SHARE_WITH_MIX_HASH uint8 = 0x14 // Agent -> Pool, for ETH
	CMD_SET_EXTRA_NONCE            uint8 = 0x22 // Pool  -> Agent, pool nonce prefix allocation result (Ethereum)
)

type SerializableExMessage interface {
	Serialize() []byte
}

type UnserializableExMessage interface {
	Unserialize(data []byte) (err error)
}

type ExMessageHeader struct {
	Type uint8
	Size uint16
}

type ExMessage struct {
	ExMessageHeader
	Body []byte
}

type ExMessageRegisterWorker struct {
	SessionID   uint16
	ClientAgent string
	WorkerName  string
}

func (msg *ExMessageRegisterWorker) Serialize() []byte {

	buf := new(bytes.Buffer)
	return buf.Bytes()
}

type ExMessageUnregisterWorker struct {
	SessionID uint16
}

func (msg *ExMessageUnregisterWorker) Serialize() []byte {
	buf := new(bytes.Buffer)
	return buf.Bytes()
}

type ExMessageSubmitShareBTC struct {
	Login          string
	JobID          string
	ExtraNonce2    string
	Time           string
	Nonce          string
	VersionMask    string
	HasVersionMask bool

	IsFakeJob bool
	SessionID uint16
}

func (msg *ExMessageSubmitShareBTC) Serialize() []byte {

	buf := new(bytes.Buffer)

	return buf.Bytes()
}

type ExMessageMiningSetDiff struct {
	Base struct {
		DiffExp uint8
		Count   uint16
	}
	SessionIDs []uint16
}

func (msg *ExMessageMiningSetDiff) Unserialize(data []byte) (err error) {
	buf := bytes.NewReader(data)

	err = binary.Read(buf, binary.LittleEndian, &msg.Base)
	if err != nil || msg.Base.Count == 0 {
		return
	}

	msg.SessionIDs = make([]uint16, msg.Base.Count)
	err = binary.Read(buf, binary.LittleEndian, msg.SessionIDs)
	return
}

type ExMessageSubmitResponse struct {
	Index  uint16
	Status StratumStatus
}

func (msg *ExMessageSubmitResponse) Unserialize(data []byte) (err error) {
	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.LittleEndian, msg)
	return
}

type ExMessageSetExtranonce struct {
	SessionID  uint16
	ExtraNonce uint32
}

func (msg *ExMessageSetExtranonce) Unserialize(data []byte) (err error) {
	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.LittleEndian, msg)
	return
}
