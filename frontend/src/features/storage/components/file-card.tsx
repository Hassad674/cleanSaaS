"use client";

import { useState } from "react";
import { formatDate } from "@/shared/lib/utils";
import type { FileItem } from "@/features/storage/types";

type FileCardProps = {
  file: FileItem;
  onDelete: (fileId: string) => Promise<{ error: string | null }>;
};

function formatFileSize(bytes: number): string {
  if (bytes === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${(bytes / Math.pow(1024, i)).toFixed(i > 0 ? 1 : 0)} ${units[i]}`;
}

function isImageType(contentType: string): boolean {
  return contentType.startsWith("image/");
}

function getFileIcon(contentType: string): string {
  if (contentType.startsWith("video/")) return "video";
  if (contentType.startsWith("audio/")) return "audio";
  if (contentType === "application/pdf") return "pdf";
  if (
    contentType.includes("spreadsheet") ||
    contentType.includes("excel") ||
    contentType === "text/csv"
  )
    return "spreadsheet";
  if (contentType.includes("document") || contentType.includes("word"))
    return "document";
  if (contentType.includes("zip") || contentType.includes("compressed"))
    return "archive";
  return "file";
}

function FileIcon({ type }: { type: string }) {
  const icon = getFileIcon(type);

  const labels: Record<string, string> = {
    video: "VID",
    audio: "AUD",
    pdf: "PDF",
    spreadsheet: "XLS",
    document: "DOC",
    archive: "ZIP",
    file: "FILE",
  };

  return (
    <div className="h-full w-full bg-accent flex flex-col items-center justify-center rounded-t-xl">
      <svg
        className="h-8 w-8 text-muted-foreground mb-1"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={1.5}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"
        />
      </svg>
      <span className="text-xs font-medium text-muted-foreground">
        {labels[icon]}
      </span>
    </div>
  );
}

export function FileCard({ file, onDelete }: FileCardProps) {
  const [showConfirm, setShowConfirm] = useState(false);
  const [deleting, setDeleting] = useState(false);

  async function handleDelete() {
    setDeleting(true);
    await onDelete(file.id);
    setDeleting(false);
    setShowConfirm(false);
  }

  return (
    <div className="bg-card border border-border rounded-xl shadow-sm overflow-hidden flex flex-col">
      {/* Thumbnail / Icon */}
      <div className="h-36 w-full overflow-hidden">
        {isImageType(file.content_type) ? (
          <img
            src={file.url}
            alt={file.name}
            className="h-full w-full object-cover"
            loading="lazy"
          />
        ) : (
          <FileIcon type={file.content_type} />
        )}
      </div>

      {/* File info */}
      <div className="p-4 flex-1 flex flex-col gap-2">
        <p
          className="text-sm font-medium text-foreground truncate"
          title={file.name}
        >
          {file.name}
        </p>
        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          <span>{formatFileSize(file.size_bytes)}</span>
          <span>&middot;</span>
          <span>{formatDate(file.created_at)}</span>
        </div>

        {/* Actions */}
        <div className="flex items-center gap-2 mt-auto pt-2">
          <a
            href={file.url}
            target="_blank"
            rel="noopener noreferrer"
            className="flex-1 text-center border border-border text-foreground rounded-lg px-3 py-1.5 text-xs font-medium hover:bg-muted transition-colors"
          >
            Download
          </a>

          {showConfirm ? (
            <div className="flex items-center gap-1">
              <button
                onClick={handleDelete}
                disabled={deleting}
                className="bg-destructive text-primary-foreground rounded-lg px-3 py-1.5 text-xs font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
              >
                {deleting ? "..." : "Yes"}
              </button>
              <button
                onClick={() => setShowConfirm(false)}
                className="text-xs text-muted-foreground hover:text-foreground transition-colors px-2 py-1.5"
              >
                No
              </button>
            </div>
          ) : (
            <button
              onClick={() => setShowConfirm(true)}
              className="text-muted-foreground hover:text-destructive transition-colors p-1.5"
              aria-label={`Delete ${file.name}`}
            >
              <svg
                className="h-4 w-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                strokeWidth={2}
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"
                />
              </svg>
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
