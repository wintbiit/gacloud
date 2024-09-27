import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import "reset-css/reset.css";
import './styles/index.css'
import {createBrowserRouter, RouterProvider} from "react-router-dom";
import {Provider} from "react-redux";
import store from "./stores"
import Error from "./pages/error/Error.tsx";
import Dashboard from "./pages/dashboard/Dashboard.tsx";
import Login from "./pages/login/Login.tsx";
import Setup from "./pages/setup/Setup.tsx";
import Maintenance from "./pages/setup/Maintenance.tsx";
import DashboardLayout from "./pages/DashboardLayout.tsx";
import StandaloneLayout from "./pages/StandaloneLayout.tsx";
import {logout} from "./api";

const router = createBrowserRouter([
    {
        path: "/",
        element: <DashboardLayout />,
        errorElement: <Error />,
        children: [
            {
                path: "/",
                element: <Dashboard />,
            },
        ]
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
                action: logout
            },
            {
                path: "/setup",
                element: <Setup />,
            },
            {
                path: "/maintenance",
                element: <Maintenance />,
            }
        ]
    },
])

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <Provider store={store}>
            <RouterProvider router={router} />
        </Provider>
    </StrictMode>
)
