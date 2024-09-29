import { Outlet } from "react-router-dom";

function DashboardLayout() {
  return (
    <div>
      <h1>DashboardLayout</h1>
      <Outlet />
    </div>
  );
}

export default DashboardLayout;
