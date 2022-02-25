import React from "react";
import { Story } from "@storybook/react";
import JoinUsSection from "./JoinUs";

interface JoinUsProps {}

export default {
  title: "Pages",
  component: JoinUsSection,
};

export const JoinUs: Story<JoinUsProps> = ({ ...props }) => (
  <JoinUsSection {...props} />
);

JoinUs.bind({});
