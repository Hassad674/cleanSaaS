import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

/// Design tokens matching the web frontend's CSS variables.
///
/// Light + dark themes with rose/pink primary, warm neutrals,
/// generous spacing, and smooth rounded corners.
class AppTheme {
  AppTheme._();

  // ---------------------------------------------------------------------------
  // Color palette
  // ---------------------------------------------------------------------------

  // Primary — Rose/pink (matches web --primary)
  static const Color _primaryLight = Color(0xFFE11D48); // rose-600
  static const Color _primaryDark = Color(0xFFFB7185); // rose-400
  static const Color _onPrimary = Color(0xFFFFFFFF);

  // Backgrounds
  static const Color _backgroundLight = Color(0xFFFFFFFF);
  static const Color _backgroundDark = Color(0xFF0A0A0A);

  // Foregrounds (body text)
  static const Color _foregroundLight = Color(0xFF0A0A0A);
  static const Color _foregroundDark = Color(0xFFFAFAFA);

  // Cards
  static const Color _cardLight = Color(0xFFFFFFFF);
  static const Color _cardDark = Color(0xFF1C1C1E);

  // Muted (subtle backgrounds, disabled)
  static const Color _mutedLight = Color(0xFFF5F5F4); // warm gray-100
  static const Color _mutedDark = Color(0xFF27272A); // zinc-800
  static const Color _mutedForegroundLight = Color(0xFF737373); // neutral-500
  static const Color _mutedForegroundDark = Color(0xFFA1A1AA); // zinc-400

  // Borders
  static const Color _borderLight = Color(0xFFE5E5E5); // neutral-200
  static const Color _borderDark = Color(0xFF3F3F46); // zinc-700

  // Semantic
  static const Color _destructive = Color(0xFFEF4444); // red-500
  static const Color _success = Color(0xFF22C55E); // green-500
  static const Color _warning = Color(0xFFF59E0B); // amber-500

  // Accent (hover states)
  static const Color _accentLight = Color(0xFFFFF1F2); // rose-50
  static const Color _accentDark = Color(0xFF2D1215); // dark rose tint

  // ---------------------------------------------------------------------------
  // Radii
  // ---------------------------------------------------------------------------

  static const double radiusSm = 8.0;
  static const double radiusMd = 12.0;
  static const double radiusLg = 16.0;

  // ---------------------------------------------------------------------------
  // Text theme (Geist Sans via Google Fonts, fallback to system)
  // ---------------------------------------------------------------------------

  static TextTheme _buildTextTheme(TextTheme base) {
    // Google Fonts "Geist" may not be available yet; fall back gracefully.
    try {
      return GoogleFonts.interTextTheme(base);
    } catch (_) {
      return base;
    }
  }

  // ---------------------------------------------------------------------------
  // Input decoration
  // ---------------------------------------------------------------------------

  static InputDecorationTheme _inputDecoration({
    required Color fillColor,
    required Color borderColor,
    required Color focusBorderColor,
    required Color hintColor,
  }) {
    final border = OutlineInputBorder(
      borderRadius: BorderRadius.circular(radiusMd),
      borderSide: BorderSide(color: borderColor),
    );

    return InputDecorationTheme(
      filled: true,
      fillColor: fillColor,
      hintStyle: TextStyle(color: hintColor),
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
      border: border,
      enabledBorder: border,
      focusedBorder: border.copyWith(
        borderSide: BorderSide(color: focusBorderColor, width: 2),
      ),
      errorBorder: border.copyWith(
        borderSide: BorderSide(color: _destructive),
      ),
      focusedErrorBorder: border.copyWith(
        borderSide: BorderSide(color: _destructive, width: 2),
      ),
    );
  }

  // ---------------------------------------------------------------------------
  // Button themes
  // ---------------------------------------------------------------------------

