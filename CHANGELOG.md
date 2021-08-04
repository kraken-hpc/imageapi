# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2021-08-04
### Added
- Support for NFS mounts
- Support for Bind mounts (local directories or sub-directories of mounts)
- Support for Attach mount type to mount any valid attachment  
- Support for Local attachment type (for already existing block devices, e.g. local disk)
- Support for Loopback attachment type to create a /dev/loop* attachment from local or mounted files
- Support (experimental) for iSCSI attachment type to attach an iSCSI lun
### Changed
- API moved to revision v0.2.0, refactored for simpler endpoints with query params
- API added data types to describe new mount/attach types
- API now follows a consistent scheme for attach and mount objects
- Internal engine refactored to use standard interfaces for mount/attach "drivers"
- Refactored central object storage, global ID tracking, and global reference tracking
- Refactored garbage collection
### Deprecated
- RBD mount type is replaced with generic "attach" mount type, allowing mount of any attachment type

## [0.1.5] - 2021-04-22
### Added
- Default settings for VSCode to run vet checks with GOOS=linux
### Fixed
- Fix issue #32 where makebb environments would call imageapi init() code with every command
- Fix incorrect link in documentation to swaggerhub

## [0.1.4] - 2021-03-22
### Added
- Set log level through environment variable `IMAGEAPI_LOGLEVEL`
- CirceCI integration (initial build recipe)
- CodeFactor integration
### Fixed
- Readme reflects current API, including nesting

## [0.1.3] - 2021-03-19
### Changed
- Migration to `kraken-hpc`

## [0.1.2] - 2021-03-19
### Added
- Garbage collection for unused resources
- Logging for all object types
### Fixed
- Various reference tracking bugs
- Mount.RefAdd should work regardless of specification type

## [0.1.1] - 2021-03-04
### Fixed
- Disallow trying to attach an RBD device that is already attached

## [0.1.0] - 2021-02-26
### Added
- Initial versioned release
