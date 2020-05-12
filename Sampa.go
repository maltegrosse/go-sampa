package sampa

import (
	"github.com/maltegrosse/go-bird"
	"github.com/maltegrosse/go-spa"
	"math"
)

///////////////////////////////////////////////
//                                           //
// Solar and Moon Position Algorithm (SAMPA) //
//                   for                     //
//        Solar Radiation Application        //
//                                           //
//              August 1, 2012               //
//                                           //
//                                           //
//   Afshin Michael Andreas                  //
//   Afshin.Andreas@NREL.gov (303)384-6383   //
//                                           //
//   Solar Resource and Forecasting Group    //
//   Solar Radiation Research Laboratory     //
//   National Renewable Energy Laboratory    //
//   15013 Denver W Pkwy, Golden, CO 80401   //
//                                           //
//  This code is based on the NREL           //
//  technical report "Solar Eclipse          //
//  Monitoring for Solar Energy Applications //
//  using the Solar and Moon Position 		 //
//  Algorithms" by Ibrahim Reda              //
///////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
//
//   NOTICE
//   Copyright (C) 2012 the Alliance for Sustainable Energy, LLC, All Rights Reserved
//
//This computer software is prepared by the Alliance for Sustainable Energy, LLC, (hereinafter
//the "Contractor"), under Contract DE-AC36-08GO28308 ("Contract") with the Department of
//Energy ("DOE"). The United States Government has been granted for itself and others acting
//on its behalf a paid-up, non-exclusive, irrevocable, worldwide license in the Software to
//reproduce, prepare derivative works, and perform publicly and display publicly. Beginning
//five (5) years after the date permission to assert copyright is obtained from DOE, and subject
//to any subsequent five (5) year renewals, the United States Government is granted for itself
//and others acting on its behalf a paid-up, non-exclusive, irrevocable, worldwide license in
//the Software to reproduce, prepare derivative works, distribute copies to the public, perform
//publicly and display publicly, and to permit others to do so. If the Contractor ceases to make
//this computer software available, it may be obtained from DOE's Office of Scientific and
//Technical Information's Energy Science and Technology Software Center (ESTSC) at PO Box 1020,
//Oak Ridge, TN 37831-1020. THIS SOFTWARE IS PROVIDED BY THE CONTRACTOR "AS IS" AND ANY EXPRESS
//OR IMPLIED WARRANTIES, INCLUDING BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
//AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE CONTRACTOR, DOE, OR
//THE U.S GOVERNMENT BE LIABLE FOR ANY SPECIAL, INDIRECT OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
//WHATSOEVER, INCLUDING BUT NOT LIMITED TO CLAIMS ASSOCIATED WITH THE LOSS OF DATA OR PROFITS,
//WHICH MAY RESULT FROM AN ACTION IN CONTRACT, NEGLIGENCE OR OTHER TORTIOUS CLAIM THAT ARISES
//OUT OF OR IN CONNECTION WITH THE ACCESS, USE OR PERFORMANCE OF THIS SOFTWARE.
//
//The software is being provided for internal, noncommercial purposes only and shall not be
//re-distributed. Please contact Jennifer Ramsey (Jennifer.Ramsey@nrel.gov) in the NREL
//Commercialization and Technology Transfer Office for information concerning a commercial
//license to use the Software.
//
//As a condition of using the software in an application, the developer of the application
//agrees to reference the use of the software and make this notice readily accessible to any
//end-user in a Help|About screen or equivalent manner.
//
///////////////////////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////////////////////
// Revised 20-SEPT-2012 Andreas
//         Modified a_sul and a_sul_pct such that values are reported when no eclipse is occurring
//         Set a_sul to zero when result was negative due to moon radius being larger than sun radius
//         Added call to SERI/NREL BIRD Clear Sky Model to estimate values for irradiances
//         Modified sampa_data structure to include values for estimated irradiance from BIRD model
//         Added a "function" input variable that allows the selecting of desired outputs
// Revised 08-SEPT-2014 Andreas
//         Changed all variables names from azimuth180 to azimuth_astro for consistency with SPA
//         See SPA.H header file for change in results of azimuth_astro
///////////////////////////////////////////////////////////////////////////////////////////////

//enumeration for function codes to select desired final outputs from SAMPA
const (
	SampaNoIrr = 0 //calculate all values except estimated solar irradiances
	SampaAll   = 1 //calculate all values
)

var COUNT = 60

const (
	TermD     = 0
	TermM     = 1
	TermMpr   = 2
	TermF     = 3
	TermLb    = 4
	TermR     = 5
	TermCount = 6
)

