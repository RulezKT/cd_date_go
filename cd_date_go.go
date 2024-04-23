package cd_date_go

import (
	"fmt"
	"math"

	. "github.com/RulezKT/cd_consts_go"
)

// максимальная и минимальные даты в файле, измеряется в секундах от J2000
// J2000 - 2000 Jan 1.5 (12h on January 1) or JD 2451545.0 - TT time or
// January 1, 2000, 11:58:55.816 UTC (Coordinated Universal Time).
// = 0 in seconds

//https://en.wikipedia.org/wiki/Julian_day
//https://www.aavso.org/jd-calculator
//http://www.onlineconversion.com/julian_date.htm

func GregToSecFromJD2000(gregDate GregDate) int64 {

	//JDN = (1461 × (Y + 4800 + (M − 14)/12))/4 +(367 × (M − 2 − 12 × ((M − 14)/12)))/12 − (3 × ((Y + 4900 + (M - 14)/12)/100))/4 + D − 32075

	//One must compute first the number of years (y) and months (m) since March 1 −4800 (March 1, 4801 BC)
	//integer, fraction := math.Modf(day)

	a := int((14 - gregDate.Month) / 12)
	y := int(gregDate.Year + 4800 - a)
	m := int(gregDate.Month + 12*a - 3)

	//All years in the BC era must be converted to astronomical years,
	//so that 1 BC is year 0, 2 BC is year −1, etc. Convert to a negative number, then increment toward zero.
	//JDN — это номер юлианского дня (англ. Julian Day Number),
	//который начинается в полдень числа, для которого производятся вычисления.
	//Then, if starting from a Gregorian calendar date compute:
	jdn := float64(gregDate.Day + (int)((153*m+2)/5) + 365*y + (int)(y/4) - (int)(y/100) + (int)(y/400) - 32045)

	// переводим в секунды от JD200 часы минуты и секунды пока не учтены
	date_in_sec := int64(JedToSecFromJd2000(jdn))

	// теперь отталкиваясь от полдня корректируемся на часы минуты и секунды
	// если время после 12, вычитаем половину суток из JDN и из времени в секундах
	// так как время нового дня переданное функции отталкивается от 00:00
	// а Юлианские дни от 12:00
	if gregDate.Hour < 12 {
		jdn -= 0.5
		date_in_sec -= 43200
	} else {
		gregDate.Hour -= 12
	}

	jdn += float64(gregDate.Hour*60*60+gregDate.Minutes*60+gregDate.Seconds) / 86400.0
	//System.out.println(" jdn = " +  jdn);

	date_in_sec += int64(gregDate.Hour*60*60 + gregDate.Minutes*60 + gregDate.Seconds)

	return date_in_sec
}

// Принимает JED и перводит его в sec_from_jd2000
func JedToSecFromJd2000(jedDate float64) int64 {
	return int64(float64(jedDate-JD2000) * float64(SEC_IN_1_DAY))
}

func SecFromJd2000ToGreg(date_in_seconds int64) GregDate {
	//даёт ошибки
	//например WorkWithTime.sec_from_jd2000_to_gregdate((long)20514081599.00);
	//DD.MM.YYYY: '24.1.2650
	//'HH.MM.Sec: '23.59.10
	// а должно быть  2650.1.25, 00:00git status

	//получаем Julian Day
	jdn := float64(JD2000 + float64(date_in_seconds)/float64(SEC_IN_1_DAY))

	//calculating necessary coeffs
	a := float64(jdn + 32044)
	b := int((4*a + 3) / 146097)
	c := a - float64(int(146097*b/4))
	d := int((4*c + 3) / 1461)
	e := float64(c - float64(int(1461*d/4)))
	m := (int)((5*e + 2) / 153)

	month := int(m + 3 - 12*((int)(m/10)))
	year := int(100*b + d - 4800 + (int)(m/10))

	//Дальше коэффициенты подбирал вручную, так как алгоритма дальше дня не нашел
	//в конце добавляем 0.5 так начало JD считается на 12:00 а нам надо на 00:00
	day := int(e - float64(int((153*m+2)/5)) + 1 + 0.5)
	hour := int(24 * ((e - float64(int((153*m+2)/5)) + 1 + 0.5) - float64(day)))
	//в конце вычитаем  0.808, нашел подгоном
	minute := int(60*((24*((e-float64(int((153*m+2)/5))+1+0.5)-float64(day)))-float64(hour)) - 0.808)
	if minute <= 0 {
		minute = 0
	}
	//в конце вычитаем 48.5 это 0.808 переведенный из десятых долей в секунды
	second := (int)(60*(60*((24*((e-float64(int((153*m+2)/5))+1+0.5)-float64(day)))-float64(hour))-float64(minute)) - 48.5)
	if second <= 0 {
		second = 0
	}

	return GregDate{Year: year, Month: month, Day: day, Hour: hour, Minutes: minute, Seconds: second}
}

