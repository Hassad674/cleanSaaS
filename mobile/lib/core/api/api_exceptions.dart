/// Base class for all API-related exceptions.
sealed class ApiException implements Exception {
  const ApiException(this.message, {this.statusCode});

  final String message;
  final int? statusCode;

  @override
  String toString() => 'ApiException($statusCode): $message';
}

/// 400 — The request was malformed or contained invalid data.
class BadRequestException extends ApiException {
  const BadRequestException([String message = 'Bad request'])
      : super(message, statusCode: 400);
}

/// 401 — Authentication is required or the token has expired.
class UnauthorizedException extends ApiException {
  const UnauthorizedException([String message = 'Unauthorized'])
      : super(message, statusCode: 401);
}

/// 403 — The user does not have permission for this action.
class ForbiddenException extends ApiException {
  const ForbiddenException([String message = 'Forbidden'])
      : super(message, statusCode: 403);
}

/// 404 — The requested resource was not found.
class NotFoundException extends ApiException {
  const NotFoundException([String message = 'Not found'])
      : super(message, statusCode: 404);
}

/// 409 — A conflicting resource already exists (e.g. duplicate email).
class ConflictException extends ApiException {
  const ConflictException([String message = 'Already exists'])
      : super(message, statusCode: 409);
}

/// 422 — Validation error.
class ValidationException extends ApiException {
  const ValidationException([String message = 'Validation error'])
      : super(message, statusCode: 422);
}

/// 429 — Too many requests, rate limit exceeded.
class RateLimitException extends ApiException {
  const RateLimitException([String message = 'Too many requests'])
      : super(message, statusCode: 429);
}

/// 500+ — An unexpected server error occurred.
class ServerException extends ApiException {
  const ServerException([String message = 'Internal server error'])
      : super(message, statusCode: 500);
}

/// Network-level failure (no internet, DNS resolution, timeout, etc.).
class NetworkException extends ApiException {
  const NetworkException([String message = 'Network error'])
      : super(message);
}

/// The request was cancelled before completing.
class CancelledException extends ApiException {
  const CancelledException([String message = 'Request cancelled'])
      : super(message);
}