///////////////////////////////////////////////////////
///  Moon's Periodic Terms for Longitude and Distance
///////////////////////////////////////////////////////
var MlTerms = [][]float64{
	{0, 0, 1, 0, 6288774, -20905355},
	{2, 0, -1, 0, 1274027, -3699111},
	{2, 0, 0, 0, 658314, -2955968},
	{0, 0, 2, 0, 213618, -569925},
	{0, 1, 0, 0, -185116, 48888},
	{0, 0, 0, 2, -114332, -3149},
	{2, 0, -2, 0, 58793, 246158},
	{2, -1, -1, 0, 57066, -152138},
	{2, 0, 1, 0, 53322, -170733},
	{2, -1, 0, 0, 45758, -204586},
	{0, 1, -1, 0, -40923, -129620},
	{1, 0, 0, 0, -34720, 108743},
	{0, 1, 1, 0, -30383, 104755},
	{2, 0, 0, -2, 15327, 10321},
	{0, 0, 1, 2, -12528, 0},
	{0, 0, 1, -2, 10980, 79661},
	{4, 0, -1, 0, 10675, -34782},
	{0, 0, 3, 0, 10034, -23210},
	{4, 0, -2, 0, 8548, -21636},
	{2, 1, -1, 0, -7888, 24208},
	{2, 1, 0, 0, -6766, 30824},
	{1, 0, -1, 0, -5163, -8379},
	{1, 1, 0, 0, 4987, -16675},
	{2, -1, 1, 0, 4036, -12831},
	{2, 0, 2, 0, 3994, -10445},
	{4, 0, 0, 0, 3861, -11650},
	{2, 0, -3, 0, 3665, 14403},
	{0, 1, -2, 0, -2689, -7003},
	{2, 0, -1, 2, -2602, 0},
	{2, -1, -2, 0, 2390, 10056},
	{1, 0, 1, 0, -2348, 6322},
	{2, -2, 0, 0, 2236, -9884},
	{0, 1, 2, 0, -2120, 5751},
	{0, 2, 0, 0, -2069, 0},
	{2, -2, -1, 0, 2048, -4950},
	{2, 0, 1, -2, -1773, 4130},
	{2, 0, 0, 2, -1595, 0},
	{4, -1, -1, 0, 1215, -3958},
	{0, 0, 2, 2, -1110, 0},
	{3, 0, -1, 0, -892, 3258},
	{2, 1, 1, 0, -810, 2616},
	{4, -1, -2, 0, 759, -1897},
	{0, 2, -1, 0, -713, -2117},
	{2, 2, -1, 0, -700, 2354},
	{2, 1, -2, 0, 691, 0},
	{2, -1, 0, -2, 596, 0},
	{4, 0, 1, 0, 549, -1423},
	{0, 0, 4, 0, 537, -1117},
	{4, -1, 0, 0, 520, -1571},
	{1, 0, -2, 0, -487, -1739},
	{2, 1, 0, -2, -399, 0},
	{0, 0, 2, -2, -381, -4421},
	{1, 1, 1, 0, 351, 0},
	{3, 0, -2, 0, -340, 0},
	{4, 0, -3, 0, 330, 0},
	{2, -1, 2, 0, 327, 0},
	{0, 2, 1, 0, -323, 1165},
	{1, 1, -1, 0, 299, 0},
	{2, 0, 3, 0, 294, 0},
	{2, 0, -1, -2, 0, 8752}}

