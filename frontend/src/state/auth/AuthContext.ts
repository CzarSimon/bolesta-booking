import { createContext } from "react";
import { User } from "../../types";

export interface AuthState {
  authenticated: boolean;
  user?: User;
  authenticate: (user: User) => void;
  logout: () => void;
}

const initalState: AuthState = {
  authenticated: false,
  authenticate: (user: User) => {}, // eslint-disable-line
  logout: () => {}, // eslint-disable-line
};

export const AuthContext = createContext<AuthState>(initalState);
