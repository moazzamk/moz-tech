package service

import (
	. "github.com/moazzamk/moz-tech/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	//"github.com/moazzamk/moz-tech/structures"

	"github.com/golang/mock/gomock"
	"github.com/moazzamk/moz-tech/mock"
	"github.com/moazzamk/moz-tech/structures"
)

var mockCtrl *gomock.Controller
var storage *mock_service.MockStorage

func TestSkillParser(t *testing.T) {

	mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	storage = mock_service.NewMockStorage(mockCtrl)

	RegisterTestingT(t)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Salary Parser Suite")
}

var _ = Describe("SkillParser", func () {
	Context("parses skills from tags", func () {
		It("parses abbreviated salaries", func () {
			slice := structures.NewUniqueSlice([]string{
				`Python`,
/*				`MySQL`,
				`DJango`,
				`javascript`,
				`ios`,
				`android`,
				`aws`,
				`iot`,*/
			})

			storage.EXPECT().HasSkill(`python`).Return(true)
			storage.EXPECT().HasSkill(`python`).Return(true)


			skillParser := NewSkillParser(storage)
			_ = skillParser.ParseFromTags(slice)
		})
	})
})

