#### gin_template.admin_menu 
管理员菜单栏表

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | id | 主键 | int(11) unsigned | PRI | NO | auto_increment |  |
| 2 | admin_id | 管理员ID | int(11) unsigned | MUL | NO |  | 0 |
| 3 | menu_id | 菜单栏ID | int(11) unsigned |  | NO |  | 0 |
| 4 | created_at | 创建时间 | timestamp |  | NO |  | current_timestamp() |
| 5 | created_user | 创建人 | varchar(60) |  | NO |  | '' |
