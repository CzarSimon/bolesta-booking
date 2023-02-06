import React from "react";
import { Select } from "antd";
import { User } from "../../types";

interface Props {
  placeholder: string;
  users: User[];
  onChange: (id: string) => void;
}

interface Option {
  label: string;
  value: string;
}

export function UserSelect({ placeholder, users, onChange }: Props) {
  return (
    <Select
      placeholder={placeholder}
      options={mapUsersToOptions(users)}
      onChange={onChange}
    />
  );
}

function mapUsersToOptions(users: User[]): Option[] {
  return users.map((u) => ({
    label: u.name,
    value: u.id,
  }));
}