///////////////////////////////////////////////////////
///  Moon's Periodic Terms for Latitude
///////////////////////////////////////////////////////
var MbTerms = [][]float64{
	{0, 0, 0, 1, 5128122, 0},
	{0, 0, 1, 1, 280602, 0},
	{0, 0, 1, -1, 277693, 0},
	{2, 0, 0, -1, 173237, 0},
	{2, 0, -1, 1, 55413, 0},
	{2, 0, -1, -1, 46271, 0},
	{2, 0, 0, 1, 32573, 0},
	{0, 0, 2, 1, 17198, 0},
	{2, 0, 1, -1, 9266, 0},
	{0, 0, 2, -1, 8822, 0},
	{2, -1, 0, -1, 8216, 0},
	{2, 0, -2, -1, 4324, 0},
	{2, 0, 1, 1, 4200, 0},
	{2, 1, 0, -1, -3359, 0},
	{2, -1, -1, 1, 2463, 0},
	{2, -1, 0, 1, 2211, 0},
	{2, -1, -1, -1, 2065, 0},
	{0, 1, -1, -1, -1870, 0},
	{4, 0, -1, -1, 1828, 0},
	{0, 1, 0, 1, -1794, 0},
	{0, 0, 0, 3, -1749, 0},
	{0, 1, -1, 1, -1565, 0},
	{1, 0, 0, 1, -1491, 0},
	{0, 1, 1, 1, -1475, 0},
	{0, 1, 1, -1, -1410, 0},
	{0, 1, 0, -1, -1344, 0},
	{1, 0, 0, -1, -1335, 0},
	{0, 0, 3, 1, 1107, 0},
	{4, 0, 0, -1, 1021, 0},
	{4, 0, -1, 1, 833, 0},
	{0, 0, 1, -3, 777, 0},
	{4, 0, -2, 1, 671, 0},
	{2, 0, 0, -3, 607, 0},
	{2, 0, 2, -1, 596, 0},
	{2, -1, 1, -1, 491, 0},
	{2, 0, -2, 1, -451, 0},
	{0, 0, 3, -1, 439, 0},
	{2, 0, 2, 1, 422, 0},
	{2, 0, -3, -1, 421, 0},
	{2, 1, -1, 1, -366, 0},
	{2, 1, 0, 1, -351, 0},
	{4, 0, 0, 1, 331, 0},
	{2, -1, 1, 1, 315, 0},
	{2, -2, 0, -1, 302, 0},
	{0, 0, 1, 3, -283, 0},
	{2, 1, 1, -1, -229, 0},
	{1, 1, 0, -1, 223, 0},
	{1, 1, 0, 1, 223, 0},
	{0, 1, -2, -1, -220, 0},
	{2, 1, -1, -1, -220, 0},
	{1, 0, 1, 1, -185, 0},
	{2, -1, -2, -1, 181, 0},
	{0, 1, 2, 1, -177, 0},
	{4, 0, -2, -1, 176, 0},
	{4, -1, -1, -1, 166, 0},
	{1, 0, 1, -1, -164, 0},
	{4, 0, 1, -1, 132, 0},
	{1, 0, -1, -1, -119, 0},
	{4, -1, 0, -1, 115, 0},
	{2, -2, 0, 1, 107, 0}}

// Sampa interface defines the public functions
type Sampa interface {
	Calculate() error

	SetSpaData(spa.Spa)
	GetSpaData() spa.Spa

	SetBirdData(bird.Bird)
	GetBirdData() bird.Bird

	GetMpaData() Mpa
	CalculateMpa() Mpa

	SetFunction(uint32)
	GetFunction() uint32

	GetEms() float64
	GetRs() float64
	GetRm() float64
	GetASul() float64
	GetASulPct() float64
	GetDni() float64
	GetDniSul() float64
	GetGhi() float64
	GetGhiSul() float64
	GetDhi() float64
	GetDhiSul() float64
}

// NewSampa creates new Sampa instance
func NewSampa(sp spa.Spa, bi bird.Bird) (Sampa, error) {

	var sa sampa
	sa.spaData = sp
	sa.birdData = bi
	sa.function = SampaAll
	return &sa, sa.Calculate()
}

type sampa struct {
	spaData spa.Spa //Enter required INPUT VALUES into SPA structure (see SPA.H)
	//spa.function will be forced to SPA_ZA, therefore slope & azm_rotation not required)

	mpaData Mpa //Moon Position Algorithm structure (defined above)

	function uint32 //Switch to choose functions for desired output (from enumeration)

	birdData bird.Bird

	//---------------------Final SAMPA OUTPUT VALUES------------------------

	ems float64 //local observed, topocentric, angular distance between sun and moon centers [degrees]
	rs  float64 //radius of sun disk [degrees]
	rm  float64 //radius of moon disk [degrees]

	aSul    float64 //area of sun's unshaded lune (SUL) during eclipse [degrees squared]
	aSulPct float64 //percent area of SUL during eclipse [percent]

	dni    float64 //estimated direct normal solar irradiance using SERI/NREL Bird Clear Sky Model [W/m^2]
	dniSul float64 //estimated direct normal solar irradiance from the sun's unshaded lune [W/m^2]

	ghi    float64 //estimated global horizontal solar irradiance using SERI/NREL Bird Clear Sky Model [W/m^2]
	ghiSul float64 //estimated global horizontal solar irradiance from the sun's unshaded lune [W/m^2]

	dhi    float64 //estimated diffuse horizontal solar irradiance using SERI/NREL Bird Clear Sky Model [W/m^2]
	dhiSul float64 //estimated diffuse horizontal solar irradiance from the sun's unshaded lune [W/m^2]
}

func (s *sampa) SetSpaData(sp spa.Spa) {
	s.spaData = sp
}

func (s *sampa) GetSpaData() spa.Spa {
	return s.spaData
}

func (s *sampa) SetBirdData(b bird.Bird) {
	s.birdData = b
}