// отличная версия калькуляции, получаем Грег. дату в ET
func SecJd2000ToGregMeeus(date_in_seconds int64) GregDate {

	//григорианский календарь введён 4 октября 1582 года
	// следующим днём после четверга 4 октября стала пятница 15 октября.
	//программа не учитывает эти 10 дней в датах до 15 октября 1582 года

	//получаем Julian Day
	jdn := float64(JD2000 + float64(date_in_seconds)/float64(SEC_IN_1_DAY))

	jdn += 0.5

	var A int

	Z := (int)(jdn)

	F := float64(jdn) - float64(Z)

	if Z < 2299161 {

		A = Z

	} else {

		alpha := (int)((float64(Z) - float64(1867216.25)) / 36524.25)
		A = Z + 1 + alpha - (int)(alpha/4)

	}

	B := int(A + 1524)

	C := (int)((float64(B) - float64(122.1)) / 365.25)

	D := (int)(float64(365.25) * float64(C))

	E := (int)(float64(B-D) / float64(30.6001))

	day := float64(B) - float64(D) - float64(int(float64(30.6001)*float64(E))) + F

	var month int

	if (E < 0) || (E > 15) {
		fmt.Println("sec_jd2000_to_greg_meeus, unacceptable value of E")
	}

	if E < 14 {
		month = E - 1
	} else {
		month = E - 13
	}

	var year int
	if month > 2 {
		year = C - 4716
	} else {
		year = C - 4715
	}

	//get fraction part
	_, fraction := math.Modf(day)
	hour := fraction * 24

	_, fraction = math.Modf(hour)
	minute := fraction * 60

	_, fraction = math.Modf(minute)
	second := fraction * 60

	//Math.ceil подходит лучше всего
	second = math.Ceil(second)
	//second =  Math.round(second);
	//second =  (int)(second);

	//если секунд больше 59, увеличиваем время
	if second > 59 {

		second -= 60
		minute++

		if minute > 59 {
			minute -= 60
			hour++
		}

		if hour > 23 {
			hour -= 24
			day++
		}

		if day > 28 {

			if month == 2 {
				if !isLeapYear(year) {
					day -= 28
					month++
				} else {
					if day > 29 {
						day -= 29
						month++
					}
				}

			} else if day > 30 {

				switch month {
				case 4:
				case 6:
				case 9:
				case 11:
					day -= 30
					month++

				case 1:
				case 3:
				case 5:
				case 7:
				case 8:
				case 10:
				case 12:
					if day > 31 {
						day -= 31
						month++
					}

				default:
					break
				}
			}

		}

		if month > 12 {

			month -= 12
			year++
		}
	}

	return GregDate{Year: year, Month: month, Day: int(day), Hour: int(hour), Minutes: int(minute), Seconds: int(second)}
}

func isLeapYear(year int) bool {
	return ((year%4 == 0) && (year%100 != 0)) || (year%400 == 0)
}

// delta_t observations started at 1620 and now it is 2017
// so we have 398 records for now
// each record is the year and a number of seconds
// [0] record corresponds to year 1620
// https://www.staff.science.uu.nl/~gent0113/deltat/deltat.htm
// ftp://maia.usno.navy.mil/ser7/deltat.data

//https://eclipse.gsfc.nasa.gov/SEhelp/deltaT.html
//This parameter is known as delta-T or ΔT (ΔT = TDT - UT).
// for delta_t calculations we use
// https://eclipse.gsfc.nasa.gov/SEcat5/deltatpoly.html
// algorithms

