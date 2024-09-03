// Package gorm
//
//	@program:	gin-template
//	@author:	[lliuhuan](https://github.com/lliuhuan)
//	@create:	2024-07-03 10:19
package gorm

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/pkg/timeutil"
	"github.com/LLiuHuan/gin-template/pkg/trace"

	"gorm.io/gorm"
	"gorm.io/gorm/utils"
)

const (
	callBackBeforeName = "core:before"
	callBackAfterName  = "core:after"
	startTime          = "_start_time"
)

type TracePlugin struct{}

func (op *TracePlugin) Name() string {
	return "tracePlugin"
}

func (op *TracePlugin) Initialize(db *gorm.DB) (err error) {
	// 查询没有数据，屏蔽 gorm v2 包中会爆出的错误
	// https://github.com/go-gorm/gorm/issues/3789  此 issue 所反映的问题就是我们本次解决掉的
	_ = db.Callback().Query().Before("gorm:query").Register("disable_raise_record_not_found", maskNotDataError)

	// 开始前 - 并不是都用相同的方法，可以自己自定义
	_ = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	_ = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	_ = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)
	// https://github.com/go-gorm/gorm/issues/4838
	_ = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, createBeforeHook)
	// 为了完美支持gorm的一系列回调函数
	_ = db.Callback().Update().Before("gorm:before_update").Register(callBackBeforeName, updateBeforeHook)

	// 结束后 - 并不是都用相同的方法，可以自己自定义
	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return
}

var _ gorm.Plugin = &TracePlugin{}

// before 在执行前，记录开始时间
func before(db *gorm.DB) {
	db.InstanceSet(startTime, time.Now())
	return
}

// after 在执行后，记录结束时间，以及执行的 SQL 语句，并将其保存到 Trace 中
func after(db *gorm.DB) {
	_ctx := db.Statement.Context
	ctx, ok := _ctx.(core.StdContext)
	if !ok {
		return
	}

	_ts, isExist := db.InstanceGet(startTime)
	if !isExist {
		return
	}

	ts, ok := _ts.(time.Time)
	if !ok {
		return
	}

	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)

	sqlInfo := new(trace.SQL)
	sqlInfo.Timestamp = timeutil.CSTLayoutString()
	sqlInfo.SQL = sql
	sqlInfo.Stack = utils.FileWithLineNum()
	sqlInfo.Rows = db.Statement.RowsAffected
	sqlInfo.CostSeconds = time.Since(ts).Seconds()

	if ctx.Trace != nil {
		ctx.Trace.AppendSQL(sqlInfo)
	}

	return
}

// maskNotDataError 屏蔽 gorm v2 包中会爆出的错误，正常没有数据会抛出异常，但是感觉不合理，没什么必要
func maskNotDataError(db *gorm.DB) {
	db.Statement.RaiseErrorOnNotFound = false
}

// createBeforeHook InterceptCreatePramsNotPtrError 拦截 create 函数参数如果是非指针类型的错误,新用户最容犯此错误
func createBeforeHook(db *gorm.DB) {
	if reflect.TypeOf(db.Statement.Dest).Kind() != reflect.Ptr {
		//db.Error = errors.New("gorm Create 函数的参数必须是一个指针")
		fmt.Println("gorm Create 函数的参数必须是一个指针")
	} else {
		destValueOf := reflect.ValueOf(db.Statement.Dest).Elem()
		if destValueOf.Type().Kind() == reflect.Slice || destValueOf.Type().Kind() == reflect.Array {
			inLen := destValueOf.Len()
			for i := 0; i < inLen; i++ {
				row := destValueOf.Index(i)
				if row.Type().Kind() == reflect.Struct {
					if b, column := structHasSpecialField("CreatedAt", row); b {
						destValueOf.Index(i).FieldByName(column).Set(reflect.ValueOf(time.Now().Format(timeutil.CSTLayout)))
					}
					if b, column := structHasSpecialField("UpdatedAt", row); b {
						destValueOf.Index(i).FieldByName(column).Set(reflect.ValueOf(time.Now().Format(timeutil.CSTLayout)))
					}

				} else if row.Type().Kind() == reflect.Map {
					if b, column := structHasSpecialField("created_at", row); b {
						row.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(timeutil.CSTLayout)))
					}
					if b, column := structHasSpecialField("updated_at", row); b {
						row.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(timeutil.CSTLayout)))
					}
				}
			}
		} else if destValueOf.Type().Kind() == reflect.Struct {
			//  if destValueOf.Type().Kind() == reflect.Struct
			// 参数校验无错误自动设置 CreatedAt、 UpdatedAt
			if b, column := structHasSpecialField("CreatedAt", db.Statement.Dest); b {
				db.Statement.SetColumn(column, time.Now().Format(timeutil.CSTLayout))
			}
			if b, column := structHasSpecialField("UpdatedAt", db.Statement.Dest); b {
				db.Statement.SetColumn(column, time.Now().Format(timeutil.CSTLayout))
			}
		} else if destValueOf.Type().Kind() == reflect.Map {
			if b, column := structHasSpecialField("created_at", db.Statement.Dest); b {
				destValueOf.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(timeutil.CSTLayout)))
			}
			if b, column := structHasSpecialField("updated_at", db.Statement.Dest); b {
				destValueOf.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(timeutil.CSTLayout)))
			}
		}
	}
}

