import { DefaultServerInfo, getServerInfo, ServerInfo } from "../api/app.ts";
import { createSlice } from "@reduxjs/toolkit";

export const serverInfoSlice = createSlice({
  name: "serverInfo",
  initialState: {
    serverInfo: DefaultServerInfo,
  },
  reducers: {
    updateServerInfo: (state) => {
      getServerInfo().then((serverInfo: ServerInfo) => {
        state.serverInfo = serverInfo;
      });
    },
  },
});

export const { updateServerInfo } = serverInfoSlice.actions;

export default serverInfoSlice.reducer;
