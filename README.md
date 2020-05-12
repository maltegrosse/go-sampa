Go-NREL's Solar and Moon Position Algorithm (SAMPA)
=======================================
[![Go Report Card](https://goreportcard.com/badge/github.com/maltegrosse/go-sampa)](https://goreportcard.com/report/github.com/maltegrosse/go-sampa)
[![GoDoc](https://godoc.org/github.com/maltegrosse/go-sampa?status.svg)](https://pkg.go.dev/github.com/maltegrosse/go-sampa)
![Go](https://github.com/maltegrosse/go-sampa/workflows/Go/badge.svg) 

NREL's Solar and Moon Position Algorithm (SAMPA) calculates  the solar and lunar zenith and azimuth angles in the period from the year -2000 to 6000, with uncertainties of +/- 0.0003 degrees for the Sun and +/- 0.003 degrees for the Moon, based on the date, time, and location on Earth. The algorithm can be used for solar eclipse monitoring and estimating the reduction in solar irradiance for many applications, such as smart grid, solar energy, etc.

(Reference: Reda, I. (2010). Solar Eclipse Monitoring for Solar Energy Applications Using the Solar and Moon Position Algorithms. 35 pp.; NREL Report No. TP-3B0-47681). 
## Installation

This packages requires Go 1.13. If you installed it and set up your GOPATH, just run:

`go get -u github.com/maltegrosse/go-sampa`

## Usage

You can find some examples in the [examples](examples) directory.

Please visit https://midcdmz.nrel.gov/sampa for additional information.

Some additional helper functions have been added to the original application logic.
## Notes


|       | NREL sampa_tester.c    | GO Sampa    |  
|---------------|-------|-------|
| Julian Day     | 2455034.564583  | 2455034.564583  | 
| L          | 299.4024  | 299.402381  | 
| B      | -0.00001308059  | -0.000013080591  | 
| R            |  1.016024   | 1.016024218757  | 
| H  |  344.999100  |   344.999099851812  | 
| Delta Psi          | 0.004441121  | 0.004441121189  | 
| Delta Epsilon         | 0.001203311  | 0.001203311382  | 
| Epsilon          | 23.439252 | 23.439252167574  | 
| Zenith     | 14.512686 | 14.512686209564 | 
| Azimuth     | 104.387917  |  104.387916743210 | 
| Angular dist      | 0.374760  | 0.374759984176  | 
| Sun Radius         | 0.262360  | 0.262359778407 | 
| Moon Radius            |  0.283341  | 0.283341456977 |
| Area unshaded            | 78.363514  | 78.363513779787 |
| DNI             |  719.099358  | 719.099358263094 |



## License
**[NREL SAMPA License](https://midcdmz.nrel.gov/sampa/#license)**

Adoption in Golang under **[MIT license](http://opensource.org/licenses/mit-license.php)** 2020 © Malte Grosse.

