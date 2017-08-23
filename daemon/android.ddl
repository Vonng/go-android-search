-----------------------------------------------------------
-- Android 母表DDL
-----------------------------------------------------------
CREATE TABLE android (
  source       TEXT NOT NULL DEFAULT 'none',
  id           TEXT, --标识
  name         TEXT, --名称
  url          TEXT, --页面
  icon         TEXT, --图标
  link         TEXT, --下载
  version      TEXT, --版本
  vendor       TEXT, --厂商
  genre        TEXT, --分类
  tags         TEXT [], --标签
  categories   TEXT [], --类目
  price        BIGINT, --价格
  system       TEXT, --系统
  platform     TEXT [], --支持的平台
  permissions  TEXT [], --所需权限
  size         BIGINT, --大小
  rating       BIGINT, --评分均值
  install_cnt  BIGINT, -- 安装/下载人数
  comment_cnt  BIGINT, --评论数
  appkey       TEXT, -- 友盟分配的appkey, 留空
  app_id       BIGINT, --平台分配的应用ID
  apk_code     BIGINT, --平台分配的Apk代码
  subtitle     TEXT, --副标题
  commentary   TEXT, --编辑评论
  description  TEXT, --应用描述,带有换行符
  reviews      JSONB, --客户评论,JSON数组,每项为三元组`(YYYY-MM-DD,user,content)`
  news         JSONB, --新闻技巧与攻略，JSON数组,每项为`(title,src,vendor)`
  extra        JSONB, -- 额外信息
  screenshots  TEXT [], --截图列表
  related_apps TEXT [], --推荐的相关应用
  sibling_apps TEXT [], --同一开发者的其他应用，豌豆荚无
  release_note TEXT, --最近更新日志,带有换行符
  release_time TIMESTAMPTZ, --最近更新时间
  crawled_time   TIMESTAMPTZ   DEFAULT CURRENT_TIMESTAMP --最近爬取时间
);



COMMENT ON TABLE android IS '安卓应用数据表(总合表)';
COMMENT ON COLUMN android.source IS '标记数据来源,wdj/sjqq/mi';
COMMENT ON COLUMN android.id IS '标识，即APK,PkgName';
COMMENT ON COLUMN android.url IS '页面URL';
COMMENT ON COLUMN android.name IS '名称';
COMMENT ON COLUMN android.icon IS '图标';
COMMENT ON COLUMN android.link IS '下载';
COMMENT ON COLUMN android.version IS '版本';
COMMENT ON COLUMN android.vendor IS '厂商';
COMMENT ON COLUMN android.genre IS '分类';
COMMENT ON COLUMN android.categories IS '类目(数组)';
COMMENT ON COLUMN android.tags IS '标签(数组)';
COMMENT ON COLUMN android.price IS '售价';
COMMENT ON COLUMN android.system IS '系统要求';
COMMENT ON COLUMN android.platform IS '支持设备';
COMMENT ON COLUMN android.permissions IS '所需权限,数组';
COMMENT ON COLUMN android.size IS '大小';
COMMENT ON COLUMN android.rating IS '评分';
COMMENT ON COLUMN android.install_cnt IS '安装数';
COMMENT ON COLUMN android.comment_cnt IS '评论数';
COMMENT ON COLUMN android.appkey IS '友盟分配的AppKey';
COMMENT ON COLUMN android.app_id IS '平台分配的应用ID,豌豆荚无';
COMMENT ON COLUMN android.apk_code IS '平台分配的Apk代码,豌豆荚无';
COMMENT ON COLUMN android.subtitle IS '副标题';
COMMENT ON COLUMN android.commentary IS '编辑评论';
COMMENT ON COLUMN android.description IS '应用描述,带有换行符';
COMMENT ON COLUMN android.reviews IS '客户评论,JSON数组,每项为三元组(YYYY-MM-DD,user,content)';
COMMENT ON COLUMN android.news IS '新闻技巧与攻略，JSON数组,每项为`(title,src,vendor)';
COMMENT ON COLUMN android.extra IS '额外扩展用字段';
COMMENT ON COLUMN android.screenshots IS '截图列表';
COMMENT ON COLUMN android.related_apps IS '推荐的相关应用';
COMMENT ON COLUMN android.sibling_apps IS '同一开发者的其他应用，豌豆荚无';
COMMENT ON COLUMN android.release_note IS '最近更新日志,带有换行符';
COMMENT ON COLUMN android.release_time IS '最近更新时间';
COMMENT ON COLUMN android.crawled_time IS '最近爬取时间';
-----------------------------------------------------------


