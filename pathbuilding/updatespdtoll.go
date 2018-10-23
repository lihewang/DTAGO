package pathbuilding

import (
	//"fmt"
	"io"
	"math"
	tp "typedef"
)

//UpdateLink update link speed and toll
func UpdateLink(par *tp.Par, links map[int]*tp.Link, wr io.Writer) {
	for _, link := range links {
		//loop time steps
		for i := 1; i <= 96; i++ {
			//total vol
			for _, cls := range par.VehClass { //loop vehicle classes
				if par.Iter > 1 { //MSA
					link.VolCls[cls][i-1] = link.VolClsPre[cls][i-1]*(1-1/float64(par.Iter)) + link.VolCls[cls][i-1]/float64(par.Iter)
				}
				link.VolClsPre[cls][i-1] = link.VolCls[cls][i-1]
				link.Vol[i-1] = link.Vol[i-1] + link.VolCls[cls][i-1]
			}
			link.CgSpeed[i-1] = link.FFSPEED / (1 + link.ALPHA*math.Pow(link.Vol[i-1]/(float64(link.CAPACITY)/4), link.BETA))
			link.CgTime[i-1] = link.DIST / link.CgSpeed[i-1]
			//update
			policy := link.TOLLPOLICY
			if policy > 0 {
				tl := par.TollPolicy[policy]["MinToll"] + (par.TollPolicy[policy]["MaxToll"]-par.TollPolicy[policy]["MinToll"])*
					math.Pow(link.Vol[i-1]/(float64(link.CAPACITY)/4)+par.TollPolicy[policy]["Offset"], par.TollPolicy[policy]["Exp"])
				if tl > par.TollPolicy[policy]["MaxToll"] {
					tl = par.TollPolicy[policy]["MaxToll"]
				}
				link.TollRate[i-1] = tl
			}
		}
	}
	par.VehClass = nil
}
