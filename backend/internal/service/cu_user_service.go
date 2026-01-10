package service

import (
	"context"
	"errors"
	"time"

	"license-manager/internal/config"
	"license-manager/internal/models"
	"license-manager/internal/repository"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
	"license-manager/pkg/utils"

	"gorm.io/gorm"
)

type CuUserService interface {
	Register(ctx context.Context, req *models.CuUserRegisterRequest) (*models.CuUser, error)
	Login(ctx context.Context, req *models.CuUserLoginRequest, ip string) (*models.CuUser, string, error)
	GetProfile(ctx context.Context, userID string) (*models.CuUser, error)
	UpdateProfile(ctx context.Context, userID string, req *models.CuUserProfileUpdateRequest) error
	UpdatePhone(ctx context.Context, userID string, req *models.CuUserPhoneUpdateRequest) error
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	SendRegisterSms(ctx context.Context, req *models.CuUserSendRegisterSmsRequest) error
	ForgotPassword(ctx context.Context, req *models.CuUserForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req *models.CuUserResetPasswordRequest) error
	GetUsersByCustomer(ctx context.Context, customerID string, offset, limit int) ([]*models.CuUser, int64, error)
	LockAccount(ctx context.Context, userID string, duration time.Duration) error
	UnlockAccount(ctx context.Context, userID string) error
}

type cuUserService struct {
	repo         repository.CuUserRepository
	customerRepo repository.CustomerRepository
	db           *gorm.DB
}

func NewCuUserService(repo repository.CuUserRepository, customerRepo repository.CustomerRepository, db *gorm.DB) CuUserService {
	return &cuUserService{
		repo:         repo,
		customerRepo: customerRepo,
		db:           db,
	}
}

func (s *cuUserService) Register(ctx context.Context, req *models.CuUserRegisterRequest) (*models.CuUser, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 设置默认值
	phoneCountryCode := req.PhoneCountryCode
	if phoneCountryCode == "" {
		phoneCountryCode = "+86" // 默认中国
	}

	// 检查手机号是否已存在
	exists, err := s.repo.CheckPhoneExists(req.Phone, phoneCountryCode, "")
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}
	if exists {
		return nil, i18n.NewI18nError("500002", lang) // 手机号已被注册
	}

	// 检查邮箱是否已存在（如果提供）
	if req.Email != nil {
		exists, err := s.repo.CheckEmailExists(*req.Email, "")
		if err != nil {
			return nil, i18n.NewI18nError("900004", lang, err.Error())
		}
		if exists {
			return nil, i18n.NewI18nError("500002", lang) // 邮箱已被注册
		}
	}

	// 处理客户ID：如果未提供，则自动创建客户
	customerID := req.CustomerID
	if customerID == "" {
		// 创建新客户
		customer, err := s.createDefaultCustomer(ctx, req.Phone, phoneCountryCode)
		if err != nil {
			return nil, err
		}
		customerID = customer.ID
	}

	// 生成密码哈希
	var passwordHash, salt string
	if req.Password != "" {
		salt = utils.GenerateSalt()
		passwordHash, err = utils.HashPasswordWithSalt(req.Password, salt)
		if err != nil {
			return nil, i18n.NewI18nError("900004", lang, err.Error())
		}
	}

	// 创建用户
	user := &models.CuUser{
		CustomerID:       customerID,
		Phone:            req.Phone,
		PhoneCountryCode: phoneCountryCode,
		Password:         &passwordHash,
		Salt:             &salt,
		UserRole:         "member", // 默认普通成员
		RealName:         req.RealName,
		Email:            req.Email,
		Status:           "active",
		PhoneVerified:    true, // 注册时已验证
		EmailVerified:    false,
		Language:         "zh-CN",
		Timezone:         "Asia/Shanghai",
	}

	// 委托给Repository层进行数据创建
	if err := s.repo.Create(user); err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	return user, nil
}

