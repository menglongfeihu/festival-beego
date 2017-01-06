package consts

const (
	//用户日志类型 0点亮ip, 1集福
	USER_LOG_IP    = 0
	USER_LOG_FAVOR = 1

	// 奖品类型
	PRODUCT_TYPE_VIP    = 0 // vip 黄金会员
	PRODUCT_TYPE_ENTITY = 1 // 实物奖品

	//用户状态：0 ：未参与 1：点亮ip  2：集福 -1：拉黑
	USER_STATUS_NOT_INCLUDE = 0
	USER_STATUS_IPING       = 1
	USER_STATUS_FAVORING    = 2
	USER_STATUS_BLACKLIST   = -1
	//领取奖品状态 0 已领取 1未领取 2已过期
	AWARD_STATUS_DRAW   = 0
	AWARD_STATUS_UNDRAW = 1
	AWARD_STATUS_EXPIRE = 2
	// VIP类型
	VIP_TYPE_3  = 3
	VIP_TYPE_7  = 7
	VIP_TYPE_14 = 14

	//活动黄金会员集福数量
	FAVOR_VIP_3  = 3
	FAVOR_VIP_7  = 5
	FAVOR_VIP_14 = 10

	//点亮ip数量
	COUNT_IP = 6

	//缓存时间
	EXPIRE_ONE_SECOND = 1
	EXPIRE_ONE_MINUTE = 60
	EXPIRE_ONE_HOUR   = 60 * 60
	EXPIRE_ONE_DAY    = 24 * 60 * 60
	EXPIRE_ONE_WEEK   = 7 * 24 * 60 * 60

	DAY_LIMIT = 1000
	IP_LIMIT  = 100
	//缓存key
	CACHE_ACCOUNT_FAVOR_LIMIT_DAY          = "favor_limit_day_%s"
	CACHE_ACCOUNT_FAVOR_LIMIT_IP           = "favor_limit_ip_%s_%s"
	CACHE_ACCOUNT_BEFAVOR_LIMIT_DAY        = "befavor_limit_day_%s"
	CACHE_ACCOUNT_BEFAVOR_LIMIT_IP         = "befavor_limit_ip_%s_%s"
	CACHE_ACCOUNT_INFO                     = "account_info_%s"
	CACHE_ACCOUNT_PHONE                    = "account_phone_%s"
	CACHE_ACCOUNT_LOGIN                    = "account_login_%s"
	CACHE_PRODUCT_LIST                     = "product_list%s"
	CACHE_PRODUCT_ID                       = "product_id_%d"
	CACHE_VIDEOINFO_PASSPORT_VNAME         = "videoinfo_%s_%s"
	CACHE_VIDEOINFO_LIST_PASSPORT          = "videoinfo_list_%s"
	CACHE_USERLOG_PASSPORT_TYPE_DAY        = "userlog_%s_%s_%d"
	CACHE_USERLOG_FAVOR_PASSPORT_FPASSPORT = "userlog_favor_%s_%s"
	CACHE_USERLOG_LIST_PASSPORT_TYPE       = "userlog_list_%s_%d"
	CACHE_USERLOG_COUNT_PASSPORT_TYPE      = "userlog_count_%s_%d"
	CACHE_USERAWARD_PASSPORT_PROID         = "useraward_%s_%d"
	CACHE_USERAWARD_LIST_PASSPORT          = "useraward_list_%s"
	CACHE_USER_ID                          = "user_id_%d"
	CACHE_USER_PASSPORT                    = "user_passport_%s"
	CACHE_RANKAWARDMAP_RANK                = "randawardmap_%s"
)
