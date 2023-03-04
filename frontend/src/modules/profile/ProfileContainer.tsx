import React from "react";
import { useNavigate } from "react-router-dom";
import { changePassword } from "../../api";
import { LoadingIndicator } from "../../components/LoadingIndicator";
import { useAuth } from "../../state/auth/hooks";
import { ChangePasswordRequest } from "../../types";
import { ProfileView } from "./components/ProfileView";

export function ProfileContainer() {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  if (!user) {
    navigate("/login");
    return <LoadingIndicator />;
  }

  const requestChangePassword = (req: ChangePasswordRequest): Promise<void> => {
    return changePassword(user.id, req).then();
  };

  return (
    <ProfileView
      user={user}
      logout={logout}
      changePassword={requestChangePassword}
    />
  );
}
