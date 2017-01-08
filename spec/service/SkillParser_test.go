package service

import (
	. "github.com/moazzamk/moz-tech/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/moazzamk/moz-tech/structures"
	"github.com/moazzamk/moz-tech/mock/service"
)


var tt *testing.T

func TestSkillParser(t *testing.T) {
	tt = t

	RegisterTestingT(t)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Skill Parser Suite")
}

var _ = Describe("SkillParser", func () {

	var mockCtrl *gomock.Controller
	var storage *mock_service.MockStorage

	BeforeEach(func () {
		mockCtrl = gomock.NewController(tt)
		defer mockCtrl.Finish()

		storage = mock_service.NewMockStorage(mockCtrl)
	})

	Context("parses skills from tags", func () {
		It("does nothing if skills already exist", func () {
			slice := structures.NewUniqueSlice([]string{
				`Python `,
				`MySQL `,
				`Software`,
				`software developer`,
				`php python`,
			})

			storage = mock_service.NewMockStorage(mockCtrl)
			storage.EXPECT().HasSkill(`python`).Return(true)
			storage.EXPECT().HasSkill(`python`).Return(true)
			storage.EXPECT().HasSkill(`python`).Return(true)
			storage.EXPECT().HasSkill(`python`).Return(true)

			storage.EXPECT().HasSkill(`mysql`).Return(true)
			storage.EXPECT().HasSkill(`mysql`).Return(true)
			storage.EXPECT().HasSkill(`php`).Return(true)


			skillParser := NewSkillParser(storage)
			_ = skillParser.ParseFromTags(slice)
		})

		It(`splits compound skills and adds any pieces known to be skills`, func () {
			slice := structures.NewUniqueSlice([]string{
				`ruby on rails`,
			})

			storage.EXPECT().HasSkill(`ruby`).Return(false)
			storage.EXPECT().HasSkill(`rails`).Return(true)

			skillParser := NewSkillParser(storage)
			val := skillParser.ParseFromTags(slice)

			Expect(val.ToSlice()).To(HaveLen(2))
		})

		It(`learns skills from tags`, func () {
			slice := structures.NewUniqueSlice([]string{
				`porter `,
				`cable `,
			})

			storage.EXPECT().HasSkill(`porter`).Return(false)
			storage.EXPECT().HasSkill(`porter`).Return(false)
			storage.EXPECT().HasSkill(`cable`).Return(true)
			storage.EXPECT().HasSkill(`cable`).Return(true)
			storage.EXPECT().AddSkill(`porter`)

			skillParser := NewSkillParser(storage)
			_ = skillParser.ParseFromTags(slice)
		})
/*
		It(`splits skills by / and learns the individual skills`, func () {

			slice := structures.NewUniqueSlice([]string{
				`abc/javascript `,
				`polo `,
			})

			storage = mock_service.NewMockStorage(mockCtrl)
			storage.EXPECT().HasSkill(`javascript`).Return(false)
			storage.EXPECT().HasSkill(`polo`).Return(false)
			storage.EXPECT().HasSkill(`polo`).Return(false)
			storage.EXPECT().HasSkill(`abc`).Return(false)
			storage.EXPECT().AddSkill(`javascript`)
			storage.EXPECT().AddSkill(`polo`)
			storage.EXPECT().AddSkill(`abc`)

			skillParser := NewSkillParser(storage)
			_ = skillParser.ParseFromTags(slice)

		})*/
	})
})

