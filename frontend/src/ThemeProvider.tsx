import React from "react";
import { ConfigProvider } from "antd";

interface Props {
  children: JSX.Element;
}

export function ThemeProvider({ children }: Props) {
  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: "#00B96B",
        },
      }}
    >
      {children}
    </ConfigProvider>
  );
}
