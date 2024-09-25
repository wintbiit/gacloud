import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import {createBrowserRouter, RouterProvider} from "react-router-dom";
import Layout from "./pages/Layout.tsx";
import Error from "./pages/error/Error.tsx";
import {dashboardLoader} from "./pages/dashboard/Dashboard.tsx";
import {Provider} from "react-redux";
import store from "./stores"

const router = createBrowserRouter([
    {
        path: "/",
        element: <Layout />,
        errorElement: <Error />,
        children: [
            {
                path: "/",
                action: () => import("./pages/dashboard/Dashboard.tsx"),
                loader: dashboardLoader
            },
            {
                path: "/login",
                action: () => import("./pages/login/Login.tsx"),
            }
        ]
    },
    {
        path: "/setup",
        action: () => import("./pages/setup/Setup.tsx"),
    }
])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Provider store={store}>
        <RouterProvider router={router} />
    </Provider>
  </StrictMode>,
)