func calculate_delta_t(year int) float64 {

	var delta_t_sec float64
	// before year 1620 (observations started from 1620, before were only estimations)
	if year < 1620 {
		if year < -500 {
			delta_t_sec = -20 + 32*math.Pow(float64(year-1820)/100, 2)
			return delta_t_sec
		} else if year >= -500 && year <= 500 {
			delta_t_sec = 10583.6 - float64(1014.41)*float64(year)/100 + 33.78311*math.Pow(float64(year)/100, 2) - 5.952053*math.Pow(float64(year)/100, 3) - 0.1798452*math.Pow(float64(year)/100, 4) + 0.022174192*math.Pow(float64(year)/100, 5) + 0.0090316521*math.Pow(float64(year)/100, 6)
			return delta_t_sec
		} else if year > 500 && year <= 1600 {
			delta_t_sec = 1574.2 - 556.01*(float64(year)-1000)/100 + 71.23472*math.Pow((float64(year)-1000)/100, 2) + 0.319781*math.Pow((float64(year)-1000)/100, 3) - 0.8503463*math.Pow((float64(year)-1000)/100, 4) - 0.005050998*math.Pow((float64(year)-1000)/100, 5) + 0.0083572073*math.Pow((float64(year)-1000)/100, 6)
			return delta_t_sec
		} else { // from 1600 to 1620
			delta_t_sec = 120 - 0.9808*(float64(year)-1600) - 0.01532*math.Pow((float64(year)-1600), 2) + math.Pow((float64(year)-1600), 3)/7129
			return delta_t_sec
		}
	}

	if year >= 1620 && year <= 1700 {
		delta_t_sec = 120 - 0.9808*(float64(year)-1600) - 0.01532*math.Pow((float64(year)-1600), 2) + math.Pow((float64(year)-1600), 3)/7129
		return delta_t_sec
	}

	if year > 1700 && year <= 1800 {
		delta_t_sec = 8.83 + 0.1603*(float64(year)-1700) - 0.0059285*math.Pow((float64(year)-1700), 2) + 0.00013336*math.Pow((float64(year)-1700), 3) - math.Pow((float64(year)-1700), 4)/1174000
		return delta_t_sec
	}

	if year > 1800 && year <= 1860 {
		delta_t_sec = 13.72 - 0.332447*(float64(year)-1800) + 0.0068612*math.Pow((float64(year)-1800), 2) + 0.0041116*math.Pow((float64(year)-1800), 3) - 0.00037436*math.Pow((float64(year)-1800), 4) + 0.0000121272*math.Pow((float64(year)-1800), 5) - 0.0000001699*math.Pow((float64(year)-1800), 6) + 0.000000000875*math.Pow((float64(year)-1800), 7)
		return delta_t_sec
	}

	if year > 1860 && year <= 1900 {
		delta_t_sec = 7.62 + 0.5737*(float64(year)-1860) - 0.251754*math.Pow((float64(year)-1860), 2) + 0.01680668*math.Pow((float64(year)-1860), 3) - 0.0004473624*math.Pow((float64(year)-1860), 4) + math.Pow((float64(year)-1860), 5)/233174
		return delta_t_sec
	}

	if year > 1900 && year <= 1920 {
		delta_t_sec = -2.79 + 1.494119*(float64(year)-1900) - 0.0598939*math.Pow((float64(year)-1900), 2) + 0.0061966*math.Pow((float64(year)-1900), 3) - 0.000197*math.Pow((float64(year)-1900), 4)
		return delta_t_sec
	}

	if year > 1920 && year <= 1941 {
		delta_t_sec = 21.20 + 0.84493*(float64(year)-1920) - 0.076100*math.Pow((float64(year)-1920), 2) + 0.0020936*math.Pow((float64(year)-1920), 3)
		return delta_t_sec
	}

	if year > 1941 && year <= 1961 {
		delta_t_sec = 29.07 + 0.407*(float64(year)-1950) - math.Pow((float64(year)-1950), 2)/233.0 + math.Pow((float64(year)-1950), 3)/2547.0
		return delta_t_sec
	}

	if year > 1961 && year <= 1986 {
		delta_t_sec = 45.45 + 1.067*(float64(year)-1975) - math.Pow((float64(year)-1975), 2)/260.0 - math.Pow((float64(year)-1975), 3)/718.0
		return delta_t_sec
	}

	if year > 1986 && year <= 2005 {
		delta_t_sec = 63.86 + 0.3345*(float64(year)-2000) - 0.060374*math.Pow((float64(year)-2000), 2) + 0.0017275*math.Pow((float64(year)-2000), 3) + 0.000651814*math.Pow((float64(year)-2000), 4) + 0.00002373599*math.Pow((float64(year)-2000), 5)
		return delta_t_sec
	}

	if year > 2005 && year <= 2050 {
		delta_t_sec = 62.92 + 0.32217*(float64(year)-2000) + 0.005589*math.Pow((float64(year)-2000), 2)
		return delta_t_sec
	}

	if year > 2050 && year <= 2150 {
		delta_t_sec = -20 + 32*math.Pow(((float64(year)-1820)/100.0), 2) - 0.5628*(2150-float64(year))
		return delta_t_sec
	}

	if year > 2150 {
		delta_t_sec = -20 + 32*math.Pow(((float64(year)-1820)/100.0), 2)
		return delta_t_sec
	}

	return -1

}

