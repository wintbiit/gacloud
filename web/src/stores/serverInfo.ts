import {getServerInfo, ServerInfo} from "../api/app.ts";
import {createSlice} from "@reduxjs/toolkit";

const serverInfo = await getServerInfo();

export const serverInfoSlice = createSlice({
    name: 'serverInfo',
    initialState: serverInfo as ServerInfo,
    reducers: {}
});

export const {} = serverInfoSlice.actions;

export default serverInfoSlice.reducer;