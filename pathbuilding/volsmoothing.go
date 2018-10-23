package pathbuilding

import (
	//"fmt"
	"io"
	tp "typedef"
)

//Volsmoothing smooth volume
func Volsmoothing(par *tp.Par, links map[int]*tp.Link, wr io.Writer) {
	percentf := float64(par.SmoothForwardPercent) / 100
	percentb := float64(par.SmoothBackworkPercent) / 100
	mid := 1 - percentf - percentb
	var tempvol *[96]float64

	for _, link := range links {
		for k := range link.VolCls {
			tempvol = link.VolCls[k]
			for j := 1; j <= par.SmoothIter; j++ {
				var newvol [96]float64
				for i := 1; i <= 96; i++ {
					var ip int
					var ia int
					if i == 1 {
						ip = 96
						ia = i + 1
					} else if i == 96 {
						ip = i - 1
						ia = 1
					} else {
						ip = i - 1
						ia = i + 1
					}
					newvol[i-1] = tempvol[ip-1]*percentf + tempvol[i-1]*mid + tempvol[ia-1]*percentb
				}
				tempvol = &newvol
			}
			for i := 1; i <= 96; i++ {
				link.VolCls[k][i-1] = tempvol[i-1]
			}
		}
		//fmt.Fprintf(wr, "link new vol %v\n", link.Vol)
	}
}
