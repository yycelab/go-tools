package tx

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"unsafe"

	"gorm.io/gorm"
	"yycelab.com/go-tools/errors"
)

var (
	db_session = sessionkey{}
)

type sessionkey struct{}

type CommitOrRollbackCall func(commit bool)

type EndTx func(call ...CommitOrRollbackCall)

func PanicErr(result any) {
	var err error
	switch t := result.(type) {
	case SessionContext:
		err = t.Err()
	case *gorm.DB:
		err = t.Error
	case error:
		err = t
	default:
	}
	if err != nil {
		panic(err)
	}
}

type SessionContext interface {
	context.Context
	AppendErr(err error)
	AddIfEmpty(err error)
	HasError() bool
	Errors() []error
	Err() error
	//上下文中如果存在err,panic
	PanicErr()
	//上下文中如果存在err,panic(err)
	PanicNewError(err error)
	LastInsertId() int
	WithLogger(logger *log.Logger, flags ...int) *log.Logger
}

var _ SessionContext = &errorContext{}

type errorContext struct {
	context context.Context
	errs    []error
}

func (ec *errorContext) PanicNewError(err error) {
	if ec.HasError() {
		panic(err)
	}
}
func (ec *errorContext) PanicErr() {
	if ec.HasError() {
		panic(ec.Err())
	}
}
func (ec *errorContext) LastInsertId() int {
	tx := ec.Value(db_session)
	db := tx.(*txdb)
	var id int
	ec.AddIfEmpty(db.db.Raw("SELECT LAST_INSERT_ID() AS ID").Scan(&id).Error)
	return id
}
func (ec *errorContext) WithLogger(logger *log.Logger, flags ...int) *log.Logger {
	v := ec.context.Value(db_session)
	id := v.(*txdb).sessionId
	prefix := fmt.Sprintf("[gtm(%+v)] ", id)
	var f int
	if len(flags) == 0 {
		if logger != nil {
			f = log.Flags()
		} else {
			f = log.Ldate | log.Ltime | log.Llongfile
		}
	} else {
		for i := range flags {
			f = f | flags[i]
		}
	}
	if logger == nil {
		return log.New(log.Writer(), prefix, f)
	}
	prefix = fmt.Sprintf("%s%s", logger.Prefix(), prefix)
	return log.New(logger.Writer(), prefix, logger.Flags())
}

func (ec *errorContext) Value(key any) any {
	return ec.context.Value(key)
}

func (ec *errorContext) Deadline() (deadline time.Time, ok bool) {

	return ec.context.Deadline()
}
func (ec *errorContext) Done() <-chan struct{} {
	return ec.context.Done()
}

func (ec *errorContext) HasError() bool {
	return len(ec.errs) > 0
}

func (ec *errorContext) Errors() []error {
	return ec.errs
}

func (ec *errorContext) Err() error {
	if len(ec.errs) > 0 {
		return ec.errs[0]
	}
	return nil
}

func (ec *errorContext) AddIfEmpty(err error) {
	if err != nil && len(ec.errs) == 0 {
		ec.errs = append(ec.errs, err)
	}
}

func (ec *errorContext) AppendErr(err error) {
	if err != nil {
		ec.errs = append(ec.errs, err)
	}
}

// 解决跨多表的事务 ,service层开启
type TransactionManager interface {
	//获取一个链接的DB,和会话id(可打印)
	Session(ctx context.Context) (db *gorm.DB, session unsafe.Pointer)
	//事务包裹执行,如果当前已有事务,加入事务中执行 ,该方法使用场景,在dao层.数据访问
	Tx(ctx context.Context, exec func(db *gorm.DB, sc SessionContext), options ...*sql.TxOptions)
	TransactionSupport
}

type TransactionSupport interface {
	//开启一个事务,返回一个可追加error的context .多个执行时,可以先判断是有错误再执行.
	//返回的EndTx钩子.请在begin后 使用defer 调用,可以传入callback针对不同的提交结果做一些事
	//请确保在callback中不出现panic. 该panic不影响本次事务的执行结果
	Begin(ctx context.Context, options ...*sql.TxOptions) (ec SessionContext, end EndTx)
}

func GormTxManager(db *gorm.DB) TransactionManager {
	return &transactionManager{db: db}
}

func GormTxDebugManager(db *gorm.DB, logger *log.Logger) TransactionManager {
	return &transactionManager{db: db, logger: logger}
}

type txdb struct {
	sessionId unsafe.Pointer
	db        *gorm.DB
}

