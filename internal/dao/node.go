package dao

import (
	"database/sql"
	"github.com/apulis/bmod/aistudio-aom/internal/dto"
	"github.com/apulis/bmod/aistudio-aom/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Nodes struct{}

func (n *Nodes) Save(nodes []model.Node) error {
	return GetDB().Session(&gorm.Session{FullSaveAssociations: true}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "machine_id"}},
			DoUpdates: clause.AssignmentColumns(n.upsertColumns()),
		}).CreateInBatches(nodes, 10).Error
}

func (n *Nodes) upsertColumns() []string {
	var columns []string
	columnTypes, _ := GetDB().Migrator().ColumnTypes(&model.Node{})
	for _, c := range columnTypes {
		if c.Name() != "org_resource_id" {
			columns = append(columns, c.Name())
		}
	}
	return columns
}

func (n *Nodes) FindNode(oldNode *model.Node, newNodes []model.Node) int {
	for k, newNode := range newNodes {
		if oldNode.MachineID == newNode.MachineID {
			return k
		}
	}
	return len(newNodes) + 1
}

func (n *Nodes) FindNodeDevice(oldDev *model.NodeDevice, newDevs []model.NodeDevice) int {
	for k, newDev := range newDevs {
		if oldDev.Type == newDev.Type &&
			oldDev.Arch == newDev.Arch &&
			oldDev.Model == newDev.Model &&
			oldDev.ComputeType == newDev.ComputeType &&
			oldDev.Series == newDev.Series {
			return k
		}
	}
	return len(newDevs) + 1
}

func (n *Nodes) Delete(node *model.Node) error {
	return GetDB().Select("Devs").Delete(&node).Error
}
func (n *Nodes) FreeNode(node *model.Node) error {
	return GetDB().Model(&node).Update("org_resource_id", nil).Error
}

func (n *Nodes) Get(id int64) (*model.Node, error) {
	var node model.Node
	db := GetDB().Where("id = ?", id).First(&node)
	return &node, db.Error
}

func (n *Nodes) Gets(ids []int64) ([]model.Node, error) {
	var nodes []model.Node
	db := GetDB().Model(&model.Node{}).Find(&nodes, ids)
	return nodes, db.Error
}

func (n *Nodes) Count(cond *dto.NodeList) (int64, error) {
	var count int64
	db := n.QueryDb(cond)
	db = n.QueryFilter(db, cond).Count(&count)

	return count, db.Error
}

func (n *Nodes) List(cond *dto.NodeList) ([]model.Node, error) {
	db := n.QueryDb(cond).Preload("Devs").Preload("OrgResource")
	var nodes []model.Node
	db = n.QueryFilter(db, cond)
	db = db.Offset((cond.BaseListDto.PageNum - 1) * cond.BaseListDto.PageSize).
		Limit(cond.BaseListDto.PageSize).Order("created_at DESC").
		Find(&nodes)
	return nodes, db.Error
}

func (n *Nodes) QueryDb(cond *dto.NodeList) *gorm.DB {
	subQuery := GetDB().Model(&model.Node{}).
		Select("nodes.*, org_resources.org_id, string_agg(key, '') as keys").
		Joins("LEFT JOIN org_resources ON nodes.org_resource_id = org_resources.id").
		Joins("LEFT JOIN node_devices Devs ON nodes.id = Devs.node_id").
		Group("nodes.id, org_resources.id")
	db := GetDB().Table("(?) as n", subQuery)
	return db
}

func (n *Nodes) QueryFilter(db *gorm.DB, cond *dto.NodeList) *gorm.DB {
	if cond.OrgId != 0 {
		db = db.Where("org_id = ?", cond.OrgId)
	}
	if cond.Type != "" {
		db = db.Where("type = ?", cond.Type)
	}
	if cond.CpuArch != "" {
		db = db.Where("cpu_arch = ?", cond.CpuArch)
	}
	if cond.CpuNum != 0 {
		db = db.Where("cpu_num >= ?", cond.CpuNum)
	}
	if cond.AvlCpuNum != 0 {
		db = db.Where("avl_cpu_num >= ?", cond.AvlCpuNum)
	}
	if cond.Mem != 0 {
		db = db.Where("mem >= ?", cond.Mem)
	}
	if cond.AvlMem != 0 {
		db = db.Where("avl_mem >= ?", cond.AvlMem)
	}
	if cond.Role != "" {
		db = db.Where("role = ?", cond.Role)
	}
	if cond.Status != "" {
		db = db.Where("status = ?", cond.Status)
	}
	if cond.Keyword != "" {
		keyword := "%" + cond.Keyword + "%"
		db = db.Where("name like ? or internal_ip like ? or keys like ?",
			keyword, keyword, keyword)
	}
	return db
}

func (n *Nodes) OrgAvlMem(orgId int64) (int64, error) {
	var sum sql.NullInt64
	db := GetDB().Model(&model.Node{}).
		Select("sum(avl_mem)").
		Joins("LEFT JOIN org_resources ON nodes.org_resource_id = org_resources.id").
		Where("org_id = ?", orgId).
		Scan(&sum)
	return sum.Int64, db.Error
}

func (n *Nodes) OrgAvlMemes() ([]model.OrgAvlMem, error) {
	var sums []model.OrgAvlMem
	db := GetDB().Model(&model.Node{}).
		Select("org_id, sum(avl_mem) as avl_mem").
		Joins("LEFT JOIN org_resources ON nodes.org_resource_id = org_resources.id").
		Group("org_id").
		Find(&sums)
	return sums, db.Error
}

func (n *Nodes) ListAll() ([]model.Node, error) {
	var nodes []model.Node
	db := GetDB().Model(&model.Node{}).Joins("OrgResource")
	db = db.Find(&nodes)
	return nodes, db.Error
}

func (n *Nodes) ListWithDevs() ([]model.Node, error) {
	var nodes []model.Node
	db := GetDB().Model(&model.Node{}).Preload("Devs")
	db = db.Find(&nodes)
	return nodes, db.Error
}

func (n *Nodes) Roles() ([]string, error) {
	var roles []string
	db := GetDB().Model(&model.Node{}).Pluck("role", &roles)
	return roles, db.Error
}
