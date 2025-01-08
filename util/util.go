package util

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zachturing/login/common/define"
)

// 指定加密密钥
var jwtSecret = []byte("mix_paper_dev")

func init() {
	tmpSecret := os.Getenv("JWT_SECRET_KEY")
	if len(tmpSecret) > 0 {
		jwtSecret = []byte(tmpSecret)
	}
}

// Claims 是一些实体（通常指的用户）的状态和额外的元数据
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 根据用户的用户名和密码产生token
func GenerateToken(userID int) (string, int64, error) {
	expiredTime := time.Now().Add(define.TokenExpireTime)
	// 生成时间戳
	expiredTimeStamp := expiredTime.Unix()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, expiredTimeStamp, err
}

// ParseToken 根据传入的token值获取到Claims对象信息（进而获取其中的用户id）
func ParseToken(token string) (*Claims, error) {
	// 用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// AlphanumericSet 自定义字符集，去除了容易混淆的字符
var AlphanumericSet = []rune{
	'2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'm', 'n', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

// GenerateInvCodeByUserId 获取指定长度的邀请码，默认为6位
func GenerateInvCodeByUserId(uid uint64) string {
	// 放大 + 加盐
	uid = uid*define.PRIME1 + define.SALT
	// 邀请码长度默认为6位
	l := define.DefaultInvCodeLength

	var code []rune
	slIdx := make([]byte, l)

	// 扩散
	for i := 0; i < l; i++ {
		slIdx[i] = byte(uid % uint64(len(AlphanumericSet)))                   // 获取 52 进制的每一位值
		slIdx[i] = (slIdx[i] + byte(i)*slIdx[0]) % byte(len(AlphanumericSet)) // 其他位与个位加和再取余（让个位的变化影响到所有位）
		uid = uid / uint64(len(AlphanumericSet))                              // 相当于右移一位（52进制）
	}

	// 混淆
	for i := 0; i < l; i++ {
		idx := (byte(i) * define.PRIME2) % byte(l)
		code = append(code, AlphanumericSet[slIdx[idx]])
	}
	return string(code)
}

func invertAlphanumericMap() map[rune]byte {
	inverseMap := make(map[rune]byte)
	for i, r := range AlphanumericSet {
		inverseMap[r] = byte(i)
	}
	return inverseMap
}

// DecodeInvCodeToUID 根据邀请码解析出userId，返回值为0说明邀请码不存在
func DecodeInvCodeToUID(code string) (uint64, error) {
	l := len(code)
	if l == 0 {
		return 0, fmt.Errorf("inv code is empty")
	}

	inverseMap := invertAlphanumericMap()

	// 反混淆
	slIdx := make([]byte, l)
	for i, r := range code {
		idx, ok := inverseMap[r]
		if !ok {
			return 0, fmt.Errorf("inv code is invalid")
		}
		slIdx[(byte(i)*define.PRIME2)%byte(l)] = idx
	}

	// 反扩散
	var uid uint64
	for i := l - 1; i >= 0; i-- {
		if i > 0 {
			slIdx[i] = (slIdx[i] + byte(len(AlphanumericSet)) - byte(i)*slIdx[0]%byte(len(AlphanumericSet))) % byte(len(AlphanumericSet))
		}
		uid = uid*uint64(len(AlphanumericSet)) + uint64(slIdx[i])
	}

	// 去盐
	if uid < define.SALT {
		return 0, fmt.Errorf("inv code is invalid")
	}
	uid = (uid - define.SALT) / define.PRIME1

	return uid, nil
}
