import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';

import '../../config/constants.dart';
import '../storage/secure_storage.dart';
import 'api_exceptions.dart';

/// Centralized HTTP client built on [Dio].
///
/// Features:
/// - Automatic JWT injection via auth interceptor
/// - Maps HTTP errors to typed [ApiException] subclasses
/// - Debug logging in non-release builds
/// - Configurable base URL and timeouts
class ApiClient {
  ApiClient({
    required SecureStorageService storage,
    String? baseUrl,
  }) : _storage = storage {
    _dio = Dio(
      BaseOptions(
        baseUrl: baseUrl ?? AppConstants.apiBaseUrl,
        connectTimeout:
            const Duration(milliseconds: AppConstants.connectTimeout),
        receiveTimeout:
            const Duration(milliseconds: AppConstants.receiveTimeout),
        sendTimeout: const Duration(milliseconds: AppConstants.sendTimeout),
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
      ),
    );

    // Auth interceptor — injects stored JWT into every request.
    _dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) async {
          final token = await _storage.getToken();
          if (token != null) {
            options.headers['Authorization'] = 'Bearer $token';
          }
          handler.next(options);
        },
        onError: (error, handler) async {
          // If the server returns 401, clear stored credentials.
          if (error.response?.statusCode == 401) {
            await _storage.clearAll();
          }
          handler.next(error);
        },
      ),
    );

    // Error mapping interceptor — converts DioExceptions to our typed hierarchy.
    _dio.interceptors.add(
      InterceptorsWrapper(
        onError: (error, handler) {
          handler.next(error);
        },
      ),
    );

    // Debug logging (only in debug / profile builds).
    if (kDebugMode) {
      _dio.interceptors.add(
        LogInterceptor(
          requestBody: true,
          responseBody: true,
          logPrint: (message) => debugPrint('[API] $message'),
        ),
      );
    }
  }

  final SecureStorageService _storage;
  late final Dio _dio;

  // ---------------------------------------------------------------------------
  // Public HTTP methods
  // ---------------------------------------------------------------------------

  Future<Response<T>> get<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    CancelToken? cancelToken,
  }) async {
    return _execute(
      () => _dio.get<T>(
        path,
        queryParameters: queryParameters,
        cancelToken: cancelToken,
      ),
    );
  }

  Future<Response<T>> post<T>(
    String path, {
    Object? data,
    Map<String, dynamic>? queryParameters,
    CancelToken? cancelToken,
  }) async {
    return _execute(
      () => _dio.post<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        cancelToken: cancelToken,
      ),
    );
  }

  Future<Response<T>> put<T>(
    String path, {
    Object? data,
    Map<String, dynamic>? queryParameters,
    CancelToken? cancelToken,
  }) async {
    return _execute(
      () => _dio.put<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        cancelToken: cancelToken,
      ),
    );
  }

  Future<Response<T>> patch<T>(
    String path, {
    Object? data,
    Map<String, dynamic>? queryParameters,
    CancelToken? cancelToken,
  }) async {
    return _execute(
      () => _dio.patch<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        cancelToken: cancelToken,
      ),
    );
  }

  Future<Response<T>> delete<T>(
    String path, {
    Object? data,
    Map<String, dynamic>? queryParameters,
    CancelToken? cancelToken,
  }) async {
    return _execute(
      () => _dio.delete<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        cancelToken: cancelToken,
      ),
    );
  }

  // ---------------------------------------------------------------------------
  // Internal helpers
  // ---------------------------------------------------------------------------

  /// Wraps every request in a try/catch that maps [DioException] to our
  /// typed [ApiException] hierarchy.
  Future<Response<T>> _execute<T>(
    Future<Response<T>> Function() request,
  ) async {
    try {
      return await request();
    } on DioException catch (e) {
      throw _mapDioException(e);
    }
  }

  /// Maps a [DioException] to the appropriate [ApiException] subclass.
  ApiException _mapDioException(DioException error) {
    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
      case DioExceptionType.connectionError:
        return const NetworkException('Connection timed out');

      case DioExceptionType.cancel:
        return const CancelledException();

      case DioExceptionType.badResponse:
        return _mapStatusCode(error.response);

      case DioExceptionType.badCertificate:
        return const NetworkException('Invalid SSL certificate');

      case DioExceptionType.unknown:
        if (error.error != null) {
          return NetworkException(error.error.toString());
        }
        return const NetworkException();
    }
  }

  /// Maps an HTTP status code to the matching [ApiException].
  ApiException _mapStatusCode(Response<dynamic>? response) {
    final statusCode = response?.statusCode ?? 500;
    final body = response?.data;

    // The backend returns `{"error": "message"}` for errors.
    String message = 'An unexpected error occurred';
    if (body is Map<String, dynamic> && body.containsKey('error')) {
      message = body['error'] as String;
    }

    switch (statusCode) {
      case 400:
        return BadRequestException(message);
      case 401:
        return UnauthorizedException(message);
      case 403:
        return ForbiddenException(message);
      case 404:
        return NotFoundException(message);
      case 409:
        return ConflictException(message);
      case 422:
        return ValidationException(message);
      case 429:
        return RateLimitException(message);
      default:
        if (statusCode >= 500) {
          return ServerException(message);
        }
        return ServerException(message);
    }
  }
}