-----------------------------------------------------------
-- wdj DDL 豌豆荚
-----------------------------------------------------------
CREATE TABLE wdj (
  PRIMARY KEY (id)
)
  INHERITS (android);

COMMENT ON TABLE wdj IS '豌豆荚应用数据表';
COMMENT ON COLUMN wdj.source IS '标记数据来源,固定为`wdj`';
COMMENT ON COLUMN wdj.id IS 'APK,PkgName';
COMMENT ON COLUMN wdj.url IS '页面URL';
COMMENT ON COLUMN wdj.name IS '名称';
COMMENT ON COLUMN wdj.icon IS '图标';
COMMENT ON COLUMN wdj.link IS '下载';
COMMENT ON COLUMN wdj.version IS '版本';
COMMENT ON COLUMN wdj.vendor IS '厂商';
COMMENT ON COLUMN wdj.genre IS '分类';
COMMENT ON COLUMN wdj.categories IS '类目(数组)';
COMMENT ON COLUMN wdj.tags IS '标签(数组)';
COMMENT ON COLUMN wdj.price IS '售价，豌豆荚无';
COMMENT ON COLUMN wdj.system IS '系统要求(安卓版本号)';
COMMENT ON COLUMN wdj.platform IS '支持设备，豌豆荚无';
COMMENT ON COLUMN wdj.permissions IS '所需权限,数组';
COMMENT ON COLUMN wdj.size IS '大小';
COMMENT ON COLUMN wdj.rating IS '评分';
COMMENT ON COLUMN wdj.install_cnt IS '安装数';
COMMENT ON COLUMN wdj.comment_cnt IS '评论数';
COMMENT ON COLUMN wdj.appkey IS '友盟分配的AppKey';
COMMENT ON COLUMN wdj.app_id IS '平台分配的应用ID,豌豆荚无';
COMMENT ON COLUMN wdj.apk_code IS '平台分配的Apk代码,豌豆荚无';
COMMENT ON COLUMN wdj.subtitle IS '副标题';
COMMENT ON COLUMN wdj.commentary IS '编辑评论';
COMMENT ON COLUMN wdj.description IS '应用描述,带有换行符';
COMMENT ON COLUMN wdj.reviews IS '客户评论,JSON数组,每项为三元组(YYYY-MM-DD,user,content)';
COMMENT ON COLUMN wdj.news IS '新闻技巧与攻略，JSON数组,每项为`(title,src,vendor)';
COMMENT ON COLUMN wdj.extra IS '额外扩展用字段';
COMMENT ON COLUMN wdj.screenshots IS '截图列表';
COMMENT ON COLUMN wdj.related_apps IS '推荐的相关应用';
COMMENT ON COLUMN wdj.sibling_apps IS '同一开发者的其他应用，豌豆荚无';
COMMENT ON COLUMN wdj.release_note IS '最近更新日志,带有换行符';
COMMENT ON COLUMN wdj.release_time IS '最近更新时间';
COMMENT ON COLUMN wdj.crawled_time IS '最近爬取时间';
-----------------------------------------------------------


-----------------------------------------------------------
-- sjqq DDL 应用宝
-----------------------------------------------------------
CREATE TABLE sjqq (
  PRIMARY KEY (id)
)
  INHERITS (android);

