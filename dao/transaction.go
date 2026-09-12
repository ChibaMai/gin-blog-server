package dao

import (
	"gin-blog/utils"
	"go.uber.org/zap"

	"gorm.io/gorm"
)

// 事务执行函数类型
type TxFunc func(*gorm.DB) error

// 在事务中执行操作，支持自动回滚
// 用法示例:
// err := dao.Transaction(func(tx *gorm.DB) error {
//     // 执行多个数据库操作
//     if err := tx.Create(&user).Error; err != nil {
//         return err  // 自动回滚
//     }
//     if err := tx.Create(&profile).Error; err != nil {
//         return err  // 自动回滚
//     }
//     return nil      // 自动提交
// })
func Transaction(fn TxFunc) error {
	tx := DB.BeginTx(DB.Statement.Context, nil)
	if err := fn(tx); err != nil {
		tx.Rollback()
		utils.Logger.Error("事务执行失败，已回滚",
			zap.Error(err),
		)
		return err
	}
	if err := tx.Commit().Error; err != nil {
		utils.Logger.Error("事务提交失败",
			zap.Error(err),
		)
		return err
	}
	return nil
}

// 在事务中执行操作，返回结果和错误
// 用法示例:
// result, err := dao.TransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
//     var user User
//     if err := tx.Create(&user).Error; err != nil {
//         return nil, err
//     }
//     return user.ID, nil
// })
func TransactionWithResult(fn func(*gorm.DB) (interface{}, error)) (interface{}, error) {
	tx := DB.BeginTx(DB.Statement.Context, nil)
	result, err := fn(tx)
	if err != nil {
		tx.Rollback()
		utils.Logger.Error("事务执行失败，已回滚",
			zap.Error(err),
		)
		return result, err
	}
	if err := tx.Commit().Error; err != nil {
		utils.Logger.Error("事务提交失败",
			zap.Error(err),
		)
		return result, err
	}
	return result, nil
}