func DeltaT(year int) float64 {

	first_year := DeltaTable[0].Year
	last_year := DeltaTable[len(DeltaTable)-1].Year
	// fmt.Println("first_year", first_year)
	// fmt.Println("last_year", last_year)

	if year < first_year || year > last_year {
		return calculate_delta_t(year)
	}

	// if year < bsp.DeltaTTable.FirstYear || year > bsp.DeltaTTable.LastYear {
	// 	return calculate_delta_t(year)
	// }

	return DeltaTable[year-first_year].Seconds
}

var DeltaTable = []DeltaTTableStructure{

	{Year: 1950, Seconds: 29.15}, {Year: 1951, Seconds: 29.57}, {Year: 1952, Seconds: 29.97}, {Year: 1953, Seconds: 30.36}, {Year: 1954, Seconds: 30.72},
	{Year: 1955, Seconds: 31.07}, {Year: 1956, Seconds: 31.35}, {Year: 1957, Seconds: 31.68}, {Year: 1958, Seconds: 32.18}, {Year: 1959, Seconds: 32.68},

	{Year: 1960, Seconds: 33.15}, {Year: 1961, Seconds: 33.59}, {Year: 1962, Seconds: 34.00}, {Year: 1963, Seconds: 34.47}, {Year: 1964, Seconds: 35.03},
	{Year: 1965, Seconds: 35.73}, {Year: 1966, Seconds: 36.54}, {Year: 1967, Seconds: 37.43}, {Year: 1968, Seconds: 38.29}, {Year: 1969, Seconds: 39.20},

	{Year: 1970, Seconds: 40.18}, {Year: 1971, Seconds: 41.17}, {Year: 1972, Seconds: 42.23}, {Year: 1973, Seconds: 43.37}, {Year: 1974, Seconds: 44.49},
	{Year: 1975, Seconds: 45.48}, {Year: 1976, Seconds: 46.46}, {Year: 1977, Seconds: 47.52}, {Year: 1978, Seconds: 48.53}, {Year: 1979, Seconds: 49.59},

	{Year: 1980, Seconds: 50.54}, {Year: 1981, Seconds: 51.38}, {Year: 1982, Seconds: 52.17}, {Year: 1983, Seconds: 52.96}, {Year: 1984, Seconds: 53.79},
	{Year: 1985, Seconds: 54.34}, {Year: 1986, Seconds: 54.87}, {Year: 1987, Seconds: 55.32}, {Year: 1988, Seconds: 55.82}, {Year: 1989, Seconds: 56.30},

	{Year: 1990, Seconds: 56.86}, {Year: 1991, Seconds: 57.57}, {Year: 1992, Seconds: 58.31}, {Year: 1993, Seconds: 59.12}, {Year: 1994, Seconds: 59.99},
	{Year: 1995, Seconds: 60.78}, {Year: 1996, Seconds: 61.63}, {Year: 1997, Seconds: 62.30}, {Year: 1998, Seconds: 62.97}, {Year: 1999, Seconds: 63.47},

	{Year: 2000, Seconds: 63.83}, {Year: 2001, Seconds: 64.09}, {Year: 2002, Seconds: 64.30}, {Year: 2003, Seconds: 64.47}, {Year: 2004, Seconds: 64.57},
	{Year: 2005, Seconds: 64.69}, {Year: 2006, Seconds: 64.85}, {Year: 2007, Seconds: 65.15}, {Year: 2008, Seconds: 65.46}, {Year: 2009, Seconds: 65.78},

	{Year: 2010, Seconds: 66.07}, {Year: 2011, Seconds: 66.32}, {Year: 2012, Seconds: 66.60}, {Year: 2013, Seconds: 66.91}, {Year: 2014, Seconds: 67.28},
	{Year: 2015, Seconds: 67.18}, {Year: 2016, Seconds: 68.18}, {Year: 2017, Seconds: 69.18}, {Year: 2018, Seconds: 69.18}, {Year: 2019, Seconds: 69.18},
	{Year: 2020, Seconds: 69.18},
}

// func DeltaTPtr() *DeltaTTable {

// 	var dt DeltaTTable
// 	dt.Table = []DeltaTTableStructure{
// 		{Year: 1620, Seconds: 124},

// 		{Year: 1621, Seconds: 119}, {Year: 1622, Seconds: 115}, {Year: 1623, Seconds: 110}, {Year: 1624, Seconds: 106},
// 		{Year: 1625, Seconds: 102}, {Year: 1626, Seconds: 98}, {Year: 1627, Seconds: 95}, {Year: 1628, Seconds: 91}, {Year: 1629, Seconds: 88},

