import { axioser } from "./index.ts";

export interface MySQLOptionsProps {
  type: string;
  host: string;
  port: number;
  username: string;
  password: string;
  database: string;
}

export interface PostgreSQLOptionsProps {
  type: string;
  host: string;
  port: number;
  username: string;
  password: string;
  database: string;
}

export interface SQLiteOptionsProps {
  type: string;
  path: string;
}

interface SetupStatus {
  currentStep: number;
}

async function getSetupStatus(): Promise<number> {
  return axioser.get<SetupStatus>("/setup").then((res) => res.data.currentStep);
}

const setupDatabase = async (
  values: MySQLOptionsProps | PostgreSQLOptionsProps | SQLiteOptionsProps,
) => {
  await axioser.post("/setup/database", {
    type: values.type,
    params: values,
  });
};

interface TestResponse {
  success: boolean;
  reason: string;
}

async function testDatabase(
  values: MySQLOptionsProps | PostgreSQLOptionsProps | SQLiteOptionsProps,
): Promise<TestResponse> {
  return (
    await axioser.post<TestResponse>("/setup/database/test", {
      type: values.type,
      params: values,
    })
  ).data;
}

export interface ElasticsearchOptionsProps {
  host: string;
  user: string;
  password: string;
}

const setupElasticsearch = async (values: ElasticsearchOptionsProps) => {
  await axioser.post("/setup/elasticsearch", values);
};

const testElasticsearch = async (
  values: ElasticsearchOptionsProps,
): Promise<TestResponse> => {
  return (await axioser.post<TestResponse>("/setup/elasticsearch/test", values))
    .data;
};

export interface StorageOptionsProps {
  name: string;
  type: string;
  credential: any;
}

export interface StorageProviderConfig {
  rootPath: string;
  providers: string[];
}

const setupStorage = async (values: StorageOptionsProps) => {
  await axioser.post("/setup/storage", {
    name: values.name,
    type: values.type,
    credential: JSON.stringify(values.credential),
  });
};

const testStorage = async (
  values: StorageOptionsProps,
): Promise<TestResponse> => {
  return (
    await axioser.post<TestResponse>("/setup/storage/test", {
      name: values.name,
      type: values.type,
      credential: JSON.stringify(values.credential),
    })
  ).data;
};

const getStorageProviders = async (): Promise<StorageProviderConfig> => {
  return (await axioser.get<StorageProviderConfig>("/setup/storage/providers"))
    .data;
};

export interface SetAdminProps {
  username: string;
  password: string;
  email: string;
}

const setAdmin = async (values: SetAdminProps) => {
  await axioser.post("/setup/admin", values);
};

const finishSetup = async () => {
  await axioser.post("/setup/finish");
};

export {
  getSetupStatus,
  setupDatabase,
  testDatabase,
  setupElasticsearch,
  testElasticsearch,
  setupStorage,
  testStorage,
  getStorageProviders,
  setAdmin,
  finishSetup,
};
