import React from "react";
import { Cabin } from "../../../../types";
import { CabinListItem } from "../CabinListItem";

import styles from "./StartPage.module.css";

interface Props {
  cabins: Cabin[];
}

export function StartPage({ cabins }: Props) {
  return (
    <div className={styles.StartPage}>
      <h1 className={styles.Title}>Bölesta booking</h1>
      <div className={styles.ListTable}>
        <h2 className={styles.ListTitle}>Stugor</h2> 
        <ul className={styles.CabinList}>
          {cabins.map((cabin) => (
            <li key={cabin.id}>
              <CabinListItem cabin={cabin} />
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
