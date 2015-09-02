package equipment

import S "bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/structures"

// GetDslamByID GET DSLAM by ID in data file
func GetDslamByID(data S.Data, id string) S.DSLAM {

	for i := 0; i < len(data.DSLAM); i++ {
		if data.DSLAM[i].ID == id {
			return data.DSLAM[i]
		}
	}
	var null S.DSLAM
	return null

}

// GetDslamPosByID GET DSLAM Position by ID in data file
func GetDslamPosByID(data S.Data, id string) int {

	for i := 0; i < len(data.DSLAM); i++ {
		if data.DSLAM[i].ID == id {
			return i
		}
	}
	return -1

}
