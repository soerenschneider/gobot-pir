# Changelog

## [1.6.3](https://github.com/soerenschneider/gobot-pir/compare/v1.6.2...v1.6.3) (2023-10-16)


### Bug Fixes

* **deps:** bump github.com/go-playground/validator/v10 ([79c980d](https://github.com/soerenschneider/gobot-pir/commit/79c980d38688e4eafb5464859410b7b45e733d80))
* **deps:** bump github.com/prometheus/client_golang ([ba016a9](https://github.com/soerenschneider/gobot-pir/commit/ba016a93d9d3e5704d69eb9242646235f66d4db0))
* fix default values ([57a4ce6](https://github.com/soerenschneider/gobot-pir/commit/57a4ce6e6a545f493d50e6926e7b649809afef40))
* fix syntax ([d5f2937](https://github.com/soerenschneider/gobot-pir/commit/d5f29376070779f2657ff55969f3295a7eec0649))

## [1.6.2](https://github.com/soerenschneider/gobot-pir/compare/v1.6.1...v1.6.2) (2023-07-11)


### Bug Fixes

* import ([f18eb2a](https://github.com/soerenschneider/gobot-pir/commit/f18eb2a420fc946318fcab6acd41815b23a94cfe))

## [1.6.1](https://github.com/soerenschneider/gobot-pir/compare/v1.6.0...v1.6.1) (2022-06-07)


### Bug Fixes

* call validate() on mqttconfig ([af945e3](https://github.com/soerenschneider/gobot-pir/commit/af945e3c9c7a75a9078eaa5f15022eaa325940fd))

## [1.6.0](https://www.github.com/soerenschneider/gobot-pir/compare/v1.5.0...v1.6.0) (2022-05-04)


### Features

* enable tls client cert auth ([e2907ff](https://www.github.com/soerenschneider/gobot-pir/commit/e2907ffcbcebab2046a2e747d6c5e584787e081c))

## [1.5.0](https://www.github.com/soerenschneider/gobot-pir/compare/v1.4.1...v1.5.0) (2021-11-22)


### Features

* Collect statistics over configurable intervals ([d61b5e0](https://www.github.com/soerenschneider/gobot-pir/commit/d61b5e0f0e6e933a3b880f4b1c3e58159217f29d))


### Bug Fixes

* Forgot to call FormatTopics ([8863854](https://www.github.com/soerenschneider/gobot-pir/commit/8863854c136998affc225522f1df1853a84c22a3))

### [1.4.1](https://www.github.com/soerenschneider/gobot-pir/compare/v1.4.0...v1.4.1) (2021-11-02)


### Miscellaneous Chores

* Trigger release ([ead1bf8](https://www.github.com/soerenschneider/gobot-pir/commit/ead1bf8581321f8bd66d95185d325bbcff588800))

## [1.4.0](https://www.github.com/soerenschneider/gobot-pir/compare/v1.3.0...v1.4.0) (2021-10-20)


### Features

* Better control over what to send with mqtt payload ([4328eb4](https://www.github.com/soerenschneider/gobot-pir/commit/4328eb47a8f1af59ae787b7fb0386269298e6b0c))


### Bug Fixes

* Enable auto-reconnect to mqtt broker ([58d1cc8](https://www.github.com/soerenschneider/gobot-pir/commit/58d1cc81b9ffc806bc12427acef7842bc1aa4b3a))

## [1.3.0](https://www.github.com/soerenschneider/gobot-pir/compare/v1.2.0...v1.3.0) (2021-10-09)


### Features

* Add version metric ([ad646e7](https://www.github.com/soerenschneider/gobot-pir/commit/ad646e738c72172439a999d44df6f09fa7f68bfc))
* Flag to print version and quit ([a569add](https://www.github.com/soerenschneider/gobot-pir/commit/a569addb62fec873836191d6f12e70f68b554f75))

## [1.2.0](https://www.github.com/soerenschneider/gobot-pir/compare/v1.1.3...v1.2.0) (2021-10-08)


### Features

* Add heartbeat metric ([b8f2eae](https://www.github.com/soerenschneider/gobot-pir/commit/b8f2eae1bf8899b64145019fcfec5e05fec8980c))
* Allow reproducible builds ([15e5a4a](https://www.github.com/soerenschneider/gobot-pir/commit/15e5a4a739d81f1c3cf663afff3a81d99b074a9f))


### Bug Fixes

* Add missing type ([f3000b4](https://www.github.com/soerenschneider/gobot-pir/commit/f3000b45a588e1ff1dec09ac00d49d6f940c93c0))

### [1.1.3](https://www.github.com/soerenschneider/gobot-pir/compare/v1.1.2...v1.1.3) (2021-09-16)


### Bug Fixes

* Build conf from env instead of default values conf ([cddc628](https://www.github.com/soerenschneider/gobot-pir/commit/cddc628595157abaf0bca4c450e3221f1e6ec90b))

### [1.1.2](https://www.github.com/soerenschneider/gobot-pir/compare/v1.1.1...v1.1.2) (2021-09-15)


### Bug Fixes

* use correct mod, change binary name ([b650851](https://www.github.com/soerenschneider/gobot-pir/commit/b65085147eba291b714670d1e66f84735a280b33))

### [1.1.1](https://www.github.com/soerenschneider/gobot-motion-detection/compare/v1.1.0...v1.1.1) (2021-09-14)


### Bug Fixes

* use correct mod, change binary name ([3aaee1e](https://www.github.com/soerenschneider/gobot-motion-detection/commit/3aaee1e1a4dc45593fda2b03a49c28c28f931a81))

## [1.1.0](https://www.github.com/soerenschneider/gobot-motion-detection/compare/v1.0.0...v1.1.0) (2021-09-14)


### Features

* Add build info to runtime ([d010070](https://www.github.com/soerenschneider/gobot-motion-detection/commit/d010070c0513e083e17a8bcade4d1538aca6676c))

## 1.0.0 (2021-09-13)


### Miscellaneous Chores

* release 1.0.0 ([2c8dba8](https://www.github.com/soerenschneider/gobot-motion-detection/commit/2c8dba80dd3738f3f432553ce2342adccc9aae5a))
