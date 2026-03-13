export type FileItem = {
  id: string;
  name: string;
  key: string;
  size_bytes: number;
  content_type: string;
  url: string;
  created_at: string;
  updated_at: string;
};

export type FilesResponse = {
  files: FileItem[];
  total: number;
  page: number;
  limit: number;
};