// 		{Year: 1630, Seconds: 85}, {Year: 1631, Seconds: 82}, {Year: 1632, Seconds: 79}, {Year: 1633, Seconds: 77}, {Year: 1634, Seconds: 74},
// 		{Year: 1635, Seconds: 72}, {Year: 1636, Seconds: 70}, {Year: 1637, Seconds: 67}, {Year: 1638, Seconds: 65}, {Year: 1639, Seconds: 63},

// 		{Year: 1640, Seconds: 62}, {Year: 1641, Seconds: 60}, {Year: 1642, Seconds: 58}, {Year: 1643, Seconds: 57}, {Year: 1644, Seconds: 55},
// 		{Year: 1645, Seconds: 54}, {Year: 1646, Seconds: 53}, {Year: 1647, Seconds: 51}, {Year: 1648, Seconds: 50}, {Year: 1649, Seconds: 49},

// 		{Year: 1650, Seconds: 48}, {Year: 1651, Seconds: 47}, {Year: 1652, Seconds: 46}, {Year: 1653, Seconds: 45}, {Year: 1654, Seconds: 44},
// 		{Year: 1655, Seconds: 43}, {Year: 1656, Seconds: 42}, {Year: 1657, Seconds: 41}, {Year: 1658, Seconds: 40}, {Year: 1659, Seconds: 38},

// 		{Year: 1660, Seconds: 37}, {Year: 1661, Seconds: 36}, {Year: 1662, Seconds: 35}, {Year: 1663, Seconds: 34}, {Year: 1664, Seconds: 33},
// 		{Year: 1665, Seconds: 32}, {Year: 1666, Seconds: 31}, {Year: 1667, Seconds: 30}, {Year: 1668, Seconds: 28}, {Year: 1669, Seconds: 27},

// 		{Year: 1670, Seconds: 26}, {Year: 1671, Seconds: 25}, {Year: 1672, Seconds: 24}, {Year: 1673, Seconds: 23}, {Year: 1674, Seconds: 22},
// 		{Year: 1675, Seconds: 21}, {Year: 1676, Seconds: 20}, {Year: 1677, Seconds: 19}, {Year: 1678, Seconds: 18}, {Year: 1679, Seconds: 17},

// 		{Year: 1680, Seconds: 16}, {Year: 1681, Seconds: 15}, {Year: 1682, Seconds: 14}, {Year: 1683, Seconds: 14}, {Year: 1684, Seconds: 13},
// 		{Year: 1685, Seconds: 12}, {Year: 1686, Seconds: 12}, {Year: 1687, Seconds: 11}, {Year: 1688, Seconds: 11}, {Year: 1689, Seconds: 10},

// 		{Year: 1690, Seconds: 10}, {Year: 1691, Seconds: 10}, {Year: 1692, Seconds: 9}, {Year: 1693, Seconds: 9}, {Year: 1694, Seconds: 9},
// 		{Year: 1695, Seconds: 9}, {Year: 1696, Seconds: 9}, {Year: 1697, Seconds: 9}, {Year: 1698, Seconds: 9}, {Year: 1699, Seconds: 9},

// 		{Year: 1700, Seconds: 9}, {Year: 1701, Seconds: 9}, {Year: 1702, Seconds: 9}, {Year: 1703, Seconds: 9}, {Year: 1704, Seconds: 9},
// 		{Year: 1705, Seconds: 9}, {Year: 1706, Seconds: 9}, {Year: 1707, Seconds: 9}, {Year: 1708, Seconds: 10}, {Year: 1709, Seconds: 10},

// 		{Year: 1710, Seconds: 10}, {Year: 1711, Seconds: 10}, {Year: 1712, Seconds: 10}, {Year: 1713, Seconds: 10}, {Year: 1714, Seconds: 10},
// 		{Year: 1715, Seconds: 10}, {Year: 1716, Seconds: 10}, {Year: 1717, Seconds: 11}, {Year: 1718, Seconds: 11}, {Year: 1719, Seconds: 11},

// 		{Year: 1720, Seconds: 11}, {Year: 1721, Seconds: 11}, {Year: 1722, Seconds: 11}, {Year: 1723, Seconds: 11}, {Year: 1724, Seconds: 11},
// 		{Year: 1725, Seconds: 11}, {Year: 1726, Seconds: 11}, {Year: 1727, Seconds: 11}, {Year: 1728, Seconds: 11}, {Year: 1729, Seconds: 11},

// 		{Year: 1730, Seconds: 11}, {Year: 1731, Seconds: 11}, {Year: 1732, Seconds: 11}, {Year: 1733, Seconds: 11}, {Year: 1734, Seconds: 12},
// 		{Year: 1735, Seconds: 12}, {Year: 1736, Seconds: 12}, {Year: 1737, Seconds: 12}, {Year: 1738, Seconds: 12}, {Year: 1739, Seconds: 12},

