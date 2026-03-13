import { formatDate } from "@/shared/lib/utils";
import type { BlogPost } from "@/features/blog/types";

type PostContentProps = {
  post: BlogPost;
};

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}

/**
 * Sanitizes a URL to prevent javascript: and data: protocol XSS attacks.
 * Only allows http:, https:, mailto:, and relative URLs.
 */
function sanitizeUrl(url: string): string {
  const trimmed = url.trim();
  // Allow relative URLs
  if (trimmed.startsWith("/") || trimmed.startsWith("#")) {
    return trimmed;
  }
  // Allow http(s) and mailto only
  if (
    trimmed.startsWith("http://") ||
    trimmed.startsWith("https://") ||
    trimmed.startsWith("mailto:")
  ) {
    return trimmed;
  }
  // Block everything else (javascript:, data:, vbscript:, etc.)
  return "#";
}

/**
 * Converts basic markdown syntax to HTML.
 *
 * Security: All content is HTML-escaped FIRST, then markdown syntax is
 * converted. This prevents stored XSS from blog content in the database.
 */
function markdownToHtml(markdown: string): string {
  // Step 1: Extract and preserve fenced code blocks before escaping.
  // We replace them with placeholders, escape everything else, then restore.
  const codeBlocks: string[] = [];
  let processed = markdown.replace(
    /```(\w*)\n([\s\S]*?)```/g,
    (_match, _lang, code) => {
      const idx = codeBlocks.length;
      codeBlocks.push(
        `<pre class="bg-muted rounded-lg p-4 overflow-x-auto my-4"><code>${escapeHtml(code.trim())}</code></pre>`
      );
      return `%%CODEBLOCK_${idx}%%`;
    }
  );

  // Step 2: Extract inline code before escaping
  const inlineCodes: string[] = [];
  processed = processed.replace(/`([^`]+)`/g, (_match, code) => {
    const idx = inlineCodes.length;
    inlineCodes.push(
      `<code class="bg-muted rounded px-1.5 py-0.5 text-sm font-mono">${escapeHtml(code)}</code>`
    );
    return `%%INLINECODE_${idx}%%`;
  });

  // Step 3: Escape ALL remaining HTML to prevent XSS
  let html = escapeHtml(processed);

  // Step 4: Now apply markdown transformations on the escaped content.
  // Since content is escaped, user-injected HTML/JS is neutralized.

  // Images — sanitize URLs
  html = html.replace(
    /!\[([^\]]*)\]\(([^)]+)\)/g,
    (_match, alt, src) =>
      `<img src="${sanitizeUrl(src)}" alt="${alt}" class="rounded-lg my-4 max-w-full" loading="lazy" />`
  );

  // Links — sanitize URLs
  html = html.replace(
    /\[([^\]]+)\]\(([^)]+)\)/g,
    (_match, text, href) =>
      `<a href="${sanitizeUrl(href)}" class="text-primary hover:underline" target="_blank" rel="noopener noreferrer">${text}</a>`
  );

  // Headings
  html = html.replace(
    /^#### (.+)$/gm,
    '<h4 class="text-base font-bold text-foreground mt-6 mb-2">$1</h4>'
  );
  html = html.replace(
    /^### (.+)$/gm,
    '<h3 class="text-lg font-bold text-foreground mt-8 mb-3">$1</h3>'
  );
  html = html.replace(
    /^## (.+)$/gm,
    '<h2 class="text-xl font-bold text-foreground mt-8 mb-3">$1</h2>'
  );
  html = html.replace(
    /^# (.+)$/gm,
    '<h1 class="text-2xl font-bold text-foreground mt-8 mb-4">$1</h1>'
  );

  // Horizontal rules
  html = html.replace(
    /^---$/gm,
    '<hr class="border-border my-8" />'
  );

  // Blockquotes
  html = html.replace(
    /^&gt; (.+)$/gm,
    '<blockquote class="border-l-4 border-primary pl-4 my-4 text-muted-foreground italic">$1</blockquote>'
  );

  // Bold and italic
  html = html.replace(/\*\*\*(.+?)\*\*\*/g, "<strong><em>$1</em></strong>");
  html = html.replace(/\*\*(.+?)\*\*/g, "<strong>$1</strong>");
  html = html.replace(/\*(.+?)\*/g, "<em>$1</em>");

  // Unordered lists
  html = html.replace(
    /^[*-] (.+)$/gm,
    '<li class="ml-4 list-disc text-foreground">$1</li>'
  );

  // Ordered lists
  html = html.replace(
    /^\d+\. (.+)$/gm,
    '<li class="ml-4 list-decimal text-foreground">$1</li>'
  );

  // Wrap consecutive <li> elements in <ul> or <ol>
  html = html.replace(
    /(<li class="ml-4 list-disc[^"]*">[\s\S]*?<\/li>\n?)+/g,
    (match) => `<ul class="my-4 space-y-1">${match}</ul>`
  );
  html = html.replace(
    /(<li class="ml-4 list-decimal[^"]*">[\s\S]*?<\/li>\n?)+/g,
    (match) => `<ol class="my-4 space-y-1">${match}</ol>`
  );

  // Paragraphs: wrap remaining lines that aren't already HTML
  const lines = html.split("\n");
  const result: string[] = [];
  for (const line of lines) {
    const trimmed = line.trim();
    if (trimmed === "") {
      result.push("");
    } else if (trimmed.startsWith("<") || trimmed.startsWith("%%")) {
      result.push(line);
    } else {
      result.push(`<p class="text-foreground leading-relaxed my-3">${trimmed}</p>`);
    }
  }

  html = result.join("\n");

  // Step 5: Restore code blocks and inline code
  for (let i = 0; i < codeBlocks.length; i++) {
    html = html.replace(`%%CODEBLOCK_${i}%%`, codeBlocks[i]);
  }
  for (let i = 0; i < inlineCodes.length; i++) {
    html = html.replace(`%%INLINECODE_${i}%%`, inlineCodes[i]);
  }

  return html;
}

export function PostContent({ post }: PostContentProps) {
  const contentHtml = markdownToHtml(post.content);

  return (
    <article className="max-w-3xl mx-auto">
      {/* Header */}
      <header className="mb-8">
        {/* Tags */}
        {post.tags.length > 0 && (
          <div className="flex flex-wrap gap-2 mb-4">
            {post.tags.map((tag) => (
              <span
                key={tag}
                className="text-xs font-medium bg-accent text-foreground rounded-md px-2.5 py-1"
              >
                {tag}
              </span>
            ))}
          </div>
        )}

        <h1 className="text-3xl sm:text-4xl font-bold text-foreground mb-4">
          {post.title}
        </h1>

        {post.excerpt && (
          <p className="text-lg text-muted-foreground mb-4">
            {post.excerpt}
          </p>
        )}

        {post.published_at && (
          <time
            dateTime={post.published_at}
            className="text-sm text-muted-foreground"
          >
            Published on {formatDate(post.published_at)}
          </time>
        )}
      </header>

      {/* Cover image */}
      {post.cover_image_url && (
        <div className="mb-8 rounded-xl overflow-hidden">
          <img
            src={post.cover_image_url}
            alt={post.title}
            className="w-full object-cover"
          />
        </div>
      )}

      {/* Content */}
      <div
        className="prose-custom text-foreground"
        dangerouslySetInnerHTML={{ __html: contentHtml }}
      />
    </article>
  );
}