func (s *sampa) GetBirdData() bird.Bird {
	return s.birdData
}

func (s *sampa) GetMpaData() Mpa {
	return s.mpaData
}

func (s *sampa) SetFunction(f uint32) {
	s.function = f
}

func (s *sampa) GetFunction() uint32 {
	return s.function
}

//local observed, topocentric, angular distance between sun and moon centers [degrees]
func (s *sampa) GetEms() float64 {
	return s.ems
}

//radius of sun disk [degrees]
func (s *sampa) GetRs() float64 {
	return s.rs
}

//radius of moon disk [degrees]
func (s *sampa) GetRm() float64 {
	return s.rm
}

//area of sun's unshaded lune (SUL) during eclipse [degrees squared]
func (s *sampa) GetASul() float64 {
	return s.aSul
}

//percent area of SUL during eclipse [percent]
func (s *sampa) GetASulPct() float64 {
	return s.aSulPct
}

//estimated direct normal solar irradiance using SERI/NREL Bird Clear Sky Model [W/m^2]
func (s *sampa) GetDni() float64 {
	return s.dni
}

//estimated direct normal solar irradiance from the sun's unshaded lune [W/m^2]
func (s *sampa) GetDniSul() float64 {
	return s.dniSul
}

//estimated global horizontal solar irradiance using SERI/NREL Bird Clear Sky Model [W/m^2]
func (s *sampa) GetGhi() float64 {
	return s.ghi
}

//estimated global horizontal solar irradiance from the sun's unshaded lune [W/m^2]
func (s *sampa) GetGhiSul() float64 {
	return s.ghiSul
}

//estimated diffuse horizontal solar irradiance using SERI/NREL Bird Clear Sky Model [W/m^2]
func (s *sampa) GetDhi() float64 {
	return s.dhi
}

//estimated diffuse horizontal solar irradiance from the sun's unshaded lune [W/m^2]
func (s *sampa) GetDhiSul() float64 {
	return s.dhiSul
}

// Mpa interface defines the public functions
type Mpa interface {
	Calculate(s *sampa)
	//moon mean longitude [degrees]
	GetLPrime() float64
	//moon mean elongation [degrees]
	GetD() float64
	//sun mean anomaly [degrees]
	GetM() float64
	//moon mean anomaly [degrees]
	GetMPrime() float64
	//moon argument of latitude [degrees]
	GetF() float64
	//term l
	GetL() float64
	//term r
	GetR() float64
	//term b
	GetB() float64
	//moon longitude [degrees]
	GetLamdaPrime() float64
	//moon latitude [degrees]
	GetBeta() float64
	//distance from earth to moon [kilometers]
	GetCapDelta() float64
	//moon equatorial horizontal parallax [degrees]
	GetPi() float64
	//apparent moon longitude [degrees]
	GetLamda() float64
	//geocentric moon right ascension [degrees]
	GetAlpha() float64
	//geocentric moon declination [degrees]
	GetDelta() float64
	//observer hour angle [degrees]
	GetH() float64
	//moon right ascension parallax [degrees]
	GetDelAlpha() float64
	//topocentric moon declination [degrees]
	GetDeltaPrime() float64
	//topocentric moon right ascension [degrees]
	GetAlphaPrime() float64
	//topocentric local hour angle [degrees]
	GetHPrime() float64
	//topocentric elevation angle (uncorrected) [degrees]
	GetE0() float64
	//atmospheric refraction correction [degrees]
	GetDelE() float64
	//topocentric elevation angle (corrected) [degrees]
	GetE() float64
	//---------------------Final MPA OUTPUT VALUES------------------------
	//topocentric zenith angle [degrees]
	GetZenith() float64
	//topocentric azimuth angle (westward from south) [for astronomers]
	GetAzimuthAstro() float64
	//topocentric azimuth angle (eastward from north) [for navigators and solar radiation]
	GetAzimuth() float64
}

// NewBird creates new Bird instance
func (s *sampa) CalculateMpa() Mpa {
	var m mpa
	m.Calculate(s)
	return &m
}

