#### gin_template.cron_task 
后台任务表

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | id | 主键 | int(11) unsigned | PRI | NO | auto_increment |  |
| 2 | name | 任务名称 | varchar(64) | MUL | NO |  | '' |
| 3 | spec | crontab 表达式 | varchar(64) |  | NO |  | '' |
| 4 | command | 执行命令 | varchar(255) |  | NO |  | '' |
| 5 | protocol | 执行方式 1:shell 2:http | tinyint(1) unsigned |  | NO |  | 1 |
| 6 | http_method | http 请求方式 1:get 2:post | tinyint(1) unsigned |  | NO |  | 1 |
| 7 | timeout | 超时时间(单位:秒) | int(11) unsigned |  | NO |  | 60 |
| 8 | retry_times | 重试次数 | tinyint(1) |  | NO |  | 3 |
| 9 | retry_interval | 重试间隔(单位:秒) | int(11) |  | NO |  | 60 |
| 10 | notify_status | 执行结束是否通知 1:不通知 2:失败通知 3:结束通知 4:结果关键字匹配通知 | tinyint(1) unsigned |  | NO |  | 0 |
| 11 | notify_type | 通知类型 1:邮件 2:webhook | tinyint(1) unsigned |  | NO |  | 0 |
| 12 | notify_receiver_email | 通知者邮箱地址(多个用,分割) | varchar(255) |  | NO |  | '' |
| 13 | notify_keyword | 通知匹配关键字(多个用,分割) | varchar(255) |  | NO |  | '' |
| 14 | remark | 备注 | varchar(100) |  | NO |  | '' |
| 15 | is_used | 是否启用 1:是  -1:否 | tinyint(1) |  | NO |  | 1 |
| 16 | created_at | 创建时间 | timestamp |  | NO |  | current_timestamp() |
| 17 | created_user | 创建人 | varchar(60) |  | NO |  | '' |
| 18 | updated_at | 更新时间 | timestamp |  | NO | on update current_timestamp() | current_timestamp() |
| 19 | updated_user | 更新人 | varchar(60) |  | NO |  | '' |