// updateBeforeHook
// InterceptUpdatePramsNotPtrError 拦截 save、update 函数参数如果是非指针类型的错误
// 对于开发者来说，以结构体形式更新数，只需要在 update 、save 函数的参数前面添加 & 即可
// 最终就可以完美兼支持、兼容 gorm 的所有回调函数
// 但是如果是指定字段更新，例如： UpdateColumn 函数则只传递值即可，不需要做校验
func updateBeforeHook(db *gorm.DB) {
	before(db)
	if reflect.TypeOf(db.Statement.Dest).Kind() == reflect.Struct {
		//_ = db.AddError(errors.New(my_errors.ErrorsGormDBUpdateParamsNotPtr))
		//variable.ZapLog.Warn(my_errors.ErrorsGormDBUpdateParamsNotPtr)
	} else if reflect.TypeOf(db.Statement.Dest).Kind() == reflect.Map {
		// 如果是调用了 gorm.Update 、updates 函数 , 在参数没有传递指针的情况下，无法触发回调函数

	} else if reflect.TypeOf(db.Statement.Dest).Kind() == reflect.Ptr && reflect.ValueOf(db.Statement.Dest).Elem().Kind() == reflect.Struct {
		// 参数校验无错误自动设置 UpdatedAt
		if b, column := structHasSpecialField("UpdatedAt", db.Statement.Dest); b {
			db.Statement.SetColumn(column, time.Now().Format(timeutil.CSTLayout))
		}
	} else if reflect.TypeOf(db.Statement.Dest).Kind() == reflect.Ptr && reflect.ValueOf(db.Statement.Dest).Elem().Kind() == reflect.Map {
		if b, column := structHasSpecialField("updated_at", db.Statement.Dest); b {
			destValueOf := reflect.ValueOf(db.Statement.Dest).Elem()
			destValueOf.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(timeutil.CSTLayout)))
		}
	}
}

// structHasSpecialField  检查结构体是否有特定字段
func structHasSpecialField(fieldName string, anyStructPtr interface{}) (bool, string) {
	var tmp reflect.Type
	if reflect.TypeOf(anyStructPtr).Kind() == reflect.Ptr && reflect.ValueOf(anyStructPtr).Elem().Kind() == reflect.Map {
		destValueOf := reflect.ValueOf(anyStructPtr).Elem()
		for _, item := range destValueOf.MapKeys() {
			if item.String() == fieldName {
				return true, fieldName
			}
		}
	} else if reflect.TypeOf(anyStructPtr).Kind() == reflect.Ptr && reflect.ValueOf(anyStructPtr).Elem().Kind() == reflect.Struct {
		destValueOf := reflect.ValueOf(anyStructPtr).Elem()
		tf := destValueOf.Type()
		for i := 0; i < tf.NumField(); i++ {
			if !tf.Field(i).Anonymous && tf.Field(i).Type.Kind() != reflect.Struct {
				if tf.Field(i).Name == fieldName {
					return true, getColumnNameFromGormTag(fieldName, tf.Field(i).Tag.Get("gorm"))
				}
			} else if tf.Field(i).Type.Kind() == reflect.Struct {
				tmp = tf.Field(i).Type
				for j := 0; j < tmp.NumField(); j++ {
					if tmp.Field(j).Name == fieldName {
						return true, getColumnNameFromGormTag(fieldName, tmp.Field(j).Tag.Get("gorm"))
					}
				}
			}
		}
	} else if reflect.Indirect(anyStructPtr.(reflect.Value)).Type().Kind() == reflect.Struct {
		// 处理结构体
		destValueOf := anyStructPtr.(reflect.Value)
		tf := destValueOf.Type()
		for i := 0; i < tf.NumField(); i++ {
			if !tf.Field(i).Anonymous && tf.Field(i).Type.Kind() != reflect.Struct {
				if tf.Field(i).Name == fieldName {
					return true, getColumnNameFromGormTag(fieldName, tf.Field(i).Tag.Get("gorm"))
				}
			} else if tf.Field(i).Type.Kind() == reflect.Struct {
				tmp = tf.Field(i).Type
				for j := 0; j < tmp.NumField(); j++ {
					if tmp.Field(j).Name == fieldName {
						return true, getColumnNameFromGormTag(fieldName, tmp.Field(j).Tag.Get("gorm"))
					}
				}
			}
		}
	} else if reflect.Indirect(anyStructPtr.(reflect.Value)).Type().Kind() == reflect.Map {
		destValueOf := anyStructPtr.(reflect.Value)
		for _, item := range destValueOf.MapKeys() {
			if item.String() == fieldName {
				return true, fieldName
			}
		}
	}
	return false, ""
}

// getColumnNameFromGormTag 从 gorm 标签中获取字段名
//
//	@defaultColumn	如果没有 gorm：column 标签为字段重命名，则使用默认字段名
//	@TagValue		字段中含有的gorm："column:created_at" 标签值，可能的格式：1. column:created_at    、2. default:null;  column:created_at  、3.  column:created_at; default:null
func getColumnNameFromGormTag(defaultColumn, TagValue string) (str string) {
	pos1 := strings.Index(TagValue, "column:")
	if pos1 == -1 {
		str = defaultColumn
		return
	} else {
		TagValue = TagValue[pos1+7:]
	}
	pos2 := strings.Index(TagValue, ";")
	if pos2 == -1 {
		str = TagValue
	} else {
		str = TagValue[:pos2]
	}
	return strings.ReplaceAll(str, " ", "")
}
