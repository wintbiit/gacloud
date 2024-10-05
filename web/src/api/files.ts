import { axioser } from "./index.ts";

export interface File {
  path: string;
  size: number;
  mime: string;
  sum: string;
  providerId: string;
  createdAt: string;
  updatedAt: string;
}

export type ListFilesResponse = File[];

const listFiles = async (path: string) => {
  const resp = await axioser.get<ListFilesResponse>(`/files/list`, {
    params: {
      path,
    },
  });
  return resp.data;
};

export { listFiles };
