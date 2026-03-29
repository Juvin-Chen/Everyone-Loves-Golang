# GORM 完整版常用标签速查表
---

## 🔴 第一类：字段类型与长度（必用）
| 标签 | 含义 | 示例 |
| ---- | ---- | ---- |
| `type` | 自定义数据库字段类型 | `type:varchar(50)`、`type:int`、`type:decimal(10,2)` |
| `size` | 简写字符串长度 | `size:100` = varchar(100) |
| `precision` | 数字/时间精度 | `precision:10,2`（小数保留两位） |
| `unsigned` | 无符号整数（仅正数） | `unsigned`（年龄/ID 不能为负） |

---

## 🟠 第二类：字段约束（必用）
| 标签 | 含义 | 示例 |
| ---- | ---- | ---- |
| `not null` | 非空，不能为空 | `not null` |
| `unique` | 唯一约束 | `unique` |
| `uniqueIndex` | 唯一索引（查询快+唯一） | `uniqueIndex` |
| `index` | 普通索引（加速查询） | `index` |
| `primaryKey` | 主键 | `primaryKey` |
| `autoIncrement` | 自增 | `autoIncrement` |

---

## 🟡 第三类：默认值 & 注释（必用）
| 标签 | 含义 | 示例 |
| ---- | ---- | ---- |
| `default` | 默认值 | `default:'available'`、`default:0` |
| `default:CURRENT_TIMESTAMP` | 默认当前时间 | `default:CURRENT_TIMESTAMP` |
| `comment` | 字段注释 | `comment:雨伞状态` |

---

## 🟢 第四类：字段控制（高频）
| 标签 | 含义 | 示例 |
| ---- | ---- | ---- |
| `-` | 忽略字段，不存入数据库 | `-`（临时变量专用） |
| `column` | 自定义数据库字段名 | `column:serial_num` |
| `embedded` | 嵌入结构体（合并字段） | `embedded` |

---

## 🔵 第五类：自动时间字段（高频）
| 标签 | 含义 | 用法 |
| ---- | ---- | ---- |
| `autoCreateTime` | 自动记录创建时间 | `autoCreateTime` |
| `autoUpdateTime` | 自动记录更新时间 | `autoUpdateTime` |

---

## 🟣 第六类：关联关系（**业务项目核心**）
做用户、订单、借还记录**必用**，比前面所有都重要！
| 标签 | 含义 | 适用场景 |
| ---- | ---- | ---- |
| `foreignKey` | 外键 | 一对多（用户借多把伞） |
| `references` | 关联主键 | 关联另一张表的ID |
| `many2many` | 多对多关联 | 多用户借多把伞 |

---

## 🟤 第七类：软删除（企业必备）
| 标签 | 含义 |
| ---- | ---- |
| `deletedAt` | 软删除标记（不真删数据） |

---

# 🔥 新手优先级（必背顺序）
1. **第一优先级（天天用）**
   `type`/`size`、`not null`、`uniqueIndex`、`default`、`comment`、`index`
2. **第二优先级（常用）**
   `column`、`unsigned`、`-`、`autoCreateTime`
3. **第三优先级（业务进阶）**
   关联标签 `foreignKey`、`many2many`

# 综合示例（雨伞项目完整版）
```go
type Umbrella struct {
	gorm.Model
	SerialNumber string `gorm:"size:50;uniqueIndex;not null;comment:雨伞编号;column:serial_num"`
	Status       string `gorm:"size:20;default:'available';comment:状态"`
	Age          int    `gorm:"unsigned;default:0;comment:使用时长"`
	Price        float64`gorm:"precision:10,2;comment:价格"`
	TempField    string `gorm:"-"` // 不存入数据库
}
```