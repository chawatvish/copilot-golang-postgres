# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Phone number field to User model (required in API, nullable in database)
- Address field to User model (optional in both API and database)
- Database migration support for adding new fields to existing tables with data
- Pointer type handling for nullable database fields while maintaining API contracts

### Changed

- User model structure to include `Phone *string` and `Address *string` fields
- API request/response structures to include phone and address fields
- Service layer methods to handle pointer conversion for new fields
- Repository implementations (both GORM and in-memory) to support new fields
- Database seeding with sample phone and address data

### Technical Details

- Used `*string` pointer types to handle nullable database columns
- Implemented proper GORM tags: `gorm:"type:text;default:null"` for phone field
- Maintained backward compatibility with existing data
- All tests pass in both in-memory and database modes

## [Previous Versions]

- Initial implementation with basic User CRUD operations
- Clean architecture with handlers, services, and repositories
- PostgreSQL database integration with GORM
- In-memory fallback repository for testing
- Comprehensive test suite
