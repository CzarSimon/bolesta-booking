import React from "react";
import { useNavigate } from "react-router-dom";
import { Cabin } from "../../../../types";

import styles from "./CabinListItem.module.css";

interface Props {
  cabin: Cabin;
}

export function CabinListItem({ cabin }: Props) {
  const navigate = useNavigate();
  const select = () => {
    navigate(`cabins/${cabin.id}`);
  };

  return (
    <div className={styles.CabinListItem} onClick={select}>
      <h2>{cabin.name}</h2>
    </div>
  );
}
