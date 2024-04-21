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
	// а должно быть  2650.1.25, 00:00

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

func DeltaT(year int, bsp BspFile) float64 {

	if year < bsp.DeltaTTable.FirstYear || year > bsp.DeltaTTable.LastYear {
		return calculate_delta_t(year)
	}

	return bsp.DeltaTTable.Table[year-bsp.DeltaTTable.FirstYear].Seconds
}

func DeltaTPtr() *DeltaTTable {

	var dt DeltaTTable
	dt.Table = []DeltaTTableStructure{
		{Year: 1620, Seconds: 124},

		{1621, 119}, {1622, 115}, {1623, 110}, {1624, 106},
		{1625, 102}, {1626, 98}, {1627, 95}, {1628, 91}, {1629, 88},

		{1630, 85}, {1631, 82}, {1632, 79}, {1633, 77}, {1634, 74},
		{1635, 72}, {1636, 70}, {1637, 67}, {1638, 65}, {1639, 63},

		{1640, 62}, {1641, 60}, {1642, 58}, {1643, 57}, {1644, 55},
		{1645, 54}, {1646, 53}, {1647, 51}, {1648, 50}, {1649, 49},

		{1650, 48}, {1651, 47}, {1652, 46}, {1653, 45}, {1654, 44},
		{1655, 43}, {1656, 42}, {1657, 41}, {1658, 40}, {1659, 38},

		{1660, 37}, {1661, 36}, {1662, 35}, {1663, 34}, {1664, 33},
		{1665, 32}, {1666, 31}, {1667, 30}, {1668, 28}, {1669, 27},

		{1670, 26}, {1671, 25}, {1672, 24}, {1673, 23}, {1674, 22},
		{1675, 21}, {1676, 20}, {1677, 19}, {1678, 18}, {1679, 17},

		{1680, 16}, {1681, 15}, {1682, 14}, {1683, 14}, {1684, 13},
		{1685, 12}, {1686, 12}, {1687, 11}, {1688, 11}, {1689, 10},

		{1690, 10}, {1691, 10}, {1692, 9}, {1693, 9}, {1694, 9},
		{1695, 9}, {1696, 9}, {1697, 9}, {1698, 9}, {1699, 9},

		{1700, 9}, {1701, 9}, {1702, 9}, {1703, 9}, {1704, 9},
		{1705, 9}, {1706, 9}, {1707, 9}, {1708, 10}, {1709, 10},

		{1710, 10}, {1711, 10}, {1712, 10}, {1713, 10}, {1714, 10},
		{1715, 10}, {1716, 10}, {1717, 11}, {1718, 11}, {1719, 11},

		{1720, 11}, {1721, 11}, {1722, 11}, {1723, 11}, {1724, 11},
		{1725, 11}, {1726, 11}, {1727, 11}, {1728, 11}, {1729, 11},

		{1730, 11}, {1731, 11}, {1732, 11}, {1733, 11}, {1734, 12},
		{1735, 12}, {1736, 12}, {1737, 12}, {1738, 12}, {1739, 12},

		{1740, 12}, {1741, 12}, {1742, 12}, {1743, 12}, {1744, 13},
		{1745, 13}, {1746, 13}, {1747, 13}, {1748, 13}, {1749, 13},

		{1750, 13}, {1751, 14}, {1752, 14}, {1753, 14}, {1754, 14},
		{1755, 14}, {1756, 14}, {1757, 14}, {1758, 15}, {1759, 15},

		{1760, 15}, {1761, 15}, {1762, 15}, {1763, 15}, {1764, 15},
		{1765, 16}, {1766, 16}, {1767, 16}, {1768, 16}, {1769, 16},

		{1770, 16}, {1771, 16}, {1772, 16}, {1773, 16}, {1774, 16},
		{1775, 17}, {1776, 17}, {1777, 17}, {1778, 17}, {1779, 17},

		{1780, 17}, {1781, 17}, {1782, 17}, {1783, 17}, {1784, 17},
		{1785, 17}, {1786, 17}, {1787, 17}, {1788, 17}, {1789, 17},

		{1790, 17}, {1791, 17}, {1792, 16}, {1793, 16}, {1794, 16},
		{1795, 16}, {1796, 15}, {1797, 15}, {1798, 14}, {1799, 14},

		{1800, 13.7}, {1801, 13.4}, {1802, 13.1}, {1803, 12.9}, {1804, 12.7},
		{1805, 12.6}, {1806, 12.5}, {1807, 12.5}, {1808, 12.5}, {1809, 12.5},

		{1810, 12.5}, {1811, 12.5}, {1812, 12.5}, {1813, 12.5}, {1814, 12.5},
		{1815, 12.5}, {1816, 12.5}, {1817, 12.4}, {1818, 12.3}, {1819, 12.3},

		{1820, 12.0}, {1821, 11.7}, {1822, 11.4}, {1823, 11.1}, {1824, 10.6},
		{1825, 10.2}, {1826, 9.6}, {1827, 9.1}, {1828, 8.6}, {1829, 8.0},

		{1830, 7.5}, {1831, 7.0}, {1832, 6.6}, {1833, 6.3}, {1834, 6.0},
		{1835, 5.8}, {1836, 5.7}, {1837, 5.6}, {1838, 5.6}, {1839, 5.6},

		{1840, 5.7}, {1841, 5.8}, {1842, 5.9}, {1843, 6.1}, {1844, 6.2},
		{1845, 6.3}, {1846, 6.5}, {1847, 6.6}, {1848, 6.8}, {1849, 6.9},

		{1850, 7.1}, {1851, 7.2}, {1852, 7.3}, {1853, 7.4}, {1854, 7.5},
		{1855, 7.6}, {1856, 7.7}, {1857, 7.7}, {1858, 7.8}, {1859, 7.8},

		{1860, 7.88}, {1861, 7.82}, {1862, 7.54}, {1863, 6.97}, {1864, 6.40},
		{1865, 6.02}, {1866, 5.41}, {1867, 4.10}, {1868, 2.92}, {1869, 1.82},

		{1870, 1.61}, {1871, 0.10}, {1872, -1.02}, {1873, -1.28}, {1874, -2.69},
		{1875, -3.24}, {1876, -3.64}, {1877, -4.54}, {1878, -4.71}, {1879, -5.1},

		{1880, -5.40}, {1881, -5.42}, {1882, -5.20}, {1883, -5.46}, {1884, -5.46},
		{1885, -5.79}, {1886, -5.63}, {1887, -5.64}, {1888, -5.80}, {1889, -5.66},

		{1890, -5.87}, {1891, -6.01}, {1892, -6.19}, {1893, -6.64}, {1894, -6.44},
		{1895, -6.47}, {1896, -6.09}, {1897, -5.76}, {1898, -4.66}, {1899, -3.74},

		{1900, -2.72}, {1901, -1.54}, {1902, -0.02}, {1903, 1.24}, {1904, 2.64},
		{1905, 3.86}, {1906, 5.37}, {1907, 6.14}, {1908, 7.75}, {1909, 9.13},

		{1910, 10.46}, {1911, 11.53}, {1912, 13.36}, {1913, 14.65}, {1914, 16.01},
		{1915, 17.20}, {1916, 18.24}, {1917, 19.06}, {1918, 20.25}, {1919, 20.95},

		{1920, 21.16}, {1921, 22.25}, {1922, 22.41}, {1923, 23.03}, {1924, 23.49},
		{1925, 23.69}, {1926, 23.86}, {1927, 24.49}, {1928, 24.34}, {1929, 24.08},

		{1930, 24.02}, {1931, 24.00}, {1932, 23.87}, {1933, 23.95}, {1934, 23.86},
		{1935, 23.93}, {1936, 23.73}, {1937, 23.92}, {1938, 23.96}, {1939, 24.02},

		{1940, 24.33}, {1941, 24.83}, {1942, 25.30}, {1943, 25.70}, {1944, 26.24},
		{1945, 26.77}, {1946, 27.28}, {1947, 27.78}, {1948, 28.25}, {1949, 28.71},

		{1950, 29.15}, {1951, 29.57}, {1952, 29.97}, {1953, 30.36}, {1954, 30.72},
		{1955, 31.07}, {1956, 31.35}, {1957, 31.68}, {1958, 32.18}, {1959, 32.68},

		{1960, 33.15}, {1961, 33.59}, {1962, 34.00}, {1963, 34.47}, {1964, 35.03},
		{1965, 35.73}, {1966, 36.54}, {1967, 37.43}, {1968, 38.29}, {1969, 39.20},

		{1970, 40.18}, {1971, 41.17}, {1972, 42.23}, {1973, 43.37}, {1974, 44.49},
		{1975, 45.48}, {1976, 46.46}, {1977, 47.52}, {1978, 48.53}, {1979, 49.59},

		{1980, 50.54}, {1981, 51.38}, {1982, 52.17}, {1983, 52.96}, {1984, 53.79},
		{1985, 54.34}, {1986, 54.87}, {1987, 55.32}, {1988, 55.82}, {1989, 56.30},

		{1990, 56.86}, {1991, 57.57}, {1992, 58.31}, {1993, 59.12}, {1994, 59.99},
		{1995, 60.78}, {1996, 61.63}, {1997, 62.30}, {1998, 62.97}, {1999, 63.47},

		{2000, 63.83}, {2001, 64.09}, {2002, 64.30}, {2003, 64.47}, {2004, 64.57},
		{2005, 64.69}, {2006, 64.85}, {2007, 65.15}, {2008, 65.46}, {2009, 65.78},

		{2010, 66.07}, {2011, 66.32}, {2012, 66.60}, {2013, 66.91}, {2014, 67.28},
		{2015, 67.18}, {2016, 68.18}, {2017, 69.18}, {2018, 69.18}, {2019, 69.18},
		{2020, 69.18},
	}

	dt.FirstYear = dt.Table[0].Year
	dt.LastYear = dt.Table[len(dt.Table)-1].Year
	return &dt

}