func (s *cuUserService) Login(ctx context.Context, req *models.CuUserLoginRequest, ip string) (*models.CuUser, string, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if req == nil {
		return nil, "", i18n.NewI18nError("900001", lang)
	}

	// 设置默认国家代码
	phoneCountryCode := req.PhoneCountryCode
	if phoneCountryCode == "" {
		phoneCountryCode = "+86" // 默认中国
	}

	// 根据手机号查找用户
	user, err := s.repo.GetByPhone(req.Phone, phoneCountryCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", i18n.NewI18nError("500003", lang) // 手机号或密码错误
		}
		return nil, "", i18n.NewI18nError("900004", lang, err.Error())
	}

	// 检查账号状态
	if user.Status != "active" {
		return nil, "", i18n.NewI18nError("500003", lang) // 账号已被禁用
	}

	// 检查是否被锁定
	if user.IsAccountLocked() {
		return nil, "", i18n.NewI18nError("500003", lang) // 账号已被锁定
	}

	// 验证密码
	if user.Password == nil {
		return nil, "", i18n.NewI18nError("500003", lang) // 请先设置密码
	}

	if !utils.CheckPassword(req.Password, *user.Password, *user.Salt) {
		// 增加登录失败次数
		if err := s.repo.IncrementLoginAttempts(user.ID); err != nil {
			// 记录错误但不影响登录流程
		}

		// 检查是否需要锁定账号
		cfg := config.GetConfig()
		if cfg != nil && user.LoginAttempts+1 >= cfg.Auth.Security.MaxLoginAttempts {
			lockUntil := time.Now().Add(time.Duration(cfg.Auth.Security.LockoutDurationMinutes) * time.Minute)
			if err := s.repo.LockAccount(user.ID, lockUntil); err != nil {
				// 记录错误但不影响登录流程
			}
			return nil, "", i18n.NewI18nError("500003", lang) // 账号已被锁定
		}

		return nil, "", i18n.NewI18nError("500003", lang) // 手机号或密码错误
	}

	// 登录成功，重置登录失败次数并更新登录信息
	if err := s.repo.ResetLoginAttempts(user.ID); err != nil {
		// 记录错误但不影响登录流程
	}
	if err := s.repo.UpdateLoginInfo(user.ID, ip); err != nil {
		// 记录错误但不影响登录流程
	}

	// 生成JWT Token
	token, err := utils.GenerateCuToken(user.ID, user.CustomerID, user.UserRole, user.Phone)
	if err != nil {
		return nil, "", i18n.NewI18nError("500013", lang, err.Error())
	}

	return user, token, nil
}

func (s *cuUserService) GetProfile(ctx context.Context, userID string) (*models.CuUser, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if userID == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}

	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, i18n.NewI18nError("500009", lang) // 用户不存在
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	return user, nil
}

func (s *cuUserService) UpdateProfile(ctx context.Context, userID string, req *models.CuUserProfileUpdateRequest) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if userID == "" || req == nil {
		return i18n.NewI18nError("900001", lang)
	}

	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i18n.NewI18nError("500009", lang) // 用户不存在
		}
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	// 检查邮箱是否已被其他用户使用
	if req.Email != nil && *req.Email != "" {
		exists, err := s.repo.CheckEmailExists(*req.Email, userID)
		if err != nil {
			return i18n.NewI18nError("900004", lang, err.Error())
		}
		if exists {
			return i18n.NewI18nError("500006", lang) // 邮箱已被其他用户使用
		}
		user.Email = req.Email
		user.EmailVerified = false // 修改邮箱后需要重新验证
	}

	// 更新其他字段
	if req.RealName != nil {
		user.RealName = req.RealName
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	if req.Language != nil {
		user.Language = *req.Language
	}
	if req.Timezone != nil {
		user.Timezone = *req.Timezone
	}
	if req.AdditionalInfo != nil {
		user.AdditionalInfo = req.AdditionalInfo
	}
	if req.Remark != nil {
		user.Remark = req.Remark
	}

	if err := s.repo.Update(user); err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	return nil
}