type mpa struct {
	//-----------------Intermediate MPA OUTPUT VALUES--------------------

	lPrime     float64 //moon mean longitude [degrees]
	d          float64 //moon mean elongation [degrees]
	m          float64 //sun mean anomaly [degrees]
	mPrime     float64 //moon mean anomaly [degrees]
	f          float64 //moon argument of latitude [degrees]
	l          float64 //term l
	r          float64 //term r
	b          float64 //term b
	lamdaPrime float64 //moon longitude [degrees]
	beta       float64 //moon latitude [degrees]
	capDelta   float64 //distance from earth to moon [kilometers]
	pi         float64 //moon equatorial horizontal parallax [degrees]
	lamda      float64 //apparent moon longitude [degrees]

	alpha float64 //geocentric moon right ascension [degrees]
	delta float64 //geocentric moon declination [degrees]

	h          float64 //observer hour angle [degrees]
	delAlpha   float64 //moon right ascension parallax [degrees]
	deltaPrime float64 //topocentric moon declination [degrees]
	alphaPrime float64 //topocentric moon right ascension [degrees]
	hPrime     float64 //topocentric local hour angle [degrees]

	e0   float64 //topocentric elevation angle (uncorrected) [degrees]
	delE float64 //atmospheric refraction correction [degrees]
	e    float64 //topocentric elevation angle (corrected) [degrees]

	//---------------------Final MPA OUTPUT VALUES------------------------

	zenith       float64 //topocentric zenith angle [degrees]
	azimuthAstro float64 //topocentric azimuth angle (westward from south) [for astronomers]
	azimuth      float64 //topocentric azimuth angle (eastward from north) [for navigators and solar radiation]
}

func (m *mpa) GetLPrime() float64 {
	return m.lPrime
}

func (m *mpa) GetD() float64 {
	return m.d
}

func (m *mpa) GetM() float64 {
	return m.m
}

func (m *mpa) GetMPrime() float64 {
	return m.mPrime
}

func (m *mpa) GetF() float64 {
	return m.f
}

func (m *mpa) GetL() float64 {
	return m.l
}

func (m *mpa) GetR() float64 {
	return m.r
}

func (m *mpa) GetB() float64 {
	return m.b
}

func (m *mpa) GetLamdaPrime() float64 {
	return m.lamdaPrime
}

func (m *mpa) GetBeta() float64 {
	return m.beta
}

func (m *mpa) GetCapDelta() float64 {
	return m.capDelta
}

func (m *mpa) GetPi() float64 {
	return m.pi
}

func (m *mpa) GetLamda() float64 {
	return m.lamda
}

func (m *mpa) GetAlpha() float64 {
	return m.alpha
}

func (m *mpa) GetDelta() float64 {
	return m.delta
}

func (m *mpa) GetH() float64 {
	return m.h
}

func (m *mpa) GetDelAlpha() float64 {
	return m.delAlpha
}

func (m *mpa) GetDeltaPrime() float64 {
	return m.deltaPrime
}

func (m *mpa) GetAlphaPrime() float64 {
	return m.alphaPrime
}

func (m *mpa) GetHPrime() float64 {
	return m.hPrime
}

func (m *mpa) GetE0() float64 {
	return m.e0
}

func (m *mpa) GetDelE() float64 {
	return m.delE
}

func (m *mpa) GetE() float64 {
	return m.e
}

func (m *mpa) GetZenith() float64 {
	return m.zenith
}

func (m *mpa) GetAzimuthAstro() float64 {
	return m.azimuthAstro
}

func (m *mpa) GetAzimuth() float64 {
	return m.azimuth
}

func (s *sampa) fourthOrderPolynomial(a float64, b float64, c float64, d float64, e float64, x float64) float64 {
	return (((a*x+b)*x+c)*x+d)*x + e
}
func (s *sampa) deg2rad(degrees float64) float64 {
	return (math.Pi / 180.0) * degrees
}
func (s *sampa) rad2deg(radians float64) float64 {
	return (180.0 / math.Pi) * radians
}
func (s *sampa) moonMeanLongitude(jce float64) float64 {
	return s.limitDegrees(s.fourthOrderPolynomial(
		-1.0/65194000, 1.0/538841, -0.0015786, 481267.88123421, 218.3164477, jce))
}

func (s *sampa) moonMeanElongation(jce float64) float64 {
	return s.limitDegrees(s.fourthOrderPolynomial(
		-1.0/113065000, 1.0/545868, -0.0018819, 445267.1114034, 297.8501921, jce))
}
func (s *sampa) limitDegrees(degrees float64) float64 {
	var limited float64
	degrees /= 360.0
	limited = 360.0 * (degrees - math.Floor(degrees))
	if limited < 0 {
		limited += 360.0
	}

	return limited
}

func (s *sampa) sunMeanAnomaly(jce float64) float64 {
	return s.limitDegrees(s.thirdOrderPolynomial(
		1.0/24490000, -0.0001536, 35999.0502909, 357.5291092, jce))
}
func (s *sampa) thirdOrderPolynomial(a float64, b float64, c float64, d float64, x float64) float64 {
	return ((a*x+b)*x+c)*x + d
}

func (s *sampa) moonMeanAnomaly(jce float64) float64 {
	return s.limitDegrees(s.fourthOrderPolynomial(
		-1.0/14712000, 1.0/69699, 0.0087414, 477198.8675055, 134.9633964, jce))
}

