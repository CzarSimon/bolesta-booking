import React from "react";
import { useAuth } from "../../state/auth/hooks";
import { LoginRequest } from "../../types";
import { Login } from "./components/Login";

export function LoginContainer() {
  const { login } = useAuth();
  return <Login submit={login} />;
}
