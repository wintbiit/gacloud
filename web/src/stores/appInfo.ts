import { createSlice } from "@reduxjs/toolkit";
import { AppInfo, DefaultAppInfo, getAppInfo } from "../api/app.ts";

export const appInfoSlice = createSlice({
  name: "appInfo",
  initialState: {
    appInfo: DefaultAppInfo,
  },
  reducers: {
    updateAppInfo: (state) => {
      getAppInfo().then((appInfo: AppInfo) => {
        state.appInfo = appInfo;
      });
    },
  },
});

export const { updateAppInfo } = appInfoSlice.actions;

export default appInfoSlice.reducer;
