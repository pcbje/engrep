package engrep

import (
	"reflect"
	"sort"
	"testing"
)


func RunTestR(t *testing.T, k int, patterns []string, probe string, expected []string, r int) {
	minlength := 99

	for _, p := range patterns {
		if len(p) < minlength {
			minlength = len(p)
		}
	}

	r = (minlength / 2) - 1

	trie := CreateEngrep(k, true, CreateDawg(k))
	trie.AddReferences(patterns)

	cache := map[string]bool{}
	actual := []string{}

	trie.Scan(probe, k, func(s int, e int, str string, pre string, suf string, d int) {
		if _, ok := cache[str]; !ok {
			actual = append(actual, str)
			cache[str] = true
		}
	})

	sort.Strings(actual)
	sort.Strings(expected)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func RunTest(t *testing.T, k int, patterns []string, probe string, expected []string) {
	RunTestR(t, k, patterns, probe, expected, 0)
}

func Test1(t *testing.T) {
	RunTest(t, 1, []string{"00222", "0011"}, "0122200", []string{"01222"})
}

func Test0_1(t *testing.T) {
	RunTest(t, 0, []string{"11001", "10000"}, "110000", []string{"10000"})
}

func Test0_2(t *testing.T) {
	RunTest(t, 0, []string{"1011", "0001", "100001"}, "1011100001011001000001101001101111", []string{"1011", "100001", "0001"})
}

func Test0_3(t *testing.T) {
	RunTest(t, 0, []string{"0111010", "11001011", "10010101"}, "0111011011010010101000001100101", []string{"10010101"})
}

func Test0_4(t *testing.T) {
	RunTest(t, 0, []string{"11011"}, "11000000000000011", []string{})
}

func Test0_5(t *testing.T) {
	RunTest(t, 0, []string{"abcdefxyz", "defghi"}, "abcdefghi....", []string{"defghi"})
}

func Test0_6(t *testing.T) {
	RunTest(t, 0, []string{"0000", "11111"}, "10000", []string{"0000"})
}

func Test0_7(t *testing.T) {
	RunTest(t, 0, []string{"11000"}, "111000", []string{"11000"})
}

func Test0_8(t *testing.T) {
	RunTest(t, 0, []string{"00111001", "11101"}, "001110100", []string{"11101"})
}

func Test0_9(t *testing.T) {
	RunTest(t, 0, []string{"000100110", "1001"}, "00010010000000", []string{"1001"})
}

func Test0_10(t *testing.T) {
	RunTest(t, 0, []string{"01111010", "101100001"}, "101111010101011", []string{"01111010"})
}

// Fails if only testing longest blue.
func Test0_11(t *testing.T) {
	RunTest(t, 0, []string{"abcdy", "bcdx", "cde"}, "abcde", []string{"cde"})
}

// Fails if mis-excluding blues
func Test0_12(t *testing.T) {
	RunTest(t, 0, []string{"bcx", "cde"}, "bcde", []string{"cde"})
}

func TestBjelland(t *testing.T) {
	RunTest(t, 2, []string{"Bjelland"}, "bjelland", []string{"cde"})
}

// Works only when we allow overlapping patterns.
func Test0_13(t *testing.T) {
	RunTest(t, 0, []string{"aaaaaabbaaa", "aabba"}, "x aaaaaabbaa x", []string{"aabba"})
}

func Test0_14(t *testing.T) {
	RunTest(t, 0, []string{"aabba"}, "a  aabbaaa", []string{"aabba"})
}

func Test1_1(t *testing.T) {
	RunTest(t, 1, []string{"00100000"}, "00100100", []string{"00100100"})
}

func Test1_2(t *testing.T) {
	RunTest(t, 1, []string{"010001000"}, "010001100", []string{"010001100"})
}

func Test1_3(t *testing.T) {
	RunTest(t, 1, []string{"11001111", "1001000"}, "1001100", []string{"1001100"})
}

func Test1_5(t *testing.T) {
	RunTest(t, 1, []string{"1011111000", "1111111"}, "111011110", []string{"1101111", "1110111", "11101111"})
}

// Delete
func Test1_6(t *testing.T) {
	RunTest(t, 1, []string{"petter chrisstian"}, "petter christian", []string{"petter christian"})
}

// Insert
func Test1_7(t *testing.T) {
	RunTest(t, 1, []string{"petter christian"}, "petter chrisstian", []string{"etter chrisstian", "petter chrisstia", "petter chrisstian"})
}

// Substitute
func Test1_8(t *testing.T) {
	RunTest(t, 1, []string{"petter christian"}, "petter chrispian", []string{"petter chrispian"})
}

func Test1_9(t *testing.T) {
	RunTest(t, 1, []string{"petter"}, "etter", []string{"etter"})
}

func Test1_10(t *testing.T) {
	RunTest(t, 1, []string{"petter"}, "pette", []string{"pette"})
}

func Test1_11(t *testing.T) {
	RunTest(t, 1, []string{"petter"}, "pett", []string{})
}

func Test1_12(t *testing.T) {
	RunTest(t, 1, []string{"abc", "cde"}, "xabcx", []string{"ab", "abc", "bc"})
}

func Test1_13(t *testing.T) {
	RunTest(t, 1, []string{"abcdefxyz", "defghi"}, "abcdefhi....", []string{"defhi"})
}

func Test1_14(t *testing.T) {
	RunTest(t, 1, []string{"xxaaaaaaaxxxxxx", "aabaaaa"}, "aaabaaaaaa", []string{"aaaaaa", "aaabaaa", "aabaaa", "aabaaaa", "aabaaaaa", "abaaaa", "abaaaaa"})
}

func Test1_15(t *testing.T) {
	RunTest(t, 1, []string{"0101010", "1010101", "00110011", "11001100"}, "0010101010101010101010101010101", []string{"001010", "010101", "0101010", "101010", "1010101"})
}

func Test1_16(t *testing.T) {
	//log.Print("Warning: (Test1_16) Actual distance == 2")
	RunTest(t, 1, []string{"1010010001"}, "11011001001", []string{"1011001001"})
}

func Test1_17(t *testing.T) {
	RunTest(t, 1, []string{"1231", "0043043"}, "220121", []string{"121"})
}

func Test1_18(t *testing.T) {
	RunTest(t, 1, []string{"0022", "0011"}, "aaa0122aaaa", []string{"0122"})
}

// Fails if blue is not checked at max error.
func Test1_19(t *testing.T) {
	RunTest(t, 1, []string{"aabbdd", "bbcc"}, "aabxbcc", []string{"bcc", "bxbc", "bxbcc"})
}

// Blue with insert between
func Test1_20(t *testing.T) {
	RunTest(t, 1, []string{"abcdxx", "cdzyy"}, "abcdyy", []string{"cdyy"})
}

func Test1_20_1(t *testing.T) {
	RunTest(t, 1, []string{"abcdxx", "cdzyy"}, "abcdyy", []string{"cdyy"})
}

// Blue with delete between
func Test1_21(t *testing.T) {
	RunTest(t, 1, []string{"abcdxx", "cdyy"}, "abcdzyy", []string{"cdzy", "cdzyy", "dzyy"})
}

// Why is this working?
func Test1_22(t *testing.T) {
	RunTest(t, 1, []string{"aabbdd", "bbccc", "bbxx"}, "aabbcxx", []string{"bbcx", "bbcxx", "bcxx"})
}

func Test1_23(t *testing.T) {
	RunTest(t, 1, []string{"1010", "0001", "0100"}, "  1100   ", []string{"100", "110", "1100"})
}

func Test1_24(t *testing.T) {
	RunTest(t, 1, []string{"0100"}, "11010", []string{"010", "1010"})
}

func Test1_25(t *testing.T) {
	RunTest(t, 1, []string{"00010000", "00100011", "10001000"}, "110001000110", []string{"0001000", "0010001", "00100011", "0100011", "1000100", "10001000"})
}

func Test1_26(t *testing.T) {
	RunTest(t, 1, []string{"101001", "11000110101", "00011111", "1011000"}, "011000011011", []string{"00011011", "011000", "0110000", "100001", "11000011011"})
}

func Test1_27(t *testing.T) {
	RunTest(t, 1, []string{"1102100201"}, "1102100201", []string{"102100201", "110210020", "1102100201"})
}

func Test1_28(t *testing.T) {
	// 														         1000010001
	RunTest(t, 1, []string{"1000010001"}, "1010010001", []string{"1010010001"})
}

func Test1_29(t *testing.T) {
	RunTest(t, 1, []string{"0022"}, "0122", []string{"0122"})
}

func _Test2_2(t *testing.T) {
	//                                  101011010
	RunTest(t, 2, []string{"1010101000", "101111110"}, "101011010", []string{"101111110", "1010101000"})
}

func Test2_3(t *testing.T) {
	RunTest(t, 2, []string{"101111110"}, "101011010", []string{"101011010"})
}

func Test2_4(t *testing.T) {
	RunTest(t, 2, []string{"petter"}, "tter", []string{"tter"})
}

func Test2_5(t *testing.T) {
	RunTest(t, 2, []string{"petter"}, "pett", []string{"pett"})
}

// Two deletes in between.
func Test2_7(t *testing.T) {
	RunTest(t, 2, []string{"aabbdedfd", "bbcccss", "bbxx"}, "aabbccxx", []string{"bb", "bbccx", "bbccxx", "bccx", "xx"})
}

func Test2_8(t *testing.T) {
	RunTest(t, 2, []string{"aabbddddd", "bbzccss", "bbxx"}, "aabbccxx", []string{"bb", "bbccx", "bbccxx", "bccx", "xx"})
}

func Test2_9(t *testing.T) {
	RunTest(t, 2, []string{"01101010", "10111001", "00010110"}, "011100000110", []string{"000011", "0000110", "000110", "011100", "0111000", "01110000"})
}

func Test2_10(t *testing.T) {
	RunTest(t, 2, []string{"3011442", "4114424"}, "01044212210312221123", []string{"010442"})
}

func Test2_11(t *testing.T) {
	RunTest(t, 2, []string{"0131343", "3040003"}, "14434310032404032101", []string{"3240403", "3431003"})
}

func Test2_12(t *testing.T) {
	RunTest(t, 2, []string{"1330112", "1443113"}, "22142311243112410133", []string{"1241013", "124311", "1243112", "142311", "1423112"})
}

func Test2_13(t *testing.T) {
	RunTest(t, 2, []string{"0331244"}, "24341101221313034244", []string{"034244"})
}

func Test2_14(t *testing.T) {
	RunTest(t, 2, []string{"3114031", "0553231"}, "3014431344", []string{"3014431"})
}

func Test2_15(t *testing.T) {
	RunTest(t, 2, []string{"0000011011"}, "10011011011", []string{"0011011011"})
}

func Test2_16(t *testing.T) {
	RunTest(t, 2, []string{"1210120011", "0221132310"}, "21220200112211121020", []string{"122020011"})
}

func Test2_17(t *testing.T) {
	// 																		 0200002110
	RunTest(t, 2, []string{"0200002110"}, "0101002110", []string{"0101002110"})
}

func Test2_18(t *testing.T) {
	//																														  01410434211
	RunTest(t, 2, []string{"01113434211"}, "01410434211", []string{"01410434211"})
}

func Test2_19(t *testing.T) {
	RunTest(t, 2, []string{"22111111011", "20200011220", "02002102201"}, "0010210210020112101122001", []string{"02102100201"})
}

func Test2_20(t *testing.T) {
	RunTest(t, 2, []string{"12210020222", "01001022102", "02222121211"}, "1020222101222101211010020", []string{"01222101211"})
}

func Test2_21(t *testing.T) {
	RunTest(t, 2, []string{"03231001000", "20131310223", "21202232200"}, "0313100300022013221103130", []string{"03131003000"})
}

func Test2_22(t *testing.T) {
	RunTestR(t, 2, []string{"02220002111"}, "02020012111", []string{"02020012111"}, 3)
}

func Test2_23(t *testing.T) {
	RunTestR(t, 2, []string{"1323313100"}, "1023343100", []string{"1023343100"}, 3)
}

func Test2_24(t *testing.T) {
	RunTestR(t, 2, []string{"Manta Hjelup Fuglstad"}, "Maa Hjelup Fuglstad", []string{"Maa Hjelup Fuglstad"}, 3)
}

func Test2_25(t *testing.T) {
	RunTestR(t, 2, []string{"2011011021"}, "2021021021", []string{"2021021021"}, 3)
}

func TestX_1(t *testing.T) {
	RunTestR(t, 1, []string{"0200010001", "1220110120", "0010121002"}, "1121011012020211012102012", []string{"1210110120"}, 3)
}

func TestX_2(t *testing.T) {
	// 022342401
	RunTestR(t, 2, []string{"0233424101"}, "022342401", []string{"022342401"}, 3)
}

func TestX_3(t *testing.T) {
	RunTestR(t, 2, []string{"41133002312"}, "41433022312", []string{"41433022312"}, 3)
}

func TestX_4(t *testing.T) {
	RunTestR(t, 2, []string{"14042143012"}, "33120223041132431442140012", []string{"1442140012"}, 4)
}

func TestX_5(t *testing.T) {
	RunTestR(t, 3, []string{"3313101311", "3130310311", "0032030303"}, "3111203323333033100003221", []string{"0332333303"}, 3)
}

func TestX_6(t *testing.T) {
	RunTestR(t, 3, []string{"2041302233", "2012221133", "0413230022"}, "3322320231014410204102344", []string{"2041023"}, 3)
}

func TestX_7(t *testing.T) {
	RunTestR(t, 3, []string{"1224103230"}, "3022034303140323030103244", []string{"1403230", "140323030"}, 3)
}

func TestX_8(t *testing.T) {
	RunTestR(t, 2, []string{"abcdefghij"}, "xx cdefghij", []string{"cdefghij"}, 2)
}

func TestX_9(t *testing.T) {
	RunTestR(t, 3, []string{"axxbxcdefghij"}, "abcdefghij", []string{"abcdefghij"}, 3)
}

func TestX_10(t *testing.T) {
	RunTestR(t, 2, []string{"xxabcdefghij"}, "abcdefghij", []string{"abcdefghij"}, 3)
	RunTestR(t, 2, []string{"axbxcdefghij"}, "abcdefghij", []string{"abcdefghij"}, 3)
	RunTestR(t, 2, []string{"axxbcdefghij"}, "abcdefghij", []string{"abcdefghij"}, 3)
	RunTestR(t, 2, []string{"abxxcdefghij"}, "abcdefghij", []string{"abcdefghij"}, 3)
	RunTestR(t, 2, []string{"abxcxdefghij"}, "abcdefghij", []string{"abcdefghij"}, 3)
}

func TestX_11(t *testing.T) {
	RunTestR(t, 1, []string{"100100", "01110"}, "1001111", []string{"0111", "01111"}, 0)
}

func TestX_12(t *testing.T) {
	RunTestR(t, 1, []string{"111"}, "101", []string{"101"}, 0)
}

func TestX_13(t *testing.T) {
	RunTestR(t, 1, []string{"011", "100"}, "000", []string{"000"}, 0)
}

func TestX_14(t *testing.T) {
	RunTestR(t, 1, []string{"011010000", "011100011", "100111110"}, "001111101011", []string{"00111110"}, 0)
}

func TestX_15(t *testing.T) {
	RunTestR(t, 2, []string{"103033031", "221022302", "302120002"}, "322211122302", []string{"211122302", "2211122302"}, 0)
}

func TestX_16(t *testing.T) {
	RunTestR(t, 1, []string{"100000101"}, "110100000001", []string{"100000001"}, 0)
}

func TestX_17(t *testing.T) {
	RunTestR(t, 2, []string{"022002220", "220001012", "221121002"}, "011000010120", []string{"0001012", "0001012", "100001012"}, 0)
}

func TestX_18(t *testing.T) {
	RunTestR(t, 3, []string{"416072670", "664426620", "996321917"}, "699638849174", []string{"9963884917"}, 0)
}

func TestX_19(t *testing.T) {
	RunTestR(t, 2, []string{"0313030133", "3012330021", "3211222101"}, "301233121", []string{"11000010", "110000100"}, 3)
}

func TestX_20(t *testing.T) {
	RunTestR(t, 3, []string{"3000343411"}, "3200314113", []string{"320031411"}, 3)
}

func TestX_21(t *testing.T) {
	RunTestR(t, 1, []string{"1000101011"}, "1000111011", []string{"1000111011"}, 3)
}

func TestX_22(t *testing.T) {
	RunTestR(t, 3, []string{"0011033230", "2332023031", "3123113001"}, "3231000011130330", []string{"323100001"}, 3)
}

func TestX_23(t *testing.T) {
	RunTestR(t, 3, []string{"2260103132", "4110306220", "6031044341"}, "2310160364441026", []string{"60364441"}, 3)
}

func TestX_24(t *testing.T) {
	RunTestR(t, 1, []string{"0010111100110", "1001100011011", "1010111100000"}, "00101010110000001000", []string{"60364441"}, 4)
}

func TestX_25(t *testing.T) {
	RunTestR(t, 3, []string{"0041202041432"}, "0040241432", []string{"0040241432"}, 4)
}

func TestX_26(t *testing.T) {
	RunTestR(t, 3, []string{"3020011334134"}, "300013331134", []string{"300013331134"}, 4)
}

func TestX_27(t *testing.T) {
	RunTestR(t, 3, []string{"0201000011110", "1012021121101", "2100110111021"}, "22102022110120112200", []string{"1020221101"}, 4)
}

func TestX_28(t *testing.T) {
	RunTestR(t, 1, []string{"0210111210220", "1101021202100", "2021120222202"}, "21101020202100012002", []string{"1020221101"}, 4)
}

func TestX_29(t *testing.T) {
	RunTestR(t, 2, []string{"0001101000111", "1000010100011", "1111100111110"}, "00010001110011111000", []string{"11100111110"}, 4)
}

func TestX_30(t *testing.T) {
	RunTestR(t, 2, []string{"0012011202020", "1110202100212", "2110002210212"}, "21201120202001021101", []string{"12011202020", "212011202020"}, 4)
}

func TestX_31(t *testing.T) {
	RunTestR(t, 3, []string{"3132220103010"}, "01100230322200000102", []string{"3032220000010"}, 4)
}

func TestX_32(t *testing.T) {
	RunTestR(t, 3, []string{"0212210020100", "1001110001011", "2211111220112"}, "02022110010000000112", []string{"020221100100"}, 4)
}

func TestX_33(t *testing.T) {
	RunTestR(t, 3, []string{"0124412334342", "1134440310413", "3104021034323"}, "20300210143234410112", []string{"30021014323"}, 4)
}

func TestX_34(t *testing.T) {
	RunTestR(t, 1, []string{"0110010010011"}, "011001100100100010", []string{"0110010010001"}, 4)
}

/*

01011001100100100010 [0000110110001 0100000100100 0101110010001 0110010010011 1011000000101]

[01011000010010] [] 00000010110000100100 [0001111001100 0011001011011 0011101100110 0101100000010 0110001011000]

[140222243033] [] 14022224303310011232 [0031302121330 0343014411340 0433422322042 1412221241033 2441420130301]

[0313044400322] [] 03130444003221330003 [0213044320322 4033040031231 4244403200112]

[30021014323] [] 20300210143234410112 [0124412334342 1134440310413 3104021034323

[23302002203] [] 40423302002203210412 [0040211333223 0301222040013 2320200242013]


[020221100100] [] 02022110010000000112 [0212210020100 1001110001011 2211111220112]

[3032220000010] [] 01100230322200000102 [0203103011031 3132220103010 3320111002101] 0



00010001110011111000 [0001101000111 1000010100011 1111100111110]

[1101020202100] [] 21101020202100012002 [0210111210220 1101021202100 2021120222202]

[1020221101] [211012011 1012011220 10120112200] 22102022110120112200 [0201000011110 1012021121101 2100110111021]

[300013331134] [] 31003000133311340010 [1330422223313 2122243130423 3020011334134]

[0110111101100] [] 10011011110110000000 [0110110101100 1000100110001 1101001010011]

02420202140040241432 [0041202041432 0302433430333 2202301414003]

[111110010010 1111100100100] [] 11111001001001011001 [0001100001000 0101001111010 1111100110000]


[1012100121102] [] 21012010121001211022 [1012000121102 1220212210100 2101001110022

[1010101100000] [] 00101010110000001000 [0010111100110 1001100011011 1010111100000]

[60364441] [23101603 1016036 41026] 2310160364441026 [2260103132 4110306220 6031044341]

[323100001] 3231000011130330 [0011033230 2332023031 3123113001]

[1000111011] [] 0011100011101110 [0000100010 1000101011 1110101011]
[320031411] [] 3222323200314113 [0241420122 3000343411 4231242403]


2200020020002112 [1112211010 2200222002 2221220121
1101010110000000 [0110000010 1010111100 1011111101]


1101101011010110 [0101100011 1101100000 1110001100


1000110110011111 [0000000111 1101000111 1110000100]

0111000001000000 [0000011010 0000101100 1100010100]


 2020111210222220 [0111220122 2110000212 2212001222]

[0101111011]
 0101110011
*/
// 0001010111101111 [0010100010 0101110011 1001110000]
//[9963884917] [] 699638849174 [416072670 664426620 996321917]

//[100000001] [100000001] [110100] 110100000001 [010100111 100000101 101111001
//[211122302 2211122302] [] 322211122302 [103033031 221022302 302120002]

//[01100111 00001100] [00001100] [0110011 01100111 1001110] 000011001110 [011001011 011100011 100001100]

//[112020021 12020021] [] 211202002112 [010210022 012220021 110011120]

// [6213314633] [6213314633] [] 006213314633 [063402316 130121363 623314633
