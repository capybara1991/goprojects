package iprange

import (
	"encoding/binary"
	"fmt"
	"net"
)

type AddressRangeList []AddressRange

type AddressRange struct {
	Min net.IP
	Max net.IP
}

type octetRange struct {
	min byte
	max byte
}

type ipSymType struct {
	yys       int
	num       byte
	octRange  octetRange
	addrRange AddressRange
	result    AddressRangeList
}

const num = 57346

var ipToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"num",
	"','",
	"' '",
	"'/'",
	"'.'",
	"'*'",
	"'-'",
}

const ipEofCode = 1
const ipErrCode = 2
const ipInitialStackSize = 16

func ParseList(in string) (AddressRangeList, error) {
	lex := &ipLex{line: []byte(in)}
	errCode := ipParse(lex)
	if errCode != 0 || lex.err != nil {
		return nil, fmt.Errorf("could not parse target: %w", lex.err)
	}
	return lex.output, nil
}

func Parse(in string) (*AddressRange, error) {
	l, err := ParseList(in)
	if err != nil {
		return nil, err
	}
	return &l[0], nil
}

var ipExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const ipPrivate = 57344

const ipLast = 22

var ipAct = [...]int{
	4, 5, 12, 20, 2, 10, 6, 18, 11, 14,
	9, 17, 16, 13, 15, 8, 1, 7, 3, 19,
	0, 21,
}
var ipPact = [...]int{
	-3, 5, -1000, -2, 0, -8, -1000, -1000, -3, 3,
	10, -3, 7, -1000, -1000, -1000, -1, -1000, -3, -5,
	-3, -1000,
}
var ipPgo = [...]int{
	0, 18, 4, 0, 17, 16, 15,
}
var ipR1 = [...]int{
	0, 5, 5, 6, 6, 2, 2, 1, 3, 3,
	3, 4,
}
var ipR2 = [...]int{
	0, 1, 3, 1, 2, 3, 1, 7, 1, 1,
	1, 3,
}
var ipChk = [...]int{
	-1000, -5, -2, -1, -3, 4, 9, -4, -6, 5,
	7, 8, 10, -2, 6, 4, -3, 4, 8, -3,
	8, -3,
}
var ipDef = [...]int{
	0, -2, 1, 6, 0, 8, 9, 10, 0, 3,
	0, 0, 0, 2, 4, 5, 0, 11, 0, 0,
	0, 7,
}

var ipTok1 = [...]int{
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 6, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 9, 3, 5, 10, 8, 7,
}
var ipTok2 = [...]int{
	2, 3, 4,
}
var ipTok3 = [...]int{0}

func ipTokname(c int) string {
	if c >= 1 && c-1 < len(ipToknames) {
		if ipToknames[c-1] != "" {
			return ipToknames[c-1]
		}
	}
	return fmt.Sprintf("tok-%v", c)
}

func ipErrorMessage(state, lookAhead int) string {
	return "syntax error"
}

type ipLexer interface {
	Lex(lval *ipSymType) int
	Error(s string)
}

type ipParser interface {
	Parse(ipLexer) int
	Lookahead() int
}

type ipParserImpl struct {
	lval  ipSymType
	stack [ipInitialStackSize]ipSymType
	char  int
}

func (p *ipParserImpl) Lookahead() int { return p.char }

func ipNewParser() ipParser { return &ipParserImpl{} }

const ipFlag = -1000

func iplex1(lex ipLexer, lval *ipSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = ipTok1[0]
		goto out
	}
	if char < len(ipTok1) {
		token = ipTok1[char]
		goto out
	}
	if char >= ipPrivate {
		if char < ipPrivate+len(ipTok2) {
			token = ipTok2[char-ipPrivate]
			goto out
		}
	}
	for i := 0; i < len(ipTok3); i += 2 {
		token = ipTok3[i+0]
		if token == char {
			token = ipTok3[i+1]
			goto out
		}
	}
out:
	if token == 0 {
		token = ipTok2[1]
	}
	return char, token
}

//nolint:gocognit
func ipParse(iplex ipLexer) int {
	return ipNewParser().Parse(iplex)
}