func (s *sampa) moonLatitudeArgument(jce float64) float64 {
	return s.limitDegrees(s.fourthOrderPolynomial(
		1.0/863310000, -1.0/3526000, -0.0036539, 483202.0175233, 93.2720950, jce))
}

func (s *sampa) moonPeriodicTermSummation(d float64, m float64, m_prime float64, f float64, jce float64, terms [][]float64, sin_sum *float64, cos_sum *float64) {
	var eMult, trigArg float64
	e := 1.0 - jce*(0.002516+jce*0.0000074)
	*sin_sum = 0
	*cos_sum = 0

	for i := 0; i < int(COUNT); i++ {

		eMult = math.Pow(e, math.Abs(terms[i][TermM]))
		trigArg = s.deg2rad(terms[i][TermD]*d + terms[i][TermM]*m + terms[i][TermF]*f + terms[i][TermMpr]*m_prime)

		*sin_sum += eMult * terms[i][TermLb] * math.Sin(trigArg)
		*cos_sum += eMult * terms[i][TermR] * math.Cos(trigArg)

	}
}

func (s *sampa) moonLongitudeAndLatitude(jce float64, lPrime float64, f float64, mPrime float64, l float64, b float64, lamdaPrime *float64, beta *float64) {
	a1 := 119.75 + 131.849*jce
	a2 := 53.09 + 479264.290*jce
	a3 := 313.45 + 481266.484*jce
	deltaL := 3958*math.Sin(s.deg2rad(a1)) + 318*math.Sin(s.deg2rad(a2)) + 1962*math.Sin(s.deg2rad(lPrime-f))
	deltaB := -2235*math.Sin(s.deg2rad(lPrime)) + 175*math.Sin(s.deg2rad(a1-f)) + 127*math.Sin(s.deg2rad(lPrime-mPrime)) + 382*math.Sin(s.deg2rad(a3)) + 175*math.Sin(s.deg2rad(a1+f)) - 115*math.Sin(s.deg2rad(lPrime+mPrime))

	*lamdaPrime = s.limitDegrees(lPrime + (l+deltaL)/1000000)
	*beta = s.limitDegrees((b + deltaB) / 1000000)
}

func (s *sampa) moonEarthDistance(r float64) float64 {
	return 385000.56 + r/1000
}

func (s *sampa) moonEquatorialHorizParallax(delta float64) float64 {
	return s.rad2deg(math.Asin(6378.14 / delta))
}

func (s *sampa) apparentMoonLongitude(lamda_prime float64, del_psi float64) float64 {
	return lamda_prime + del_psi
}

func (s *sampa) angularDistanceSunMoon(zen_sun float64, azm_sun float64, zen_moon float64, azm_moon float64) float64 {
	zs := s.deg2rad(zen_sun)
	zm := s.deg2rad(zen_moon)
	return s.rad2deg(math.Acos(math.Cos(zs)*math.Cos(zm) + math.Sin(zs)*math.Sin(zm)*math.Cos(s.deg2rad(azm_sun-azm_moon))))
}

func (s *sampa) sunDiskRadius(r float64) float64 {
	return 959.63 / (3600.0 * r)
}

func (s *sampa) moonDiskRadius(e float64, pi float64, cap_delta float64) float64 {
	return 358473400 * (1 + math.Sin(s.deg2rad(e))*math.Sin(s.deg2rad(pi))) / (3600.0 * cap_delta)
}

