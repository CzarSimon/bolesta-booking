import React, { useState } from "react";
import { EditOutlined } from "@ant-design/icons";
import { Button, Descriptions } from "antd";
import { NavTitle } from "../../../../components/NavTitle";
import { ChangePasswordRequest, User } from "../../../../types";
import { ChangePasswordForm } from "../ChangePasswordForm";

import styles from "./ProfileView.module.css";

interface Props {
  user: User;
  logout: () => void;
  changePassword: (req: ChangePasswordRequest) => Promise<void>;
}

export function ProfileView({ user, logout, changePassword }: Props) {
  const [changePasswordOpen, setChangePasswordOpen] = useState<Boolean>(false);

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
      {changePasswordOpen && (
        <ChangePasswordForm
          submit={changePassword}
          close={() => setChangePasswordOpen(false)}
        />
      )}
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
