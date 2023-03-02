import React, { useState } from "react";
import { EditOutlined } from "@ant-design/icons";
import { Button, Descriptions } from "antd";
import { NavTitle } from "../../../../components/NavTitle";
import { User } from "../../../../types";

import styles from "./ProfileView.module.css";
import { ChangePasswordForm } from "../ChangePasswordForm";

interface Props {
  user: User;
  logout: () => void;
}

export function ProfileView({ user, logout }: Props) {
  const [changePasswordOpen, setChangePasswordOpen] = useState<Boolean>(false);
  const chpwd = () => {};

  return (
    <div className={styles.Content}>
      <NavTitle title="Profil" />
      <Descriptions>
        <Descriptions.Item label="Namn">{user.name}</Descriptions.Item>
        <Descriptions.Item label="Email">{user.email}</Descriptions.Item>
      </Descriptions>
      {!changePasswordOpen && (
        <Button
          type="default"
          size="large"
          block
          icon={<EditOutlined />}
          onClick={() => setChangePasswordOpen(true)}
        >
          Ändra lösenord
        </Button>
      )}
      {changePasswordOpen && <ChangePasswordForm submit={chpwd} />}
      <Button
        type="text"
        size="large"
        block
        onClick={logout}
        className={styles.LogoutButtton}
      >
        Logga ut
      </Button>
    </div>
  );
}
