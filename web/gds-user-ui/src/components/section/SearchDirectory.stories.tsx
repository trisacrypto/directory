import React from "react";
import { Story } from "@storybook/react";
import SearchDirectory from "./SearchDirectory";

interface SearchDirectoryProps {}

export default {
  title: "Components/SearchDirectory",
  component: SearchDirectory,
};

export const standard: Story<SearchDirectoryProps> = ({ ...props }) => (
  <SearchDirectory {...props} />
);

standard.bind({});
