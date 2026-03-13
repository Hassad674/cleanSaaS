import 'package:flutter/material.dart';

/// A reusable settings list tile with icon, label, and optional trailing widget.
///
/// Used throughout the settings screen to display individual setting options.
/// The [trailing] can be a switch, arrow icon, badge, or any custom widget.
class SettingsTile extends StatelessWidget {
  /// Leading icon displayed before the label.
  final IconData icon;

  /// Primary text label for the setting.
  final String label;

  /// Optional subtitle text displayed below the label.
  final String? subtitle;

  /// Optional trailing widget (e.g., [Switch], [Icon], badge).
  /// Defaults to a chevron-right icon if [onTap] is provided and no trailing
  /// widget is specified.
  final Widget? trailing;

  /// Callback when the tile is tapped. If null, the tile is not tappable.
  final VoidCallback? onTap;

  /// Optional color override for the icon. Defaults to theme primary.
  final Color? iconColor;

  /// Whether the tile is destructive (e.g., delete account).
  /// When true, uses error color for icon and label.
  final bool isDestructive;

  const SettingsTile({
    super.key,
    required this.icon,
    required this.label,
    this.subtitle,
    this.trailing,
    this.onTap,
    this.iconColor,
    this.isDestructive = false,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final effectiveIconColor = isDestructive
        ? theme.colorScheme.error
        : iconColor ?? theme.colorScheme.primary;
    final effectiveLabelColor = isDestructive
        ? theme.colorScheme.error
        : theme.colorScheme.onSurface;

    final effectiveTrailing = trailing ??
        (onTap != null
            ? Icon(
                Icons.chevron_right,
                color: theme.colorScheme.onSurface.withOpacity(0.4),
              )
            : null);

    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
          child: Row(
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: effectiveIconColor.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(10),
                ),
                child: Icon(
                  icon,
                  size: 20,
                  color: effectiveIconColor,
                ),
              ),
              const SizedBox(width: 14),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      label,
                      style: theme.textTheme.bodyLarge?.copyWith(
                        fontWeight: FontWeight.w500,
                        color: effectiveLabelColor,
                      ),
                    ),
                    if (subtitle != null) ...[
                      const SizedBox(height: 2),
                      Text(
                        subtitle!,
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: theme.colorScheme.onSurface.withOpacity(0.5),
                        ),
                      ),
                    ],
                  ],
                ),
              ),
              if (effectiveTrailing != null) effectiveTrailing,
            ],
          ),
        ),
      ),
    );
  }
}
