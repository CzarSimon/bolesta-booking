import React, { useEffect } from "react";
import { Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../state/auth/hooks";
import { readUser } from "../init";

export function ProtectedRoutes() {
  const navigate = useNavigate();
  const { authenticated } = useAuth();

  useEffect(() => {
    const user = readUser();
    if (!user && !authenticated) {
      navigate("/login");
    }
  }, [authenticated]);

  return <Outlet />;
}
