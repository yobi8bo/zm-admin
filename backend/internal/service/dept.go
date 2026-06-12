package service

import (
	"errors"

	"gorm.io/gorm"
	"zm-project/backend/internal/dto"
	"zm-project/backend/internal/model"
	"zm-project/backend/internal/repository"
	"zm-project/backend/pkg/response"
)

type DeptService struct {
	deptRepo *repository.DeptRepo
}

func NewDeptService(deptRepo *repository.DeptRepo) *DeptService {
	return &DeptService{deptRepo: deptRepo}
}

func (s *DeptService) List() ([]dto.DeptResp, error) {
	depts, err := s.deptRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return buildDeptTree(depts, 0), nil
}

func (s *DeptService) Get(id uint) (*dto.DeptResp, error) {
	d, err := s.deptRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &BizError{Code: response.CodeDeptNotFound}
		}
		return nil, err
	}
	r := toDeptResp(d)
	return &r, nil
}

func (s *DeptService) Create(req *dto.CreateDeptReq) error {
	d := &model.SysDept{
		ParentID: req.ParentID,
		Name:     req.Name,
		Sort:     req.Sort,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	if d.Status == 0 {
		d.Status = 1
	}
	return s.deptRepo.Create(d)
}

func (s *DeptService) Update(id uint, req *dto.UpdateDeptReq) error {
	d, err := s.deptRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BizError{Code: response.CodeDeptNotFound}
		}
		return err
	}
	d.ParentID = req.ParentID
	d.Name = req.Name
	d.Sort = req.Sort
	d.Leader = req.Leader
	d.Phone = req.Phone
	d.Email = req.Email
	d.Status = req.Status
	d.Remark = req.Remark
	return s.deptRepo.Update(d)
}

func (s *DeptService) Delete(id uint) error {
	if _, err := s.deptRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BizError{Code: response.CodeDeptNotFound}
		}
		return err
	}
	hasChildren, err := s.deptRepo.HasChildren(id)
	if err != nil {
		return err
	}
	if hasChildren {
		return &BizError{Code: response.CodeDeptHasChildren}
	}
	hasUsers, err := s.deptRepo.HasUsers(id)
	if err != nil {
		return err
	}
	if hasUsers {
		return &BizError{Code: response.CodeDeptHasUsers}
	}
	return s.deptRepo.Delete(id)
}

func toDeptResp(d *model.SysDept) dto.DeptResp {
	return dto.DeptResp{
		ID:       d.ID,
		ParentID: d.ParentID,
		Name:     d.Name,
		Sort:     d.Sort,
		Leader:   d.Leader,
		Phone:    d.Phone,
		Email:    d.Email,
		Status:   d.Status,
	}
}

func buildDeptTree(depts []model.SysDept, parentID uint) []dto.DeptResp {
	var nodes []dto.DeptResp
	for _, d := range depts {
		if d.ParentID == parentID {
			node := toDeptResp(&d)
			node.Children = buildDeptTree(depts, d.ID)
			nodes = append(nodes, node)
		}
	}
	return nodes
}