COMMENT ON TABLE sjqq IS '应用宝应用数据表';
COMMENT ON COLUMN sjqq.source IS '标记数据来源，固定为`sjqq`';
COMMENT ON COLUMN sjqq.id IS '标识，即APK,PkgName';
COMMENT ON COLUMN sjqq.url IS '页面';
COMMENT ON COLUMN sjqq.name IS '名称';
COMMENT ON COLUMN sjqq.icon IS '图标';
COMMENT ON COLUMN sjqq.link IS '下载';
COMMENT ON COLUMN sjqq.version IS '版本';
COMMENT ON COLUMN sjqq.vendor IS '厂商';
COMMENT ON COLUMN sjqq.genre IS '分类';
COMMENT ON COLUMN sjqq.categories IS '类目(数组)';
COMMENT ON COLUMN sjqq.tags IS '标签(数组)';
COMMENT ON COLUMN sjqq.price IS '售价';
COMMENT ON COLUMN sjqq.system IS '系统要求';
COMMENT ON COLUMN sjqq.platform IS '支持设备';
COMMENT ON COLUMN sjqq.permissions IS '所需权限,数组';
COMMENT ON COLUMN sjqq.size IS '大小';
COMMENT ON COLUMN sjqq.rating IS '评分';
COMMENT ON COLUMN sjqq.install_cnt IS '安装数';
COMMENT ON COLUMN sjqq.comment_cnt IS '评论数';
COMMENT ON COLUMN sjqq.appkey IS '友盟分配的AppKey';
COMMENT ON COLUMN sjqq.app_id IS '平台分配的应用ID,豌豆荚无';
COMMENT ON COLUMN sjqq.apk_code IS '平台分配的Apk代码,豌豆荚无';
COMMENT ON COLUMN sjqq.subtitle IS '副标题';
COMMENT ON COLUMN sjqq.commentary IS '编辑评论';
COMMENT ON COLUMN sjqq.description IS '应用描述,带有换行符';
COMMENT ON COLUMN sjqq.reviews IS '客户评论,JSON数组,每项为三元组(YYYY-MM-DD,user,content)';
COMMENT ON COLUMN sjqq.news IS '新闻技巧与攻略，JSON数组,每项为`(title,src,vendor)';
COMMENT ON COLUMN sjqq.extra IS '额外扩展用字段';
COMMENT ON COLUMN sjqq.screenshots IS '截图列表';
COMMENT ON COLUMN sjqq.related_apps IS '推荐的相关应用';
COMMENT ON COLUMN sjqq.sibling_apps IS '同一开发者的其他应用，豌豆荚无';
COMMENT ON COLUMN sjqq.release_note IS '最近更新日志,带有换行符';
COMMENT ON COLUMN sjqq.release_time IS '最近更新时间';
COMMENT ON COLUMN sjqq.crawled_time IS '最近爬取时间';
-----------------------------------------------------------


---------------------------------------------------------------
-- Task Queue
---------------------------------------------------------------
-- DROP TABLE android_queue;
CREATE TABLE IF NOT EXISTS android_queue (
  id TEXT PRIMARY KEY
);
COMMENT ON TABLE android_queue IS 'Apple Task Queue';
-----------------------------------------
-- Function: add android id to queue
CREATE OR REPLACE FUNCTION android_apk(_apk TEXT)
  RETURNS VOID AS
$$BEGIN INSERT INTO android_queue (id) VALUES ('!' || _apk);
END;$$
LANGUAGE plpgsql VOLATILE;
COMMENT ON FUNCTION android_apk(BIGINT) IS '向安卓队列中添加apk任务';
-- SELECT android_aid(1031569344)
-----------------------------------------
-- Function: add search keyword to queue
CREATE OR REPLACE FUNCTION android_key(keyword TEXT)
  RETURNS VOID AS
$$BEGIN INSERT INTO android_queue (id) VALUES ('#' || keyword);
END;$$
LANGUAGE plpgsql VOLATILE;
COMMENT ON FUNCTION android_key(TEXT) IS '向安卓队列中添加关键词任务';
-- SELECT android_key('蛤蛤');
-----------------------------------------