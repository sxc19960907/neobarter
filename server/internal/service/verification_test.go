package service

import (
	"testing"

	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserSvcDB(t *testing.T) (*UserService, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}))
	return NewUserService(repository.NewUserRepository(db)), db
}

func TestVerifyRealName_Success(t *testing.T) {
	svc, db := setupUserSvcDB(t)
	u := &model.User{Phone: "13800000001", UserType: "personal"}
	require.NoError(t, db.Create(u).Error)

	err := svc.VerifyRealName(u.ID, "张三", "11010119900307721X")
	require.NoError(t, err)

	var after model.User
	db.First(&after, u.ID)
	assert.True(t, after.RealNameVerified)
	assert.Equal(t, "张三", after.RealName)
}

func TestVerifyRealName_InvalidIDCard(t *testing.T) {
	svc, db := setupUserSvcDB(t)
	u := &model.User{Phone: "13800000001", UserType: "personal"}
	require.NoError(t, db.Create(u).Error)

	err := svc.VerifyRealName(u.ID, "张三", "123") // 非法身份证
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "身份证")
}

func TestVerifyRealName_NoDuplicate(t *testing.T) {
	svc, db := setupUserSvcDB(t)
	u := &model.User{Phone: "13800000001", UserType: "personal", RealNameVerified: true}
	require.NoError(t, db.Create(u).Error)

	err := svc.VerifyRealName(u.ID, "张三", "11010119900307721X")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "重复")
}

func TestVerifyEnterprise_Success(t *testing.T) {
	svc, db := setupUserSvcDB(t)
	u := &model.User{Phone: "13800000002", UserType: "enterprise"}
	require.NoError(t, db.Create(u).Error)

	err := svc.VerifyEnterprise(u.ID, "测试公司", "https://oss/license.png")
	require.NoError(t, err)

	var after model.User
	db.First(&after, u.ID)
	assert.True(t, after.EnterpriseVerified)
	assert.Equal(t, "测试公司", after.EnterpriseName)
}

func TestVerifyEnterprise_PersonalRejected(t *testing.T) {
	svc, db := setupUserSvcDB(t)
	u := &model.User{Phone: "13800000003", UserType: "personal"} // 个人用户
	require.NoError(t, db.Create(u).Error)

	err := svc.VerifyEnterprise(u.ID, "测试公司", "https://oss/license.png")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "企业用户")
}
