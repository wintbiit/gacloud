import { axioser } from "./index.ts";

export interface AppInfo {
  site_name: string;
  external_url: string;
  site_logo: string;
}

export const DefaultAppInfo: AppInfo = {
  site_name: "GaCloud",
  external_url: "http://localhost:8080",
  site_logo: "",
};

export interface ServerInfo {
  version: string;
  build_revision: string;
  build_time: string;
  go_version: string;
  data_dir: string;
  log_dir: string;
  addr: string;
}

export const DefaultServerInfo: ServerInfo = {
  version: "0.0.0",
  build_revision: "0000000",
  build_time: "0000-00-00T00:00:00Z",
  go_version: "0.0.0",
  data_dir: "",
  log_dir: "",
  addr: "",
};

export const getAppInfo = async (): Promise<AppInfo> => {
  const resp = await axioser.get<AppInfo>("/appinfo");
  return resp.data;
};

export const getServerInfo = async (): Promise<ServerInfo> => {
  const resp = await axioser.get<ServerInfo>("/serverinfo");

  return resp.data;
};
