#### gin_template.menu_action 
功能权限表

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | id | 主键 | int(11) unsigned | PRI | NO | auto_increment |  |
| 2 | menu_id | 菜单栏ID | int(11) unsigned | MUL | NO |  | 0 |
| 3 | method | 请求方式 | varchar(30) |  | NO |  | '' |
| 4 | api | 请求地址 | varchar(100) |  | NO |  | '' |
| 5 | is_deleted | 是否删除 1:是  -1:否 | tinyint(1) |  | NO |  | -1 |
| 6 | created_at | 创建时间 | timestamp |  | NO |  | current_timestamp() |
| 7 | created_user | 创建人 | varchar(60) |  | NO |  | '' |
| 8 | updated_at | 更新时间 | timestamp |  | NO | on update current_timestamp() | current_timestamp() |
| 9 | updated_user | 更新人 | varchar(60) |  | NO |  | '' |
