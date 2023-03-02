import React from "react";
import { Select } from "antd";

export interface Option {
  label: string;
  value: string;
}

export interface Item {
  id: string;
  name: string;
}

interface Props {
  placeholder?: string;
  items: Item[];
  onChange: (id: string) => void;
  anyOption?: Option;
  large?: boolean;
}

export function ItemSelect({
  placeholder,
  items,
  onChange,
  anyOption,
  large,
}: Props) {
  return (
    <Select
      placeholder={placeholder}
      defaultValue={anyOption?.value}
      options={mapItemsToOptions(items, anyOption)}
      onChange={onChange}
      size={large ? "large" : "middle"}
      allowClear
    />
  );
}

function mapItemsToOptions(items: Item[], anyOption?: Option): Option[] {
  const optionsMap = items.map((i) => ({
    label: i.name,
    value: i.id,
  }));

  return anyOption ? [anyOption, ...optionsMap] : optionsMap;
}
