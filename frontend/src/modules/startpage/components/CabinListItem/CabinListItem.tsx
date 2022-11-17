import React from "react";
import { Cabin } from "../../../../types";

import styles from "./CabinListItem.module.css";

interface Props {
  cabin: Cabin;
}

export function CabinListItem({ cabin }: Props) {
  const select = () => {
    alert(`You selected ${cabin.name}`);
  };

  return (
    <div className={styles.CabinListItem} onClick={select}>
      <h2>{cabin.name}</h2>
    </div>
  );
}