// 		{Year: 1740, Seconds: 12}, {Year: 1741, Seconds: 12}, {Year: 1742, Seconds: 12}, {Year: 1743, Seconds: 12}, {Year: 1744, Seconds: 13},
// 		{Year: 1745, Seconds: 13}, {Year: 1746, Seconds: 13}, {Year: 1747, Seconds: 13}, {Year: 1748, Seconds: 13}, {Year: 1749, Seconds: 13},

// 		{Year: 1750, Seconds: 13}, {Year: 1751, Seconds: 14}, {Year: 1752, Seconds: 14}, {Year: 1753, Seconds: 14}, {Year: 1754, Seconds: 14},
// 		{Year: 1755, Seconds: 14}, {Year: 1756, Seconds: 14}, {Year: 1757, Seconds: 14}, {Year: 1758, Seconds: 15}, {Year: 1759, Seconds: 15},

// 		{Year: 1760, Seconds: 15}, {Year: 1761, Seconds: 15}, {Year: 1762, Seconds: 15}, {Year: 1763, Seconds: 15}, {Year: 1764, Seconds: 15},
// 		{Year: 1765, Seconds: 16}, {Year: 1766, Seconds: 16}, {Year: 1767, Seconds: 16}, {Year: 1768, Seconds: 16}, {Year: 1769, Seconds: 16},

// 		{Year: 1770, Seconds: 16}, {Year: 1771, Seconds: 16}, {Year: 1772, Seconds: 16}, {Year: 1773, Seconds: 16}, {Year: 1774, Seconds: 16},
// 		{Year: 1775, Seconds: 17}, {Year: 1776, Seconds: 17}, {Year: 1777, Seconds: 17}, {Year: 1778, Seconds: 17}, {Year: 1779, Seconds: 17},

// 		{Year: 1780, Seconds: 17}, {Year: 1781, Seconds: 17}, {Year: 1782, Seconds: 17}, {Year: 1783, Seconds: 17}, {Year: 1784, Seconds: 17},
// 		{Year: 1785, Seconds: 17}, {Year: 1786, Seconds: 17}, {Year: 1787, Seconds: 17}, {Year: 1788, Seconds: 17}, {Year: 1789, Seconds: 17},

// 		{Year: 1790, Seconds: 17}, {Year: 1791, Seconds: 17}, {Year: 1792, Seconds: 16}, {Year: 1793, Seconds: 16}, {Year: 1794, Seconds: 16},
// 		{Year: 1795, Seconds: 16}, {Year: 1796, Seconds: 15}, {Year: 1797, Seconds: 15}, {Year: 1798, Seconds: 14}, {Year: 1799, Seconds: 14},

// 		{Year: 1800, Seconds: 13.7}, {Year: 1801, Seconds: 13.4}, {Year: 1802, Seconds: 13.1}, {Year: 1803, Seconds: 12.9}, {Year: 1804, Seconds: 12.7},
// 		{Year: 1805, Seconds: 12.6}, {Year: 1806, Seconds: 12.5}, {Year: 1807, Seconds: 12.5}, {Year: 1808, Seconds: 12.5}, {Year: 1809, Seconds: 12.5},

// 		{Year: 1810, Seconds: 12.5}, {Year: 1811, Seconds: 12.5}, {Year: 1812, Seconds: 12.5}, {Year: 1813, Seconds: 12.5}, {Year: 1814, Seconds: 12.5},
// 		{Year: 1815, Seconds: 12.5}, {Year: 1816, Seconds: 12.5}, {Year: 1817, Seconds: 12.4}, {Year: 1818, Seconds: 12.3}, {Year: 1819, Seconds: 12.3},

// 		{Year: 1820, Seconds: 12.0}, {Year: 1821, Seconds: 11.7}, {Year: 1822, Seconds: 11.4}, {Year: 1823, Seconds: 11.1}, {Year: 1824, Seconds: 10.6},
// 		{Year: 1825, Seconds: 10.2}, {Year: 1826, Seconds: 9.6}, {Year: 1827, Seconds: 9.1}, {Year: 1828, Seconds: 8.6}, {Year: 1829, Seconds: 8.0},

// 		{Year: 1830, Seconds: 7.5}, {Year: 1831, Seconds: 7.0}, {Year: 1832, Seconds: 6.6}, {Year: 1833, Seconds: 6.3}, {Year: 1834, Seconds: 6.0},
// 		{Year: 1835, Seconds: 5.8}, {Year: 1836, Seconds: 5.7}, {Year: 1837, Seconds: 5.6}, {Year: 1838, Seconds: 5.6}, {Year: 1839, Seconds: 5.6},