type transactionManager struct {
	db     *gorm.DB
	logger *log.Logger
}

// 请传入真确的context ,如果想在事务中执行;请使用 Begin返回的ErrorContext作为参数
func (tm *transactionManager) Session(ctx context.Context) (*gorm.DB, unsafe.Pointer) {
	v := ctx.Value(db_session)
	if v != nil {
		t := v.(*txdb)
		return t.db, t.sessionId
	}
	return tm.db, unsafe.Pointer(tm.db)
}

func (tm *transactionManager) log(session unsafe.Pointer) (*log.Logger, bool) {
	if tm.logger != nil {
		logger := log.New(
			tm.logger.Writer(),
			fmt.Sprintf("%s[gtm(%+v)] ", tm.logger.Prefix(), session),
			tm.logger.Flags(),
		)
		return logger, true
	}
	return nil, false
}

func (tm *transactionManager) Tx(ctx context.Context, exec func(db *gorm.DB, ec SessionContext), options ...*sql.TxOptions) {
	db, ec, end, inTx := tm.doBegin(ctx, options...)
	logger, debug := tm.log(db.sessionId)
	if inTx {
		if debug {
			logger.Printf("join tansaction before exec")
		}
		exec(db.db, ec)
		if debug {
			logger.Printf("join tansaction after exec")
		}
		return
	}
	defer end(func(commit bool) {
		if debug {
			if commit {
				logger.Println("commit exec ")
			} else {
				logger.Printf("rollback exec ,err:%s ", ec.Err().Error())
			}
		}
	})
	if debug {
		logger.Printf("new tansaction before exec")
	}
	exec(db.db, ec)
	if debug {
		logger.Printf("new tansaction after exec")
	}
}

// 执行完,记得检查error ,ErrorContext 包含出错信息
func (tm *transactionManager) Begin(ctx context.Context, options ...*sql.TxOptions) (ec SessionContext, end EndTx) {
	_, ec, end, _ = tm.doBegin(ctx, options...)
	return
}

func (tm *transactionManager) doBegin(ctx context.Context, options ...*sql.TxOptions) (db *txdb, ec SessionContext, end EndTx, inTx bool) {
	v := ctx.Value(db_session)
	inTx = v != nil
	if v != nil {
		db = v.(*txdb)
		ec = ctx.(SessionContext)
	} else {
		c := ctx
		if c == nil {
			c = context.Background()
		}
		s := &gorm.Session{Context: c}
		gdb := tm.db.Session(s).Begin(options...)
		db = &txdb{db: gdb, sessionId: unsafe.Pointer(s)}

		val := context.WithValue(c, db_session, db)
		ec = &errorContext{
			context: val,
			errs:    make([]error, 0, 2),
		}
	}
	logger, debug := tm.log(db.sessionId)
	end = func(callbacks ...CommitOrRollbackCall) {
		err := recover()
		result := errors.RecoverError(err)
		if result.HasError {
			ec.AddIfEmpty(result.Err)
		}
		commit := tm.end(ec)
		for i := range callbacks {
			callbacks[i](commit)
		}
	}
	if debug {
		logger.Printf("try begin new transaction,opened:%t,occrur err:%t", db != nil, ec == nil || ec.HasError())
	}
	return
}

func (tm *transactionManager) end(ctx SessionContext) (commit bool) {
	var tx *txdb
	commit = !ctx.HasError()
	err := ctx.Err()
	defer func() {
		var sid unsafe.Pointer
		if tx != nil {
			sid = tx.sessionId
		}
		logger, debug := tm.log(sid)
		if debug {
			hasErr := err != ctx.Err()
			if commit {
				logger.Printf("try commit transaction,occur err:%t", hasErr)
			} else {
				logger.Printf("try rollback transaction cause:%s,occur err:%t,", err.Error(), hasErr)
			}
		}
	}()
	v := ctx.Value(db_session)
	if v == nil {
		ctx.AddIfEmpty(&errors.ResourceError{Msg: "无效的事务", Kind: errors.RESOURCE_DATABASE})
		return false
	}
	tx = v.(*txdb)
	if commit {
		ctx.AddIfEmpty(tx.db.Commit().Error)
	} else {
		ctx.AddIfEmpty(tx.db.Rollback().Error)
	}
	PanicErr(err)
	return
}

func Begin(ctx context.Context, tm TransactionSupport, opts ...*sql.TxOptions) (ec SessionContext, end EndTx) {
	return tm.Begin(ctx, opts...)
}
