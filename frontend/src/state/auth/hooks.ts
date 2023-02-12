import { useContext } from "react";
import { useNavigate } from "react-router-dom";
import { AuthenticatedResponse, LoginRequest, User } from "../../types";
import { AuthContext } from "./AuthContext";
import { requestLogin } from "../../api";
import { AUTH_TOKEN_KEY, CURRENT_USER_KEY } from "../../constants";
import { removeHeader, setHeader } from "../../api/httpclient";

interface UseAuthResult {
  login: (req: LoginRequest) => void;
  user?: User;
  authenticated: boolean;
  authenticate: (user: User) => void;
  logout: () => void;
}

export function useAuth(): UseAuthResult {
  const { user, authenticated, authenticate, logout } = useContext(AuthContext);
  const navigate = useNavigate();

  const login = async (req: LoginRequest) => {
    const res = await requestLogin(req);
    storeAuthInfo(res);
    setTimeout(() => {
      authenticate(res.user);
      navigate("/");
    }, 100);
  };

  const onLogout = () => {
    logout();
    removeAuthInfo();
    navigate("/login");
  };

  return {
    login,
    user,
    authenticated,
    authenticate,
    logout: onLogout,
  };
}

export function useIsAuthenticated(): boolean {
  const { authenticated } = useContext(AuthContext);
  return authenticated;
}

function storeAuthInfo({ user, token }: AuthenticatedResponse) {
  setHeader("Authorization", `Bearer ${token}`);
  localStorage.setItem(AUTH_TOKEN_KEY, token);
  localStorage.setItem(CURRENT_USER_KEY, JSON.stringify(user));
}

function removeAuthInfo() {
  removeHeader("Authorization");
  localStorage.removeItem(AUTH_TOKEN_KEY);
  localStorage.removeItem(CURRENT_USER_KEY);
}
