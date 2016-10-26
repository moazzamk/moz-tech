package service

import (
	"github.com/moazzamk/moz-tech/structures"
	"strings"
	"fmt"
)

type SkillParser struct {
	Storage Storage
}

func NewSkillParser(storage Storage) SkillParser {
	return SkillParser{storage}
}

func (r *SkillParser) ParseFromTags(tags *structures.UniqueSlice) *structures.UniqueSlice {
	fmt.Println(`Parsing `, tags)
	return r.processJobSkill(tags)
}

func (r *SkillParser) ParseFromDescription(description string) *structures.UniqueSlice {
	skills := structures.UniqueSlice{}
	description = strings.ToLower(description)
	descriptionSentences := strings.Split(description, `. `)
	for i := range descriptionSentences {
		tmpSkill := make(map[string]int)
		tmp := strings.Split(descriptionSentences[i], ` `)
		for j := range tmp {
			tmp1 := strings.Trim(strings.Replace(tmp[j], `,`, ` `, -1), ` `)
			if len([]rune(tmp1)) >= 3 {
				tmp1 = strings.Trim(r.getNormalizedSkillSynonym(tmp1), ` `)
				tmpSkill[tmp1] = 1
			}
		}

		for j := range tmpSkill {
			if !strings.Contains(j, ` `) && r.Storage.HasSkill(j) {
				skills.Append(j)
			}
		}
	}

	return &skills
}

func (r *SkillParser) processJobSkill(skills *structures.UniqueSlice) *structures.UniqueSlice {
	ret := structures.NewUniqueSlice([]string{})

	for _, value := range skills.ToSlice() {
		tmp := strings.ToLower(strings.Trim(value, ` `))
		tmp = r.getNormalizedSkillSynonym(tmp)
		ret.Append(tmp)

		// If skill is more than 1 word, then check if it has multiple skills listed
		tmpSlice := strings.Split(tmp, ` `)
		tmpSlice = r.removeStopWords(tmpSlice)
		tmpSliceLen := len(tmpSlice)

		for i := range tmpSlice {
			tmp1 := r.getNormalizedSkillSynonym(tmpSlice[i])
			if r.Storage.HasSkill(tmp1) {
				ret.Append(tmp1)
			}
		}

		// If the skill is one word and not present in our storage then add it

		println(`SSSS`, tmp, tmpSlice)
		if tmpSliceLen == 1 && !r.Storage.HasSkill(tmp) {
			_, err := r.Storage.AddSkill(tmp)
			if err != nil {
				panic(err)
			}
		}
	}

	return ret
}

/**
 * Filter out stop words  from skills
 */
func (r *SkillParser) removeStopWords(skills []string) []string {
	ret := []string{}
	stopWords := []string{
		`developer`,
		`development`,
		`programmer`,
		`senior`,
		`software`,
		`on`,
		`and`,
		`or`,
	}

	stopWordsMap := make(map[string]bool)
	for _, val := range stopWords {
		stopWordsMap[val] = true
	}

	for i := 0; i < len(skills); i++ {
		if _, ok := stopWordsMap[skills[i]]; !ok {
			ret = append(ret, skills[i])
		}
	}

	return ret
}


// Correct all spellings, etc of the skill and normalize synonyms to 1 name
func (r *SkillParser) getNormalizedSkillSynonym(skill string) string {
	ret := skill
	synonyms := map[string][]string{
		`mongo`: []string{
			`mongodb`,
			`mongo db`,
		},
		`redhat`: []string{
			`red hat`,
		},
		`javascript`: []string{
			`java script`,
			`jafascript`,
		},
		`angular`: []string{
			`angularjs`,
			`angular.js`,
			`angular js`,
		},
		`ember`: []string{
			`ember.js`,
			`emberjs`,
		},
		`mysql`: []string{
			`my sql`,
		},
		`mssql`: []string{
			`sql server`,
			`ms server`,
		},
		`aws`: []string{
			`amazon web services`,
		},
		`java`: []string{
			`corejava`,
			`core java`,
			`java8`,
		},
		`nodejs`: []string{
			`node js`,
			`node.js`,
		},
		`bootstrap`: []string{
			`boot strap`,
		},
		`bigdata`: []string{
			`big data`,
		},
		`elasticsearch`: []string{
			`elastic search`,
		},
		`machine_learning`: []string{
			`machine learning`,
		},
		`cognitive_computing`: []string{
			`cognitive computing`,
		},
		`cloud_computing`: []string{
			`cloud computing`,
		},
		`data_warehouse`: []string{
			`data warehouse design`,
			`data warehouse`,
			`data warehousing`,
		},
		`automated_testing`: []string{
			`automation test`,
		},
		`data_mining`: []string{
			`data mining`,
		},

		`predictive_analytics`: []string{
			`predictive analytics`,
		},
		`version_control`: []string{
			`version control`,
			`vcs`,
		},
		`business_intelligence`: []string{
			`business_intelligence`,
			` bi `,
			`bi `,
		},
		`azure`: []string{
			`ms azure`,
		},
		`business_analysis`: []string{
			`business analysis`,
			`business analyst`,
		},
		`data_science`: []string{
			`data science`,
			`data scientist`,
		},
	}
	for key, values := range synonyms {
		for i := range values {
			ret = strings.Replace(ret, values[i], key, -1)
		}
	}

	return ret
}