func (s *sampa) sulArea(ems float64, rs float64, rm float64, a_sul *float64, a_sul_pct *float64) {

	ems2 := ems * ems
	rs2 := rs * rs
	rm2 := rm * rm
	var snum, ai, m, sa, h float64

	if ems < (rs + rm) {
		if ems <= math.Abs(rs-rm) {
			ai = math.Pi * rm2
		} else {
			snum = ems2 + rs2 - rm2
			m = (ems2 - rs2 + rm2) / (2 * ems)
			sa = snum / (2 * ems)
			h = math.Sqrt(4*ems2*rs2-snum*snum) / (2 * ems)
			ai = rs2*math.Acos(sa/rs) - h*sa + rm2*math.Acos(m/rm) - h*m
		}
	} else {
		ai = 0
	}

	*a_sul = math.Pi*rs2 - ai

	if *a_sul < 0 {
		*a_sul = 0
	}
	*a_sul_pct = *a_sul * 100.0 / (math.Pi * rs2)

}
func (s *sampa) geocentricRightAscension(lamda float64, epsilon float64, beta float64) float64 {
	lamdaRad := s.deg2rad(lamda)
	epsilonRad := s.deg2rad(epsilon)

	return s.limitDegrees(s.rad2deg(math.Atan2(math.Sin(lamdaRad)*math.Cos(epsilonRad)-
		math.Tan(s.deg2rad(beta))*math.Sin(epsilonRad), math.Cos(lamdaRad))))
}
func (s *sampa) geocentricDeclination(beta float64, epsilon float64, lambda float64) float64 {
	betaRad := s.deg2rad(beta)
	epsilonRad := s.deg2rad(epsilon)

	return s.rad2deg(math.Asin(math.Sin(betaRad)*math.Cos(epsilonRad) +
		math.Cos(betaRad)*math.Sin(epsilonRad)*math.Sin(s.deg2rad(lambda))))
}
func (s *sampa) observerHourAngle(nu float64, longitude float64, alphaDeg float64) float64 {
	return s.limitDegrees(nu + longitude - alphaDeg)
}
func (s *sampa) rightAscensionParallaxAndTopocentricDec(latitude float64, elevation float64, xi float64, h float64, delta float64, delAlpha *float64, deltaPrime *float64) {
	var deltaAlphaRad float64
	latRad := s.deg2rad(latitude)
	xiRad := s.deg2rad(xi)
	hRad := s.deg2rad(h)
	deltaRad := s.deg2rad(delta)
	u := math.Atan(0.99664719 * math.Tan(latRad))
	y := 0.99664719*math.Sin(u) + elevation*math.Sin(latRad)/6378140.0
	x := math.Cos(u) + elevation*math.Cos(latRad)/6378140.0

	deltaAlphaRad = math.Atan2(-x*math.Sin(xiRad)*math.Sin(hRad), math.Cos(deltaRad)-x*math.Sin(xiRad)*math.Cos(hRad))

	*deltaPrime = s.rad2deg(math.Atan2((math.Sin(deltaRad)-y*math.Sin(xiRad))*math.Cos(deltaAlphaRad),
		math.Cos(deltaRad)-x*math.Sin(xiRad)*math.Cos(hRad)))

	*delAlpha = s.rad2deg(deltaAlphaRad)

}
func (s *sampa) topocentricRightAscension(alphaDeg float64, deltaAlpha float64) float64 {
	return alphaDeg + deltaAlpha
}
func (s *sampa) topocentricLocalHourAngle(h float64, deltaAlpha float64) float64 {
	return h - deltaAlpha
}
func (s *sampa) topocentricElevationAngle(latitude float64, deltaPrime float64, hPrime float64) float64 {
	latRad := s.deg2rad(latitude)
	deltaPrimeRad := s.deg2rad(deltaPrime)

	return s.rad2deg(math.Asin(math.Sin(latRad)*math.Sin(deltaPrimeRad) +
		math.Cos(latRad)*math.Cos(deltaPrimeRad)*math.Cos(s.deg2rad(hPrime))))
}
func (s *sampa) atmosphericRefractionCorrection(pressure float64, temperature float64, atmosRefract float64, e0 float64) float64 {
	delE := 0.

	if e0 >= -1*(spa.SunRadius+atmosRefract) {
		delE = (pressure / 1010.0) * (283.0 / (273.0 + temperature)) * 1.02 / (60.0 * math.Tan(s.deg2rad(e0+10.3/(e0+5.11))))
	}
	return delE
}
func (s *sampa) topocentricElevationAngleCorrected(e0 float64, deltaE float64) float64 {
	return e0 + deltaE
}

func (s *sampa) topocentricZenithAngle(e float64) float64 {
	return 90.0 - e
}

func (s *sampa) topocentricAzimuthAngleAstro(hPrime float64, latitude float64, deltaPrime float64) float64 {
	hPrimeRad := s.deg2rad(hPrime)
	latRad := s.deg2rad(latitude)

	return s.limitDegrees(s.rad2deg(math.Atan2(math.Sin(hPrimeRad),
		math.Cos(hPrimeRad)*math.Sin(latRad)-math.Tan(s.deg2rad(deltaPrime))*math.Cos(latRad))))
}