  static ElevatedButtonThemeData _elevatedButton(Color primary) {
    return ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: primary,
        foregroundColor: _onPrimary,
        minimumSize: const Size(double.infinity, 50),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(radiusMd),
        ),
        textStyle: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w600,
        ),
        elevation: 0,
      ),
    );
  }

  static OutlinedButtonThemeData _outlinedButton(Color borderColor) {
    return OutlinedButtonThemeData(
      style: OutlinedButton.styleFrom(
        minimumSize: const Size(double.infinity, 50),
        side: BorderSide(color: borderColor),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(radiusMd),
        ),
        textStyle: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }

  static TextButtonThemeData _textButton(Color primary) {
    return TextButtonThemeData(
      style: TextButton.styleFrom(
        foregroundColor: primary,
        textStyle: const TextStyle(
          fontSize: 14,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }

  // ---------------------------------------------------------------------------
  // Card theme
  // ---------------------------------------------------------------------------

  static CardThemeData _card(Color color, Color borderColor) {
    return CardThemeData(
      color: color,
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(radiusLg),
        side: BorderSide(color: borderColor),
      ),
      margin: EdgeInsets.zero,
    );
  }

  // ---------------------------------------------------------------------------
  // App bar
  // ---------------------------------------------------------------------------

  static AppBarTheme _appBar({
    required Color background,
    required Color foreground,
    required Color borderColor,
  }) {
    return AppBarTheme(
      backgroundColor: background,
      foregroundColor: foreground,
      elevation: 0,
      scrolledUnderElevation: 0,
      surfaceTintColor: Colors.transparent,
      shape: Border(bottom: BorderSide(color: borderColor, width: 0.5)),
      titleTextStyle: TextStyle(
        color: foreground,
        fontSize: 18,
        fontWeight: FontWeight.w600,
      ),
    );
  }

  // ---------------------------------------------------------------------------
  // Public theme getters
  // ---------------------------------------------------------------------------

  static ThemeData get light {
    final base = ThemeData.light(useMaterial3: true);
    final textTheme = _buildTextTheme(base.textTheme);

    return base.copyWith(
      colorScheme: ColorScheme.light(
        primary: _primaryLight,
        onPrimary: _onPrimary,
        secondary: _accentLight,
        surface: _cardLight,
        onSurface: _foregroundLight,
        error: _destructive,
      ),
      scaffoldBackgroundColor: _backgroundLight,
      textTheme: textTheme,
      appBarTheme: _appBar(
        background: _backgroundLight,
        foreground: _foregroundLight,
        borderColor: _borderLight,
      ),
      cardTheme: _card(_cardLight, _borderLight),
      elevatedButtonTheme: _elevatedButton(_primaryLight),
      outlinedButtonTheme: _outlinedButton(_borderLight),
      textButtonTheme: _textButton(_primaryLight),
      inputDecorationTheme: _inputDecoration(
        fillColor: _backgroundLight,
        borderColor: _borderLight,
        focusBorderColor: _primaryLight,
        hintColor: _mutedForegroundLight,
      ),
      dividerColor: _borderLight,
      dividerTheme: DividerThemeData(color: _borderLight, thickness: 0.5),
      bottomNavigationBarTheme: BottomNavigationBarThemeData(
        backgroundColor: _backgroundLight,
        selectedItemColor: _primaryLight,
        unselectedItemColor: _mutedForegroundLight,
        type: BottomNavigationBarType.fixed,
        elevation: 0,
      ),
      chipTheme: ChipThemeData(
        backgroundColor: _mutedLight,
        labelStyle: TextStyle(color: _foregroundLight),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(radiusSm),
        ),
      ),
      snackBarTheme: SnackBarThemeData(
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(radiusMd),
        ),
      ),
      extensions: <ThemeExtension<dynamic>>[
        AppColors(
          muted: _mutedLight,
          mutedForeground: _mutedForegroundLight,
          accent: _accentLight,
          border: _borderLight,
          success: _success,
          warning: _warning,
        ),
      ],
    );
  }

  static ThemeData get dark {
    final base = ThemeData.dark(useMaterial3: true);
    final textTheme = _buildTextTheme(base.textTheme);

    return base.copyWith(
      colorScheme: ColorScheme.dark(
        primary: _primaryDark,
        onPrimary: _onPrimary,
        secondary: _accentDark,
        surface: _cardDark,
        onSurface: _foregroundDark,
        error: _destructive,
      ),
      scaffoldBackgroundColor: _backgroundDark,
      textTheme: textTheme,
      appBarTheme: _appBar(
        background: _backgroundDark,
        foreground: _foregroundDark,
        borderColor: _borderDark,
      ),
      cardTheme: _card(_cardDark, _borderDark),
      elevatedButtonTheme: _elevatedButton(_primaryDark),
      outlinedButtonTheme: _outlinedButton(_borderDark),
      textButtonTheme: _textButton(_primaryDark),
      inputDecorationTheme: _inputDecoration(
        fillColor: _cardDark,
        borderColor: _borderDark,
        focusBorderColor: _primaryDark,
        hintColor: _mutedForegroundDark,
      ),
      dividerColor: _borderDark,
      dividerTheme: DividerThemeData(color: _borderDark, thickness: 0.5),
      bottomNavigationBarTheme: BottomNavigationBarThemeData(
        backgroundColor: _backgroundDark,
        selectedItemColor: _primaryDark,
        unselectedItemColor: _mutedForegroundDark,
        type: BottomNavigationBarType.fixed,
        elevation: 0,
      ),
      chipTheme: ChipThemeData(
        backgroundColor: _mutedDark,
        labelStyle: TextStyle(color: _foregroundDark),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(radiusSm),
        ),
      ),
      snackBarTheme: SnackBarThemeData(
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(radiusMd),
        ),
      ),
      extensions: <ThemeExtension<dynamic>>[
        AppColors(
          muted: _mutedDark,
          mutedForeground: _mutedForegroundDark,
          accent: _accentDark,
          border: _borderDark,
          success: _success,
          warning: _warning,
        ),
      ],
    );
  }
}

/// Custom theme extension for colors not covered by [ColorScheme].
///
/// Access via `Theme.of(context).extension<AppColors>()!`.
@immutable
class AppColors extends ThemeExtension<AppColors> {
  const AppColors({
    required this.muted,
    required this.mutedForeground,
    required this.accent,
    required this.border,
    required this.success,
    required this.warning,
  });

  final Color muted;
  final Color mutedForeground;
  final Color accent;
  final Color border;
  final Color success;
  final Color warning;

  @override
  AppColors copyWith({
    Color? muted,
    Color? mutedForeground,
    Color? accent,
    Color? border,
    Color? success,
    Color? warning,
  }) {
    return AppColors(
      muted: muted ?? this.muted,
      mutedForeground: mutedForeground ?? this.mutedForeground,
      accent: accent ?? this.accent,
      border: border ?? this.border,
      success: success ?? this.success,
      warning: warning ?? this.warning,
    );
  }

  @override
  AppColors lerp(ThemeExtension<AppColors>? other, double t) {
    if (other is! AppColors) return this;
    return AppColors(
      muted: Color.lerp(muted, other.muted, t)!,
      mutedForeground: Color.lerp(mutedForeground, other.mutedForeground, t)!,
      accent: Color.lerp(accent, other.accent, t)!,
      border: Color.lerp(border, other.border, t)!,
      success: Color.lerp(success, other.success, t)!,
      warning: Color.lerp(warning, other.warning, t)!,
    );
  }
}
