// Copyright 2020-2025 Consensys Software Inc.
// Licensed under the Apache License, Version 2.0. See the LICENSE file for details.

// Code generated by gnark DO NOT EDIT

package cs

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/blang/semver/v4"
)

// WriteTo encodes R1CS into provided io.Writer using cbor
func (cs *system) WriteTo(w io.Writer) (int64, error) {
	b, err := cs.System.ToBytes()
	if err != nil {
		return 0, err
	}

	c := cs.CoeffTable.toBytes()

	totalLen := uint64(len(b) + len(c))
	gnarkVersion := semver.MustParse(cs.GnarkVersion)
	// write totalLen, gnarkVersion.Major, gnarkVersion.Minor, gnarkVersion.Patch using
	// binary.LittleEndian
	if err := binary.Write(w, binary.LittleEndian, totalLen); err != nil {
		return 0, err
	}
	if err := binary.Write(w, binary.LittleEndian, gnarkVersion.Major); err != nil {
		return 0, err
	}
	if err := binary.Write(w, binary.LittleEndian, gnarkVersion.Minor); err != nil {
		return 0, err
	}
	if err := binary.Write(w, binary.LittleEndian, gnarkVersion.Patch); err != nil {
		return 0, err
	}

	// write the system
	n, err := w.Write(b)
	if err != nil {
		return int64(n), err
	}

	// write the coeff table
	m, err := w.Write(c)
	return int64(n+m) + 4*8, err
}

// ReadFrom attempts to decode R1CS from io.Reader using cbor
func (cs *system) ReadFrom(r io.Reader) (int64, error) {
	var totalLen uint64
	if err := binary.Read(r, binary.LittleEndian, &totalLen); err != nil {
		return 0, err
	}

	var major, minor, patch uint64
	if err := binary.Read(r, binary.LittleEndian, &major); err != nil {
		return 0, err
	}
	if err := binary.Read(r, binary.LittleEndian, &minor); err != nil {
		return 0, err
	}
	if err := binary.Read(r, binary.LittleEndian, &patch); err != nil {
		return 0, err
	}
	// TODO @gbotrel validate version, duplicate logic with core.go CheckSerializationHeader
	if major != 0 || minor < 10 {
		return 0, fmt.Errorf("unsupported gnark version %d.%d.%d", major, minor, patch)
	}

	data := make([]byte, totalLen)
	if _, err := io.ReadFull(r, data); err != nil {
		return 0, err
	}
	n, err := cs.System.FromBytes(data)
	if err != nil {
		return 0, err
	}
	data = data[n:]

	if err := cs.CoeffTable.fromBytes(data); err != nil {
		return 0, err
	}

	return int64(totalLen) + 4*8, nil
}
