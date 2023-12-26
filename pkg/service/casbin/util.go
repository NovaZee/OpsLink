package casbin

import "github.com/denovo/permission/protoc/model"

func filterRoles(groups [][]string) (roles []string) {
	for _, group := range groups {
		roles = append(roles, group[1])
	}
	return
}

func filterGModel(data [][]string) (gModes []*model.GModel) {
	for _, row := range data {
		if len(row) >= 2 {
			m := &model.GModel{
				User:  row[0],
				Role:  row[1],
				PType: "g",
			}
			gModes = append(gModes, m)
		}
	}
	return gModes
}

func filterPModel(data [][]string) (PModes []*model.PModel) {
	for _, row := range data {
		if len(row) >= 4 {
			p := &model.PModel{
				UserRole:  row[0],
				Namespace: row[1],
				Source:    row[2],
				Action:    row[3],
				PType:     "p",
			}
			PModes = append(PModes, p)
		}
	}
	return PModes
}
