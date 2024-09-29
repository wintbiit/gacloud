import {
  getSetupStatus,
  getStorageProviders,
  StorageProviderConfig,
} from "../../api/setup.ts";
import { defer, redirect } from "react-router-dom";

const setupLoader = async () => {
  // const status = await getSetupStatus();
  // const storageProviders = await getStorageProviders();
  let status = 0;
  let storageProviders: StorageProviderConfig = {
    providers: [],
    rootPath: "",
  };

  await Promise.all([
    getSetupStatus().then((res) => (status = res)),
    getStorageProviders().then((res) => (storageProviders = res)),
  ]).catch(() => {
    status = 5;
  });

  if (status > 4) {
    return redirect("/login");
  }

  return defer({ status, storageProviders });
};

export default setupLoader;
