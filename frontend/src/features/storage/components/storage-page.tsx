"use client";

import { useStorage } from "@/features/storage/hooks/use-storage";
import { FileUpload } from "@/features/storage/components/file-upload";
import { FileList } from "@/features/storage/components/file-list";

export function StoragePage() {
  const {
    files,
    total,
    page,
    totalPages,
    loading,
    uploading,
    error,
    hasNext,
    hasPrev,
    uploadFile,
    deleteFile,
    goToNextPage,
    goToPrevPage,
  } = useStorage();

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Files</h1>
        <p className="text-muted-foreground mt-1">
          Upload and manage your files.
        </p>
      </div>

      {error && (
        <div className="bg-destructive/10 border border-destructive/20 rounded-lg px-4 py-3">
          <p className="text-sm text-destructive">{error}</p>
        </div>
      )}

      <FileUpload onUpload={uploadFile} uploading={uploading} />

      <FileList
        files={files}
        total={total}
        page={page}
        totalPages={totalPages}
        loading={loading}
        hasNext={hasNext}
        hasPrev={hasPrev}
        onDelete={deleteFile}
        onNextPage={goToNextPage}
        onPrevPage={goToPrevPage}
      />
    </div>
  );
}
