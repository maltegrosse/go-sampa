package main

import (
	"fmt"
	"github.com/maltegrosse/go-bird"
	"github.com/maltegrosse/go-sampa"
	"github.com/maltegrosse/go-spa"
	"os"
	"text/tabwriter"
	"time"
)

func main() {
	deltaUt1 := 0.
	deltaT := 66.4
	longitude := 143.36167
	latitude := 24.61167
	elevation := 0.
	pressure := 1000.
	temperature := 11.
	slope := 30.
	azmRotation := -10.
	atmosRefract := 0.5667

	dt := time.Date(2009, 7, 22, 1, 33, 0, 0, time.FixedZone("ManualTimeZone", int(0*3600)))
	sp, err := spa.NewSpa(dt, latitude, longitude, elevation, pressure, temperature, deltaT, deltaUt1, slope, azmRotation, atmosRefract)
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := bird.NewBird(sp.GetZenith(), sp.GetR(), sp.GetPressure(), 0.3, 1.5, 0.07637, 0.85, 0.2, 0.783635)

	if err != nil {
		fmt.Println(err)
		return
	}
	s, err := sampa.NewSampa(sp, b)
	if err != nil {
		fmt.Println(err)
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	_, err = fmt.Fprintln(writer, "- \tNREL sampa_tester.c\tGO Sampa")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fmt.Fprintln(writer, "Julian Day", "\t", "2455034.564583", "\t", fmt.Sprintf("%.6f", s.GetSpaData().GetJd()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fmt.Fprintln(writer, "L", "\t", "299.4024", "\t", fmt.Sprintf("%.6f", s.GetSpaData().GetL()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(writer, "B", "\t", "-0.00001308059", "\t", fmt.Sprintf("%.12f", s.GetSpaData().GetB()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(writer, "R", "\t", "1.016024", "\t", fmt.Sprintf("%.12f", s.GetSpaData().GetR()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(writer, "H", "\t", "344.999100", "\t", fmt.Sprintf("%.12f", s.GetSpaData().GetH()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fmt.Fprintln(writer, "Delta Psi", "\t", "0.004441121", "\t", fmt.Sprintf("%.12f", s.GetSpaData().GetDelPsi()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(writer, "Delta Epsilon", "\t", "0.001203311", "\t", fmt.Sprintf("%.12f", s.GetSpaData().GetDelEpsilon()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(writer, "Epsilon", "\t", "23.439252", "\t", fmt.Sprintf("%.12f", s.GetSpaData().GetEpsilon()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(writer, "Zenith", "\t", "14.512686", "\t", fmt.Sprintf("%.12f", s.GetSpaData().GetZenith()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fmt.Fprintln(writer, "Azimuth", "\t", "104.387917", "\t", fmt.Sprintf("%.12f", s.GetSpaData().GetAzimuth()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fmt.Fprintln(writer, "Angular dist", "\t", "0.374760", "\t", fmt.Sprintf("%.12f", s.GetEms()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(writer, "Sun Radius", "\t", "0.262360", "\t", fmt.Sprintf("%.12f", s.GetRs()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(writer, "Moon Radius", "\t", "0.283341", "\t", fmt.Sprintf("%.12f", s.GetRm()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(writer, "Area unshaded", "\t", "78.363514", "\t", fmt.Sprintf("%.12f", s.GetASulPct()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(writer, "DNI", "\t", "719.099358", "\t", fmt.Sprintf("%.12f", s.GetDniSul()), "\t")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println(err)
		return
	}

}
