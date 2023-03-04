import React, { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../state/auth/hooks";
import { Login } from "./components/Login";

export function LoginContainer() {
  const navigate = useNavigate();
  const { login, authenticated } = useAuth();

  useEffect(() => {
    if (authenticated) {
      navigate("/");
    }
  }, [authenticated]); // eslint-disable-line react-hooks/exhaustive-deps
  return <Login submit={login} />;
}