func (s *cuUserService) UpdatePhone(ctx context.Context, userID string, req *models.CuUserPhoneUpdateRequest) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if userID == "" || req == nil {
		return i18n.NewI18nError("900001", lang)
	}

	// 检查新手机号是否已被使用
	exists, err := s.repo.CheckPhoneExists(req.NewPhone, req.NewPhoneCountryCode, userID)
	if err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}
	if exists {
		return i18n.NewI18nError("500007", lang) // 新手机号已被其他用户使用
	}

	// 验证新手机号（这里需要调用短信验证服务，暂时跳过）

	// 更新手机号
	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i18n.NewI18nError("500009", lang) // 用户不存在
		}
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	user.Phone = req.NewPhone
	user.PhoneCountryCode = req.NewPhoneCountryCode
	user.PhoneVerified = true // 手机号已验证

	if err := s.repo.Update(user); err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	return nil
}

func (s *cuUserService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if userID == "" || oldPassword == "" || newPassword == "" {
		return i18n.NewI18nError("900001", lang)
	}

	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i18n.NewI18nError("500009", lang) // 用户不存在
		}
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	// 验证旧密码
	if user.Password == nil {
		return i18n.NewI18nError("500008", lang) // 请先设置密码
	}

	if !utils.CheckPassword(oldPassword, *user.Password, *user.Salt) {
		return i18n.NewI18nError("500008", lang) // 旧密码错误
	}

	// 生成新密码哈希
	salt := utils.GenerateSalt()
	passwordHash, err := utils.HashPasswordWithSalt(newPassword, salt)
	if err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	user.Password = &passwordHash
	user.Salt = &salt

	if err := s.repo.Update(user); err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	return nil
}

func (s *cuUserService) SendRegisterSms(ctx context.Context, req *models.CuUserSendRegisterSmsRequest) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if req == nil {
		return i18n.NewI18nError("900001", lang)
	}

	// 设置默认国家代码
	phoneCountryCode := req.PhoneCountryCode
	if phoneCountryCode == "" {
		phoneCountryCode = "+86" // 默认中国
	}

	// 检查手机号是否已被注册
	exists, err := s.repo.CheckPhoneExists(req.Phone, phoneCountryCode, "")
	if err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}
	if exists {
		return i18n.NewI18nError("200002", lang) // 客户已存在
	}

	// TODO: 实现频率限制检查
	// 这里应该检查Redis中该手机号的发送频率
	// 例如：1分钟内最多发送1次，1小时内最多发送5次
	// 如果超出限制，返回错误码 200003

	// TODO: 生成6位随机验证码并发送短信
	// 1. 生成随机验证码
	// 2. 调用短信服务发送验证码
	// 3. 将验证码存储到Redis中，设置过期时间（5分钟）
	// 例如:
	// verificationCode := generateRandomCode(6)
	// smsService.SendSms(phoneCountryCode + req.Phone, fmt.Sprintf("您的注册验证码是：%s", verificationCode))
	// redis.Set(fmt.Sprintf("register_sms:%s:%s", phoneCountryCode, req.Phone), verificationCode, 5*time.Minute)
	// 如果发送失败，返回错误码 200004

	// 暂时返回成功，等待短信服务实现
	return nil
}

func (s *cuUserService) ForgotPassword(ctx context.Context, req *models.CuUserForgotPasswordRequest) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if req == nil {
		return i18n.NewI18nError("900001", lang)
	}

	// 设置默认国家代码
	phoneCountryCode := req.PhoneCountryCode
	if phoneCountryCode == "" {
		phoneCountryCode = "+86" // 默认中国
	}

	// 查找用户
	user, err := s.repo.GetByPhone(req.Phone, phoneCountryCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i18n.NewI18nError("500009", lang) // 用户不存在
		}
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	// 检查账号状态
	if user.Status != "active" {
		return i18n.NewI18nError("500004", lang) // 账号已被禁用
	}

	// TODO: 实现短信验证码发送逻辑
	// 1. 生成6位随机验证码
	// 2. 调用短信服务发送验证码
	// 3. 将验证码存储到Redis中，设置过期时间（5分钟）
	// 例如:
	// verificationCode := generateRandomCode(6)
	// smsService.SendSms(user.Phone, fmt.Sprintf("您的密码重置验证码是：%s", verificationCode))
	// redis.Set(fmt.Sprintf("reset_pwd:%s", user.Phone), verificationCode, 5*time.Minute)
	// 暂时返回成功，等待短信服务实现

	return nil
}

