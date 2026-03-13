"use client";

import { FileCard } from "@/features/storage/components/file-card";
import type { FileItem } from "@/features/storage/types";

type FileListProps = {
  files: FileItem[];
  total: number;
  page: number;
  totalPages: number;
  loading: boolean;
  hasNext: boolean;
  hasPrev: boolean;
  onDelete: (fileId: string) => Promise<{ error: string | null }>;
  onNextPage: () => void;
  onPrevPage: () => void;
};

function LoadingSkeleton() {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {[1, 2, 3, 4].map((i) => (
        <div
          key={i}
          className="bg-card border border-border rounded-xl shadow-sm animate-pulse h-64"
        />
      ))}
    </div>
  );
}

function EmptyState() {
  return (
    <div className="bg-card border border-border rounded-xl p-12 shadow-sm text-center">
      <svg
        className="mx-auto h-12 w-12 text-muted-foreground mb-4"
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
      <h3 className="text-base font-medium text-foreground mb-1">
        No files yet
      </h3>
      <p className="text-sm text-muted-foreground">
        Upload your first file using the form above.
      </p>
    </div>
  );
}

export function FileList({
  files,
  total,
  page,
  totalPages,
  loading,
  hasNext,
  hasPrev,
  onDelete,
  onNextPage,
  onPrevPage,
}: FileListProps) {
  if (loading) {
    return <LoadingSkeleton />;
  }

  if (files.length === 0 && page === 1) {
    return <EmptyState />;
  }

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {files.map((file) => (
          <FileCard key={file.id} file={file} onDelete={onDelete} />
        ))}
      </div>

      {/* Pagination */}
      {(hasPrev || hasNext) && (
        <div className="flex items-center justify-between pt-2">
          <button
            onClick={onPrevPage}
            disabled={!hasPrev}
            className="text-sm text-muted-foreground hover:text-foreground transition-colors disabled:opacity-50"
          >
            Previous
          </button>
          <span className="text-xs text-muted-foreground">
            Page {page} of {totalPages} &middot; {total} file
            {total !== 1 ? "s" : ""}
          </span>
          <button
            onClick={onNextPage}
            disabled={!hasNext}
            className="text-sm text-muted-foreground hover:text-foreground transition-colors disabled:opacity-50"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}
