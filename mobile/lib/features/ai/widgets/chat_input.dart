import 'package:flutter/material.dart';

/// Text input bar for composing and sending chat messages.
///
/// Includes a text field and send button. The send button is disabled
/// when the input is empty or while a message is being sent.
class ChatInput extends StatefulWidget {
  final Future<void> Function(String content) onSend;
  final bool isSending;

  const ChatInput({
    super.key,
    required this.onSend,
    this.isSending = false,
  });

  @override
  State<ChatInput> createState() => _ChatInputState();
}

class _ChatInputState extends State<ChatInput> {
  final TextEditingController _controller = TextEditingController();
  final FocusNode _focusNode = FocusNode();
  bool _hasText = false;

  @override
  void initState() {
    super.initState();
    _controller.addListener(() {
      final hasText = _controller.text.trim().isNotEmpty;
      if (hasText != _hasText) {
        setState(() => _hasText = hasText);
      }
    });
  }

  @override
  void dispose() {
    _controller.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  Future<void> _handleSend() async {
    final content = _controller.text.trim();
    if (content.isEmpty || widget.isSending) return;

    _controller.clear();
    await widget.onSend(content);
    _focusNode.requestFocus();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final bottomPadding = MediaQuery.of(context).padding.bottom;

    return Container(
      padding: EdgeInsets.only(
        left: 12,
        right: 8,
        top: 8,
        bottom: 8 + bottomPadding,
      ),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(
          top: BorderSide(
            color: theme.colorScheme.outline.withOpacity(0.1),
          ),
        ),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          Expanded(
            child: TextField(
              controller: _controller,
              focusNode: _focusNode,
              textCapitalization: TextCapitalization.sentences,
              maxLines: 5,
              minLines: 1,
              enabled: !widget.isSending,
              decoration: InputDecoration(
                hintText: 'Type a message...',
                hintStyle: TextStyle(
                  color: theme.colorScheme.onSurface.withOpacity(0.4),
                ),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(24),
                  borderSide: BorderSide(
                    color: theme.colorScheme.outline.withOpacity(0.2),
                  ),
                ),
                enabledBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(24),
                  borderSide: BorderSide(
                    color: theme.colorScheme.outline.withOpacity(0.2),
                  ),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(24),
                  borderSide: BorderSide(
                    color: theme.colorScheme.primary,
                  ),
                ),
                contentPadding: const EdgeInsets.symmetric(
                  horizontal: 16,
                  vertical: 10,
                ),
                filled: true,
                fillColor: theme.colorScheme.surfaceContainerHighest
                    .withOpacity(0.3),
              ),
              onSubmitted: (_) => _handleSend(),
              textInputAction: TextInputAction.send,
            ),
          ),
          const SizedBox(width: 8),
          _buildSendButton(theme),
        ],
      ),
    );
  }

  Widget _buildSendButton(ThemeData theme) {
    final canSend = _hasText && !widget.isSending;

    return AnimatedContainer(
      duration: const Duration(milliseconds: 200),
      child: Material(
        color: canSend
            ? theme.colorScheme.primary
            : theme.colorScheme.onSurface.withOpacity(0.1),
        borderRadius: BorderRadius.circular(24),
        child: InkWell(
          onTap: canSend ? _handleSend : null,
          borderRadius: BorderRadius.circular(24),
          child: Container(
            width: 44,
            height: 44,
            alignment: Alignment.center,
            child: widget.isSending
                ? SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(
                      strokeWidth: 2,
                      color: theme.colorScheme.onPrimary,
                    ),
                  )
                : Icon(
                    Icons.send_rounded,
                    size: 20,
                    color: canSend
                        ? theme.colorScheme.onPrimary
                        : theme.colorScheme.onSurface.withOpacity(0.3),
                  ),
          ),
        ),
      ),
    );
  }
}
