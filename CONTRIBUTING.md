# Git 提交信息规范

## 1. 提交信息格式
提交信息应具体明了，遵循以下格式：
```
<类型>: <简短描述>
<空行>
<详细描述>
```

## 2. 类型
以下是允许使用的提交类型，每种类型对应不同的代码更改目的：

- **feat**: 引入新功能。
- **fix**: 修复 bug。
- **pref**: 提升性能的代码更改。
- **docs**: 仅文档新增或变更。
- **style**: 不影响代码逻辑的样式调整（如格式化）。
- **refactor**: 代码重构，无新增功能或修复 bug。
- **test**: 添加或修改测试代码。
- **chore**: 其他更新，如构建流程或辅助工具的变动。
- **ci**: 持续集成相关文件和脚本的变更。

## 3. 简短描述
简短描述应简明扼要地说明提交的主要变更，通常不超过50个字符。

## 4. 详细描述（可选）
如果需要，可在第二段提供更详尽的描述，说明更改的动机和与之前行为的对比等。

## 示例

```
feat: add user login functionality

Add user login functionality using JWT for authentication. This includes
models, controllers, and tests.
```

## 使用规范的好处
遵守这些规范将帮助团队成员理解历史提交的目的，便于代码审查和版本控制。此外，这也有助于自动化工具生成更清晰的版本日志。

## 提交前的自查清单
在提交前，请确保您的代码已经通过了所有测试，并且所有的代码更改都符合项目的编码规范。

## 总结
通过遵循这些简单的规则，我们可以保持项目的历史记录清晰、有序。团队成员应定期回顾这些规范以确保持续遵守。

