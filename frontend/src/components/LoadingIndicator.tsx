import React from "react";
import { Spin } from "antd";

export function LoadingIndicator() {
  return <Spin tip="Sidan laddas" size="large"></Spin>;
}
