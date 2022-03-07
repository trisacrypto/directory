import React from "react";
import { Story } from "@storybook/react";
import JoinUsSection from "./JoinUs";

interface JoinUsProps {}

export default {
  title: "Components/JoinUs",
  component: JoinUsSection,
};

export const Default: Story<JoinUsProps> = ({ ...props }) => (
  <JoinUsSection {...props} />
);

Default.bind({});