// 		{Year: 1840, Seconds: 5.7}, {Year: 1841, Seconds: 5.8}, {Year: 1842, Seconds: 5.9}, {Year: 1843, Seconds: 6.1}, {Year: 1844, Seconds: 6.2},
// 		{Year: 1845, Seconds: 6.3}, {Year: 1846, Seconds: 6.5}, {Year: 1847, Seconds: 6.6}, {Year: 1848, Seconds: 6.8}, {Year: 1849, Seconds: 6.9},

// 		{Year: 1850, Seconds: 7.1}, {Year: 1851, Seconds: 7.2}, {Year: 1852, Seconds: 7.3}, {Year: 1853, Seconds: 7.4}, {Year: 1854, Seconds: 7.5},
// 		{Year: 1855, Seconds: 7.6}, {Year: 1856, Seconds: 7.7}, {Year: 1857, Seconds: 7.7}, {Year: 1858, Seconds: 7.8}, {Year: 1859, Seconds: 7.8},

// 		{Year: 1860, Seconds: 7.88}, {Year: 1861, Seconds: 7.82}, {Year: 1862, Seconds: 7.54}, {Year: 1863, Seconds: 6.97}, {Year: 1864, Seconds: 6.40},
// 		{Year: 1865, Seconds: 6.02}, {Year: 1866, Seconds: 5.41}, {Year: 1867, Seconds: 4.10}, {Year: 1868, Seconds: 2.92}, {Year: 1869, Seconds: 1.82},

// 		{Year: 1870, Seconds: 1.61}, {Year: 1871, Seconds: 0.10}, {Year: 1872, Seconds: -1.02}, {Year: 1873, Seconds: -1.28}, {Year: 1874, Seconds: -2.69},
// 		{Year: 1875, Seconds: -3.24}, {Year: 1876, Seconds: -3.64}, {Year: 1877, Seconds: -4.54}, {Year: 1878, Seconds: -4.71}, {Year: 1879, Seconds: -5.1},

// 		{Year: 1880, Seconds: -5.40}, {Year: 1881, Seconds: -5.42}, {Year: 1882, Seconds: -5.20}, {Year: 1883, Seconds: -5.46}, {Year: 1884, Seconds: -5.46},
// 		{Year: 1885, Seconds: -5.79}, {Year: 1886, Seconds: -5.63}, {Year: 1887, Seconds: -5.64}, {Year: 1888, Seconds: -5.80}, {Year: 1889, Seconds: -5.66},

// 		{Year: 1890, Seconds: -5.87}, {Year: 1891, Seconds: -6.01}, {Year: 1892, Seconds: -6.19}, {Year: 1893, Seconds: -6.64}, {Year: 1894, Seconds: -6.44},
// 		{Year: 1895, Seconds: -6.47}, {Year: 1896, Seconds: -6.09}, {Year: 1897, Seconds: -5.76}, {Year: 1898, Seconds: -4.66}, {Year: 1899, Seconds: -3.74},

// 		{Year: 1900, Seconds: -2.72}, {Year: 1901, Seconds: -1.54}, {Year: 1902, Seconds: -0.02}, {Year: 1903, Seconds: 1.24}, {Year: 1904, Seconds: 2.64},
// 		{Year: 1905, Seconds: 3.86}, {Year: 1906, Seconds: 5.37}, {Year: 1907, Seconds: 6.14}, {Year: 1908, Seconds: 7.75}, {Year: 1909, Seconds: 9.13},

// 		{Year: 1910, Seconds: 10.46}, {Year: 1911, Seconds: 11.53}, {Year: 1912, Seconds: 13.36}, {Year: 1913, Seconds: 14.65}, {Year: 1914, Seconds: 16.01},
// 		{Year: 1915, Seconds: 17.20}, {Year: 1916, Seconds: 18.24}, {Year: 1917, Seconds: 19.06}, {Year: 1918, Seconds: 20.25}, {Year: 1919, Seconds: 20.95},

// 		{Year: 1920, Seconds: 21.16}, {Year: 1921, Seconds: 22.25}, {Year: 1922, Seconds: 22.41}, {Year: 1923, Seconds: 23.03}, {Year: 1924, Seconds: 23.49},
// 		{Year: 1925, Seconds: 23.69}, {Year: 1926, Seconds: 23.86}, {Year: 1927, Seconds: 24.49}, {Year: 1928, Seconds: 24.34}, {Year: 1929, Seconds: 24.08},

