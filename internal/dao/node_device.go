package dao

import (
	"github.com/apulis/bmod/aistudio-aom/internal/model"
)

type NodeDevices struct{}

func (d *NodeDevices) List() ([]model.NodeDeviceType, error) {
	var types []model.NodeDeviceType
	db := GetDB().Model(&model.NodeDevice{}).
		Select("key, type, arch, model, compute_type, series, sum(num) as num, sum(avl_num) as avl_num").
		Group("key, type, arch, model, compute_type, series").
		Scan(&types)
	return types, db.Error
}

func (d *NodeDevices) OrgDevs(orgId int64) ([]model.NodeDeviceType, error) {
	var types []model.NodeDeviceType
	db := GetDB().Model(&model.NodeDevice{}).
		Select("key, node_devices.type, arch, model, compute_type, series, sum(num) as num, sum(avl_num) as avl_num").
		Joins("LEFT JOIN nodes ON node_devices.node_id = nodes.id").
		Joins("LEFT JOIN org_resources ON nodes.org_resource_id = org_resources.id").
		Where("org_id = ?", orgId).
		Group("key, node_devices.type, arch, model, compute_type, series").
		Find(&types)
	return types, db.Error
}

func (d *NodeDevices) OrgDevises() ([]model.NodeDeviceTypeOrg, error) {
	var types []model.NodeDeviceTypeOrg
	db := GetDB().Model(&model.NodeDevice{}).
		Select("org_id, key, node_devices.type, arch, model, compute_type, series, sum(num) as num, sum(avl_num) as avl_num").
		Joins("LEFT JOIN nodes ON node_devices.node_id = nodes.id").
		Joins("LEFT JOIN org_resources ON nodes.org_resource_id = org_resources.id").
		Group("org_id, key, node_devices.type, arch, model, compute_type, series").
		Find(&types)
	return types, db.Error
}

func (d *NodeDevices) Get(t, key string) (*model.NodeDeviceType, error) {
	var _type model.NodeDeviceType
	db := GetDB().Model(&model.NodeDevice{}).
		Select("key, type, arch, model, compute_type, series, sum(num) as num, sum(avl_num) as avl_num").
		Where("type = ?", t).
		Where("key = ?", key).
		Group("key, type, arch, model, compute_type, series").
		Find(&_type)
	return &_type, db.Error
}

func (d *NodeDevices) Deletes(ids []int64) error {
	return GetDB().Delete(&model.NodeDevice{}, ids).Error
}
