import React from "react";
import { useNavigate } from "react-router-dom";
import { LoadingIndicator } from "../../components/LoadingIndicator";
import { useAuth } from "../../state/auth/hooks";
import { ProfileView } from "./components/ProfileView";

export function ProfileContainer() {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  if (!user) {
    navigate("/login");
    return <LoadingIndicator />;
  }

  return <ProfileView user={user} logout={logout} />;
}
