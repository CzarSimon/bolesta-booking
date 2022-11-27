import React from "react";
import { Optional } from "../../types";

import styles from "./ErrorText.module.css";

interface Props {
  error: Optional<string>;
}

export function ErrorText({ error }: Props) {
  return error ? <p className={styles.Text}>{error}</p> : null;
}
