# CHANGELOG

## [1.0.4] 
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