import switchDarkMode from "./loaders/darkmode.ts";
switchDarkMode(true)

import { lazy, StrictMode, Suspense } from "react";
import { createRoot } from "react-dom/client";
import "reset-css/reset.css";
import "./styles/index.css";
import {createBrowserRouter, redirect, RouterProvider} from "react-router-dom";
import { Provider } from "react-redux";
import store from "./stores";
import Loading from "./pages/Loading.tsx";
const DashboardLayout = lazy(() => import("./pages/DashboardLayout.tsx"));
const Files = lazy(() => import("./pages/files/Files.tsx"));
const Error = lazy(() => import("./pages/error/Error.tsx"));
const Login = lazy(() => import("./pages/login/Login.tsx"));
const Setup = lazy(() => import("./pages/setup/Setup.tsx"));
const Maintenance = lazy(() => import("./pages/setup/Maintenance.tsx"));
const StandaloneLayout = lazy(() => import("./pages/StandaloneLayout.tsx"));
const Likes = lazy(() => import("./pages/likes/Likes.tsx"));
const Shares = lazy(() => import("./pages/shares/Shares.tsx"));
const GroupFolders = lazy(() => import("./pages/groups/GroupFolders.tsx"));
import { logout } from "./api";
import setupLoader from "./loaders/setupLoader.ts";
import filesLoader from "./loaders/filesLoader.ts";

const router = createBrowserRouter([
  {
    path: "/",
    element: <DashboardLayout />,
    errorElement: <Error />,
    children: [
      {
        path: "/",
        loader: () => redirect("/files/")
      },
      {
        path: "/likes",
        element: <Likes />,
      },
      {
        path: "/shares",
        element: <Shares />
      },
      {
        path: "/groups",
        element: <GroupFolders />
      },
      {
        path: "/files/*",
        element: <Files />,
        loader: filesLoader
      },
    ],
  },
  {
    path: "/",
    element: <StandaloneLayout />,
    children: [
      {
        path: "/login",
        element: <Login />,
      },
      {
        path: "/logout",
        action: logout,
      },
      {
        path: "/setup",
        element: <Setup />,
        loader: setupLoader,
      },
      {
        path: "/maintenance",
        element: <Maintenance />,
      },
    ],
  },
]);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Provider store={store}>
      <Suspense fallback={<Loading />}>
        <RouterProvider router={router} fallbackElement={<Loading />} />
      </Suspense>
    </Provider>
  </StrictMode>,
);
