package dictionary

import (
	"context"

	readerClient "github.com/web-rabis/reader-client/client"
	readerModel "github.com/web-rabis/reader-client/model"
)

type IManager interface {
	AcademicDegreeList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error)
	EducationList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error)
	FacultyList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error)
	InstitutionList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error)
	NationalityList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error)
	SexList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error)
	SocialStatusList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.SocialStatus, error)
	SpecialityList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error)
	TypeCardList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.TypeCard, error)
	TypeCardById(ctx context.Context, id int64) (*readerModel.TypeCard, error)
}
type Manager struct {
	dictClient readerClient.DictionaryService
}

func NewManager(dictClient readerClient.DictionaryService) *Manager {
	return &Manager{
		dictClient: dictClient,
	}
}

func (m *Manager) AcademicDegreeList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error) {
	return m.dictClient.AcademicDegreeList(ctx, paging)
}
func (m *Manager) EducationList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error) {
	return m.dictClient.EducationList(ctx, paging)
}
func (m *Manager) FacultyList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error) {
	return m.dictClient.FacultyList(ctx, paging)
}
func (m *Manager) InstitutionList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error) {
	return m.dictClient.InstitutionList(ctx, paging)
}
func (m *Manager) NationalityList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error) {
	return m.dictClient.NationalityList(ctx, paging)
}
func (m *Manager) SexList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error) {
	return m.dictClient.SexList(ctx, paging)
}
func (m *Manager) SocialStatusList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.SocialStatus, error) {
	return m.dictClient.SocialStatusList(ctx, paging)
}
func (m *Manager) SpecialityList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.Dictionary, error) {
	return m.dictClient.SpecialityList(ctx, paging)
}
func (m *Manager) TypeCardList(ctx context.Context, paging *readerModel.Paging) (int64, []*readerModel.TypeCard, error) {
	return m.dictClient.TypeCardList(ctx, paging)
}
func (m *Manager) TypeCardById(ctx context.Context, id int64) (*readerModel.TypeCard, error) {
	return m.dictClient.TypeCardById(ctx, id)
}
