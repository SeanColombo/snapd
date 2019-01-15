// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2018 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package hotplug

import (
	"fmt"

	"github.com/snapcore/snapd/interfaces/utils"

	"github.com/snapcore/snapd/snap"
)

// Definer can be implemented by interfaces that need to create slots in response to hotplug events.
type Definer interface {
	HotplugDeviceDetected(di *HotplugDeviceInfo, spec *Specification) error
}

// HotplugKeyHandler can be implemented by interfaces that need to provide a non-standard key for hotplug devices.
type HotplugKeyHandler interface {
	HotplugKey(di *HotplugDeviceInfo) (string, error)
}

// HandledByGadgetPredicate can be implemented by hotplug interfaces to decide whether a device is already handled by given gadget slot.
type HandledByGadgetPredicate interface {
	HandledByGadget(di *HotplugDeviceInfo, slot *snap.SlotInfo) bool
}

// RequestedSlotSpec is a definition of the slot to create in response to a hotplug event.
type RequestedSlotSpec struct {
	// Name is how the interface wants to name the slot. When left empty,
	// one will be generated on demand. The hotplug machinery appends a
	// suffix to ensure uniqueness of the name.
	Name  string
	Label string
	Attrs map[string]interface{}
}

// Specification contains a slot definition to create in response to a hotplug event.
type Specification struct {
	slot *RequestedSlotSpec
}

// NewSpecification creates an empty hotplug Specification.
func NewSpecification() *Specification {
	return &Specification{}
}

// SetSlot sets the slot specification.
func (h *Specification) SetSlot(slotSpec *RequestedSlotSpec) error {
	if h.slot != nil {
		return fmt.Errorf("slot specification already created")
	}
	// only validate name if not empty, otherwise name is created by hotplug
	// subsystem later on when the spec is processed.
	if slotSpec.Name != "" {
		if err := snap.ValidateSlotName(slotSpec.Name); err != nil {
			return err
		}
	}
	attrs := slotSpec.Attrs
	if attrs == nil {
		attrs = make(map[string]interface{})
	} else {
		attrs = utils.CopyAttributes(slotSpec.Attrs)
	}
	h.slot = &RequestedSlotSpec{
		Name:  slotSpec.Name,
		Label: slotSpec.Label,
		Attrs: utils.NormalizeInterfaceAttributes(attrs).(map[string]interface{}),
	}
	return nil
}

// Slot returns the specification of the requested slot.
func (h *Specification) Slot() *RequestedSlotSpec {
	return h.slot
}