func (s *sampa) topocentricAzimuthAngle(azimuthAstro float64) float64 {
	return s.limitDegrees(azimuthAstro + 180.0)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Calculate all MPA parameters and put into structure
// Note: All inputs values (listed in SPA header file) must already be in structure
///////////////////////////////////////////////////////////////////////////////////////////
func (m *mpa) Calculate(s *sampa) {
	m.lPrime = s.moonMeanLongitude(s.spaData.GetJce())
	m.d = s.moonMeanElongation(s.spaData.GetJce())
	m.m = s.sunMeanAnomaly(s.spaData.GetJce())
	m.mPrime = s.moonMeanAnomaly(s.spaData.GetJce())
	m.f = s.moonLatitudeArgument(s.spaData.GetJce())

	s.moonPeriodicTermSummation(m.d, m.m, m.mPrime, m.f, s.spaData.GetJce(), MlTerms, &m.l, &m.r)
	var tmpCos float64
	tmpCos = 0
	s.moonPeriodicTermSummation(m.d, m.m, m.mPrime, m.f, s.spaData.GetJce(), MbTerms, &m.b, &tmpCos)

	s.moonLongitudeAndLatitude(s.spaData.GetJce(), m.lPrime, m.f, m.mPrime, m.l, m.b, &m.lamdaPrime, &m.beta)

	m.capDelta = s.moonEarthDistance(m.r)
	m.pi = s.moonEquatorialHorizParallax(m.capDelta)

	m.lamda = s.apparentMoonLongitude(m.lamdaPrime, s.spaData.GetDelPsi())

	m.alpha = s.geocentricRightAscension(m.lamda, s.spaData.GetEpsilon(), m.beta)
	m.delta = s.geocentricDeclination(m.beta, s.spaData.GetEpsilon(), m.lamda)

	m.h = s.observerHourAngle(s.spaData.GetNu(), s.spaData.GetLongitude(), m.alpha)

	s.rightAscensionParallaxAndTopocentricDec(s.spaData.GetLatitude(), s.spaData.GetElevation(), m.pi, m.h, m.delta, &m.delAlpha, &m.deltaPrime)

	m.alphaPrime = s.topocentricRightAscension(m.alpha, m.delAlpha)
	m.hPrime = s.topocentricLocalHourAngle(m.h, m.delAlpha)

	m.e0 = s.topocentricElevationAngle(s.spaData.GetLatitude(), m.deltaPrime, m.hPrime)
	m.delE = s.atmosphericRefractionCorrection(s.spaData.GetPressure(), s.spaData.GetTemperature(),
		s.spaData.GetAtmosRefract(), m.e0)
	m.e = s.topocentricElevationAngleCorrected(m.e0, m.delE)

	m.zenith = s.topocentricZenithAngle(m.e)
	m.azimuthAstro = s.topocentricAzimuthAngleAstro(m.hPrime, s.spaData.GetLatitude(), m.deltaPrime)
	m.azimuth = s.topocentricAzimuthAngle(m.azimuthAstro)

}

///////////////////////////////////////////////////////////////////////////////////////////
// Estimate solar irradiances using the SERI/NREL's Bird Clear Sky Model
///////////////////////////////////////////////////////////////////////////////////////////
func (s *sampa) estimateIrr() error {

	zenith := s.spaData.GetZenith()
	r := s.spaData.GetR()

	pressure := s.spaData.GetPressure()
	ozone := s.birdData.GetOzone()
	water := s.birdData.GetWater()
	taua := s.birdData.GetTaua()
	ba := s.birdData.GetBa()
	albedo := s.birdData.GetAlbedo()
	dni_mod := s.aSulPct / 100.0

	b, err := bird.NewBird(zenith, r, pressure, ozone, water, taua, ba, albedo, dni_mod)
	if err != nil {
		return err
	}
	s.birdData = b
	s.dni = b.GetDirectNormal()
	s.dniSul = b.GetDirectNormalMod()
	s.ghi = b.GetGlobalHoriz()
	s.ghiSul = b.GetGlobalHorizMod()
	s.dhi = b.GetDiffuseHoriz()
	s.dhiSul = b.GetDiffuseHorizMod()
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////
// Calculate all SAMPA parameters and put into structure
// Note: All inputs values (listed in SPA header file) must already be in structure
///////////////////////////////////////////////////////////////////////////////////////////
func (s *sampa) Calculate() error {
	s.spaData.SetSPAFunction(0)
	s.spaData.SetSPAFunction(spa.SpaZa)
	err := s.spaData.Calculate()
	if err != nil {
		return err
	}
	s.mpaData = s.CalculateMpa()

	s.ems = s.angularDistanceSunMoon(s.spaData.GetZenith(), s.spaData.GetAzimuth(), s.mpaData.GetZenith(), s.mpaData.GetAzimuth())
	s.rs = s.sunDiskRadius(s.spaData.GetR())
	s.rm = s.moonDiskRadius(s.mpaData.GetE(), s.mpaData.GetPi(), s.mpaData.GetCapDelta())

	s.sulArea(s.ems, s.rs, s.rm, &s.aSul, &s.aSulPct)

	if s.function == SampaAll {
		err = s.estimateIrr()
		if err != nil {
			return err
		}
	}
	return nil

}
