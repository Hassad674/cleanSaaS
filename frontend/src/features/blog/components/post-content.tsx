import { formatDate } from "@/shared/lib/utils";
import type { BlogPost } from "@/features/blog/types";

type PostContentProps = {
  post: BlogPost;
};

/**
 * Converts basic markdown syntax to HTML.
 * Handles: headings, bold, italic, links, code blocks, inline code,
 * unordered lists, ordered lists, blockquotes, images, horizontal rules,
 * and paragraphs.
 */
function markdownToHtml(markdown: string): string {
  let html = markdown;

  // Fenced code blocks (```...```)
  html = html.replace(
    /```(\w*)\n([\s\S]*?)```/g,
    (_match, _lang, code) =>
      `<pre class="bg-muted rounded-lg p-4 overflow-x-auto my-4"><code>${escapeHtml(code.trim())}</code></pre>`
  );

  // Inline code
  html = html.replace(
    /`([^`]+)`/g,
    '<code class="bg-muted rounded px-1.5 py-0.5 text-sm font-mono">$1</code>'
  );

  // Images
  html = html.replace(
    /!\[([^\]]*)\]\(([^)]+)\)/g,
    '<img src="$2" alt="$1" class="rounded-lg my-4 max-w-full" loading="lazy" />'
  );

  // Links
  html = html.replace(
    /\[([^\]]+)\]\(([^)]+)\)/g,
    '<a href="$2" class="text-primary hover:underline" target="_blank" rel="noopener noreferrer">$1</a>'
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
    /^> (.+)$/gm,
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
    } else if (trimmed.startsWith("<")) {
      result.push(line);
    } else {
      result.push(`<p class="text-foreground leading-relaxed my-3">${trimmed}</p>`);
    }
  }

  return result.join("\n");
}

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
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
