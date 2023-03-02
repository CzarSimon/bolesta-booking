import React from "react";
import { Button } from "antd";
import { User } from "../../../../types";
import { useNavigate } from "react-router-dom";

import styles from "./ProfileLink.module.css";

interface Props {
  user: User;
}

export function ProfileLink({ user }: Props) {
  const navigate = useNavigate();
  const goToProfile = () => {
    navigate("profile");
  };

  return (
    <Button type="text" block className={styles.Button} onClick={goToProfile}>
      Inloggad som: {user.name}
    </Button>
  );
}