//nolint:gocognit
func (iprcvr *ipParserImpl) Parse(iplex ipLexer) int {
	var ipn int
	var ipVAL ipSymType
	var ipDollar []ipSymType
	ipS := iprcvr.stack[:]

	Nerrs := 0
	Errflag := 0
	ipstate := 0
	iprcvr.char = -1
	iptoken := -1
	defer func() {
		ipstate = -1
		iprcvr.char = -1
		iptoken = -1
	}()
	ipp := -1
	goto ipstack

ret0:
	return 0
ret1:
	return 1

ipstack:
	ipp++
	if ipp >= len(ipS) {
		nyys := make([]ipSymType, len(ipS)*2)
		copy(nyys, ipS)
		ipS = nyys
	}
	ipS[ipp] = ipVAL
	ipS[ipp].yys = ipstate

ipnewstate:
	ipn = ipPact[ipstate]
	if ipn <= ipFlag {
		goto ipdefault
	}
	if iprcvr.char < 0 {
		iprcvr.char, iptoken = iplex1(iplex, &iprcvr.lval)
	}
	ipn += iptoken
	if ipn < 0 || ipn >= ipLast {
		goto ipdefault
	}
	ipn = ipAct[ipn]
	if ipChk[ipn] == iptoken {
		iprcvr.char = -1
		iptoken = -1
		ipVAL = iprcvr.lval
		ipstate = ipn
		if Errflag > 0 {
			Errflag--
		}
		goto ipstack
	}

ipdefault:
	ipn = ipDef[ipstate]
	if ipn == -2 {
		if iprcvr.char < 0 {
			iprcvr.char, iptoken = iplex1(iplex, &iprcvr.lval)
		}
		xi := 0
		for {
			if ipExca[xi+0] == -1 && ipExca[xi+1] == ipstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			ipn = ipExca[xi+0]
			if ipn < 0 || ipn == iptoken {
				break
			}
		}
		ipn = ipExca[xi+1]
		if ipn < 0 {
			goto ret0
		}
	}
	if ipn == 0 {
		switch Errflag {
		case 0:
			iplex.Error(ipErrorMessage(ipstate, iptoken))
			Nerrs++
			fallthrough
		case 1, 2:
			Errflag = 3
			for ipp >= 0 {
				ipn = ipPact[ipS[ipp].yys] + ipErrCode
				if ipn >= 0 && ipn < ipLast {
					ipstate = ipAct[ipn]
					if ipChk[ipstate] == ipErrCode {
						goto ipstack
					}
				}
				ipp--
			}
			goto ret1
		case 3:
			if iptoken == ipEofCode {
				goto ret1
			}
			iprcvr.char = -1
			iptoken = -1
			goto ipnewstate
		}
	}

	ipnt := ipn
	ippt := ipp
	ipp -= ipR2[ipn]
	if ipp+1 >= len(ipS) {
		nyys := make([]ipSymType, len(ipS)*2)
		copy(nyys, ipS)
		ipS = nyys
	}
	ipVAL = ipS[ipp+1]

	ipn = ipR1[ipn]
	ipg := ipPgo[ipn]
	ipj := ipg + ipS[ipp].yys + 1

	if ipj >= ipLast {
		ipstate = ipAct[ipg]
	} else {
		ipstate = ipAct[ipj]
		if ipChk[ipstate] != -ipn {
			ipstate = ipAct[ipg]
		}
	}
	switch ipnt {
	case 1:
		ipDollar = ipS[ippt-1 : ippt+1]
		{
			tmp := append(AddressRangeList(nil), ipDollar[1].addrRange)
			ipVAL.result = tmp
			iplex.(*ipLex).output = ipVAL.result
		}
	case 2:
		ipDollar = ipS[ippt-3 : ippt+1]
		{
			tmp := append(append(AddressRangeList(nil), ipDollar[1].result...), ipDollar[3].addrRange)
			ipVAL.result = tmp
			iplex.(*ipLex).output = ipVAL.result
		}
	case 5:
		ipDollar = ipS[ippt-3 : ippt+1]
		{
			if int(ipDollar[3].num) > 32 {
				iplex.(*ipLex).Error("invalid cidr size")
				ipVAL.addrRange = AddressRange{}
			} else {
				mask := net.CIDRMask(int(ipDollar[3].num), 32)
				min := ipDollar[1].addrRange.Min.Mask(mask)
				maxInt := binary.BigEndian.Uint32([]byte(min)) + 0xffffffff - binary.BigEndian.Uint32([]byte(mask))
				maxBytes := make([]byte, 4)
				binary.BigEndian.PutUint32(maxBytes, maxInt)
				maxBytes = maxBytes[len(maxBytes)-4:]
				max := net.IP(maxBytes)
				ipVAL.addrRange = AddressRange{
					Min: min.To4(),
					Max: max.To4(),
				}
			}
		}
	case 6:
		ipDollar = ipS[ippt-1 : ippt+1]
		{
			ipVAL.addrRange = ipDollar[1].addrRange
		}
	case 7:
		ipDollar = ipS[ippt-7 : ippt+1]
		{
			ipVAL.addrRange = AddressRange{
				Min: net.IPv4(ipDollar[1].octRange.min, ipDollar[3].octRange.min, ipDollar[5].octRange.min, ipDollar[7].octRange.min).To4(),
				Max: net.IPv4(ipDollar[1].octRange.max, ipDollar[3].octRange.max, ipDollar[5].octRange.max, ipDollar[7].octRange.max).To4(),
			}
		}
	case 8:
		ipDollar = ipS[ippt-1 : ippt+1]
		{
			ipVAL.octRange = octetRange{ipDollar[1].num, ipDollar[1].num}
		}
	case 9:
		ipDollar = ipS[ippt-1 : ippt+1]
		{
			ipVAL.octRange = octetRange{0, 255}
		}
	case 10:
		ipDollar = ipS[ippt-1 : ippt+1]
		{
			ipVAL.octRange = ipDollar[1].octRange
		}
	case 11:
		ipDollar = ipS[ippt-3 : ippt+1]
		{
			if ipDollar[1].num > ipDollar[3].num {
				iplex.(*ipLex).Error("invalid octet range")
				ipVAL.octRange = octetRange{0, 0}
			} else {
				ipVAL.octRange = octetRange{ipDollar[1].num, ipDollar[3].num}
			}
		}
	}
	goto ipstack
}
