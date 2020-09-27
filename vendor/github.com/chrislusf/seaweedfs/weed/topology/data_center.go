package topology

import "github.com/chrislusf/seaweedfs/weed/pb/master_pb"

type DataCenter struct {
	NodeImpl
}

func NewDataCenter(id string) *DataCenter {
	dc := &DataCenter{}
	dc.id = NodeId(id)
	dc.nodeType = "DataCenter"
	dc.children = make(map[NodeId]Node)
	dc.NodeImpl.value = dc
	return dc
}

func (dc *DataCenter) GetOrCreateRack(rackName string) *Rack {
	for _, c := range dc.Children() {
		rack := c.(*Rack)
		if string(rack.Id()) == rackName {
			return rack
		}
	}
	rack := NewRack(rackName)
	dc.LinkChildNode(rack)
	return rack
}

func (dc *DataCenter) ToMap() interface{} {
	m := make(map[string]interface{})
	m["Id"] = dc.Id()
	m["Max"] = dc.GetMaxVolumeCount()
	m["Free"] = dc.FreeSpace()
	var racks []interface{}
	for _, c := range dc.Children() {
		rack := c.(*Rack)
		racks = append(racks, rack.ToMap())
	}
	m["Racks"] = racks
	return m
}

func (dc *DataCenter) ToDataCenterInfo() *master_pb.DataCenterInfo {
	m := &master_pb.DataCenterInfo{
		Id:                string(dc.Id()),
		VolumeCount:       uint64(dc.GetVolumeCount()),
		MaxVolumeCount:    uint64(dc.GetMaxVolumeCount()),
		FreeVolumeCount:   uint64(dc.FreeSpace()),
		ActiveVolumeCount: uint64(dc.GetActiveVolumeCount()),
		RemoteVolumeCount: uint64(dc.GetRemoteVolumeCount()),
	}
	for _, c := range dc.Children() {
		rack := c.(*Rack)
		m.RackInfos = append(m.RackInfos, rack.ToRackInfo())
	}
	return m
}
