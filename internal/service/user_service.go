package service

import (
	"errors"

	"github.com/KeLes-Coding/gopress/internal/dao"
	"github.com/KeLes-Coding/gopress/internal/model"
	"github.com/KeLes-Coding/gopress/internal/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 结构体封装了所有与用户相关的业务逻辑。
type UserService struct{}

// NewUserService 是一个工厂函数，用于创建 UserService 的实例。
// 这样做可以方便未来的扩展，例如，如果 service 需要依赖其他组件（如 Redis 客户端），
// 可以在这里注入。
func NewUserService() *UserService {
	return &UserService{}
}

// SignUp 处理用户注册的核心逻辑。
func (s *UserService) SignUp(username, password string) error {
	// 1. 参数校验
	if len(username) < 4 || len(password) < 6 {
		return errors.New("用户名长度不能少于4位，密码长度不能少于6位")
	}

	// 2. 检查用户名是否存在
	// 我们直接调用 dao 层的 GetDB() 来获取数据库连接，并查询 users 表。
	db := dao.GetDB()
	var user model.User
	// First 方法会查询第一条满足条件的记录。
	// 如果记录不存在，GORM V2 会返回 gorm.ErrRecordNotFound 错误。
	err := db.Where("username = ?", username).First(&user).Error
	if err == nil {
		// 如果 err 为 nil，说明找到了记录，用户名已存在。
		return errors.New("用户名已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果是一个非 "记录未找到" 的其他数据库错误，直接返回。
		return err
	}

	// 3. 密码加密
	// 使用 bcrypt 算法对用户密码进行哈希处理。
	// GenerateFromPassword 的第二个参数是 cost，值越高，哈希计算越慢，密码也就越安全。
	// bcrypt.DefaultCost (值为10) 是一个推荐的默认值。
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// 如果加密过程中出错，返回错误。
		return err
	}

	// 4. 创建用户实例并存入数据库
	newUser := model.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		// Nickname, Email, Role 等字段会使用其零值或数据库定义的默认值。
	}

	// 调用 GORM 的 Create 方法将新用户记录插入数据库。
	if err := db.Create(&newUser).Error; err != nil {
		return err
	}

	// 注册成功，返回 nil。
	return nil
}

// Login 处理用户登录的业务逻辑
// 成功是返回生成的 JWT ，失败时返回错误
func (s *UserService) Login(username, password string) (string, error) {
	// 1. 根据用户名查询用户
	db := dao.GetDB()
	var user model.User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果记录未找到，返回一个对用户更友好的错误信息
			return "", errors.New("用户名或密码错误")
		}
		// 其他数据库错误
		return "", err
	}

	// 2. 校验密码
	// 使用 bcrypt.CompareHashAndPassword 来比较哈希后的密码和用户输入的明文密码。
	// 这个函数可以有效防止时序攻击 (Timing Attack)。
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// 如果密码不匹配，err 会是 bcrypt.ErrMismatchedHashAndPassword。
		// 为了安全，我们同样返回一个模糊的错误提示。
		return "", errors.New("用户名或密码错误")
	}

	// 3. 生成 JWT
	// 登陆成功，调用 util 包中的 GenerateToken 函数生成 token。
	token, err := util.GenerateToken(user.ID, user.Username)
	if err != nil {
		// 如果 token 生成失败，这是一个服务端内部错误
		return "", err
	}

	// 4. 返回 JWT
	return token, nil
}
