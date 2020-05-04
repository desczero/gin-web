package service

import (
	"errors"
	"fmt"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/request"
	"go-shipment-api/pkg/utils"
	"strings"
)

// 注意golang坑, 返回值有数组时, 不要为了return方便使用命名返回值
// 如func GetApis(req *request.ApiListRequestStruct) (apis []models.SysApi, err error) {
// 这会导致apis被初始化无穷大的集合, 代码无法断点调试: collecting data..., 几秒后程序异常退出:
// error layer=rpc writing response:write tcp xxx write: broken pipe
// error layer=rpc rpc:read tcp xxx read: connection reset by peer
// exit code 0
// 我曾一度以为调试工具安装配置错误, 使用其他项目代码却能稳定调试, 最终还是定位到代码本身. 踩过的坑希望大家不要再踩
// 获取所有接口
func GetApis(req *request.ApiListRequestStruct) ([]models.SysApi, error) {
	var err error
	list := make([]models.SysApi, 0)
	db := global.Mysql
	method := strings.TrimSpace(req.Method)
	if method != "" {
		db = db.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		db = db.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	// 查询条数
	err = db.Find(&list).Count(&req.PageInfo.Total).Error
	if err == nil {
		// 获取分页参数
		limit, offset := req.GetLimit()
		err = db.Limit(limit).Offset(offset).Find(&list).Error
	}
	return list, err
}

// 创建接口
func CreateApi(req *request.CreateApiRequestStruct) (err error) {
	var api models.SysApi
	utils.Struct2StructByJson(req, &api)
	// 创建数据
	err = global.Mysql.Create(&api).Error
	return
}

// 更新接口
func UpdateApiById(id uint, req *request.CreateApiRequestStruct) (err error) {
	var oldApi models.SysApi
	if global.Mysql.Where("id = ?", id).First(&oldApi).RecordNotFound() {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	var api models.SysApi
	utils.CompareDifferenceStructByJson(req, oldApi, &api)

	// 更新指定列
	err = global.Mysql.Model(&oldApi).UpdateColumns(api).Error
	return
}

// 批量删除接口
func DeleteApiByIds(ids []uint) (err error) {
	return global.Mysql.Where("id IN (?)", ids).Delete(models.SysApi{}).Error
}
