package captcha

import (
	"github.com/lexkong/log"
	"github.com/satori/go.uuid"
	"groot/comm/constk"
	"groot/proto/cserr"
	"groot/proto/ssmsg"
	"groot/service/captcha/model"
	"groot/sfw/crypto"
	"groot/sfw/util"
	"strconv"
	"strings"
	"time"
)

type (
	ApplyCodeOptions struct {
		Data         string
		MaxLength    int
		Timeout      int
		Charset      int
		CheckTimes   int
		LimitGenTime int
	}
)

func CaptchaGetSMSCodeDefaultApplyOptions() *ApplyCodeOptions {
	return &ApplyCodeOptions{
		MaxLength:    4,
		Charset:      int(util.RAND_CHARSET_NUMBER),
		Timeout:      constk.CAPTCHA_SMS_DEFAULT_TIMEOUT_SECOND,
		CheckTimes:   5,
		LimitGenTime: constk.CAPTCHA_SMS_DEFAULT_LIMIT_TIME_SECOND,
	}
}



func VerifyCodeWithTicket(target, ticket, key string, scene constk.CaptchaSceneType) bool {
	vs := strings.Split(key, ".")
	if len(vs) != 2 {
		log.Errorf(nil, "error ticket key")
		return false
	}
	secret := vs[0]
	check_magic := vs[1]

	ctx := ssmsg.CaptchaMsgTicketCtx{}
	if crypto.VerifyMsgTicket(ticket, []byte(secret), &ctx) == false {
		log.Debugf("ticket verify or prse error ")
		return false
	}

	if ctx.Target != target || ctx.Scene != uint32(scene) || strconv.Itoa(int(ctx.CheckMagic)) != check_magic {
		log.Debugf("error app ctx not match")
		return false
	}

	return true
}

//apply for one captch , return code or error
func ApplyCode(target string, scene constk.CaptchaSceneType, options *ApplyCodeOptions) (code string,
	e cserr.ICSErrCodeError) {

	state := model.DbCaptchaState{
		Target: target,
		Scene:  uint32(scene),
	}
	timeNow := time.Now().Unix()
	db := CaptchaDb
	notExists := false
	if db.Table(DbCaptchaStateTbNameBase).First(&state, "target=? AND scene=?", target, scene).RecordNotFound() {
		//not exist ok , insert
		notExists = true
	} else {
		if state.LastGenTime+int64(options.LimitGenTime) > timeNow {
			log.Warnf("error target:%s gen code hz too high", target)
			e = cserr.ErrCaptcha
			return
		}
		//update
		notExists = false
	}

	charset := util.GetRandCharset(options.Charset)
	code = util.RandomStringWithCharset(options.MaxLength, charset)

	state.ExpiredTime = timeNow + int64(options.Timeout)
	state.LastGenTime = timeNow
	state.Code = code
	state.CheckedTimes = options.CheckTimes

	var err error
	if notExists {
		state.Uuid = uuid.NewV4().String()
		err = CaptchaDb.Table(DbCaptchaStateTbNameBase).Create(&state).Error
	} else {
		err = CaptchaDb.Table(DbCaptchaStateTbNameBase).Where("Target=? AND Scene=? ",
			target, scene).Updates(model.DbCaptchaState{
			Code:         code,
			ExpiredTime:  state.ExpiredTime,
			LastGenTime:  state.LastGenTime,
			CheckedTimes: state.CheckedTimes,
		}).Error
	}
	if err != nil {
		log.Errorf(err, "db save captcha sate fail for target:%s", target)
		e = cserr.ErrCaptcha
		return
	}

	return
}

//return data, code, or err
func CheckCodeByTargetScene(target, code string, scene constk.CaptchaSceneType) error {
	capc, err := GetCaptchaByTargetScene(target, scene)
	if !util.ErrOK(err) {
		return err
	}
	return CheckDbCaptcha(code, capc)
}

func CheckCodeByUuid(uuid, code string) error {
	capc, err := GetCaptchaByUuid(uuid)
	if !util.ErrOK(err) {
		return err
	}
	return CheckDbCaptcha(code, capc)
}

func CheckDbCaptcha(code string, capc *model.DbCaptchaState) error {
	now := time.Now()
	if now.Unix() > capc.ExpiredTime || capc.CheckedTimes <= 0 {
		log.Warnf("captcha verify timeout or check times:%d", capc.CheckedTimes)
		return cserr.ErrCaptcha
	}
	if code != capc.Code {
		dberr := CaptchaDb.Table(DbCaptchaStateTbNameBase).Where("check_times = ?",
			capc.CheckedTimes).Update("check_times",
			capc.CheckedTimes-1).Error
		if dberr != nil {
			log.Errorf(dberr, "captcha db check times error !")
		}
		log.Warnf("captcha verify fail check times:%d", capc.CheckedTimes)
		return cserr.ErrCaptcha
	}
	log.Debugf("captcha verify success target:%s scene:%d uuid:%s", capc.Target, capc.Scene, capc.Uuid)
	//real remove (no need , for will check agian too quick)
	CaptchaDb.Table(DbCaptchaStateTbNameBase).Where( "Uuid=?", capc.Uuid).Update("last_gen_time",
			capc.LastGenTime - int64(constk.CAPTCHA_SMS_DEFAULT_LIMIT_TIME_SECOND))
	return nil
}

func GetCaptchaByTargetScene(target string, scene constk.CaptchaSceneType) (*model.DbCaptchaState, error) {
	capct := model.DbCaptchaState{
		Target: target,
		Scene:  uint32(scene),
	}
	var err = CaptchaDb.Table(DbCaptchaStateTbNameBase).First(&capct, "Target=? AND Scene=?", target, scene).Error
	if err != nil {
		log.Errorf(err, "get captcha error for target:%s scene:%d", target, scene)
		return nil, cserr.ErrDb
	}
	return &capct, nil
}

func GetCaptchaByUuid(uuid string) (*model.DbCaptchaState, error) {
	capct := model.DbCaptchaState{
		Uuid: uuid,
	}
	var err = CaptchaDb.Table(DbCaptchaStateTbNameBase).First(&capct, "Uuid=?", uuid).Error
	if err != nil {
		log.Errorf(err, "get captcha error for uuid:%s", uuid)
		return nil, cserr.ErrDb
	}
	return &capct, nil
}
