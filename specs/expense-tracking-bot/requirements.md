# Requirements Document

## Introduction

This feature implements a backend system for an expense tracking application that integrates with messaging platforms (Telegram, WhatsApp, Discord). Users can log their income and expenses through natural language commands sent via bots, which are then parsed by AI and stored in a PostgreSQL database. The system uses Gin framework for the REST API and GORM for database operations.

## Requirements

### Requirement 1

**User Story:** As a user, I want to send expense/income commands through messaging bots, so that I can quickly log my financial transactions without opening a separate app.

#### Acceptance Criteria

1. WHEN a user sends a command like "/in salary 3000000" or "/out kost 2jt" THEN the system SHALL receive and process the message from Discord/Telegram/Whatsapp bot
2. WHEN a bot receives a message THEN the system SHALL forward it to the backend API endpoint
3. WHEN the backend receives a bot message THEN the system SHALL include the user identifier and platform information

### Requirement 2

**User Story:** As a system, I want to parse natural language expense commands using AI, so that I can extract structured data from user messages.

#### Acceptance Criteria

1. WHEN the system receives a message like "/in freelance 5.000.000" THEN the AI service SHALL parse it into JSON format with category="job", name="freelance", value=5000000, type="income"
2. WHEN the system receives a message like "/out kost 2jt" THEN the AI service SHALL parse it into JSON format with category="housing", name="kost", value=2000000, type="outcome"
3. WHEN the system receives an unparseable message THEN the AI service SHALL return an error response
4. WHEN the AI parsing succeeds THEN the system SHALL validate the extracted data format
5. WHEN the AI parsing fails THEN the system SHALL return a user-friendly error message to the bot

### Requirement 3

**User Story:** As a user, I want my financial transactions to be stored securely in a database, so that I can track my expenses and income over time.

#### Acceptance Criteria

1. WHEN the AI successfully parses a transaction THEN the system SHALL store it in PostgreSQL database with user association
2. WHEN storing a transaction THEN the system SHALL include timestamp, user_id, category, name, value, and transaction type
3. WHEN a database error occurs THEN the system SHALL handle it gracefully and return appropriate error response
4. WHEN storing a transaction THEN the system SHALL ensure data integrity and validation
5. WHEN a transaction is stored successfully THEN the system SHALL return confirmation to the user via the bot

### Requirement 4

**User Story:** As a developer, I want a RESTful API built with Gin framework, so that the system can handle bot requests efficiently and be easily maintainable.

#### Acceptance Criteria

1. WHEN the system starts THEN it SHALL initialize a Gin HTTP server on a configurable port
2. WHEN a POST request is made to /api/transactions THEN the system SHALL process the bot message
3. WHEN the API receives a request THEN it SHALL validate the request format and authentication
4. WHEN processing is complete THEN the system SHALL return appropriate HTTP status codes and JSON responses
5. WHEN an error occurs THEN the system SHALL log the error and return structured error responses

### Requirement 5

**User Story:** As a system administrator, I want proper database management with GORM, so that the system can efficiently handle data operations and migrations.

#### Acceptance Criteria

1. WHEN the system starts THEN it SHALL connect to PostgreSQL database using GORM
2. WHEN the system initializes THEN it SHALL run database migrations for transaction and user tables
3. WHEN performing database operations THEN the system SHALL use GORM models and methods
4. WHEN a database connection fails THEN the system SHALL handle the error and provide meaningful feedback
5. WHEN the system shuts down THEN it SHALL properly close database connections

### Requirement 6

**User Story:** As a user, I want my transactions to be associated with my account, so that my data is kept separate from other users.

#### Acceptance Criteria

1. WHEN a bot message is received THEN the system SHALL identify the user from the platform-specific user ID
2. WHEN storing a transaction THEN the system SHALL associate it with the correct user account
3. WHEN a new user sends their first message THEN the system SHALL create a user record if it doesn't exist
4. WHEN querying transactions THEN the system SHALL only return data for the authenticated user
5. WHEN user identification fails THEN the system SHALL return an authentication error