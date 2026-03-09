export type ApiError = {
  error: string;
};

export type PaginatedResponse<T> = {
  data: T[];
  total: number;
  offset: number;
  limit: number;
};
