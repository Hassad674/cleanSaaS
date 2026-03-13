import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:cleansaas_mobile/features/billing/models/subscription.dart';

/// A list tile displaying a single invoice with date, amount, status badge,
/// and an optional download button.
class InvoiceTile extends StatelessWidget {
  final Invoice invoice;
  final VoidCallback? onDownload;

  const InvoiceTile({
    super.key,
    required this.invoice,
    this.onDownload,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final dateFormat = DateFormat('MMM d, yyyy');

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
      child: Row(
        children: [
          // Invoice icon
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: theme.colorScheme.primaryContainer.withOpacity(0.3),
              borderRadius: BorderRadius.circular(10),
            ),
            child: Icon(
              Icons.receipt_outlined,
              size: 20,
              color: theme.colorScheme.primary,
            ),
          ),
          const SizedBox(width: 14),

          // Date and amount
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  dateFormat.format(invoice.createdAt),
                  style: theme.textTheme.bodyMedium?.copyWith(
                    fontWeight: FontWeight.w500,
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  invoice.formattedAmount,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: theme.colorScheme.onSurface.withOpacity(0.6),
                  ),
                ),
              ],
            ),
          ),

          // Status badge
          _buildStatusBadge(theme),
          const SizedBox(width: 8),

          // Download button
          if (invoice.downloadUrl != null)
            IconButton(
              onPressed: onDownload,
              icon: Icon(
                Icons.download_outlined,
                color: theme.colorScheme.primary,
                size: 20,
              ),
              tooltip: 'Download invoice',
              style: IconButton.styleFrom(
                backgroundColor:
                    theme.colorScheme.primary.withOpacity(0.1),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildStatusBadge(ThemeData theme) {
    Color backgroundColor;
    Color foregroundColor;

    switch (invoice.status) {
      case 'paid':
        backgroundColor = Colors.green.withOpacity(0.1);
        foregroundColor = Colors.green;
        break;
      case 'pending':
        backgroundColor = Colors.orange.withOpacity(0.1);
        foregroundColor = Colors.orange;
        break;
      case 'failed':
        backgroundColor = theme.colorScheme.error.withOpacity(0.1);
        foregroundColor = theme.colorScheme.error;
        break;
      default:
        backgroundColor = theme.colorScheme.surfaceContainerHighest;
        foregroundColor = theme.colorScheme.onSurface;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(20),
      ),
      child: Text(
        invoice.statusLabel,
        style: theme.textTheme.labelSmall?.copyWith(
          color: foregroundColor,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }
}
