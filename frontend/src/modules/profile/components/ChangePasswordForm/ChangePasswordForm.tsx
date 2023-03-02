import React, { ChangeEvent, useState } from "react";
import { Button, Input } from "antd";
import { ErrorText } from "../../../../components/ErrorText";
import {
  ChangePasswordRequest,
  Failure,
  Optional,
  Result,
  Success,
} from "../../../../types";

import styles from "./ChangePasswordForm.module.css";

interface Props {
  submit: (req: ChangePasswordRequest) => void;
}

export function ChangePasswordForm({ submit }: Props) {
  const [oldPassword, setOldPassword] = useState<Optional<string>>();
  const [newPassword, setNewPassword] = useState<Optional<string>>();
  const [confirmPassword, setConfirmPassword] = useState<Optional<string>>();
  const [err, setErr] = useState<Optional<string>>();

  const updateOldPassword = (e: ChangeEvent<HTMLInputElement>) => {
    setOldPassword(e.target.value);
  };

  const updateNewPassword = (e: ChangeEvent<HTMLInputElement>) => {
    setNewPassword(e.target.value);
  };

  const updateConfirmPassword = (e: ChangeEvent<HTMLInputElement>) => {
    setConfirmPassword(e.target.value);
  };

  const onSubmit = () => {
    parseChangeRequest(oldPassword, newPassword, confirmPassword)
      .then((req) => submit(req))
      .catch((e) => setErr(e));
  };

  return (
    <div>
      <Input.Password
        placeholder="Nuvarande lösenord"
        onChange={updateOldPassword}
        size="large"
        className={styles.OldPassword}
      />
      <Input.Password
        placeholder="Nytt lösenord"
        onChange={updateNewPassword}
        size="large"
        className={styles.NewPassword}
      />
      <Input.Password
        placeholder="Bekräfta lösenord"
        onChange={updateConfirmPassword}
        size="large"
        className={styles.NewPassword}
      />
      <Button block type="primary" onClick={onSubmit} size="large">
        Ändra lösenord
      </Button>
      <ErrorText error={err} />
    </div>
  );
}

function parseChangeRequest(
  oldPassword?: string,
  newPassword?: string,
  confirmPassword?: string
): Result<ChangePasswordRequest, string> {
  if (!oldPassword || !newPassword || !confirmPassword) {
    return Failure("Alla fält måste fyllas i");
  }

  if (newPassword !== confirmPassword) {
    return Failure("Lösenorden stämmer inte överens");
  }

  return Success({
    oldPassword,
    newPassword,
  });
}
