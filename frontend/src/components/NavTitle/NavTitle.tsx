import React from "react";
import { LeftOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";

import styles from "./NavTitle.module.css";

interface Props {
  title: string;
}

export function NavTitle({ title }: Props) {
  const navigate = useNavigate();
  const goBack = () => {
    navigate(-1);
  };

  return (
    <div className={styles.Content}>
      <LeftOutlined onClick={goBack} className={styles.BackButton} />
      <h1>{title}</h1>
    </div>
  );
}