func (s *cuUserService) ResetPassword(ctx context.Context, req *models.CuUserResetPasswordRequest) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if req == nil {
		return i18n.NewI18nError("900001", lang)
	}

	// 设置默认国家代码
	phoneCountryCode := req.PhoneCountryCode
	if phoneCountryCode == "" {
		phoneCountryCode = "+86" // 默认中国
	}

	// 查找用户
	user, err := s.repo.GetByPhone(req.Phone, phoneCountryCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i18n.NewI18nError("500009", lang) // 用户不存在
		}
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	// 检查账号状态
	if user.Status != "active" {
		return i18n.NewI18nError("500005", lang) // 账号已被禁用
	}

	// TODO: 实现短信验证码验证逻辑
	// 这里应该验证用户提交的验证码是否正确且未过期
	// 例如:
	// storedCode := redis.Get(fmt.Sprintf("reset_pwd:%s", user.Phone))
	// if storedCode == "" { return i18n.NewI18nError("500021", lang) } // 验证码不存在
	// if storedCode != req.SmsCode { return i18n.NewI18nError("500012", lang) } // 验证码错误
	// redis.Del(fmt.Sprintf("reset_pwd:%s", user.Phone)) // 验证成功后删除验证码
	// 暂时跳过验证码验证，等待短信服务实现

	// 生成新密码哈希
	salt := utils.GenerateSalt()
	passwordHash, err := utils.HashPasswordWithSalt(req.NewPassword, salt)
	if err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	user.Password = &passwordHash
	user.Salt = &salt
	user.LockedUntil = nil // 解锁账号
	user.LoginAttempts = 0 // 重置登录失败次数

	if err := s.repo.Update(user); err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	return nil
}

func (s *cuUserService) GetUsersByCustomer(ctx context.Context, customerID string, offset, limit int) ([]*models.CuUser, int64, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if customerID == "" {
		return nil, 0, i18n.NewI18nError("900001", lang)
	}

	users, total, err := s.repo.GetByCustomerID(customerID, offset, limit)
	if err != nil {
		return nil, 0, i18n.NewI18nError("900004", lang, err.Error())
	}

	return users, total, nil
}

func (s *cuUserService) LockAccount(ctx context.Context, userID string, duration time.Duration) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if userID == "" {
		return i18n.NewI18nError("900001", lang)
	}

	lockUntil := time.Now().Add(duration)
	if err := s.repo.LockAccount(userID, lockUntil); err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	return nil
}

func (s *cuUserService) UnlockAccount(ctx context.Context, userID string) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if userID == "" {
		return i18n.NewI18nError("900001", lang)
	}

	if err := s.repo.LockAccount(userID, nil); err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	return nil
}

// createDefaultCustomer 创建默认客户（用于自动注册场景）
func (s *cuUserService) createDefaultCustomer(ctx context.Context, phone, phoneCountryCode string) (*models.Customer, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 使用手机号作为客户名称
	customerName := phoneCountryCode + phone
	companySize := "small"
	contactPhone := phoneCountryCode + phone

	// 创建默认客户
	customer := &models.Customer{
		CustomerName:  customerName,  // 客户名称默认为手机号
		CustomerType:  "individual",  // 客户类型默认为个人客户
		ContactPerson: contactPhone,  // 联系人默认为手机号
		CustomerLevel: "basic",       // 等级默认为基础客户
		Status:        "active",      // 状态为激活
		CompanySize:   &companySize,  // 公司规模为小型企业
		CreatedBy:     "system",      // 系统自动创建
		Phone:         &contactPhone, // 设置联系电话
	}

	// 委托给Repository层进行数据创建
	if err := s.customerRepo.CreateCustomer(ctx, customer); err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	return customer, nil
}
