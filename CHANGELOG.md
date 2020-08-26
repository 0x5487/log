# CHANGELOG
## [2.0.0-beta.4] 2020-08-26
- add `StackTrace()` fn
- `error`, `panic`, and `fatal` level add stack_trace into entry by default, but it can be turn off by `log.AutoStackTrace = false` 
- add task runner (Taskfile.yml)
- update github workflow to v2

## [2.0.0-beta.3] 2020-08-07
- add SaveToDefault feature
- add json handler
- centralized error handling 

## [2.0.0-beta.2] 2020-05-16
- gelf will auto flush every 10 second
- redesign hook func
- add more func into context
- fix go module v2 path issue

## [2.0.0-beta1] 2020-05-09
- refactoring architecture
- use JSON as default format
- replace WithDafaultFields to hook
- replace WithFields to strongly type field type
- rename RegisterHandler to AddHandler
- add Hook function
- add WithContext func
- handler interface has been changed
- bulit-in handlers have been redesigned
- performance has been improved
- add more unit tests

## [1.0.4] 2020-04-30
- fix gelf handler race condition issue
- add standard field type

## [1.0.3] 2020-03-13
- use slice fields to improve performance
- use write buffer in gelf handler to improve performance
- use cacheLeveledHandler to improve performance (reduce map loop up)

## [1.0.2] 2020-02-23
- remove lock and improve performance
- add benchmark suite
- add code coverage