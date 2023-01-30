import React, { useState } from "react";
import { AuthContext } from "./AuthContext";
import { Optional, User } from "../../types";

interface Props {
  children: JSX.Element;
}

export function AuthProvider({ children }: Props) {
  const [authenticated, setAuthenticated] = useState<boolean>(false);
  const [user, setUser] = useState<Optional<User>>(undefined);

  const authenticate = (user: User) => {
    setAuthenticated(true);
    setUser(user);
  };

  const logout = () => {
    setAuthenticated(false);
    setUser(undefined);
  };

  return (
    <AuthContext.Provider value={{ authenticated, user, authenticate, logout }}>
      {children}
    </AuthContext.Provider>
  );
}
