import {createSlice} from "@reduxjs/toolkit";
import {AppInfo, DefaultAppInfo, getAppInfo} from "../api/app.ts";

export const appInfoSlice = createSlice({
    name: 'appInfo',
    initialState: DefaultAppInfo,
    reducers: {
    },
});

export const {
    updateAppInfo
} = appInfoSlice.actions;

export default appInfoSlice.reducer;