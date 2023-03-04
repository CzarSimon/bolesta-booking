import React, { useState } from "react";
import { Button } from "antd";
import { DeleteOutlined } from "@ant-design/icons";

interface Props {
  onConfirm: () => void;
}

export function DeleteButton({ onConfirm }: Props) {
  const [clicked, setClicked] = useState<boolean>(false);
  const onClick = () => {
    if (clicked) {
      onConfirm();
    }
    setClicked(true);
  };

  return (
    <Button onClick={onClick} type="text" danger={clicked}>
      {clicked ? "Ta bort bokning" : <DeleteOutlined />}
    </Button>
  );
}
