"use client";

import { useCallback, useRef, useState } from "react";
import Link from "next/link";
import { cn } from "@/shared/lib/utils";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

type FileType = "pdf" | "image" | "spreadsheet" | "video" | "text" | "other";

interface MockFile {
  id: string;
  name: string;
  size: number;
  type: FileType;
  uploadedAt: Date;
}

interface UploadingFile {
  id: string;
  name: string;
  size: number;
  type: FileType;
  progress: number;
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function generateId(): string {
  return Math.random().toString(36).substring(2, 10);
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(0)} KB`;
  if (bytes < 1024 * 1024 * 1024)
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`;
}

function formatRelativeTime(date: Date): string {
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMin = Math.floor(diffMs / 60_000);
  const diffHours = Math.floor(diffMs / 3_600_000);
  const diffDays = Math.floor(diffMs / 86_400_000);
  const diffWeeks = Math.floor(diffDays / 7);
  const diffMonths = Math.floor(diffDays / 30);

  if (diffMin < 1) return "just now";
  if (diffMin < 60) return `${diffMin}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;
  if (diffWeeks < 5) return `${diffWeeks}w ago`;
  return `${diffMonths}mo ago`;
}

function inferFileType(fileName: string): FileType {
  const ext = fileName.split(".").pop()?.toLowerCase() ?? "";
  if (["pdf"].includes(ext)) return "pdf";
  if (["png", "jpg", "jpeg", "gif", "svg", "webp"].includes(ext))
    return "image";
  if (["xls", "xlsx", "csv"].includes(ext)) return "spreadsheet";
  if (["mp4", "mov", "avi", "webm", "mkv"].includes(ext)) return "video";
  if (["txt", "md", "doc", "docx", "rtf"].includes(ext)) return "text";
  return "other";
}

// ---------------------------------------------------------------------------
// Icons
// ---------------------------------------------------------------------------

function FileTypeIcon({
  type,
  className,
}: {
  type: FileType;
  className?: string;
}) {
  const base = cn("shrink-0", className);

  switch (type) {
    case "pdf":
      return (
        <svg
          className={base}
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
          <text
            x="8.5"
            y="18"
            fontSize="5"
            fill="currentColor"
            stroke="none"
            fontWeight="bold"
          >
            PDF
          </text>
        </svg>
      );
    case "image":
      return (
        <svg
          className={base}
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={1.5}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z"
          />
        </svg>
      );
    case "spreadsheet":
      return (
        <svg
          className={base}
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={1.5}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M3.375 19.5h17.25m-17.25 0a1.125 1.125 0 0 1-1.125-1.125M3.375 19.5h7.5c.621 0 1.125-.504 1.125-1.125m-9.75 0V5.625m0 12.75v-1.5c0-.621.504-1.125 1.125-1.125m18.375 2.625V5.625m0 12.75c0 .621-.504 1.125-1.125 1.125m1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125m0 3.75h-7.5A1.125 1.125 0 0 1 12 18.375m9.75-12.75c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125m19.5 0v1.5c0 .621-.504 1.125-1.125 1.125M2.25 5.625v1.5c0 .621.504 1.125 1.125 1.125m0 0h17.25m-17.25 0h7.5c.621 0 1.125.504 1.125 1.125M3.375 8.25c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125m17.25-3.75h-7.5c-.621 0-1.125.504-1.125 1.125m8.625-1.125c.621 0 1.125.504 1.125 1.125v1.5c0 .621-.504 1.125-1.125 1.125m-17.25 0h7.5m-7.5 0c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125M12 10.875v-1.5m0 1.5c0 .621-.504 1.125-1.125 1.125M12 10.875c0 .621.504 1.125 1.125 1.125m-2.25 0c.621 0 1.125.504 1.125 1.125M10.875 12h-7.5m0 0c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125m17.25-3.75h-7.5m0 0c.621 0 1.125.504 1.125 1.125"
          />
        </svg>
      );
    case "video":
      return (
        <svg
          className={base}
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={1.5}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="m15.75 10.5 4.72-4.72a.75.75 0 0 1 1.28.53v11.38a.75.75 0 0 1-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 0 0 2.25-2.25v-9a2.25 2.25 0 0 0-2.25-2.25h-9A2.25 2.25 0 0 0 2.25 7.5v9a2.25 2.25 0 0 0 2.25 2.25Z"
          />
        </svg>
      );
    case "text":
      return (
        <svg
          className={base}
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={1.5}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
          />
        </svg>
      );
    default:
      return (
        <svg
          className={base}
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={1.5}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
          />
        </svg>
      );
  }
}

// ---------------------------------------------------------------------------
// File type color mapping
// ---------------------------------------------------------------------------

function fileTypeColor(type: FileType): string {
  switch (type) {
    case "pdf":
      return "text-destructive";
    case "image":
      return "text-primary";
    case "spreadsheet":
      return "text-success";
    case "video":
      return "text-warning";
    case "text":
      return "text-muted-foreground";
    default:
      return "text-muted-foreground";
  }
}

// ---------------------------------------------------------------------------
// Mock data
// ---------------------------------------------------------------------------

const now = new Date();

function daysAgo(days: number): Date {
  return new Date(now.getTime() - days * 86_400_000);
}

const INITIAL_FILES: MockFile[] = [
  {
    id: "f1",
    name: "presentation.pdf",
    size: 2_516_582,
    type: "pdf",
    uploadedAt: daysAgo(2),
  },
  {
    id: "f2",
    name: "screenshot.png",
    size: 867_328,
    type: "image",
    uploadedAt: daysAgo(5),
  },
  {
    id: "f3",
    name: "report.xlsx",
    size: 1_153_434,
    type: "spreadsheet",
    uploadedAt: daysAgo(7),
  },
  {
    id: "f4",
    name: "avatar.jpg",
    size: 126_976,
    type: "image",
    uploadedAt: daysAgo(14),
  },
  {
    id: "f5",
    name: "notes.txt",
    size: 12_288,
    type: "text",
    uploadedAt: daysAgo(21),
  },
  {
    id: "f6",
    name: "demo-video.mp4",
    size: 47_395_635,
    type: "video",
    uploadedAt: daysAgo(30),
  },
];

const TOTAL_STORAGE = 10 * 1024 * 1024 * 1024; // 10 GB
const USED_STORAGE = 2.3 * 1024 * 1024 * 1024; // 2.3 GB

// ---------------------------------------------------------------------------
// Main component
// ---------------------------------------------------------------------------

export function StorageDemo() {
  const [files, setFiles] = useState<MockFile[]>(INITIAL_FILES);
  const [uploading, setUploading] = useState<UploadingFile[]>([]);
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");
  const [isDragOver, setIsDragOver] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // -------------------------------------------------------------------------
  // Upload simulation
  // -------------------------------------------------------------------------

  const simulateUpload = useCallback(
    (fileList: { name: string; size: number }[]) => {
      const newUploads: UploadingFile[] = fileList.map((f) => ({
        id: generateId(),
        name: f.name,
        size: f.size,
        type: inferFileType(f.name),
        progress: 0,
      }));

      setUploading((prev) => [...prev, ...newUploads]);

      newUploads.forEach((upload) => {
        const steps = 20;
        const interval = 2000 / steps;
        let step = 0;

        const timer = setInterval(() => {
          step++;
          const progress = Math.min(100, Math.round((step / steps) * 100));

          setUploading((prev) =>
            prev.map((u) => (u.id === upload.id ? { ...u, progress } : u))
          );

          if (step >= steps) {
            clearInterval(timer);

            // Move from uploading to files
            setTimeout(() => {
              setUploading((prev) => prev.filter((u) => u.id !== upload.id));
              setFiles((prev) => [
                {
                  id: upload.id,
                  name: upload.name,
                  size: upload.size,
                  type: upload.type,
                  uploadedAt: new Date(),
                },
                ...prev,
              ]);
            }, 300);
          }
        }, interval);
      });
    },
    []
  );

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      setIsDragOver(false);

      const droppedFiles = Array.from(e.dataTransfer.files).map((f) => ({
        name: f.name,
        size: f.size,
      }));

      if (droppedFiles.length > 0) {
        simulateUpload(droppedFiles);
      }
    },
    [simulateUpload]
  );

  const handleFileSelect = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const selected = e.target.files;
      if (!selected) return;

      const fileList = Array.from(selected).map((f) => ({
        name: f.name,
        size: f.size,
      }));

      if (fileList.length > 0) {
        simulateUpload(fileList);
      }

      // Reset input so the same file can be selected again
      e.target.value = "";
    },
    [simulateUpload]
  );

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);
  }, []);

  const handleDelete = useCallback((id: string) => {
    setFiles((prev) => prev.filter((f) => f.id !== id));
  }, []);

  const handleDownload = useCallback((fileName: string) => {
    // Simulated download — just show an alert in demo mode
    alert(`Download simulated for "${fileName}"`);
  }, []);

  // -------------------------------------------------------------------------
  // Render
  // -------------------------------------------------------------------------

  const usedPercent = (USED_STORAGE / TOTAL_STORAGE) * 100;

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 max-w-5xl">
      {/* Back link */}
      <Link
        href="/demo"
        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-6"
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
            d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"
          />
        </svg>
        Back to demos
      </Link>

      {/* Header */}
      <div className="mb-8">
        <span className="inline-block bg-primary/10 text-primary text-sm font-medium px-3 py-1 rounded-full mb-3">
          File Storage
        </span>
        <h1 className="text-3xl sm:text-4xl font-bold text-foreground tracking-tight">
          File Storage Demo
        </h1>
        <p className="text-muted-foreground mt-2 max-w-2xl">
          Upload, manage, and preview files with Cloudflare R2. This demo runs
          entirely in your browser with simulated data.
        </p>
      </div>

      {/* Storage usage bar */}
      <div className="bg-card border border-border rounded-xl p-4 sm:p-6 mb-6">
        <div className="flex items-center justify-between mb-2">
          <h2 className="text-sm font-medium text-foreground">
            Storage Usage
          </h2>
          <span className="text-sm text-muted-foreground">
            {formatFileSize(USED_STORAGE)} of {formatFileSize(TOTAL_STORAGE)}{" "}
            used
          </span>
        </div>
        <div className="h-2.5 bg-muted rounded-full overflow-hidden">
          <div
            className="h-full bg-primary rounded-full transition-all duration-500"
            style={{ width: `${usedPercent}%` }}
          />
        </div>
        <p className="text-xs text-muted-foreground mt-1.5">
          {usedPercent.toFixed(0)}% of your storage is in use
        </p>
      </div>

      {/* Drop zone */}
      <div
        onDrop={handleDrop}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        className={cn(
          "relative border-2 border-dashed rounded-xl p-8 sm:p-12 text-center transition-colors duration-200 mb-6",
          isDragOver
            ? "border-primary bg-primary/5"
            : "border-border hover:border-muted-foreground/50"
        )}
      >
        <input
          ref={fileInputRef}
          type="file"
          multiple
          onChange={handleFileSelect}
          className="hidden"
          aria-label="File upload input"
        />

        <div className="flex flex-col items-center gap-3">
          <div
            className={cn(
              "h-12 w-12 rounded-lg flex items-center justify-center transition-colors",
              isDragOver
                ? "bg-primary/10 text-primary"
                : "bg-muted text-muted-foreground"
            )}
          >
            <svg
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={1.5}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5"
              />
            </svg>
          </div>
          <div>
            <p className="text-sm font-medium text-foreground">
              {isDragOver
                ? "Drop files here..."
                : "Drag and drop files here"}
            </p>
            <p className="text-xs text-muted-foreground mt-1">
              or use the button below to browse
            </p>
          </div>
          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            className="mt-1 inline-flex items-center gap-2 bg-primary text-primary-foreground text-sm font-medium px-4 py-2 rounded-lg hover:opacity-90 transition-opacity focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
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
                d="M12 10.5v6m3-3H9m4.06-7.19l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z"
              />
            </svg>
            Browse files
          </button>
        </div>
      </div>

      {/* Uploading files */}
      {uploading.length > 0 && (
        <div className="space-y-3 mb-6">
          {uploading.map((file) => (
            <div
              key={file.id}
              className="bg-card border border-border rounded-xl p-4 flex items-center gap-4"
            >
              <div
                className={cn(
                  "h-10 w-10 rounded-lg bg-muted flex items-center justify-center",
                  fileTypeColor(file.type)
                )}
              >
                <FileTypeIcon type={file.type} className="h-5 w-5" />
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-foreground truncate">
                  {file.name}
                </p>
                <div className="flex items-center gap-3 mt-1.5">
                  <div className="flex-1 h-1.5 bg-muted rounded-full overflow-hidden">
                    <div
                      className="h-full bg-primary rounded-full transition-all duration-100"
                      style={{ width: `${file.progress}%` }}
                    />
                  </div>
                  <span className="text-xs text-muted-foreground tabular-nums w-8 text-right">
                    {file.progress}%
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Toolbar: view toggle + file count */}
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm text-muted-foreground">
          {files.length} {files.length === 1 ? "file" : "files"}
        </p>
        <div className="flex items-center gap-1 bg-muted rounded-lg p-0.5">
          <button
            type="button"
            onClick={() => setViewMode("grid")}
            className={cn(
              "p-1.5 rounded-md transition-colors",
              viewMode === "grid"
                ? "bg-card text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            )}
            aria-label="Grid view"
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
                d="M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z"
              />
            </svg>
          </button>
          <button
            type="button"
            onClick={() => setViewMode("list")}
            className={cn(
              "p-1.5 rounded-md transition-colors",
              viewMode === "list"
                ? "bg-card text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            )}
            aria-label="List view"
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
                d="M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.625 4.5h12.75a1.875 1.875 0 010 3.75H5.625a1.875 1.875 0 010-3.75z"
              />
            </svg>
          </button>
        </div>
      </div>

      {/* File grid */}
      {viewMode === "grid" && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {files.map((file) => (
            <div
              key={file.id}
              className="bg-card border border-border rounded-xl p-4 group hover:shadow-md hover:border-primary/30 transition-all duration-200"
            >
              {/* Icon + badges */}
              <div className="flex items-start justify-between mb-3">
                <div
                  className={cn(
                    "h-10 w-10 rounded-lg bg-muted flex items-center justify-center",
                    fileTypeColor(file.type)
                  )}
                >
                  <FileTypeIcon type={file.type} className="h-5 w-5" />
                </div>
                {file.type === "image" && (
                  <span className="text-xs bg-primary/10 text-primary font-medium px-2 py-0.5 rounded-md">
                    Preview
                  </span>
                )}
              </div>

              {/* File info */}
              <p className="text-sm font-medium text-foreground truncate">
                {file.name}
              </p>
              <div className="flex items-center gap-2 mt-1">
                <span className="text-xs text-muted-foreground">
                  {formatFileSize(file.size)}
                </span>
                <span className="text-muted-foreground/40">·</span>
                <span className="text-xs text-muted-foreground">
                  {formatRelativeTime(file.uploadedAt)}
                </span>
              </div>

              {/* Actions */}
              <div className="flex items-center gap-2 mt-3 pt-3 border-t border-border">
                <button
                  type="button"
                  onClick={() => handleDownload(file.name)}
                  className="inline-flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
                >
                  <svg
                    className="h-3.5 w-3.5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    strokeWidth={2}
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3"
                    />
                  </svg>
                  Download
                </button>
                <button
                  type="button"
                  onClick={() => handleDelete(file.id)}
                  className="inline-flex items-center gap-1.5 text-xs text-muted-foreground hover:text-destructive transition-colors ml-auto"
                >
                  <svg
                    className="h-3.5 w-3.5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    strokeWidth={2}
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"
                    />
                  </svg>
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* File list (table) */}
      {viewMode === "list" && (
        <div className="bg-card border border-border rounded-xl overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="border-b border-border">
                <th className="text-left text-xs font-medium text-muted-foreground px-4 py-3">
                  Name
                </th>
                <th className="text-left text-xs font-medium text-muted-foreground px-4 py-3 hidden sm:table-cell">
                  Size
                </th>
                <th className="text-left text-xs font-medium text-muted-foreground px-4 py-3 hidden md:table-cell">
                  Uploaded
                </th>
                <th className="text-right text-xs font-medium text-muted-foreground px-4 py-3">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody>
              {files.map((file, index) => (
                <tr
                  key={file.id}
                  className={cn(
                    "hover:bg-muted/50 transition-colors",
                    index !== files.length - 1 && "border-b border-border"
                  )}
                >
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-3">
                      <div
                        className={cn(
                          "h-8 w-8 rounded-lg bg-muted flex items-center justify-center shrink-0",
                          fileTypeColor(file.type)
                        )}
                      >
                        <FileTypeIcon
                          type={file.type}
                          className="h-4 w-4"
                        />
                      </div>
                      <div className="min-w-0">
                        <p className="text-sm font-medium text-foreground truncate">
                          {file.name}
                        </p>
                        <div className="flex items-center gap-2 sm:hidden">
                          <span className="text-xs text-muted-foreground">
                            {formatFileSize(file.size)}
                          </span>
                          {file.type === "image" && (
                            <span className="text-xs bg-primary/10 text-primary font-medium px-1.5 py-0.5 rounded">
                              Preview
                            </span>
                          )}
                        </div>
                      </div>
                      {file.type === "image" && (
                        <span className="hidden sm:inline-block text-xs bg-primary/10 text-primary font-medium px-2 py-0.5 rounded-md">
                          Preview
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="px-4 py-3 hidden sm:table-cell">
                    <span className="text-sm text-muted-foreground">
                      {formatFileSize(file.size)}
                    </span>
                  </td>
                  <td className="px-4 py-3 hidden md:table-cell">
                    <span className="text-sm text-muted-foreground">
                      {formatRelativeTime(file.uploadedAt)}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        type="button"
                        onClick={() => handleDownload(file.name)}
                        className="p-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
                        aria-label={`Download ${file.name}`}
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
                            d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3"
                          />
                        </svg>
                      </button>
                      <button
                        type="button"
                        onClick={() => handleDelete(file.id)}
                        className="p-1.5 rounded-md text-muted-foreground hover:text-destructive hover:bg-muted transition-colors"
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
                            d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"
                          />
                        </svg>
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Empty state */}
      {files.length === 0 && uploading.length === 0 && (
        <div className="text-center py-16">
          <div className="h-12 w-12 rounded-lg bg-muted text-muted-foreground flex items-center justify-center mx-auto mb-4">
            <svg
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={1.5}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M2.25 12.75V12A2.25 2.25 0 014.5 9.75h15A2.25 2.25 0 0121.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z"
              />
            </svg>
          </div>
          <p className="text-sm font-medium text-foreground">
            No files uploaded
          </p>
          <p className="text-xs text-muted-foreground mt-1">
            Drop files above or click &quot;Browse files&quot; to get started.
          </p>
        </div>
      )}

      {/* Footer note */}
      <p className="text-center text-sm text-muted-foreground mt-10">
        This demo runs entirely in your browser. In production, files are stored
        in Cloudflare R2 via the Go backend.
      </p>
    </div>
  );
}
