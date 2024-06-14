CREATE TABLE IF NOT EXISTS `user`
(
    `id`            bigint primary key comment "用户id",
    `user_name`     varchar(256)               NOT NULL COMMENT "用户名不允许重复",
    `avatar_url`    varchar(1024)              NULL COMMENT "用户头像",
    `user_password` varchar(256)               NOT NULL COMMENT "用户密码",
    `create_time`   varchar(256)               NOT NULL COMMENT "创建时间",
    `update_time`   varchar(256)               NOT NULL NULL COMMENT "更新时间",
    `is_delete`     tinyint     default 0      NOT NULL COMMENT "是否删除,0为不删除，1为删除",
    `user_role`     varchar(64) default 'user' NOT NULL COMMENT "用户类型，有user,admin,ban"
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC COMMENT "用户表";


CREATE TABLE IF NOT EXISTS `question`
(
    `id`           bigint comment "id" primary key,
    `title`        varchar(512)      null comment "标题",
    `content`      text              null comment "内容",
    `tag`         varchar(1024)     null comment "标签列表json数组",
    `answer`       text              null comment "题目答案",
    `submit_num`   int     default 0 not null comment "题目提交数",
    `accept_num`   int     default 0 not null comment "题目通过数",
    `judge_case`   text              null comment "判题用例json数组",
    `judge_config` text              null comment "判题配置json对象",
    `user_id`      bigint            not null comment "创建用户id",
    `create_time`  varchar(256)      NOT NULL COMMENT "创建时间",
    `update_time`  varchar(256)      NOT NULL NULL COMMENT "更新时间",
    `is_delete`    tinyint default 0 not null comment "是否删除",
    index idx_userId (user_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC COMMENT "题目信息表";


CREATE TABLE IF NOT EXISTS `submit`
(
    `id`          bigint comment "id" primary key,
    `language`    varchar(128)      not null comment "编程语言",
    `code`        text              not null comment "用户代码",
    `judge_result`  text              null comment "判题信息json对象,包含判题人（系统或ai或第三方判题系统）",
    `status`      varchar(256)      not null comment "判题状态（待判题,判题中,成功,失败)",
    `question_id` bigint            not null comment "判题id",
    `user_id`     bigint            not null comment "创建用户id",
    `create_time` varchar(256)      NOT NULL COMMENT "创建时间",
    `update_time` varchar(256)      NOT NULL NULL COMMENT "更新时间",
    `is_delete`   tinyint default 0 not null comment "是否删除",
    index idx_question_id (question_id),
    index idx_user_id (user_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC COMMENT "题目提交信息表";