// 		{Year: 1930, Seconds: 24.02}, {Year: 1931, Seconds: 24.00}, {Year: 1932, Seconds: 23.87}, {Year: 1933, Seconds: 23.95}, {Year: 1934, Seconds: 23.86},
// 		{Year: 1935, Seconds: 23.93}, {Year: 1936, Seconds: 23.73}, {Year: 1937, Seconds: 23.92}, {Year: 1938, Seconds: 23.96}, {Year: 1939, Seconds: 24.02},

// 		{Year: 1940, Seconds: 24.33}, {Year: 1941, Seconds: 24.83}, {Year: 1942, Seconds: 25.30}, {Year: 1943, Seconds: 25.70}, {Year: 1944, Seconds: 26.24},
// 		{Year: 1945, Seconds: 26.77}, {Year: 1946, Seconds: 27.28}, {Year: 1947, Seconds: 27.78}, {Year: 1948, Seconds: 28.25}, {Year: 1949, Seconds: 28.71},

// 		{Year: 1950, Seconds: 29.15}, {Year: 1951, Seconds: 29.57}, {Year: 1952, Seconds: 29.97}, {Year: 1953, Seconds: 30.36}, {Year: 1954, Seconds: 30.72},
// 		{Year: 1955, Seconds: 31.07}, {Year: 1956, Seconds: 31.35}, {Year: 1957, Seconds: 31.68}, {Year: 1958, Seconds: 32.18}, {Year: 1959, Seconds: 32.68},

// 		{Year: 1960, Seconds: 33.15}, {Year: 1961, Seconds: 33.59}, {Year: 1962, Seconds: 34.00}, {Year: 1963, Seconds: 34.47}, {Year: 1964, Seconds: 35.03},
// 		{Year: 1965, Seconds: 35.73}, {Year: 1966, Seconds: 36.54}, {Year: 1967, Seconds: 37.43}, {Year: 1968, Seconds: 38.29}, {Year: 1969, Seconds: 39.20},

// 		{Year: 1970, Seconds: 40.18}, {Year: 1971, Seconds: 41.17}, {Year: 1972, Seconds: 42.23}, {Year: 1973, Seconds: 43.37}, {Year: 1974, Seconds: 44.49},
// 		{Year: 1975, Seconds: 45.48}, {Year: 1976, Seconds: 46.46}, {Year: 1977, Seconds: 47.52}, {Year: 1978, Seconds: 48.53}, {Year: 1979, Seconds: 49.59},

// 		{Year: 1980, Seconds: 50.54}, {Year: 1981, Seconds: 51.38}, {Year: 1982, Seconds: 52.17}, {Year: 1983, Seconds: 52.96}, {Year: 1984, Seconds: 53.79},
// 		{Year: 1985, Seconds: 54.34}, {Year: 1986, Seconds: 54.87}, {Year: 1987, Seconds: 55.32}, {Year: 1988, Seconds: 55.82}, {Year: 1989, Seconds: 56.30},

// 		{Year: 1990, Seconds: 56.86}, {Year: 1991, Seconds: 57.57}, {Year: 1992, Seconds: 58.31}, {Year: 1993, Seconds: 59.12}, {Year: 1994, Seconds: 59.99},
// 		{Year: 1995, Seconds: 60.78}, {Year: 1996, Seconds: 61.63}, {Year: 1997, Seconds: 62.30}, {Year: 1998, Seconds: 62.97}, {Year: 1999, Seconds: 63.47},

// 		{Year: 2000, Seconds: 63.83}, {Year: 2001, Seconds: 64.09}, {Year: 2002, Seconds: 64.30}, {Year: 2003, Seconds: 64.47}, {Year: 2004, Seconds: 64.57},
// 		{Year: 2005, Seconds: 64.69}, {Year: 2006, Seconds: 64.85}, {Year: 2007, Seconds: 65.15}, {Year: 2008, Seconds: 65.46}, {Year: 2009, Seconds: 65.78},

// 		{Year: 2010, Seconds: 66.07}, {Year: 2011, Seconds: 66.32}, {Year: 2012, Seconds: 66.60}, {Year: 2013, Seconds: 66.91}, {Year: 2014, Seconds: 67.28},
// 		{Year: 2015, Seconds: 67.18}, {Year: 2016, Seconds: 68.18}, {Year: 2017, Seconds: 69.18}, {Year: 2018, Seconds: 69.18}, {Year: 2019, Seconds: 69.18},
// 		{Year: 2020, Seconds: 69.18},
// 	}

// 	dt.FirstYear = dt.Table[0].Year
// 	dt.LastYear = dt.Table[len(dt.Table)-1].Year
// 	return &dt

// }
