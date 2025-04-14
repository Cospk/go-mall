package domain

import (
	"fmt"
	"github.com/Cospk/go-mall/pkg/auth"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/dgrijalva/jwt-go"
	"sync"
	"time"
)

//// 生成访问令牌（有效期短）
//func GenerateAccessToken(userID string) (string, error) {
//	claims := auth.MapClaims{
//		"user_id": userID,
//		"exp":     time.Now().Add(time.Minute * 15).Unix(), // 15分钟有效期
//	}
//	token := auth.NewWithClaims(auth.SigningMethodHS256, claims)
//	return token.SignedString([]byte("access-secret"))
//}
//
//// 生成刷新令牌（有效期长）
//func GenerateRefreshToken(userID string) (string, error) {
//	claims := auth.MapClaims{
//		"user_id": userID,
//		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7天有效期
//	}
//	token := auth.NewWithClaims(auth.SigningMethodHS256, claims)
//	return token.SignedString([]byte("refresh-secret"))
//}

type AuthService interface {
	Login(userName string, password string) (string, string, error)
	RefreshToken(token string) (string, error)
}

func GenerateToken(userID int64, platform string) (string, error) {
	claim := auth.CustomClaims{
		UserId:         userID,
		Platform:       platform,
		BufferTime:     0,
		StandardClaims: jwt.StandardClaims{},
	}

	token, err := auth.CreateToken(claim)
	if err != nil {
		return "", errcode.NewError(50001, err.Error())
	}
	return token, nil
}
func PrintNum(n int) {

	var wg sync.WaitGroup
	// 使用channel控制并发
	ch1 := make(chan struct{}, 1)
	ch2 := make(chan struct{}, 1)

	// 打印奇数
	wg.Add(1)
	go func(n int) {
		defer wg.Done()
		for i := 1; i <= n; i++ {
			if i%2 == 1 {
				fmt.Printf("奇数协程输出：%d \n", i)
			}
			ch1 <- struct{}{} // 通知ch1可以执行
			<-ch2
		}
	}(n)

	// 打印偶数
	wg.Add(1)
	go func(n int) {
		defer wg.Done()
		for i := 1; i <= n; i++ {
			<-ch1 // 等待ch1执行
			if i%2 == 0 {
				fmt.Printf("偶数协程输出：%d \n", i)
			}
			ch2 <- struct{}{} // 通知ch2可以执行
		}
	}(n)
	time.Sleep(time.Second * 2)
	wg.Wait()
}

func gToPrintNum(n int, wg *sync.WaitGroup, ch1, ch2 chan struct{}, g int) {
	defer wg.Done()

	if g == 1 {
		for i := 1; i <= n; i++ {
			if i%g == 0 {
				fmt.Printf("协程 %d 输出：%d \n", g, i)
			}
		}
	} else {
		for i := 1; i <= n; i++ {
			<-ch1 // 等待ch1执行
			if i%2 == 0 {
				fmt.Printf("协程2输出：%d \n", i)
			}
		}
	}

}
